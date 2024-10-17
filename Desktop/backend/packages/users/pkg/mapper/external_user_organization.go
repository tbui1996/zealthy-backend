package mapper

import (
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Initialize with NewExteranlUserOrganization, not struct initialization
type ExternalUserOrganization struct {
	identityMap map[int]*model.ExternalUserOrganization
	db          *gorm.DB
	logger      *zap.Logger
}

type newExternalUserOrganizationInput struct {
	db     *gorm.DB
	logger *zap.Logger
}

func newExternalUserOrganization(input *newExternalUserOrganizationInput) *ExternalUserOrganization {
	return &ExternalUserOrganization{
		identityMap: make(map[int]*model.ExternalUserOrganization),
		db:          input.db,
		logger:      input.logger.With(zap.String("mapper", "ExternalUserOrganization")),
	}
}

func (m *ExternalUserOrganization) FindAll() ([]*model.ExternalUserOrganization, error) {
	m.logger.Debug("FindAll")
	records := make([]*externalUserOrganizationRecord, 0)

	result := m.db.Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	organizations := make([]*model.ExternalUserOrganization, len(records))
	for i, record := range records {
		organizations[i] = model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
			ID: record.ID,
		}).WithName(record.Name).Value()

		m.identityMap[record.ID] = organizations[i]
	}

	return organizations, nil
}

func (m *ExternalUserOrganization) Find(id int) (*model.ExternalUserOrganization, error) {
	m.logger.Debug("Find")
	existing, ok := m.identityMap[id]

	if ok {
		if existing != nil {
			return existing, nil
		}

		// remove from map if id was nil
		delete(m.identityMap, id)
	}

	record := &externalUserOrganizationRecord{}
	result := m.db.Find(&record, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	output := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: id,
	}).WithName(record.Name).Value()
	m.identityMap[id] = output
	return output, nil
}

// Returns the ID of the created organization, or an error if any occurred
// In order to create a model, use the model builder
// Guaranteed to not be inserted if an error occurred
func (m *ExternalUserOrganization) Insert(input *iface.ExternalUserOrganizationInsertInput) (*model.ExternalUserOrganization, error) {
	m.logger.Debug("Insert")
	record := externalUserOrganizationRecordFromInsertInput(input)
	result := m.db.Create(&record)

	if result.Error != nil {
		return nil, result.Error
	}

	output := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: record.ID,
	}).WithName(record.Name).Value()

	m.identityMap[output.ID] = output
	return output, nil
}

// Returns a copy with updated values
// If an error occurs, no change is made to model and nil is returned as first argument
// Guaranteed to be a transactional update
func (m *ExternalUserOrganization) Update(dm *model.ExternalUserOrganization) (*model.ExternalUserOrganization, error) {
	m.logger.Debug("Update")
	updates := make(map[string]interface{})

	if dm.NameChanged() {
		updates["name"] = dm.Name()
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := m.db.Model(&externalUserOrganizationRecord{ID: dm.ID}).Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	// Clone and mark all attributes as unchanged (since they now persist to database)
	clone := dm.Clone(true)
	m.identityMap[dm.ID] = clone
	return clone, nil
}
