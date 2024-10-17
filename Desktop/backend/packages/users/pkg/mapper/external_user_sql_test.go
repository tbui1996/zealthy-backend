package mapper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ExternalUserSQLSuite struct {
	suite.Suite
}

func (s *ExternalUserSQLSuite) TestFind_QueriesForUser() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationID := 1
	recordID := "test"

	row := sqlmock.
		NewRows([]string{"id", "external_user_organization_id"}).
		AddRow(recordID, organizationID)

	theMock.
		ExpectQuery("SELECT * FROM \"users\".\"external_users\" WHERE id = $1").
		WithArgs(recordID).
		WillReturnRows(row)

	m := &externalUserSQL{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	}

	record, err := m.find(recordID)

	expected := &externalUserSQLRecord{
		ID:                         recordID,
		ExternalUserOrganizationID: &organizationID,
	}

	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	s.NoError(err)
	s.Equal(expected, record)
}

func (s *ExternalUserSQLSuite) TestUpdate_UpdatesOrganizationIfChanged() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationID := 1
	recordID := "test"

	theMock.ExpectBegin()
	theMock.
		ExpectExec("UPDATE \"users\".\"external_users\" SET \"external_user_organization_id\"=$1 WHERE \"id\" = $2").
		WithArgs(organizationID, recordID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	theMock.ExpectCommit()

	m := &externalUserSQL{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	}

	updater := externalUserSQLRecordUpdater{
		id:                                recordID,
		externalUserOrganizationID:        &organizationID,
		externalUserOrganizationIDChanged: true,
	}

	err = m.update(updater)

	s.NoError(err)
	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *ExternalUserSQLSuite) TestUpdate_SetsOrganizationNull() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	recordID := "test"

	theMock.ExpectBegin()
	theMock.
		ExpectExec("UPDATE \"users\".\"external_users\" SET \"external_user_organization_id\"=$1 WHERE \"id\" = $2").
		WithArgs(nil, recordID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	theMock.ExpectCommit()

	m := &externalUserSQL{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	}

	updater := externalUserSQLRecordUpdater{
		id:                                recordID,
		externalUserOrganizationID:        nil,
		externalUserOrganizationIDChanged: true,
	}

	err = m.update(updater)

	s.NoError(err)
	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *ExternalUserSQLSuite) TestUpdate_DoesNothingIfNoChanges() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	recordID := "test"

	m := &externalUserSQL{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	}

	updater := externalUserSQLRecordUpdater{
		id:                                recordID,
		externalUserOrganizationID:        nil,
		externalUserOrganizationIDChanged: false,
	}

	err = m.update(updater)

	s.NoError(err)

	// No expectations
	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *ExternalUserSQLSuite) TestFindAll_ShouldFetchAllUsers() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationIDs := []int{1, 2, 3}
	recordIDs := []string{"test1", "test2", "test3"}

	rows := sqlmock.
		NewRows([]string{"id", "external_user_organization_id"})

	for i := range organizationIDs {
		rows.AddRow(recordIDs[i], organizationIDs[i])
	}

	theMock.
		ExpectQuery("SELECT * FROM \"users\".\"external_users\"").
		WillReturnRows(rows)

	m := &externalUserSQL{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	}

	record, err := m.findAll()

	expected := make([]*externalUserSQLRecord, len(organizationIDs))

	for i := range organizationIDs {
		expected[i] = &externalUserSQLRecord{
			ID:                         recordIDs[i],
			ExternalUserOrganizationID: &organizationIDs[i],
		}
	}

	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	s.NoError(err)
	s.Equal(expected, record)
}

func TestExternalUserSQLSuite(t *testing.T) {
	suite.Run(t, new(ExternalUserSQLSuite))
}
