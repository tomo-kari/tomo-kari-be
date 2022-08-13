package integration_test

import (
	. "github.com/Eun/go-hit"
	"net/http"
	"testing"

	"tomokari/pkg/rabbitmq/rmq_rpc/client"
)

// HTTP POST: /translation/do-translate.
func TestHTTPDoTranslate(t *testing.T) {
	body := `{
		"destination": "en",
		"original": "текст для перевода",
		"source": "auto"
	}`
	Test(t,
		Description("DoTranslate Success"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".translation").Equal("text to translate"),
	)

	body = `{
		"destination": "en",
		"original": "текст для перевода"
	}`
	Test(t,
		Description("DoTranslate Fail"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".message").Equal("invalid request body"),
	)
}

// HTTP GET: /translation/history.
func TestHTTPHistory(t *testing.T) {
	Test(t,
		Description("History Success"),
		Get(basePath+"/translation/history"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"history":[{`),
	)
}

// RabbitMQ RPC Client: getHistory.
func TestRMQClientRPC(t *testing.T) {
	rmqClient, err := client.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatal("RabbitMQ RPC Client - init error - client.New")
	}

	defer func() {
		err = rmqClient.Shutdown()
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - shutdown error - rmqClient.RemoteCall", err)
		}
	}()

	type Translation struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Original    string `json:"original"`
		Translation string `json:"translation"`
	}

	type historyResponse struct {
		History []Translation `json:"history"`
	}

	for i := 0; i < requests; i++ {
		var history historyResponse

		err = rmqClient.RemoteCall("getHistory", nil, &history)
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall", err)
		}

		if history.History[0].Original != "текст для перевода" {
			t.Fatal("Original != текст для перевода")
		}
	}
}
