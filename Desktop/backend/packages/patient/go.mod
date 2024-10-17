module github.com/circulohealth/sonar-backend/packages/patient

go 1.16

replace github.com/circulohealth/sonar-backend/packages/common => ../common

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-lambda-go v1.28.0 // indirect
	github.com/circulohealth/sonar-backend/packages/common v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v4 v4.14.1 // indirect
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.20.0
	golang.org/x/crypto v0.0.0-20220126234351-aa10faf2a1f8 // indirect
	gorm.io/driver/postgres v1.2.3 // indirect
	gorm.io/gorm v1.22.5
)
