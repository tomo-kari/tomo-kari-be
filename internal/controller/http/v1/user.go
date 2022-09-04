package v1

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"tomokari/internal/constants"
	"tomokari/pkg/postgres"

	"github.com/gin-gonic/gin"

	"tomokari/internal/entity"
	"tomokari/internal/usecase"
	"tomokari/pkg/logger"
)

type UserRoutes struct {
	U usecase.User
	L logger.Interface
	V *validator.Validate
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.User, l logger.Interface, v *validator.Validate) {
	r := &UserRoutes{u, l, v}

	h := handler.Group("/auth")
	{
		h.POST("/Register", r.Register)
		h.POST("/Login", r.Login)
		h.POST("/candidates", r.getCandidates)
	}
}

// Register @Summary     Register
// @Description Register with email and password
// @ID          do-Register
// @Tags  	    registration
// @Accept      json
// @Produce     json
// @Success     200 {object} registerResponse
// @Failure     500 {object} response
// @Router      /auth/Register [post]
func (ur *UserRoutes) Register(c *gin.Context) {
	var body entity.CreateUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ur.L.Error(err, "http - v1 - Register")
		errorResponse(c, http.StatusBadRequest, constants.UserInfoErrorMessage)
		return
	}
	err := ur.V.Struct(body)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	authUser, status, err := ur.U.Register(c.Request.Context(), body)
	if err != nil {
		ur.L.Error(err, "http - v1 - Register")
		errMsg := err.Error()
		errorResponse(c, status, errMsg)
		return
	}

	responseWithData(c, status, authUser, "")
}

// Login @Summary     Login
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
func (ur *UserRoutes) Login(c *gin.Context) {
	var body entity.LoginUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ur.L.Error(err, "http - v1 - Login")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	authUser, status, err := ur.U.Login(c.Request.Context(), body)
	if err != nil {
		ur.L.Error(err, "http - v1 - doTranslate")
		errorResponse(c, status, err.Error())
		return
	}

	responseWithData(c, status, authUser, "")
}

func (ur *UserRoutes) getCandidates(c *gin.Context) {
	var filter postgres.GetManyRequestBody
	if err := c.ShouldBindJSON(&filter); err != nil {
		ur.L.Error(err, "http - v1 - getCandidates")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	candidates, err := ur.U.GetCandidates(c.Request.Context(), filter)
	if err != nil {
		ur.L.Error(err, "http - v1 - getCandidates")
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if candidates == nil {
		candidates = []entity.Map{}
	}
	responseWithData(c, http.StatusOK, candidates, "")
}
