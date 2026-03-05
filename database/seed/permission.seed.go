package seed

import (
	"go-cariir-api/model/entity"
	"log"

	"gorm.io/gorm"
)

func PermissionSeed(db *gorm.DB) {
	permissions := []entity.Permission{
		{Alias: "role.create", Name: "Tambah Role", IsActive: true},
		{Alias: "role.read", Name: "Lihat Role", IsActive: true},
		{Alias: "role.update", Name: "Ubah Role", IsActive: true},
		{Alias: "role.delete", Name: "Hapus Role", IsActive: true},
		{Alias: "permission.create", Name: "Tambah Permission", IsActive: true},
		{Alias: "permission.read", Name: "Lihat Permission", IsActive: true},
		{Alias: "permission.update", Name: "Ubah Permission", IsActive: true},
		{Alias: "permission.delete", Name: "Hapus Permission", IsActive: true},
		{Alias: "user.create", Name: "Tambah User", IsActive: true},
		{Alias: "user.read", Name: "Lihat User", IsActive: true},
		{Alias: "user.update", Name: "Ubah User", IsActive: true},
	}

	for _, permission := range permissions {
		var existingPermission entity.Permission
		err := db.Where("alias = ?", permission.Alias).First(&existingPermission).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					log.Printf("Could not seed permission %s: %v", permission.Alias, err)
				}
			} else {
				log.Printf("Error checking permission %s: %v", permission.Alias, err)
			}
		}
	}
}
