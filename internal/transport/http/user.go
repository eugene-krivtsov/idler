package http

import (
	"github.com/eugene-krivtsov/idler/internal/model/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	user := api.Group("/user")
	{
		user.POST("/sign-up", h.signUp)
		user.POST("/sign-in", h.signIn)
		user.GET("/verify", h.getVerifyCode)
		user.POST("/verify/:code", h.verifyUser)

		authenticated := user.Group("/", h.userIdentity)
		{
			authenticated.GET("/profile", h.getUserProfile)
			authenticated.GET("/", h.getAllUsers)
			authenticated.GET("/:email", h.getUserByEmail)
		}
	}
}

// @Summary		User SignUp
// @Tags        user-auth
// @Description registering new user
// @ID          userSignUp
// @Accept      json
// @Produce     json
// @Param       input body dto.SignUpDTO  true  "account information"
// @Success     201
// @Failure     400  	  		{object}  errorResponse
// @Failure     500      		{object}  errorResponse
// @Failure     default  		{object}  errorResponse
// @Router      /user/sign-up 	[post]
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

// @Summary 	User SignIn
// @Tags 		user-auth
// @Description authentication new user
// @ID 			userSignIn
// @Accept  	json
// @Produce  	json
// @Param input body dto.SignInDTO true "credentials"
// @Success 	200 {string} string "token"
// @Failure 	400,404 {object} 	errorResponse
// @Failure 	500 {object} 		errorResponse
// @Failure 	default {object} 	errorResponse
// @Router 		/user/sign-in 		[post]
func (h *Handler) signIn(c *gin.Context) {
	var signInDTO dto.SignInDTO
	if err := c.BindJSON(&signInDTO); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.userService.SignIn(c.Request.Context(), signInDTO)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}

// @Summary 	Verify User
// @Security 	Bearer
// @Tags 		user
// @Description verifies user email
// @ID 			verifyUser
// @Accept  	json
// @Produce  	json
// @Success 	200
// @Failure 	400 {object} 		errorResponse
// @Failure 	500 {object} 		errorResponse
// @Failure 	default {object} 	errorResponse
// @Router 		/user/verify 		[post]
func (h *Handler) verifyUser(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newErrorResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	id, err := uuid.Parse(code)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "unsupported code format")
	}

	if err := h.userService.VerifyUser(c, id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Status(http.StatusOK)
}

// @Summary 	Get Verify Code
// @Security 	Bearer
// @Tags 		user
// @Description gets verify code to email
// @ID 			getVerifyCode
// @Accept  	json
// @Produce  	json
// @Success 	200
// @Failure 	400 {object} 		errorResponse
// @Failure 	500 {object} 		errorResponse
// @Failure 	default {object} 	errorResponse
// @Router 		/user/verify 		[get]
func (h *Handler) getVerifyCode(c *gin.Context) {

}

// @Summary 	User Profile
// @Security 	Bearer
// @Tags 		user
// @Description gets current profile user
// @ID 			currentUser
// @Accept  	json
// @Produce  	json
// @Success 	200 {object} 		dto.UserDTO
// @Failure 	500 {object} 		errorResponse
// @Failure 	default {object} 	errorResponse
// @Router 		/user/profile 		[get]
func (h *Handler) getUserProfile(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.userService.GetById(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary 	Users
// @Security 	Bearer
// @Tags 		user
// @Description gets all users
// @ID 			allUsers
// @Accept  	json
// @Produce  	json
// @Success 	200 {array} 		dto.UserDTO
// @Failure 	500 {object} 		errorResponse
// @Failure 	default {object} 	errorResponse
// @Router 		/user 				[get]
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

// @Summary User By Email
// @Security Bearer
// @Tags user
// @Description gets user by email
// @ID userByEmail
// @Accept  json
// @Produce  json
// @Param email path string true "User email"
// @Success 200 {object} dto.UserDTO
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/{email} [get]
func (h *Handler) getUserByEmail(c *gin.Context) {
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
