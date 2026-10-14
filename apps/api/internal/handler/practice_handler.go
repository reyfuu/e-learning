package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"learn-platform/backend/internal/service"
)

type PracticeHandler struct{ service *service.PracticeService }

func NewPracticeHandler(s *service.PracticeService) *PracticeHandler {
	return &PracticeHandler{service: s}
}

func (h *PracticeHandler) GetExercise(c echo.Context) error {
	if c.Param("id") != "employee-email-normalizer" {
		return echo.NewHTTPError(http.StatusNotFound, "latihan tidak ditemukan")
	}
	exercise, err := h.service.GetExercise(c.Request().Context(), c.Get("user_id").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "gagal memuat latihan")
	}
	return c.JSON(http.StatusOK, exercise)
}

func (h *PracticeHandler) UnlockHint(c echo.Context) error {
	if c.Param("id") != "employee-email-normalizer" {
		return echo.NewHTTPError(http.StatusNotFound, "latihan tidak ditemukan")
	}
	hint, unlocked, err := h.service.UnlockHint(c.Request().Context(), c.Get("user_id").(string))
	if err != nil {
		if errors.Is(err, service.ErrHintsExhausted) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "gagal membuka hint")
	}
	return c.JSON(http.StatusOK, map[string]any{"hint": hint, "hints_unlocked": unlocked})
}

func (h *PracticeHandler) Run(c echo.Context) error {
	if c.Param("id") != "employee-email-normalizer" {
		return echo.NewHTTPError(http.StatusNotFound, "latihan tidak ditemukan")
	}
	var req struct {
		Language string `json:"language"`
		Source   string `json:"source"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "request tidak valid")
	}
	if len(req.Source) == 0 || len(req.Source) > 16000 || (req.Language != "go" && req.Language != "javascript") {
		return echo.NewHTTPError(http.StatusBadRequest, "source wajib diisi (maksimal 16 KB) dan bahasa harus go atau javascript")
	}
	if strings.TrimSpace(req.Source) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "source code tidak boleh kosong")
	}
	result, err := h.service.Run(c.Request().Context(), req.Language, req.Source)
	if err != nil {
		if err.Error() == "code runner belum dikonfigurasi; atur CODE_RUNNER_URL" {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
