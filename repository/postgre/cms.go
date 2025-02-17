package postgre

import (
	"crypto/md5"
	"strconv"

	"fmt"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	LayoutDefault  = time.RFC3339Nano
	Loc, _         = time.LoadLocation("Asia/Jakarta")
	LayoutDateTime = "2006-01-02 15:04:05"
)

func (self *postgreRepo) Login(username, password string) (data model.UserCMS, err error) {

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	cond := map[string]interface{}{
		"username": username,
		"password": password,
	}

	err = self.DB.Where(cond).First(&data).Error
	if err != nil {
		err = errors.New("Invalid username or password. Please try again.")
	}

	return
}

func (self *postgreRepo) CheckToken(token string) (err error) {

	var data model.UserCMS

	cond := map[string]interface{}{
		"token": token,
	}

	err = self.DB.Where(cond).First(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CheckToken]")
	}

	return
}

func (self *postgreRepo) SetTokenLogin(id uint, token string) (err error) {

	err = self.DB.Model(&model.UserCMS{}).Where(map[string]interface{}{"id": id}).Updates(map[string]interface{}{"token": token}).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetTokenLogin]")
	}
	return

}

func (self *postgreRepo) ReportHistoryValidation(req api.Report) (data map[string]interface{}, err error) {

	type tmp struct {
		Rnum           string `json:"rnum"`
		Msisdn         string `json:"msisdn"`
		Name           string `json:"name"`
		Prize          string `json:"prize"`
		LotteryNumber  string `json:"lotteryNumber"`
		DateValidation string `json:"dateValidation"`
		Approved       bool   `json:"approved"`
		RedeemID       int    `json:"redeemID"`
		County         string `json:"county"`
		Profession     string `json:"profession"`
		NIK            string `json:"nik"`
		DateRedeem     string `json:"dateRedeem"`
		Amount         string `json:"amount"`
		Notes          string `json:"notes"`
		Receipt        string `json:"receipt"`
	}

	var (
		res   []tmp
		cond  map[string]interface{}
		datas []map[string]interface{}
		rows  int64
	)

	q := self.DB.Table("detailed_prize_redemptions").Select("*, row_number() OVER () as rnum").Where(cond)
	if req.From != "" {
		q = q.Where("date(date_validation) BETWEEN ? AND ?", req.From, req.To)
	}
	q = q.Where("date_validation != ''")

	if req.Keyword != "" {
		column := fmt.Sprintf("%s::Text like ?", req.Column)
		q = q.Where(column, "%"+req.Keyword+"%")
	}

	q.Count(&rows)
	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	}
	err = q.Order("sequence_number desc").Offset(req.Offset).Find(&res).Error
	if err != nil {
		return
	}

	for _, v := range res {
		x := map[string]interface{}{
			"rNum":             v.Rnum,
			"msisdn":           fmt.Sprintf("`%s", v.Msisdn),
			"name":             v.Name,
			"prize":            v.Prize,
			"lotteryNumber":    v.LotteryNumber,
			"dateValidation":   v.DateValidation,
			"statusValidation": pkg.StatusValidation(v.Approved),
			"id":               v.RedeemID,
			"county":           v.County,
			"nik":              v.NIK,
			"profession":       v.Profession,
			"dateRedeem":       v.DateRedeem,
			"amount":           v.Amount,
			"notes":            v.Notes,
			"receipt":          v.Receipt,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": rows,
		"data": datas,
	}

	return
}

func (r *postgreRepo) ReportSummary() (map[string]interface{}, error) {

	// Struct to hold query results
	type Result struct {
		Prize     string
		Available int
		Used      int
		Total     int
	}

	var results []Result

	// Query to group by prize and calculate counts for available, used, and total
	err := r.DB.Table("prizes").
		Select(`prize, 
			SUM(CASE WHEN is_used = false THEN 1 ELSE 0 END) AS available, 
			SUM(CASE WHEN is_used = true THEN 1 ELSE 0 END) AS used, 
			COUNT(*) AS total`).
		Group("prize").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prize summary: %w", err)
	}

	// Prepare the output map
	data := map[string]interface{}{
		"rows": len(results),
		"data": results,
	}

	return data, nil
}

