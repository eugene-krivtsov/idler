package http

import (
	"github.com/eugene-krivtsov/idler/internal/model/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUserRoutes(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.POST("/sign-up", h.signUp)
		user.POST("/sign-in", h.signIn)

		authenticated := user.Group("/", h.userIdentity)
		{
			authenticated.GET("/", h.getAllUsers)
			authenticated.GET("/:mail", h.GetByEmail)
		}
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var signUpDTO dto.SignUpDTO
	if err := c.BindJSON(&signUpDTO); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var err = h.userService.SignUp(c.Request.Context(), signUpDTO)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {
	var signInDTO dto.SignInDTO
	if err := c.BindJSON(&signInDTO); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.userService.SignIn(c.Request.Context(), signInDTO)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

func (h *Handler) getAllUsers(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	users, err := h.userService.GetAll(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetByEmail(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.userService.GetByEmail(c, c.Param("email"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}