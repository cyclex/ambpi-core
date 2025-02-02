package postgre

import (
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
)

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
