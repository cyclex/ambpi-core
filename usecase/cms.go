package usecase

import (
	"context"
	"fmt"
	"net/http"

	"time"

	"github.com/apex/log"
	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/domain/repository"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pkg/errors"
)

type cmsUcase struct {
	m              repository.ModelRepository
	q              repository.QueueRepository
	contextTimeout time.Duration
	chatUcase      domain.ChatUcase
}

func NewCmsUcase(m repository.ModelRepository, ctxTimeout time.Duration, chatUcase domain.ChatUcase, q repository.QueueRepository) domain.CmsUcase {

	return &cmsUcase{
		m:              m,
		contextTimeout: ctxTimeout,
		chatUcase:      chatUcase,
		q:              q,
	}
}

func (self *cmsUcase) Login(c context.Context, req api.Login) (data map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	res, err := self.m.Login(req.Username, req.Password)
	if err != nil {
		return
	}

	tokenCms := pkg.TokenGenerator(16)
	err = self.m.SetTokenLogin(res.ID, tokenCms)
	if err != nil {
		err = errors.Wrap(err, "[usecase.Login]")
	}

	data = map[string]interface{}{
		"username": res.Username,
		"user_id":  res.ID,
		"level":    res.Level,
		"token":    tokenCms,
	}

	return
}

func (self *cmsUcase) Report(c context.Context, req api.Report, category string) (data map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	switch category {
	case "redeem":
		data, err = self.m.ReportRedeem(req)
		break

	case "history_validation":
		data, err = self.m.ReportHistoryValidation(req)
		break

	case "availability":
		data, err = self.m.ReportAvailability()
		break

	case "prize":
		data, err = self.m.ReportPrize(req)
		break

	case "usage":
		data, err = self.m.ReportUsage(req)
		break

	case "summary":
		data, err = self.m.ReportSummary()
		break
	// case "messages":
	// 	data, err = self.m.ReportMessages(req)
	// 	break
	default:
		// data, err = self.m.ReportRedeem(req, "")
		break
	}

	return
}

func (self *cmsUcase) SendPushNotif(c context.Context, req api.SendPushNotif) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	cond := map[string]interface{}{
		"id": req.RedeemID,
	}
	dataRP, err := self.m.FindFirstRedeemPrizes(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SendPushNotif] FindFirstRedeemPrizes")
		return
	}

	cond = map[string]interface{}{
		"id": dataRP.PrizeID,
	}
	dataP, err := self.m.FindFirstPrizes(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SendPushNotif] FindFirstPrizes")
		return
	}

	res, statusCode, err := self.chatUcase.ChatToUser(dataRP.Msisdn, req.Param, "push", req.TemplateName)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SendPushNotif] ChatToUser")
		return
	}

	if statusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("%s: %s", http.StatusText(statusCode), res))
		err = errors.Wrap(err, "[usecase.SendPushNotif] ChatToUser")
		return
	}

	clog := model.ConversationsLog{
		SessionID: uuid.NewString(),
		Incoming:  "push-notif-by-admin",
		WAID:      dataRP.Msisdn,
		Outgouing: dataP.Prize,
	}
	err = self.m.CreateConversationsLog(clog)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SendPushNotif] CreateConversationsLog")
		log.Error(err.Error())
	}

	return
}

func (self *cmsUcase) GetProgram(c context.Context) (data []map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	return self.m.FindProgram()

}

func (self *cmsUcase) SetProgram(c context.Context, cond, data map[string]interface{}) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	_, err = self.m.SetProgram(cond, data)

	return
}

func (self *cmsUcase) ImportPrize(c context.Context, req api.Job) (status bool, totalRows int, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	rows, err := pkg.ReadFromFile(req.File)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	totalRows, err = self.m.CreatePrize(rows)

	return
}

func (self *cmsUcase) CreateJob(c context.Context, req api.Job) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data := model.QueueJob{
		JobType:   req.JobType,
		JobStatus: "1",
		Author:    req.Author,
		File:      req.File,
		ID:        primitive.NewObjectID(),
		StartAt:   req.StartDate,
		EndAt:     req.EndDate,
	}
	err = self.q.CreateJob(data)

	return
}

func (self *cmsUcase) CreateUser(c context.Context, req api.User) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	if req.Username == "" && req.Level == "" {
		return errors.New("invalid request")
	}

	var dataUser model.UserCMS

	dataUser = model.UserCMS{
		Username: req.Username,
		Level:    req.Level,
		Password: req.Password,
	}

	err = self.m.CreateUser(dataUser)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateUser]")
	}

	return
}

