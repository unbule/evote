package dao

import "github.com/jinzhu/gorm"

type SelectDao struct {
}

func (selectdao *SelectDao) SelectByUsername(username string) func (db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username=?",username)
	}
}

