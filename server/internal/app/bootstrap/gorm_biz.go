package bootstrap

import "gorm.io/gorm"

func bizModel(db *gorm.DB) error {
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
