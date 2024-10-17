package mapper

import (
	"fmt"
	"strings"
	"sync"

	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"go.uber.org/zap"
)

// Implements IExternalUser
// Initalize with NewExternalUser, not with struct initialization
// External to package
type ExternalUser struct {
	// private to package
	identityMap map[string]*model.ExternalUser

	logger *zap.Logger

	registry RegistryAPI
}

type externalUserContainer struct {
	cognito *externalUserCognitoRecord
	sql     *externalUserSQLRecord
}

type ExternalUserRecord struct {
	ID       string
	Username string
	Email    string
	Enabled  string
	Status   string
}

type newExternalUserInput struct {
	registry RegistryAPI
	logger   *zap.Logger
}

func newExternalUser(input *newExternalUserInput) *ExternalUser {
	return &ExternalUser{
		identityMap: make(map[string]*model.ExternalUser),
		registry:    input.registry,
		logger:      input.logger,
	}
}

func (m *ExternalUser) buildExternalUserFromRecords(cognitoRecord *externalUserCognitoRecord, sqlRecord *externalUserSQLRecord) *model.ExternalUser {
	builder := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       cognitoRecord.username,
		Username: cognitoRecord.username,
		Email:    cognitoRecord.email,
		Status:   cognitoRecord.status,
	}).WithEnabled(cognitoRecord.enabled)

	if cognitoRecord.firstName != nil {
		builder.WithFirstName(*cognitoRecord.firstName)
	}

	if cognitoRecord.lastName != nil {
		builder.WithLastName(*cognitoRecord.lastName)
	}

	if cognitoRecord.hasGroup {
		builder.WithGroup(cognitoRecord.group)
	}

	if sqlRecord.ExternalUserOrganizationID != nil {
		organization, err := m.registry.ExternalUserOrganization().Find(*sqlRecord.ExternalUserOrganizationID)

		if err != nil {
			m.logger.Error(fmt.Sprintf("ExternalUser has organization %d but it could not be found due to %s", *sqlRecord.ExternalUserOrganizationID, err.Error()))
		} else if organization != nil {
			builder.WithOrganization(organization)
		}
	}

	return builder.Value()
}

func (m *ExternalUser) FindAll() ([]*model.ExternalUser, error) {
	m.logger.Debug("FindAll")
	var wg sync.WaitGroup
	wg.Add(2) // nolint

	var recordsSQL []*externalUserSQLRecord
	var recordsCognito []*externalUserCognitoRecord
	errs := make([]error, 2) // nolint

	go func() {
		defer wg.Done()
		recordsCognito, errs[0] = m.registry.externalUserCognito().findAll()
	}()

	go func() {
		defer wg.Done()
		recordsSQL, errs[1] = m.registry.externalUserSQL().findAll()
	}()

	wg.Wait()

	recordsMap := make(map[string]*externalUserContainer, len(recordsCognito))

	for _, recordCognito := range recordsCognito {
		recordsMap[recordCognito.username] = &externalUserContainer{
			cognito: recordCognito,
		}
	}

	for _, recordSQL := range recordsSQL {
		record, ok := recordsMap[recordSQL.ID]

		if ok {
			record.sql = recordSQL
		}
	}

	values := make([]*model.ExternalUser, 0, len(recordsMap))
	for _, record := range recordsMap {
		if record.cognito == nil {
			// TODO log or error
			continue
		}

		if record.sql == nil {
			// TODO log or error
			continue
		}

		value := m.buildExternalUserFromRecords(record.cognito, record.sql)

		values = append(values, value)
		m.identityMap[value.ID] = value
	}

	return values, nil
}

