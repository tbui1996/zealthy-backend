package mapper

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type newExternalUserSQLInput struct {
	db     *gorm.DB
	logger *zap.Logger
}

func newExternalUserSQL(input *newExternalUserSQLInput) *externalUserSQL {
	return &externalUserSQL{
		db:     input.db,
		logger: input.logger.With(zap.String("mapper", "externalUserSQL")),
	}
}

// implements externalUserSQLAPI
type externalUserSQL struct {
	db     *gorm.DB
	logger *zap.Logger
}

// internal to package
func (m *externalUserSQL) find(id string) (*externalUserSQLRecord, error) {
	m.logger.Debug("find")
	record := &externalUserSQLRecord{}

	result := m.db.Find(&record, "id = ?", id)

	if result.Error != nil {
		m.logger.Debug("find error")
		return nil, result.Error
	}

	m.logger.Debug("found record")
	return record, nil
}

func (m *externalUserSQL) findAll() ([]*externalUserSQLRecord, error) {
	m.logger.Debug("findAll")
	records := make([]*externalUserSQLRecord, 0)

	result := m.db.Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func (m *externalUserSQL) update(updater externalUserSQLRecordUpdater) error {
	m.logger.Debug("update")
	updates := make(map[string]interface{})

	if updater.externalUserOrganizationIDChanged {
		if updater.externalUserOrganizationID == nil {
			updates["external_user_organization_id"] = nil
		} else {
			updates["external_user_organization_id"] = *updater.externalUserOrganizationID
		}
	}

	if len(updates) == 0 {
		return nil
	}

	result := m.db.Model(&externalUserSQLRecord{ID: updater.id}).Updates(updates)
	return result.Error
}
