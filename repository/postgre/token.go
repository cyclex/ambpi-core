package postgre

import (
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/pkg/errors"
)

func (self *postgreRepo) SetToken(updated map[string]interface{}) (err error) {

	err = self.DB.Model(&model.Token{}).Where("id = ?", "1").Updates(updated).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetToken]")
	}

	return
}

func (self *postgreRepo) FindToken() (data model.Token, err error) {

	err = self.DB.Model(&model.Token{}).First(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.FindToken]")
	}

	return
}