func (self *postgreRepo) ReportRedeem(req api.Report) (data map[string]interface{}, err error) {

	type tmp struct {
		Rnum       string `json:"rnum"`
		Msisdn     string `json:"msisdn"`
		Name       string `json:"name"`
		Nik        string `json:"nik"`
		DateRedeem string `json:"dateRedeem"`
		Prize      string `json:"prize"`
		County     string `json:"county"`
		Profession string `json:"profession"`
		RedeemID   string `json:"redeem_id"`
	}

	var (
		res   []tmp
		datas []map[string]interface{}
		rows  int64
	)

	q := self.DB.Table("detailed_prize_redemptions").Select("*, row_number() OVER () as rnum")
	if req.From != "" {
		q = q.Where("date(date_redeem) BETWEEN ? AND ?", req.From, req.To)
	}
	q = q.Where("date_validation = ''")

	if req.Keyword != "" {
		column := fmt.Sprintf("%s ilike ?", req.Column)
		q = q.Where(column, "%"+req.Keyword+"%")
	}

	q.Count(&rows)
	q = q.Order("date_redeem desc")

	if req.Limit > 0 {
		q = q.Limit(req.Limit)
	}

	err = q.Offset(req.Offset).Find(&res).Error
	if err != nil {
		return
	}

	for _, v := range res {
		dateRedeem, _ := time.Parse(LayoutDefault, v.DateRedeem)
		x := map[string]interface{}{
			"rNum":       v.Rnum,
			"msisdn":     fmt.Sprintf("`%s", v.Msisdn),
			"name":       v.Name,
			"nik":        fmt.Sprintf("`%s", v.Nik),
			"dateRedeem": dateRedeem.In(Loc).Format(LayoutDateTime),
			"prize":      v.Prize,
			"county":     v.County,
			"profession": v.Profession,
			"redeem_id":  v.RedeemID,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"rows": rows,
		"data": datas,
	}

	return
}

func (self *postgreRepo) FindRedeemID(id string) (data map[string]interface{}, err error) {

	type tmp struct {
		Rnum           string `json:"rnum"`
		Msisdn         string `json:"msisdn"`
		Name           string `json:"name"`
		County         string `json:"county"`
		Nik            string `json:"nik"`
		Profession     string `json:"profession"`
		Prize          string `json:"prize"`
		LotteryNumber  string `json:"lotteryNumber"`
		Amount         string `json:"amount"`
		Notes          string `json:"notes"`
		DateRedeem     string `json:"dateRedeem"`
		DateValidation string `json:"dateValidation"`
		Approved       string `json:"approved"`
		RedeemID       string `json:"redeem_id"`
		Receipt        string `json:"receipt"`
	}

	var (
		res tmp
	)

	err = self.DB.Table("detailed_prize_redemptions").Where("redeem_id = ?", id).Find(&res).Error
	if err != nil {
		return
	}

	if res.Name == "" {
		return
	}

	dateRedeem, _ := time.Parse(LayoutDefault, res.DateRedeem)
	approved, _ := strconv.ParseBool(res.Approved)
	data = map[string]interface{}{
		"msisdn":         fmt.Sprintf("`%s", res.Msisdn),
		"name":           res.Name,
		"nik":            fmt.Sprintf("`%s", res.Nik),
		"dateRedeem":     dateRedeem.In(Loc).Format(LayoutDateTime),
		"prize":          res.Prize,
		"county":         res.County,
		"profession":     res.Profession,
		"notes":          res.Notes,
		"approved":       pkg.StatusValidation(approved),
		"dateValidation": res.DateValidation,
		"amount":         res.Amount,
		"lotteryNumber":  res.LotteryNumber,
		"id":             res.RedeemID,
		"receipt":        res.Receipt,
	}

	return
}

func (r *postgreRepo) ReportUsage(req api.Report) (map[string]interface{}, error) {
	// Initialize the result map
	data := make(map[string]interface{})

	// Initialize counters
	var totalSubmit, totalValid, totalInvalid, totalPending int64

	// Define a reusable base query function
	baseQuery := func() *gorm.DB {
		return r.DB.Table("redeem_prizes").Where("date(created_at) BETWEEN ? AND ?", req.From, req.To)
	}

	// Count total submissions
	if err := baseQuery().Count(&totalSubmit).Error; err != nil {
		return nil, fmt.Errorf("failed to count total submissions: %w", err)
	}

	// Count valid submissions
	if err := baseQuery().Where("approved = ?", "true").Count(&totalValid).Error; err != nil {
		return nil, fmt.Errorf("failed to count valid submissions: %w", err)
	}

	// Count invalid submissions
	if err := baseQuery().Where("approved = ?", "false").Where("date_validation != ''").Count(&totalInvalid).Error; err != nil {
		return nil, fmt.Errorf("failed to count invalid submissions: %w", err)
	}

	// Count pending submissions
	if err := baseQuery().Where("date_validation = ''").Count(&totalPending).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending submissions: %w", err)
	}

	// Populate the result map
	data["totalSubmit"] = totalSubmit
	data["totalValid"] = totalValid
	data["totalInvalid"] = totalInvalid
	data["totalPending"] = totalPending

	return data, nil
}

