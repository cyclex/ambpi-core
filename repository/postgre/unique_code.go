package postgre

import (
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
)

func (self *postgreRepo) CreateUsersUniqueCodes(new model.UsersUniqueCode) (affected int64, err error, pk uint) {

	new.CreatedAt = time.Now().Local()
	q := self.DB.Create(&new)

	return q.RowsAffected, q.Error, new.ID
}
