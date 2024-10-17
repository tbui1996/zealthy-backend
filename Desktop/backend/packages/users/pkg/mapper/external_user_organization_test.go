package mapper

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ExternalUserOrganizationSuite struct {
	suite.Suite
}

func (s *ExternalUserOrganizationSuite) TestFind() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationID := 1
	organizationName := "reddoor"

	row := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(organizationID, organizationName)

	theMock.
		ExpectQuery("SELECT * FROM \"users\".\"external_user_organizations\" WHERE id = $1").
		WithArgs(organizationID).
		WillReturnRows(row)

	m := newExternalUserOrganization(&newExternalUserOrganizationInput{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	})

	record, err := m.Find(organizationID)

	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	expected := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: organizationID,
	}).WithName(organizationName).Value()

	s.NoError(err)
	s.Equal(expected, record)
}

func (s *ExternalUserOrganizationSuite) TestFindAll() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationIDs := []int{1, 2, 3}
	organizationNames := []string{"test0", "test1", "test2"}

	row := sqlmock.
		NewRows([]string{"id", "name"})

	for i := range organizationIDs {
		row.AddRow(organizationIDs[i], organizationNames[i])
	}

	theMock.
		ExpectQuery("SELECT * FROM \"users\".\"external_user_organizations\"").
		WillReturnRows(row)

	m := newExternalUserOrganization(&newExternalUserOrganizationInput{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	})

	record, err := m.FindAll()

	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	expected := make([]*model.ExternalUserOrganization, len(organizationIDs))

	for i := range organizationIDs {
		expected[i] = model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
			ID: organizationIDs[i],
		}).WithName(organizationNames[i]).Value()
	}

	s.NoError(err)
	s.Equal(expected, record)
}

func (s *ExternalUserOrganizationSuite) TestInsert() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationID := 1
	organizationName := "reddoor"

	row := sqlmock.
		NewRows([]string{"id"}).
		AddRow(organizationID)

	theMock.ExpectBegin()
	theMock.
		ExpectQuery("INSERT INTO \"users\".\"external_user_organizations\" (\"name\") VALUES ($1) RETURNING \"id\"").
		WithArgs(organizationName).
		WillReturnRows(row)
	theMock.ExpectCommit()

	m := newExternalUserOrganization(&newExternalUserOrganizationInput{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	})

	record, err := m.Insert(&iface.ExternalUserOrganizationInsertInput{
		Name: organizationName,
	})

	s.NoError(err)

	expected := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: organizationID,
	}).WithName(organizationName).Value()

	s.Equal(expected, record)

	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *ExternalUserOrganizationSuite) TestUpdate() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	organizationID := 1
	organizationName := "reddoor"
	organizationNameNext := "reddoor2"

	original := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: organizationID,
	}).WithName(organizationName).Value()

	original.SetName(organizationNameNext)

	theMock.ExpectBegin()
	theMock.
		ExpectExec("UPDATE \"users\".\"external_user_organizations\" SET \"name\"=$1 WHERE \"id\" = $2").
		WithArgs(organizationNameNext, organizationID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	theMock.ExpectCommit()

	m := newExternalUserOrganization(&newExternalUserOrganizationInput{
		db:     theDB,
		logger: zaptest.NewLogger(s.T()),
	})

	record, err := m.Update(original)

	s.NoError(err)
	s.Equal(organizationID, record.ID)
	s.Equal(organizationNameNext, record.Name())
	s.False(record.NameChanged())

	if err := theMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestExternalUserOrganizationSuite(t *testing.T) {
	suite.Run(t, new(ExternalUserOrganizationSuite))
}
