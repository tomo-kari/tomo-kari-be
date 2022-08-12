package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tomokari/internal/entity"
	"tomokari/internal/usecase"
	"tomokari/pkg/logger"
)

type userRoutes struct {
	u usecase.User
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.User, l logger.Interface) {
	r := &userRoutes{t, l}

	h := handler.Group("/auth")
	{
		h.GET("/register", r.register)
		h.POST("/login", r.login)
	}
}

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} registerResponse
// @Failure     500 {object} response
// @Router      /translation/history [get]
func (r *userRoutes) register(c *gin.Context) {
	var body entity.CreateUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - register")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}
	authUser, status, err := r.u.Register(c.Request.Context(), body)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		errorResponse(c, status, "database problems")

		return
	}

	responseWithData(c, status, authUser)
}

// @Summary     Login
// @Description Login with email and password
// @ID          do-login
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body entity.LoginUserRequestBody true "Set up translation"
// @Success     200 {object} entity.User
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *userRoutes) login(c *gin.Context) {
	var body entity.LoginUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - login")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	authUser, status, err := r.u.Login(c.Request.Context(), body)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, status, err.Error())
		return
	}

	responseWithData(c, status, authUser)
}
