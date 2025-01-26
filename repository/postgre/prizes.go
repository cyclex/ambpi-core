package postgre

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/pkg/errors"
)

func (self *postgreRepo) FindFirstPrizes(cond map[string]interface{}) (data model.Prizes, err error) {

	err = self.DB.Table("prizes").Where(cond).First(&data).Error

	return
}

func (self *postgreRepo) SetPrizes(cond, updated map[string]interface{}) (affected int64, err error) {

	q := self.DB.Table("prizes").Where(cond).Updates(updated)

	return q.RowsAffected, q.Error
}

func (self *postgreRepo) CreateRedeemPrizes(new model.RedeemPrizes) (affected int64, err error) {

	new.CreatedAt = time.Now().Local()
	q := self.DB.Create(&new)

	return q.RowsAffected, q.Error
}

func (self *postgreRepo) FindFirstRedeemPrizes(cond map[string]interface{}) (data model.RedeemPrizes, err error) {

	err = self.DB.Table("redeem_prizes").Where(cond).First(&data).Error

	return
}

func (self *postgreRepo) SetRedeemPrizes(cond, updated map[string]interface{}) (err error) {

	return self.DB.Table("redeem_prizes").Where(cond).Updates(updated).Error

}

func (self *postgreRepo) ReportPrize(req api.Report) (data map[string]interface{}, err error) {

	var (
		res   []model.Prizes
		datas []map[string]interface{}
	)

	err = self.DB.Table("v_prize_total").Where("prize_type = ?", req.PrizeType).Order("id asc").Find(&res).Error
	if err != nil {
		return
	}

	r := 1
	for _, v := range res {
		x := map[string]interface{}{
			"no":          r,
			"prize":       v.Prize,
			"quota":       0,
			"activeQuota": 0,
			"action":      "",
			"id":          v.ID,
		}

		datas = append(datas, x)
		r++
	}

	data = map[string]interface{}{
		"rows": len(datas),
		"data": datas,
	}

	return
}

func (r *postgreRepo) CreatePrize(rows [][]string) (totalRows int, err error) {
	// Start a transaction
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	if err := tx.Error; err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Validate rows are not empty and contain the header
	if len(rows) <= 1 {
		tx.Rollback()
		return 0, errors.New("no valid rows to process")
	}

	// Fetch the last sequence number from the database
	var lastSequenceNumber int
	err = tx.Table("prizes").
		Select("COALESCE(MAX(sequence_number), 0)").
		Row().
		Scan(&lastSequenceNumber)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch last sequence number: %w", err)
	}

	// Process each row, starting from the first data row (after the header)
	for i, row := range rows[1:] {
		if len(row) < 2 {
			tx.Rollback()
			return 0, fmt.Errorf("invalid row at index %d: not enough columns", i+2)
		}

		// Parse the sequence number from the row
		currentSequenceNumber, parseErr := strconv.Atoi(row[0])
		if parseErr != nil {
			tx.Rollback()
			return 0, fmt.Errorf("invalid sequence number at row #%d: %w", i+2, parseErr)
		}

		// Validate sequence order
		if currentSequenceNumber != lastSequenceNumber+1 {
			tx.Rollback()
			return 0, fmt.Errorf("sequence number out of order at row #%d: expected %d, got %d", i+2, lastSequenceNumber+1, currentSequenceNumber)
		}

		var seq = 0
		seq, _ = strconv.Atoi(row[0])
		// Create a prize entry
		prize := model.Prizes{
			SequenceNumber: seq,
			Prize:          row[1],
			IsUsed:         false,
			PrizeType:      pkg.PrizeData(row[1]),
		}
		if err := tx.Debug().Create(&prize).Error; err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to create prize at row #%d: %w", i+2, err)
		}

		// Update the last sequence number
		lastSequenceNumber = currentSequenceNumber
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return len(rows) - 1, nil // Exclude the header row from the count
}

func (self *postgreRepo) FindActivePrizes(cond map[string]interface{}, isActive bool) (data model.Prizes, err error) {

	q := self.DB.Table("prizes").Where(cond).Where("is_used = ?", false)

	err = q.First(&data).Error

	return
}
