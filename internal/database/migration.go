package database

import (
	"github.com/Delaram-Gholampoor-Sagha/production_ready_rest_api/internal/comment"
	"github.com/jinzhu/gorm"
)


func MigrationDB(db *gorm.DB) error {
      if result := db.AutoMigrate(&comment.Comment{}) ; result.Error != nil {
		  return result.Error
	  }
	  return nil 
}

