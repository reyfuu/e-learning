package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"learn-platform/backend/internal/model"
	"learn-platform/backend/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c echo.Context) error {
	limit := 10
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	users, total, err := h.service.GetAll(c.Request().Context(), limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch users")
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToResponse()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   userResponses,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID := c.Get("user_id").(string)
	user, err := h.service.GetByID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("user_id").(string)

	// Security: Users can only update their own profile (unless admin)
	userRole := c.Get("user_role").(string)
	if userID != id && userRole != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "unauthorized")
	}

	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	userRole := c.Get("user_role").(string)

	// Security: Only admin can delete users
	if userRole != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "only admin can delete users")
	}

	err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user deleted successfully",
	})
}
