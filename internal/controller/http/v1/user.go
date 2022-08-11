package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tomokari/internal/entity"
	"tomokari/internal/usecase"
	"tomokari/pkg/logger"
)

type userRoutes struct {
	t usecase.User
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

type registerResponse struct {
	History []entity.Translation `json:"history"`
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
	err := r.t.Register(c.Request.Context(), body)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, registerResponse{})
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *userRoutes) login(c *gin.Context) {
	var body entity.LoginUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	user, err := r.t.Login(c.Request.Context(), body)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")

		return
	}

	c.JSON(http.StatusOK, user)
}
