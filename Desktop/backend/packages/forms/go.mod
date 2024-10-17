module github.com/circulohealth/sonar-backend/packages/forms

go 1.16

replace github.com/circulohealth/sonar-backend/packages/common => ../common/

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-lambda-go v1.28.0
	github.com/aws/aws-sdk-go v1.40.20
	github.com/circulohealth/sonar-backend/packages/common v0.0.0-00010101000000-000000000000
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.10.2
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.20.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gorm.io/driver/postgres v1.2.3
	gorm.io/gorm v1.22.5
)
