package patients

import (
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"gorm.io/gorm"
)

type PatientRepository struct {
	DB *gorm.DB
}

func NewPatientRepository() (*PatientRepository, error) {
	db, err := dao.OpenConnectionWithTablePrefix(dao.Chat)
	if err != nil {
		return nil, err
	}
	return &PatientRepository{
		DB: db,
	}, nil
}

func (repo *PatientRepository) FindAll(filter interface{}) ([]model.Patient, error) {
	var patients []model.Patient
	if filter != nil {
		if err := repo.DB.Where(filter).Find(&patients).Error; err != nil {
			return nil, err
		}
		return patients, nil
	}

	if err := repo.DB.Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}

func (repo *PatientRepository) Find(filter interface{}) (*model.Patient, error) {
	var patient model.Patient
	if err := repo.DB.Where(filter).Take(&patient).Error; err != nil {
		return nil, err
	}

	return &patient, nil
}
