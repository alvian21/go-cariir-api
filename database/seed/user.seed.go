package seed

import (
	"go-cariir-api/model/entity"
	"go-cariir-api/utils"
	"log"

	"gorm.io/gorm"
)

func UserSeed(db *gorm.DB) {
	password, _ := utils.HashingPassword("password123")

	var adminRole entity.Role
	db.Where("alias = ?", "admin").First(&adminRole)

	var userRole entity.Role
	db.Where("alias = ?", "user").First(&userRole)

	users := []entity.User{
		{
			FullName: "Admin User",
			Email:    "admin@example.com",
			Password: password,
			IsActive: true,
			RoleID:   &adminRole.ID,
		},
		{
			FullName: "John Doe",
			Email:    "john@example.com",
			Password: password,
			IsActive: true,
			RoleID:   &userRole.ID,
		},
	}

	for _, user := range users {
		var existingUser entity.User
		err := db.Where("email = ?", user.Email).First(&existingUser).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					log.Printf("Could not seed user %s: %v", user.Email, err)
				}
			} else {
				log.Printf("Error checking user %s: %v", user.Email, err)
			}
		}
	}
}
