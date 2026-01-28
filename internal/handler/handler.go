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

// CreateUser создает нового пользователя
// @Summary      📝 Регистрация нового пользователя
// @Description  Создает нового пользователя в системе. Логин должен быть уникальным.
// @Description  **Пример запроса curl:**
// @Description  ```bash
// @Description  curl -X POST "http://localhost:8080/api/users" \
// @Description    -H "Content-Type: application/json" \
// @Description    -d '{
// @Description      "first_name": "Иван",
// @Description      "last_name": "Петров",
// @Description      "login": "ivan@example.com",
// @Description      "password": "Password123!"
// @Description    }'
// @Description  ```
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        request body models.CreateUserRequest true "Данные нового пользователя"
// @Success      201  {object}  models.ServerResponse{data=models.User}  "Пользователь успешно создан"
// @Failure      400  {object}  models.ServerResponse  "Неверный формат данных"
// @Failure      409  {object}  models.ServerResponse  "Пользователь с таким логином уже существует"
// @Failure      500  {object}  models.ServerResponse  "Внутренняя ошибка сервера"
// @Router       /users [post]
// @Example request
//
//	{
//	  "first_name": "Иван",
//	  "last_name": "Петров",
//	  "login": "ivan@example.com",
//	  "password": "Password123!"
//	}
//
// @Example response 201
//
//	{
//	  "status": "success",
//	  "data": {
//	    "user_id": 1,
//	    "first_name": "Иван",
//	    "last_name": "Петров",
//	    "login": "ivan@example.com",
//	    "last_login": "2026-01-28T12:30:45+05:00"
//	  }
//	}
//
// @Example response 409
//
//	{
//	  "status": "error",
//	  "error": "Пользователь с таким логином уже существует"
//	}
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

// GetUser возвращает пользователя по ID
// @Summary      🔍 Получить пользователя по ID
// @Description  Возвращает информацию о конкретном пользователе. Требует JWT аутентификации.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID пользователя" Format(int64)
// @Success      200  {object}  models.ServerResponse{data=models.User}  "Данные пользователя"
// @Failure      400  {object}  models.ServerResponse  "Неверный ID"
// @Failure      401  {object}  models.ServerResponse  "Требуется аутентификация"
// @Failure      404  {object}  models.ServerResponse  "Пользователь не найден"
// @Router       /users/{id} [get]
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

// GetAllUsers возвращает всех пользователей
// @Summary      👥 Получить всех пользователей
// @Description  Возвращает список всех зарегистрированных пользователей. Требует JWT аутентификации.
// @Description
// @Description  **Пример запроса curl:**
// @Description  ```bash
// @Description  curl -X GET "http://localhost:8080/api/users" \
// @Description    -H "Authorization: Bearer $TOKEN"
// @Description  ```
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.ServerResponse{data=[]models.User}  "Список пользователей"
// @Failure      401  {object}  models.ServerResponse  "Требуется аутентификация"
// @Failure      500  {object}  models.ServerResponse  "Внутренняя ошибка сервера"
// @Router       /users [get]
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

// UpdateUser обновляет данные пользователя
// @Summary      ✏️ Обновить пользователя
// @Description  Обновляет информацию о пользователе. Все поля опциональны. Требует JWT аутентификации.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID пользователя" Format(int64)
// @Param        request body models.UpdateUserRequest true "Новые данные пользователя"
// @Success      200  {object}  models.ServerResponse{data=models.User}  "Пользователь обновлен"
// @Failure      400  {object}  models.ServerResponse  "Неверные данные"
// @Failure      401  {object}  models.ServerResponse  "Требуется аутентификация"
// @Failure      404  {object}  models.ServerResponse  "Пользователь не найден"
// @Failure      409  {object}  models.ServerResponse  "Логин уже занят"
// @Router       /users/{id} [put]
// @Example request
//
//	{
//	  "first_name": "Александр",
//	  "last_name": "Иванов",
//	  "login": "alex@example.com",
//	  "password": "NewPassword456!"
//	}
//
// @Example response 200
//
//	{
//	  "status": "success",
//	  "data": {
//	    "user_id": 1,
//	    "first_name": "Александр",
//	    "last_name": "Иванов",
//	    "login": "alex@example.com",
//	    "last_login": "2026-01-28T12:30:45+05:00"
//	  }
//	}
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
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.ServerResponse{
				Status: "error",
				Error:  "User not found",
			})
			return
		}
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

// DeleteUser удаляет пользователя
// @Summary      🗑️ Удалить пользователя
// @Description  Удаляет пользователя из системы. Требует JWT аутентификации.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID пользователя" Format(int64)
// @Success      200  {object}  models.ServerResponse  "Пользователь удален"
// @Failure      400  {object}  models.ServerResponse  "Неверный ID"
// @Failure      401  {object}  models.ServerResponse  "Требуется аутентификация"
// @Failure      404  {object}  models.ServerResponse  "Пользователь не найден"
// @Router       /users/{id} [delete]
// @Example response 200
//
//	{
//	  "status": "success",
//	  "message": "User deleted successfully"
//	}
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
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, models.ServerResponse{
				Status: "error",
				Error:  "User not found",
			})
			return
		}
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

// Login аутентифицирует пользователя
// @Summary      🔐 Вход в систему
// @Description  Аутентифицирует пользователя по логину и паролю, возвращает JWT токен
// @Tags         Аутентификация
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Учетные данные"
// @Success      200  {object}  models.ServerResponse{data=models.LoginResponse}  "Успешный вход"
// @Failure      400  {object}  models.ServerResponse  "Неверный формат данных"
// @Failure      401  {object}  models.ServerResponse  "Неверные учетные данные"
// @Router       /users/login [post]
// @Example request
//
//	{
//	  "login": "ivan@example.com",
//	  "password": "Password123!"
//	}
//
// @Example response 200
//
//	{
//	  "status": "success",
//	  "data": {
//	    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "user": {
//	      "user_id": 1,
//	      "first_name": "Иван",
//	      "last_name": "Петров",
//	      "login": "ivan@example.com",
//	      "last_login": "2026-01-28T12:30:45+05:00"
//	    }
//	  }
//	}
//
// @Example curl
//
//	curl -X POST "http://localhost:8080/api/users/login" \
//	  -H "Content-Type: application/json" \
//	  -d '{
//	    "login": "ivan@example.com",
//	    "password": "Password123!"
//	  }'
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
