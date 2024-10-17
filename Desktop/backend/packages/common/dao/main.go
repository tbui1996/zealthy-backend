package dao

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func OpenConnectionWithTablePrefix(tablePrefix string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("NAME"))

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// e.g. "forms.", "router.", etc.
			TablePrefix: tablePrefix,
		},
	})
}

// This function opens a connection to Doppler. To use this, make sure you set the DOPPLER env variables in your Lambda Terraform.
// 	environment {
// 		variables = {
// 			"DOPPLERHOST"     = var.live_env ? var.doppler_host : var.db_host
// 			"DOPPLERPORT"     = var.live_env ? var.doppler_port : var.db_port
// 			"DOPPLERUSER"     = var.live_env ? var.doppler_user : var.db_username
// 			"DOPPLERPASSWORD" = var.live_env ? var.doppler_pw : var.db_password
// 			"DOPPLERDBNAME"   = var.live_env ? var.doppler_dbname : "external"
// 		}
// 	}
func OpenConnectionToDoppler() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DOPPLERHOST"), os.Getenv("DOPPLERPORT"), os.Getenv("DOPPLERUSER"), os.Getenv("DOPPLERPASSWORD"), os.Getenv("DOPPLERDBNAME"))

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
}

func IsUniqueConstraintViolation(err error, constraint string) bool {
	msg := err.Error()
	return strings.Contains(msg, "(SQLSTATE 23505)") && strings.Contains(msg, constraint)
}
