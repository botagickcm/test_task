package handler

import (
	"net/http"
	"strconv"

	"test_task/internal/models"
	"test_task/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  "Invalid request data: " + err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ServerResponse{
		Status: "success",
		Data:   user,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, models.ServerResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   user,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   users,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  "Invalid user ID",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  "Invalid request data: " + err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  "Invalid user ID",
		})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ServerResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Status: "error",
			Error:  "Invalid request data: " + err.Error(),
		})
		return
	}

	response, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ServerResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ServerResponse{
		Status: "success",
		Data:   response,
	})
}