func (self *cmsUcase) DeleteUser(c context.Context, deletedID int64) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.m.RemoveUser([]int64{deletedID})
	if err != nil {
		err = errors.Wrap(err, "[usecase.DeleteUser]")
	}

	return
}

func (self *cmsUcase) SetUser(c context.Context, req api.User) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataUser := model.UserCMS{
		Username: req.Username,
		Level:    req.Level,
		Password: req.Password,
	}

	err = self.m.SetUser(req.ID, dataUser)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetUser]")
	}

	return
}

func (self *cmsUcase) SetUserPassword(c context.Context, req api.User) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	dataUser := model.UserCMS{
		Password: req.Password,
	}

	err = self.m.SetUserPassword(req.Username, dataUser)
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetUser]")
	}

	return
}

func (self *cmsUcase) CheckToken(c context.Context, req api.CheckToken) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	err = self.m.CheckToken(req.Token)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CheckToken]")
	}

	return
}

func (self *cmsUcase) DownloadRedeem(c context.Context, req api.Job) (files string, status bool, totalRows int, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data, err := self.m.ReportRedeem(api.Report{From: req.StartDate, To: req.EndDate})
	if err != nil {
		err = errors.Wrap(err, "[usecase.DownloadRedeem]")
		return
	}

	if datas, ok := data["data"].([]map[string]interface{}); ok {
		files, totalRows, err = pkg.WriteXLS(datas)
		if err != nil {
			err = errors.Wrap(err, "[usecase.DownloadRedeem]")
		} else {
			status = true
		}
	}

	return
}

func (self *cmsUcase) ValidateRedeem(c context.Context, req api.ValidateRedeem) (err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	var lotteryNumber string
	cond := map[string]interface{}{
		"id": req.ID,
	}

	data, err := self.m.FindFirstRedeemPrizes(cond)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ValidateRedeem]")
		return
	}

	templateName := "invalid"
	var param []string
	if req.Approved {
		cond := map[string]interface{}{
			"id": data.PrizeID,
		}
		p, errs := self.m.FindFirstPrizes(cond)
		if errs != nil {
			err = errors.Wrap(errs, "[usecase.ValidateRedeem]")
			return
		}
		templateName = "valid_belumberuntung"
		lotteryNumber, errs = pkg.GenerateRandomCode(8)
		if errs != nil {
			err = errors.Wrap(errs, "[usecase.ValidateRedeem]")
			return
		}
		param = append(param, lotteryNumber)
		if p.PrizeType != "zonk" {
			templateName = "valid_hadiahlangsung"
			param = []string{p.Prize, lotteryNumber}
		}
	}

	pn := api.SendPushNotif{RedeemID: data.ID, PushBy: req.Author, TemplateName: templateName, Param: param}
	err = self.SendPushNotif(c, pn)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ValidateRedeem]")
		return
	}

	updated := map[string]interface{}{
		"amount":          req.Amount,
		"notes":           req.Notes,
		"approved":        req.Approved,
		"lottery_number":  lotteryNumber,
		"date_validation": time.Now().Local(),
	}

	err = self.m.SetRedeemPrizes(cond, updated)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ValidateRedeem]")
		return
	}

	if !req.Approved {
		cond := map[string]interface{}{
			"id": data.PrizeID,
		}
		updated := map[string]interface{}{
			"is_used": req.Approved,
		}
		_, err = self.m.SetPrizes(cond, updated)
		if err != nil {
			err = errors.Wrap(err, "[usecase.ValidateRedeem]")
			return
		}
	}

	return
}

func (self *cmsUcase) FindDetailRedeem(c context.Context, id string) (data map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	data, err = self.m.FindRedeemID(id)
	if err != nil {
		err = errors.Wrap(err, "[usecase.FindDetailRedeem]")
		return
	}

	return
}

func (self *cmsUcase) ListJob(c context.Context, category string) (data []map[string]interface{}, err error) {

	_, cancel := context.WithTimeout(c, self.contextTimeout)
	defer cancel()

	jobs, err := self.q.GetJob(category)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ListJob]")
		return
	}

	// Convert each QueueJob to a map[string]interface{}
	for _, job := range jobs {
		jobMap := map[string]interface{}{
			"id":         job.ID.Hex(), // Convert ObjectID to string
			"job_type":   job.JobType,
			"job_status": job.JobStatus,
			"author":     job.Author,
			"file":       job.File,
			"total_rows": job.TotalRows,
			"expired_at": job.ExpiredAt,
			"start_at":   job.StartAt,
			"end_at":     job.EndAt,
			"created_at": job.CreatedAt,
		}
		data = append(data, jobMap)
	}

	return
}
