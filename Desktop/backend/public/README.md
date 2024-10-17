This directory is for housing our godoc static site dependencies. See the `pages` step in our .gitlab-ci.yml file.

To add documentation here, add comments above important or complex functions that you think developers need to see.
For example, in common/dao/main.go, there is this:

```
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
```

godoc-static parses this comment when creating a static page: https://circulohealth.gitlab.io/sonar/backend/github.com/circulohealth/sonar-backend/packages/common/dao/#OpenConnectionToDoppler