func (r *postgreRepo) ReportAvailability() (data map[string]interface{}, err error) {
	// Initialize the result map
	data = make(map[string]interface{})

	var active, inactive int
	type Result struct {
		IsUsed bool
		Count  int
	}

	var results []Result
	err = r.DB.Table("prizes").
		Select("is_used, COUNT(*) AS count").
		Group("is_used").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prize availability: %w", err)
	}

	// Process results to determine active and inactive counts
	for _, result := range results {
		if result.IsUsed {
			active = result.Count
		} else {
			inactive = result.Count
		}
	}

	// Avoid division by zero
	total := active + inactive
	if total == 0 {
		data["percentage"] = 0
		data["message"] = "No data available for prizes"
		return data, nil
	}

	// Calculate the percentage of active prizes
	percentage := (float64(active) / float64(total)) * 100

	data["active"] = active
	data["inactive"] = inactive
	data["percentage"] = percentage

	return data, nil
}

func (self *postgreRepo) ReportUsers() (data map[string]interface{}, err error) {

	type tmp struct {
		Rnum      string `json:"rnum"`
		Username  string `json:"username"`
		ID        string `json:"id"`
		Flag      string `json:"flag"`
		Level     string `json:"level"`
		CreatedAt string `json:"createdAt"`
	}

	var (
		datas []map[string]interface{}
		res   []tmp
	)

	err = self.DB.Table("user_cms").Select("*, row_number() OVER () as rnum").Order("id desc").Find(&res).Error
	if err != nil {
		return
	}

	for _, v := range res {
		createdAt, _ := strconv.Atoi(v.CreatedAt)
		createdAts := time.Unix(int64(createdAt), 0).Format("2006-01-02 15:04:05")
		x := map[string]interface{}{
			"id":        v.ID,
			"username":  v.Username,
			"flag":      v.Flag,
			"level":     v.Level,
			"rNum":      v.Rnum,
			"createdAt": createdAts,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"data": datas,
		"rows": len(datas),
	}

	return
}

func (self *postgreRepo) ReportCampaign() (data map[string]interface{}, err error) {

	type tmp struct {
		Retail    string `json:"retail"`
		Status    bool   `json:"status"`
		StartDate int    `json:"start_date"`
		EndDate   int    `json:"end_date"`
		ID        int    `json:"id"`
	}

	var (
		datas []map[string]interface{}
		res   []tmp
	)

	err = self.DB.Table("programs").Find(&res).Error
	if err != nil {
		return
	}

	for _, v := range res {

		x := map[string]interface{}{
			"retail":    v.Retail,
			"startDate": time.Unix(int64(v.StartDate), 0).Local().Format("2006-01-02 15:04:05"),
			"endDate":   time.Unix(int64(v.EndDate), 0).Local().Format("2006-01-02 15:04:05"),
			"id":        v.ID,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"data": datas,
		"rows": len(datas),
	}

	return
}

func (self *postgreRepo) ReportPrizes() (data map[string]interface{}, err error) {

	type tmp struct {
		Prize          string `json:"prize"`
		SequenceNumber string `json:"sequence_number"`
	}

	var (
		datas []map[string]interface{}
		res   []tmp
	)

	err = self.DB.Table("prizes").Order("sequence_number asc").Find(&res).Error
	if err != nil {
		return
	}

	for _, v := range res {

		x := map[string]interface{}{
			"prize":          v.Prize,
			"sequenceNumber": v.SequenceNumber,
		}

		datas = append(datas, x)
	}

	data = map[string]interface{}{
		"data": datas,
		"rows": len(datas),
	}

	return
}
