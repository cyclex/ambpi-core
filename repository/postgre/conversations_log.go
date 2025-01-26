package postgre

import (
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
)

func (self *postgreRepo) CreateConversationsLog(new model.ConversationsLog) (err error) {

	new.CreatedAt = time.Now().Local()
	q := self.DB.Create(&new)

	return q.Error
}
