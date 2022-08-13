package integration_test

import (
	"fmt"
	. "github.com/Eun/go-hit"
	"net/http"
	"testing"
	"time"
)

// HTTP POST: /auth/register.
func TestHTTPRegister(t *testing.T) {
	fakeEmail := fmt.Sprintf("email%d@gmail.com", time.Now().Unix())
	fakePhone := fmt.Sprintf("0383934929%d", time.Now().Unix())
	body := fmt.Sprintf(`{
        "email": "%s",
        "phone": "%s",
        "dateOfBirth": "%s",
        "password": "%s",
        "termsOfServiceId": 1
    }`, fakeEmail, fakePhone, "08/13/2022", "123456")
	Test(t,
		Description("DoRegister Success"),
		Post(basePath+"/auth/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".success").Equal(true),
	)

	body = `{}`
	Test(t,
		Description("DoRegister Fail"),
		Post(basePath+"/auth/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".success").Equal(false),
	)
}
