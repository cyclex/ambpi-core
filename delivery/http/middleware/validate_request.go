package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cyclex/ambpi-core/api"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ReqLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")

		// Read the content
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
		}

		var req api.Login
		err := json.Unmarshal(bodyBytes, &req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		// Restore the io.ReadCloser to its original state
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return next(c)
	}
}

func ReqReport(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")

		// Read the content
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		}

		var req api.Report
		err := json.Unmarshal(bodyBytes, &req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		// Restore the io.ReadCloser to its original state
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return next(c)
	}
}

func ReqSendPush(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")

		// Read the content
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		}

		var req api.SendPushNotif
		err := json.Unmarshal(bodyBytes, &req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		// Restore the io.ReadCloser to its original state
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return next(c)
	}
}

func ReqSetProgram(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")

		// Read the content
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		}

		var req api.SetProgram
		err := json.Unmarshal(bodyBytes, &req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		// Restore the io.ReadCloser to its original state
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return next(c)
	}
}

func ReqValidateRedeem(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/json")

		// Read the content
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		}

		var req api.ValidateRedeem
		err := json.Unmarshal(bodyBytes, &req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.ResponseError{
				Status:  false,
				Message: "Parameter error",
			})
		}

		// Restore the io.ReadCloser to its original state
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return next(c)
	}
}
