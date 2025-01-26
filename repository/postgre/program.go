package postgre

import (
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
)

func (self *postgreRepo) FindProgram() (data []map[string]interface{}, err error) {

	type tmp struct {
		Retail    string `json:"retail"`
		Status    bool   `json:"status"`
		StartDate int    `json:"start_date"`
		EndDate   int    `json:"end_date"`
		ID        int    `json:"id"`
	}

	var res []tmp

	err = self.DB.Table("programs").Find(&res).Error
	if err != nil {
		return
	}

	for _, v := range res {

		// start, _ := time.Parse(time.RFC3339, "02 Sep 15 08:00 WIB")

		x := map[string]interface{}{
			"retail":    v.Retail,
			"startDate": time.Unix(int64(v.StartDate), 0).Local().Format("2006-01-02 15:04:05"),
			"endDate":   time.Unix(int64(v.EndDate), 0).Local().Format("2006-01-02 15:04:05"),
			"id":        v.ID,
		}

		data = append(data, x)
	}

	return
}

func (self *postgreRepo) SetProgram(cond, updated map[string]interface{}) (affected int64, err error) {

	q := self.DB.Table("programs").Where(cond).Updates(updated)

	return q.RowsAffected, q.Error
}

func (self *postgreRepo) IsProgramActive(retail string) (status int) {

	var program model.Program
	now := time.Now().Local().Unix()

	self.DB.Table("programs").Where(map[string]interface{}{"retail": retail}).Find(&program)

	if program.StartDate >= now {
		return 1
	} else if program.EndDate <= now {
		return 2
	} else {
		return 3
	}

}
