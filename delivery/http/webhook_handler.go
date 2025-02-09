package http

import (
	"encoding/json"
	"io"
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
	e.GET("/v1/webhooks/whatsapp", handler.health)
}

func (self *OrderHandler) webhooksWhatsapp(c echo.Context) (err error) {
	var (
		request api.CproPayload
		code    = 200
	)

	defer func(code *int) {
		res := api.ResponseChatbot{
			Code:       *code,
			Message:    http.StatusText(*code),
			ServerTime: time.Now().Local().Unix(),
		}
		c.JSON(*code, res)
	}(&code)

	// Read raw body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		err = errors.Wrap(err, "[webhooksWhatsapp] Failed to read request")
		appLog.Error(err)
		return
	}

	// Log raw body to file
	appLog.Debug(string(body))

	// Check if request body is empty
	if len(body) == 0 {
		err = errors.Wrap(err, "[webhooksWhatsapp] Received empty request body")
		appLog.Error(err)
		return
	}

	// Try parsing JSON
	err = json.Unmarshal(body, &request)
	if err != nil {
		err = errors.Wrap(err, "[webhooksWhatsapp] Failed to parse JSON")
		appLog.Error(err)
		return
	}

	var inbound api.CproMessage

	if request.ID == "" {
		err = errors.New("[webhooksWhatsapp] invalid request: missing ID")
		appLog.Debug(err)
		return
	}

	// Ensure Entry exists and is not empty
	if len(request.Data.Entry) == 0 {
		err = errors.New("[webhooksWhatsapp] invalid request: missing Entry array")
		appLog.Debug(err)
		return
	}

	// Ensure Changes exists and is not empty
	if len(request.Data.Entry[0].Changes) == 0 {
		err = errors.New("[webhooksWhatsapp] invalid request: missing Changes array")
		appLog.Debug(err)
		return
	}

	// Ensure Messages exists and is not empty
	if len(request.Data.Entry[0].Changes[0].Value.Messages) == 0 {
		err = errors.New("[webhooksWhatsapp] invalid request: missing Messages array")
		appLog.Debug(err)
		return
	}

	// Extract the message safely
	inbound = request.Data.Entry[0].Changes[0].Value.Messages[0]

	code = 200
	_, err = self.Ch.IncomingMessages(inbound)
	if err != nil {
		err = errors.Wrap(err, "[webhooksWhatsapp] IncomingMessages")
		appLog.Error(err)
	}

	return
}

func (self *OrderHandler) health(c echo.Context) (err error) {

	res := api.ResponseChatbot{
		Code:       200,
		Message:    http.StatusText(200),
		ServerTime: time.Now().Local().Unix(),
	}
	c.JSON(200, res)

	return
}
