package http

import (
	"net/http"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var appLog *logrus.Logger

type OrderHandler struct {
	Ch domain.ChatUcase
}

func NewOrderHandler(e *echo.Echo, chatUcase domain.ChatUcase, debug bool) {

	appLog = pkg.New("app", debug)

	handler := &OrderHandler{
		Ch: chatUcase,
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: pkg.New("middleware", debug).Out,
	}))

	e.POST("/v1/webhooks/whatsapp", handler.webhooksWhatsapp)
	e.GET("/v1/webhooks/whatsapp", handler.webhooksWhatsapp)
}

func (self *OrderHandler) webhooksWhatsapp(c echo.Context) (err error) {

	var (
		request api.CproWebhookPayload
		code    = 400
	)

	defer func(code *int) {
		res := api.ResponseChatbot{
			Code:       *code,
			Message:    http.StatusText(*code),
			ServerTime: time.Now().Local().Unix(),
		}
		c.JSON(*code, res)
	}(&code)

	err = c.Bind(&request)
	if err != nil {
		err = errors.Wrap(err, "[webhooksWhatsapp] Bind")
		appLog.Error(err)
	}

	code = 200
	_, err = self.Ch.IncomingMessages(request.Entry[0].Changes[0].Value.Messages[0])
	if err != nil {
		err = errors.Wrap(err, "[webhooksWhatsapp] IncomingMessages")
		appLog.Error(err)
	}

	return
}
