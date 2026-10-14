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
	return &PracticeService{db: db, runnerURL: strings.TrimRight(runnerURL, "/"), runnerToken: runnerToken, client: &http.Client{Timeout: 70 * time.Second}}
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
	if language != "javascript" && language != "go" {
		return nil, errors.New("bahasa hanya mendukung JavaScript dan Go")
	}
	inputs := []struct{ input, expected string }{
		{"  Rani@Example.COM  ", "rani@example.com\n"},
		{"BUDI+HR@CORP.ID", "budi+hr@corp.id\n"},
		{"  Support@LearnHub.ID\nADMIN@LEARNHUB.ID  ", "support@learnhub.id\nadmin@learnhub.id\n"},
	}
	var stdin, expected strings.Builder
	for _, tc := range inputs {
		stdin.WriteString(tc.input)
		if !strings.HasSuffix(tc.input, "\n") {
			stdin.WriteByte('\n')
		}
		expected.WriteString(tc.expected)
	}
	actual, status, err := s.submit(ctx, source, language, stdin.String())
	if err != nil {
		return nil, err
	}
	result := &RunResult{Total: len(inputs), Status: "Belum lolos semua test"}
	remaining := actual
	if status == "Accepted" {
		for _, tc := range inputs {
			if !strings.HasPrefix(remaining, tc.expected) {
				break
			}
			result.Passed++
			remaining = strings.TrimPrefix(remaining, tc.expected)
		}
	}
	if status == "Accepted" && actual == expected.String() {
		result.Status = "Semua test lolos"
	} else {
		result.Output = statusMessage(status, actual)
	}
	return result, nil
}

type runnerResult struct {
	Status string            `json:"status"`
	Error  string            `json:"error"`
	Files  map[string]string `json:"files"`
}

func (s *PracticeService) submit(ctx context.Context, source, language, stdin string) (string, string, error) {
	filename, args := "main.js", []string{"/usr/bin/node", "main.js"}
	if language == "go" {
		filename, args = "main.go", []string{"/bin/sh", "-c", "go build -p 2 -o /tmp/main main.go && /tmp/main"}
	}
	payload, _ := json.Marshal(map[string]any{"cmd": []any{map[string]any{
		"args":     args,
		"env":      []string{"PATH=/usr/local/go/bin:/usr/bin:/bin", "GOCACHE=/tmp/.gocache", "GOTMPDIR=/tmp", "GOMAXPROCS=2"},
		"files":    []any{map[string]string{"content": stdin}, map[string]any{"name": "stdout", "max": 8192}, map[string]any{"name": "stderr", "max": 8192}},
		"cpuLimit": 30_000_000_000, "clockLimit": 60_000_000_000, "memoryLimit": 536_870_912, "procLimit": 50,
		"copyIn": map[string]any{filename: map[string]string{"content": source}},
	}}})
	var results []runnerResult
	if err := s.request(ctx, http.MethodPost, s.runnerURL+"/run", payload, &results); err != nil {
		return "", "", err
	}
	if len(results) != 1 {
		return "", "", errors.New("runner mengembalikan hasil yang tidak valid")
	}
	res := results[0]
	if res.Error != "" || res.Status == "Internal Error" {
		return "", "", fmt.Errorf("runner gagal menjalankan sandbox: %s", firstNonEmpty(res.Error, res.Status))
	}
	return firstNonEmpty(res.Files["stdout"], res.Files["stderr"]), res.Status, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (s *PracticeService) request(ctx context.Context, method, url string, body []byte, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.runnerToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.runnerToken)
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

func statusMessage(status, output string) string {
	switch status {
	case "Accepted":
		return "Output belum sesuai expected result. " + output
	case "Nonzero Exit Status":
		return "Program gagal. Periksa kode dan pesan error. " + output
	case "Memory Limit Exceeded":
		return "Program kehabisan memori."
	case "Time Limit Exceeded":
		return "Program gagal saat dijalankan atau melewati batas waktu."
	default:
		return "Program belum lolos (" + status + "). Periksa kembali hasilnya. " + output
	}
}
