module github.com/circulohealth/sonar-backend/packages/users

go 1.16

replace github.com/circulohealth/sonar-backend/packages/common => ../common/

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-lambda-go v1.28.0
	github.com/aws/aws-sdk-go v1.40.20
	github.com/circulohealth/sonar-backend/packages/common v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/lestrrat-go/jwx v1.2.6 // indirect
	github.com/sethvargo/go-password v0.2.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.20.0
	gorm.io/driver/postgres v1.2.3
	gorm.io/gorm v1.22.5
)
