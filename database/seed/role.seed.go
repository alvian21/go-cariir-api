package seed

import (
	"go-cariir-api/model/entity"
	"log"

	"gorm.io/gorm"
)

func RoleSeed(db *gorm.DB) {
	roles := []entity.Role{
		{
			Alias: "admin",
			Name:  "Administrator",
		},
		{
			Alias: "user",
			Name:  "User",
		},
	}

	for _, role := range roles {
		var existingRole entity.Role
		err := db.Where("alias = ?", role.Alias).First(&existingRole).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					log.Printf("Could not seed role %s: %v", role.Alias, err)
				}
			} else {
				log.Printf("Error checking role %s: %v", role.Alias, err)
			}
		}
	}
}
