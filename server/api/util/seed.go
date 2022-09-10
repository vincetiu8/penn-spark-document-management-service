package util

import (
	"log"

	"gorm.io/gorm"

	"github.com/invincibot/penn-spark-server/api/models"
)

func Load(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.User{}, &models.UserRole{}, &models.AccessRole{})
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.UserRole{}, &models.AccessRole{})
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	userRole, err := models.CreateUserRole(db, models.UserRole{
		Name: "default",
	})
	if err != nil {
		log.Fatalf("cannot seed userRole: %v", err)
	}

	_, err = models.CreateUser(db, models.User{
		Username:  "admin",
		FirstName: "admin",
		LastName:  "admin",
		Password:  "password",
		IsAdmin:   true,
		UserRoles: []models.UserRole{
			userRole,
		},
	})
	if err != nil {
		log.Fatalf("cannot seed user: %v", err)
	}
}