func (m *ExternalUser) Find(id string) (*model.ExternalUser, error) {
	m.logger.Debug("Find")
	existing, ok := m.identityMap[id]

	if ok {
		if existing != nil {
			return existing, nil
		}

		delete(m.identityMap, id)
	}

	var cognitoRecord *externalUserCognitoRecord
	var sqlRecord *externalUserSQLRecord
	errs := make(chan (error))
	var wg sync.WaitGroup
	wg.Add(2) // nolint

	go func() {
		defer wg.Done()
		record, err := m.registry.externalUserCognito().find(id)
		cognitoRecord = record

		if err != nil {
			errs <- err
		}
	}()

	go func() {
		defer wg.Done()
		record, err := m.registry.externalUserSQL().find(id)
		sqlRecord = record

		if err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	updateErrs := make([]string, 0)
	for err := range errs {
		updateErrs = append(updateErrs, err.Error())
	}

	if len(updateErrs) > 0 {
		errMessage := strings.Join(updateErrs, ", ")
		return nil, fmt.Errorf(errMessage)
	}

	output := m.buildExternalUserFromRecords(cognitoRecord, sqlRecord)
	m.identityMap[id] = output
	return output, nil
}

func (m *ExternalUser) updateSQL(dm *model.ExternalUser) error {
	doSqlUpdate := false
	updater := externalUserSQLRecordUpdater{
		id: dm.ID,
	}

	if dm.OrganizationChanged() {
		if dm.Organization() == nil {
			updater.externalUserOrganizationID = nil
		} else {
			updater.externalUserOrganizationID = &dm.Organization().ID
		}

		updater.externalUserOrganizationIDChanged = true
		doSqlUpdate = true
	}

	var wg sync.WaitGroup
	errs := make(chan (error))
	if doSqlUpdate {
		wg.Add(1)
		go func() {
			defer wg.Done()
			updateErr := m.registry.externalUserSQL().update(updater)
			if updateErr != nil {
				errs <- updateErr
			}
		}()
	}

	// delegate child changes to child mappers
	if dm.Organization() != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, updateErr := m.registry.ExternalUserOrganization().Update(dm.Organization())
			if updateErr != nil {
				errs <- updateErr
			}
		}()
	}

	wg.Wait()
	close(errs)

	updateErrs := make([]string, 0)
	for err := range errs {
		updateErrs = append(updateErrs, err.Error())
	}

	if len(updateErrs) > 0 {
		errMessage := strings.Join(updateErrs, ", ")
		return fmt.Errorf(errMessage)
	}

	return nil
}

func (m *ExternalUser) updateCognito(dm *model.ExternalUser) error {
	doUpdate := false
	updater := externalUserCognitoRecordUpdater{
		username: dm.Username,
	}

	if dm.EnabledChanged() {
		doUpdate = true
		value := dm.Enabled()
		updater.enabled = value
		updater.enabledChanged = true
	}

	if dm.FirstNameChanged() {
		doUpdate = true
		value := dm.FirstName()
		updater.firstName = value
		updater.firstNameChanged = true
	}

	if dm.LastNameChanged() {
		doUpdate = true
		value := dm.LastName()
		updater.lastName = value
		updater.lastNameChanged = true
	}

	if dm.GroupChanged() {
		doUpdate = true
		value := dm.Group()
		updater.group = value
		updater.groupChanged = true
	}

	if dm.HasGroupChanged() {
		doUpdate = true
		value := dm.HasGroup()
		updater.hasGroup = value
		updater.hasGroupChanged = true
	}

	if doUpdate {
		return m.registry.externalUserCognito().update(updater)
	}

	return nil
}

// Returns a copy with updated values
func (m *ExternalUser) Update(dm *model.ExternalUser) (*model.ExternalUser, error) {
	m.logger.Debug("Update")

	var wg sync.WaitGroup
	wg.Add(2) // nolint
	errs := make(chan (error))

	go func() {
		defer wg.Done()
		// update cognito
		err := m.updateCognito(dm)

		if err != nil {
			errs <- err
		}
	}()

	go func() {
		defer wg.Done()
		// update sql
		err := m.updateSQL(dm)

		if err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	updateErrs := make([]string, 0)
	for err := range errs {
		updateErrs = append(updateErrs, err.Error())
	}

	if len(updateErrs) > 0 {
		errMessage := strings.Join(updateErrs, ", ")
		return nil, fmt.Errorf(errMessage)
	}

	clone := dm.Clone(true)
	m.identityMap[dm.ID] = clone
	return clone, nil
}
