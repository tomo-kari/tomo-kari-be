package v1

import (
	"encoding/json"
	"net/http"
	"tomokari/internal/constants"
	"tomokari/pkg/postgres"

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
		h.POST("/register", gin.WrapF(r.register))
		h.POST("/login", r.login)
		h.POST("/candidates", r.getCandidates)
	}
}

// @Summary     Register
// @Description Register with email and password
// @ID          do-register
// @Tags  	    registration
// @Accept      json
// @Produce     json
// @Success     200 {object} registerResponse
// @Failure     500 {object} response
// @Router      /auth/register [post]
func (ur *userRoutes) register(res http.ResponseWriter, req *http.Request) {
	var body entity.CreateUserRequestBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		ur.l.Error(err, "http - v1 - register")
		errorResponse2(res, http.StatusBadRequest, constants.UserInfoErrorMessage)
		return
	}
	authUser, status, err := ur.u.Register(req.Context(), body)
	if err != nil {
		ur.l.Error(err, "http - v1 - register")
		errMsg := err.Error()
		errorResponse2(res, status, errMsg)
		return
	}

	responseWithData2(res, status, authUser, "")
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
func (ur *userRoutes) login(c *gin.Context) {
	var body entity.LoginUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		ur.l.Error(err, "http - v1 - login")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	authUser, status, err := ur.u.Login(c.Request.Context(), body)
	if err != nil {
		ur.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, status, err.Error())
		return
	}

	responseWithData(c, status, authUser, "")
}

func (ur *userRoutes) getCandidates(c *gin.Context) {
	var filter postgres.GetManyRequestBody
	if err := c.ShouldBindJSON(&filter); err != nil {
		ur.l.Error(err, "http - v1 - getCandidates")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	candidates, err := ur.u.GetCandidates(c.Request.Context(), filter)
	if err != nil {
		ur.l.Error(err, "http - v1 - getCandidates")
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if candidates == nil {
		candidates = []entity.Map{}
	}
	responseWithData(c, http.StatusOK, candidates, "")
}
