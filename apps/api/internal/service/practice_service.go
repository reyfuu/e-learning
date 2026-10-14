package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learn-platform/backend/internal/model"
)

const exerciseID = "employee-email-normalizer"

var ErrHintsExhausted = errors.New("semua hint sudah digunakan")

type PracticeService struct {
	db          *gorm.DB
	runnerURL   string
	runnerToken string
	client      *http.Client
}

type Exercise struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Scenario      string            `json:"scenario"`
	Prompt        string            `json:"prompt"`
	Starter       map[string]string `json:"starter"`
	Hints         []string          `json:"hints,omitempty"`
	HintsUnlocked int               `json:"hints_unlocked"`
}

type RunResult struct {
	Passed int    `json:"passed"`
	Total  int    `json:"total"`
	Status string `json:"status"`
	Output string `json:"output,omitempty"`
}

func NewPracticeService(db *gorm.DB, runnerURL, runnerToken string) *PracticeService {
	return &PracticeService{db: db, runnerURL: strings.TrimRight(runnerURL, "/"), runnerToken: runnerToken, client: &http.Client{Timeout: 15 * time.Second}}
}

var emailExercise = Exercise{
	ID:       exerciseID,
	Title:    "Normalisasi email karyawan",
	Scenario: "Tim HR menerima data karyawan dari beberapa sistem. Sebelum disimpan, backend harus menormalkan alamat email agar pencarian akun konsisten.",
	Prompt:   "Buat program yang membaca satu email per baris dari stdin dan mencetak versi yang sudah di-trim dan huruf kecil. Pertahankan satu baris output untuk setiap input. Contoh:  Rani@Example.COM  menjadi rani@example.com.",
	Starter: map[string]string{
		"javascript": "const fs = require('fs');\nconst lines = fs.readFileSync(0, 'utf8').split(/\\r?\\n/);\nfor (const line of lines) {\n  if (line.length === 0) continue;\n  // TODO: trim spasi lalu ubah menjadi huruf kecil\n  console.log(line);\n}",
		"go":         "package main\n\nimport (\n  \"bufio\"\n  \"fmt\"\n  \"os\"\n)\n\nfunc main() {\n  scanner := bufio.NewScanner(os.Stdin)\n  for scanner.Scan() {\n    // TODO: trim spasi lalu ubah menjadi huruf kecil\n    fmt.Println(scanner.Text())\n  }\n}",
	},
	Hints: []string{
		"Hilangkan spasi di awal dan akhir string sebelum melakukan transformasi lain.",
		"Di JavaScript, string memiliki method trim() dan toLowerCase().",
		"Di Go, gunakan strings.TrimSpace dan strings.ToLower.",
		"Pastikan setiap baris diproses sendiri; jangan menggabungkan alamat email yang berbeda.",
		"Urutannya: trim whitespace, lowercase, lalu tulis satu hasil per baris.",
	},
}

func (s *PracticeService) GetExercise(ctx context.Context, userID string) (*Exercise, error) {
	var progress model.PracticeProgress
	err := s.db.WithContext(ctx).First(&progress, "user_id = ? AND exercise_id = ?", userID, exerciseID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	out := emailExercise
	unlocked := max(0, min(progress.HintsUnlocked, len(emailExercise.Hints)))
	out.Hints = append([]string(nil), emailExercise.Hints[:unlocked]...)
	out.HintsUnlocked = unlocked
	return &out, nil
}

func (s *PracticeService) UnlockHint(ctx context.Context, userID string) (string, int, error) {
	row := model.PracticeProgress{UserID: userID, ExerciseID: exerciseID}
	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
		return "", 0, err
	}
	updated := s.db.WithContext(ctx).Model(&model.PracticeProgress{}).
		Where("user_id = ? AND exercise_id = ? AND hints_unlocked < ?", userID, exerciseID, len(emailExercise.Hints)).
		UpdateColumn("hints_unlocked", gorm.Expr("hints_unlocked + 1"))
	if updated.Error != nil {
		return "", 0, updated.Error
	}
	if updated.RowsAffected == 0 {
		return "", len(emailExercise.Hints), ErrHintsExhausted
	}
	if err := s.db.WithContext(ctx).First(&row, "user_id = ? AND exercise_id = ?", userID, exerciseID).Error; err != nil {
		return "", 0, err
	}
	return emailExercise.Hints[row.HintsUnlocked-1], row.HintsUnlocked, nil
}

