package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/pkg"

	_AppMW "github.com/cyclex/ambpi-core/delivery/http/middleware"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var cmsLog *logrus.Logger

type CmsHandler struct {
	CmsGw domain.CmsUcase
}

func NewCmsHandler(e *echo.Echo, gw domain.CmsUcase, debug bool) {

	cmsLog = pkg.New("cms", debug)
	handler := &CmsHandler{
		CmsGw: gw,
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: pkg.New("middleware", false).Out,
	}))

	e.POST("/v1/login", handler.login, _AppMW.ReqLogin)
	e.GET("/v1/checkToken/:token", handler.checkToken)

	e.POST("/v1/report/:type", handler.report)
	e.GET("/v1/redeem/:id", handler.detailHistory)

	e.POST("/v1/push", handler.sendPush, _AppMW.ReqSendPush)

	e.PUT("/v1/campaign/:id", handler.setProgram)

	e.POST("/v1/job/create", handler.createJob)
	e.GET("/v1/job/:type", handler.listJob)

	e.POST("/v1/redeem/validation", handler.validateRedeem, _AppMW.ReqValidateRedeem)

	e.POST("/v1/user", handler.createUser)
	e.DELETE("/v1/user/:id", handler.deleteUser)
	e.PUT("/v1/user/:id", handler.setUser)
	e.PUT("/v1/user/password", handler.setUserPassword)
}

func (self *CmsHandler) checkToken(c echo.Context) error {

	var (
		request api.CheckToken
		res     interface{}
		code    = http.StatusForbidden
		ctx     = c.Request().Context()
	)

	request.Token = c.Param("token")
	err := self.CmsGw.CheckToken(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseError{
			Status: true,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) login(c echo.Context) error {

	var (
		request api.Login
		res     interface{}
	)

	code := http.StatusForbidden

	c.Bind(&request)

	ctx := c.Request().Context()

	data, err := self.CmsGw.Login(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
			Data:    data,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) report(c echo.Context) error {

	var (
		request api.Report
		res     interface{}
		code    = http.StatusInternalServerError
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	data, err := self.CmsGw.Report(ctx, request, c.Param("type"))
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
			Data:    data,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) sendPush(c echo.Context) error {

	var (
		request api.SendPushNotif
		res     interface{}
	)

	code := http.StatusInternalServerError

	c.Bind(&request)

	ctx := c.Request().Context()

	err := self.CmsGw.SendPushNotif(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) setProgram(c echo.Context) error {

	var (
		request api.SetProgram
		res     interface{}
	)

	code := http.StatusInternalServerError

	c.Bind(&request)

	ctx := c.Request().Context()

	cond := map[string]interface{}{"id": c.Param("id")}
	data := map[string]interface{}{"updated_at": time.Now().Local(), "start_date": request.StartDate, "end_date": request.EndDate}

	err := self.CmsGw.SetProgram(ctx, cond, data)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) createJob(c echo.Context) error {

	var (
		request api.Job
		res     interface{}
	)

	code := http.StatusInternalServerError

	ctx := c.Request().Context()
	c.Bind(&request)

	err := self.CmsGw.CreateJob(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseReport{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) importPrize(c echo.Context) error {

	var (
		request api.Job
		res     interface{}
	)

	code := http.StatusInternalServerError

	ctx := c.Request().Context()
	c.Bind(&request)

	_, _, err := self.CmsGw.ImportPrize(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseReport{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) createUser(c echo.Context) error {

	var (
		request api.User
		res     interface{}
		code    = http.StatusInternalServerError
		ctx     = c.Request().Context()
	)

	c.Bind(&request)
	err := self.CmsGw.CreateUser(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: "Invalid account",
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) deleteUser(c echo.Context) error {

	var (
		res  interface{}
		code = http.StatusInternalServerError
		ctx  = c.Request().Context()
	)

	n, _ := strconv.Atoi(c.Param("id"))
	x := int64(n)

	err := self.CmsGw.DeleteUser(ctx, x)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) setUser(c echo.Context) error {

	var (
		request api.User
		res     interface{}
		ctx     = c.Request().Context()
		code    = http.StatusInternalServerError
	)

	c.Bind(&request)
	n, _ := strconv.Atoi(c.Param("id"))
	request.ID = int64(n)

	err := self.CmsGw.SetUser(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) setUserPassword(c echo.Context) error {

	var (
		request api.User
		res     interface{}
		ctx     = c.Request().Context()
		code    = http.StatusInternalServerError
	)

	c.Bind(&request)

	err := self.CmsGw.SetUserPassword(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) validateRedeem(c echo.Context) error {

	var (
		request api.ValidateRedeem
		res     interface{}
		ctx     = c.Request().Context()
		code    = http.StatusInternalServerError
	)

	c.Bind(&request)

	err := self.CmsGw.ValidateRedeem(ctx, request)
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseReport{
			Status:  true,
			Message: "success",
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) detailHistory(c echo.Context) error {

	var res interface{}

	code := http.StatusInternalServerError

	ctx := c.Request().Context()

	data, err := self.CmsGw.FindDetailRedeem(ctx, c.Param("id"))
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseSuccess{
			Status:  true,
			Message: "success",
			Data:    data,
		}
	}

	return c.JSON(code, res)
}

func (self *CmsHandler) listJob(c echo.Context) error {

	var (
		res  interface{}
		ctx  = c.Request().Context()
		code = http.StatusInternalServerError
	)

	data, err := self.CmsGw.ListJob(ctx, c.Param("type"))
	if err != nil {
		cmsLog.Error(err)
		res = api.ResponseError{
			Status:  false,
			Message: err.Error(),
		}

	} else {
		code = http.StatusOK
		res = api.ResponseReport{
			Status:  true,
			Message: "success",
			Data:    data,
		}
	}

	return c.JSON(code, res)
}
