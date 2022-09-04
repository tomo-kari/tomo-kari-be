package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tomokari/internal/entity"
)

func TestRegister(t *testing.T) {
	//_ = refreshUserTable()
	_ = deleteData((&entity.User{}).Table())
	seedOneTos()

	samples := []struct {
		inputJson  entity.CreateUserRequestBody
		statusCode int
		message    string
	}{
		{
			inputJson: entity.CreateUserRequestBody{
				BasicInfo: entity.BasicInfo{
					Email: "user",
					Phone: "",
				},
				DateOfBirth:      "",
				Password:         "",
				TermsOfServiceId: 0,
			},
			statusCode: http.StatusBadRequest,
			message:    "Error:Field validation",
		},
		{
			inputJson: entity.CreateUserRequestBody{
				BasicInfo: entity.BasicInfo{
					Email: "user1@gmail.com",
					Phone: "user@1234",
				},
				DateOfBirth:      "2022-08-23T10:41:20.513Z",
				Password:         "pass12",
				TermsOfServiceId: 1,
			},
			statusCode: http.StatusCreated,
			message:    "",
		},
	}

	for _, sample := range samples {
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(sample.inputJson)
		req, err := http.NewRequest("POST", "/register", &buf)
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		recorder := httptest.NewRecorder()
		handler := gin.New()
		handler.POST("/register", userRoute.Register)
		handler.ServeHTTP(recorder, req)

		body := recorder.Body
		var res map[string]interface{}
		_ = json.Unmarshal(body.Bytes(), &res)
		assert.Equal(t, strings.Contains(res["message"].(string), sample.message), true)
		assert.Equal(t, recorder.Code, sample.statusCode)
	}
}

func TestLogin(t *testing.T) {
	samples := []struct {
		inputJson  entity.LoginUserRequestBody
		statusCode int
		message    string
	}{
		{
			inputJson: entity.LoginUserRequestBody{
				Email:    "user1@gmail.com",
				Password: "pass12",
			},
			statusCode: http.StatusOK,
			message:    "",
		},
	}
	for _, sample := range samples {
		var buf bytes.Buffer
		_ = json.NewEncoder(&buf).Encode(sample.inputJson)
		req, err := http.NewRequest("POST", "/login", &buf)
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		recorder := httptest.NewRecorder()
		handler := gin.New()
		handler.POST("/login", userRoute.Login)
		handler.ServeHTTP(recorder, req)

		body := recorder.Body
		var res map[string]interface{}
		_ = json.Unmarshal(body.Bytes(), &res)
		assert.Equal(t, strings.Contains(res["message"].(string), sample.message), true)
		assert.Equal(t, recorder.Code, sample.statusCode)
	}
}