func (s *PracticeService) Run(ctx context.Context, language, source string) (*RunResult, error) {
	if s.runnerURL == "" {
		return nil, errors.New("code runner belum dikonfigurasi; atur CODE_RUNNER_URL")
	}
	if len(source) == 0 || len(source) > 16000 {
		return nil, errors.New("source code wajib diisi dan maksimal 16 KB")
	}
	languageID := map[string]int{"javascript": 63, "go": 60}[language]
	if languageID == 0 {
		return nil, errors.New("bahasa hanya mendukung JavaScript dan Go")
	}
	inputs := []struct{ input, expected string }{
		{"  Rani@Example.COM  ", "rani@example.com\n"},
		{"BUDI+HR@CORP.ID", "budi+hr@corp.id\n"},
		{"  Support@LearnHub.ID\nADMIN@LEARNHUB.ID  ", "support@learnhub.id\nadmin@learnhub.id\n"},
	}
	result := &RunResult{Total: len(inputs), Status: "Belum lolos semua test"}
	for _, tc := range inputs {
		actual, status, err := s.submit(ctx, source, languageID, tc.input, tc.expected)
		if err != nil {
			return nil, err
		}
		if status == 3 && actual == tc.expected {
			result.Passed++
		}
		if result.Passed < result.Total {
			result.Output = statusMessage(status, actual)
		}
	}
	if result.Passed == result.Total {
		result.Status = "Semua test lolos"
		result.Output = ""
	}
	return result, nil
}

type judgeResult struct {
	Token         string  `json:"token"`
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

func (s *PracticeService) submit(ctx context.Context, source string, languageID int, stdin, expected string) (string, int, error) {
	payload, _ := json.Marshal(map[string]any{"source_code": source, "language_id": languageID, "stdin": stdin, "expected_output": expected, "cpu_time_limit": 2, "wall_time_limit": 5, "memory_limit": 128000, "max_file_size": 1024})
	url := s.runnerURL + "/submissions?base64_encoded=false&wait=true"
	var res judgeResult
	if err := s.request(ctx, http.MethodPost, url, payload, &res); err != nil {
		return "", 0, err
	}
	if res.Status.ID == 0 {
		if res.Token == "" {
			return "", 0, errors.New("runner tidak mengembalikan token submission")
		}
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			time.Sleep(250 * time.Millisecond)
			if err := s.request(ctx, http.MethodGet, s.runnerURL+"/submissions/"+res.Token+"?base64_encoded=false", nil, &res); err != nil {
				return "", 0, err
			}
			if res.Status.ID > 2 {
				break
			}
		}
	}
	actual := ""
	if res.Stdout != nil {
		actual = *res.Stdout
	}
	if res.Status.ID <= 2 {
		return actual, res.Status.ID, errors.New("runner timeout menunggu hasil eksekusi")
	}
	return actual, res.Status.ID, nil
}

func (s *PracticeService) request(ctx context.Context, method, url string, body []byte, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.runnerToken != "" {
		req.Header.Set("X-Auth-Token", s.runnerToken)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("runner tidak dapat dijangkau: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("runner merespons HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 64*1024)).Decode(out)
}

func statusMessage(status int, output string) string {
	switch status {
	case 3:
		return "Output belum sesuai expected result. " + output
	case 4:
		return "Jawaban salah. Periksa whitespace dan huruf besar/kecil."
	case 5:
		return "Program kehabisan memori."
	case 6:
		return "Program gagal dikompilasi."
	case 7, 8, 9, 10, 11, 12:
		return "Program gagal saat dijalankan atau melewati batas waktu."
	default:
		return "Program belum lolos. Periksa kembali hasilnya."
	}
}
