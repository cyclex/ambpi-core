package postgre

import (
	"time"

	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/pkg/errors"
)

func (self *postgreRepo) FindUserBy(cond map[string]interface{}) (data []model.UserCMS, err error) {

	err = self.DB.Where(cond).Find(&data).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.FindUserBy]")
	}

	return
}

func (self *postgreRepo) SetUser(id int64, kol model.UserCMS) (err error) {

	kol.UpdatedAt = time.Now().Local().Unix()
	kol.Password = pkg.HashPassword(kol.Password)
	err = self.DB.Where("id = ?", id).Updates(kol).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetUser]")
	}

	return

}

func (self *postgreRepo) RemoveUser(id []int64) (err error) {

	err = self.DB.Delete(&model.UserCMS{}, id).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.RemoveUser]")
	}

	return

}

func (self *postgreRepo) CreateUser(new model.UserCMS) (err error) {

	new.CreatedAt = time.Now().Local().Unix()
	new.Password = pkg.HashPassword(new.Password)
	err = self.DB.Create(&new).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.CreateUser]")
	}

	return
}

func (self *postgreRepo) SetUserPassword(username string, kol model.UserCMS) (err error) {

	kol.UpdatedAt = time.Now().Local().Unix()
	kol.Password = pkg.HashPassword(kol.Password)
	err = self.DB.Where("username = ?", username).Updates(kol).Error
	if err != nil {
		err = errors.Wrap(err, "[postgre.SetUser]")
	}

	return

}
