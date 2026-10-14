package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"learn-platform/backend/internal/model"
	"learn-platform/backend/internal/service"
)

type CourseHandler struct {
	service *service.CourseService
}

func NewCourseHandler(service *service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) CreateCourse(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	var req model.CreateCourseRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	resp, err := h.service.CreateCourse(c.Request().Context(), userID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *CourseHandler) GetUserCourses(c echo.Context) error {
	userID := c.Get("user_id").(string)

	resp, err := h.service.GetUserCourses(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch courses")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) GetCourse(c echo.Context) error {
	userID := c.Get("user_id").(string)
	courseID := c.Param("id")

	resp, err := h.service.GetCourse(c.Request().Context(), courseID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "course not found")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) GetAllCourses(c echo.Context) error {
	resp, err := h.service.GetAllCourses(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch catalog")
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) EnrollCourse(c echo.Context) error {
	userID := c.Get("user_id").(string)
	courseID := c.Param("id")

	if err := h.service.EnrollCourse(c.Request().Context(), userID, courseID); err != nil {
		if err.Error() == "already enrolled" {
			return echo.NewHTTPError(http.StatusConflict, "already enrolled")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "enrolled successfully"})
}

func (h *CourseHandler) GetEnrolledCourses(c echo.Context) error {
	userID := c.Get("user_id").(string)
	resp, err := h.service.GetEnrolledCourses(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch enrolled courses")
	}
	return c.JSON(http.StatusOK, resp)
}


func (h *CourseHandler) CreateModule(c echo.Context) error {
	userID := c.Get("user_id").(string)
	courseID := c.Param("id")

	var req model.CreateModuleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	resp, err := h.service.CreateModule(c.Request().Context(), userID, courseID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *CourseHandler) ToggleModuleCompletion(c echo.Context) error {
	userID := c.Get("user_id").(string)
	moduleID := c.Param("id")

	resp, err := h.service.ToggleModuleCompletion(c.Request().Context(), userID, moduleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) UpdateCourseStatus(c echo.Context) error {
	userID := c.Get("user_id").(string)
	courseID := c.Param("id")

	var req model.UpdateCourseStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	resp, err := h.service.UpdateCourseStatus(c.Request().Context(), userID, courseID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
