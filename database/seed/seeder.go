package seed

import (
	"fmt"
	"go-cariir-api/database"
)

func RunSeed() {
	RoleSeed(database.DB)
	PermissionSeed(database.DB)
	RolePermissionSeed(database.DB)
	UserSeed(database.DB)

	fmt.Println("Seeding complete")
}
