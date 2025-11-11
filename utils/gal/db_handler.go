package gal

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	model "seaotterms-api/model/galgame"
)

func WriteTmpData(dataType string, dataContent string, expirationAt time.Time, db *gorm.DB) error {
	data := model.TmpData{
		Type:         dataType,
		Content:      dataContent,
		ExpirationAt: expirationAt,
	}

	err := db.Create(&data).Error
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func WriteLog(dataType string, message string, severity uint, db *gorm.DB) error {
	data := model.Log{
		Type:     dataType,
		Message:  message,
		Severity: severity,
	}

	err := db.Create(&data).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
