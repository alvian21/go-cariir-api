package seed

import (
	"go-cariir-api/model/entity"
	"log"

	"gorm.io/gorm"
)

func RolePermissionSeed(db *gorm.DB) {
	var adminRole entity.Role
	db.Where("alias = ?", "admin").First(&adminRole)

	var userRole entity.Role
	db.Where("alias = ?", "user").First(&userRole)

	var permissions []entity.Permission
	db.Find(&permissions)

	// ADMIN gets all permissions
	for _, permission := range permissions {
		var existingRolePermission entity.RolePermission
		err := db.Where("role_id = ? AND permission_id = ?", adminRole.ID, permission.ID).First(&existingRolePermission).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			if err := db.Create(&entity.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: permission.ID,
			}).Error; err != nil {
				log.Printf("Could not link admin to permission %s: %v", permission.Alias, err)
			}
		}
	}

	// USER gets only read permissions
	for _, permission := range permissions {
		if permission.Alias == "role.read" || permission.Alias == "permission.read" || permission.Alias == "user.read" {
			var existingRolePermission entity.RolePermission
			err := db.Where("role_id = ? AND permission_id = ?", userRole.ID, permission.ID).First(&existingRolePermission).Error
			if err != nil && err == gorm.ErrRecordNotFound {
				if err := db.Create(&entity.RolePermission{
					RoleID:       userRole.ID,
					PermissionID: permission.ID,
				}).Error; err != nil {
					log.Printf("Could not link user to permission %s: %v", permission.Alias, err)
				}
			}
		}
	}
}
