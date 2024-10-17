package mapper

import (
	"testing"

	"github.com/circulohealth/sonar-backend/packages/users/pkg/mapper/iface"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type ExternalUserSuite struct {
	suite.Suite
	externalUserSQLAPI          *mockExternalUserSQLAPI
	externalUserCognitoAPI      *mockExternalUserCognitoAPI
	externalUserAPI             *iface.MockExternalUser
	externalUserOrganizationAPI *iface.MockExternalUserOrganization
	registry                    *MockRegistryAPI
}

func (suite *ExternalUserSuite) SetupTest() {
	suite.externalUserSQLAPI = new(mockExternalUserSQLAPI)
	suite.externalUserCognitoAPI = new(mockExternalUserCognitoAPI)
	suite.externalUserAPI = new(iface.MockExternalUser)
	suite.externalUserOrganizationAPI = new(iface.MockExternalUserOrganization)

	suite.registry = new(MockRegistryAPI)
	suite.registry.On("ExternalUser").Return(suite.externalUserAPI)
	suite.registry.On("ExternalUserOrganization").Return(suite.externalUserOrganizationAPI)
	suite.registry.On("externalUserCognito").Return(suite.externalUserCognitoAPI)
	suite.registry.On("externalUserSQL").Return(suite.externalUserSQLAPI)
}

func (suite *ExternalUserSuite) TestFind_QueriesCognitoAndSQLAndMergesResults() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	expected := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithEnabled(true).Value()

	suite.externalUserSQLAPI.On("find", mock.Anything).Return(&externalUserSQLRecord{
		ID:                         expected.ID,
		ExternalUserOrganizationID: nil,
	}, nil)

	suite.externalUserCognitoAPI.On("find", mock.Anything).Return(&externalUserCognitoRecord{
		username:  expected.Username,
		status:    expected.Status,
		email:     expected.Email,
		enabled:   expected.Enabled(),
		firstName: nil,
		lastName:  nil,
		group:     "",
		hasGroup:  false,
	}, nil)

	actual, err := mapper.Find("1")
	suite.NoError(err)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "find", "1")
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "find", "1")
	suite.Equal(expected, actual)
}

func (suite *ExternalUserSuite) TestFind_IncludesStandardAttributes() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	expected := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithEnabled(false).WithFirstName("firstName").WithLastName("lastName").Value()

	suite.externalUserSQLAPI.On("find", mock.Anything).Return(&externalUserSQLRecord{
		ID:                         expected.ID,
		ExternalUserOrganizationID: nil,
	}, nil)

	firstName := expected.FirstName()
	lastName := expected.LastName()
	suite.externalUserCognitoAPI.On("find", mock.Anything).Return(&externalUserCognitoRecord{
		username:  expected.Username,
		status:    expected.Status,
		email:     expected.Email,
		enabled:   expected.Enabled(),
		firstName: &firstName,
		lastName:  &lastName,
		group:     "",
		hasGroup:  false,
	}, nil)

	actual, err := mapper.Find("1")
	suite.NoError(err)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "find", "1")
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "find", "1")
	suite.Equal(expected, actual)
}

func (suite *ExternalUserSuite) TestFind_IncludesGroup() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	expected := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithEnabled(false).WithGroup("test").Value()

	suite.externalUserSQLAPI.On("find", mock.Anything).Return(&externalUserSQLRecord{
		ID:                         expected.ID,
		ExternalUserOrganizationID: nil,
	}, nil)

	suite.externalUserCognitoAPI.On("find", mock.Anything).Return(&externalUserCognitoRecord{
		username:  expected.Username,
		status:    expected.Status,
		email:     expected.Email,
		enabled:   expected.Enabled(),
		firstName: nil,
		lastName:  nil,
		group:     expected.Group(),
		hasGroup:  true,
	}, nil)

	actual, err := mapper.Find("1")
	suite.NoError(err)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "find", "1")
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "find", "1")
	suite.Equal(expected, actual)
}

func (suite *ExternalUserSuite) TestFind_IncludesOrganization() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	organization := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: 0,
	}).WithName("test").Value()

	suite.externalUserOrganizationAPI.On("Find", organization.ID).Return(organization, nil)

	expected := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithEnabled(false).WithOrganization(organization).Value()

	suite.externalUserSQLAPI.On("find", mock.Anything).Return(&externalUserSQLRecord{
		ID:                         expected.ID,
		ExternalUserOrganizationID: &organization.ID,
	}, nil)

	suite.externalUserCognitoAPI.On("find", mock.Anything).Return(&externalUserCognitoRecord{
		username:  expected.Username,
		status:    expected.Status,
		email:     expected.Email,
		enabled:   expected.Enabled(),
		firstName: nil,
		lastName:  nil,
		group:     "",
		hasGroup:  false,
	}, nil)

	actual, err := mapper.Find("1")
	suite.NoError(err)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "find", "1")
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "find", "1")
	suite.Equal(expected, actual)
}

func (suite *ExternalUserSuite) TestUpdate_DoesNotCallNestedMappersIfNoChange() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	original := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).Value()

	suite.externalUserCognitoAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserSQLAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserOrganizationAPI.On("Update", mock.Anything).Return(nil, nil)

	actual, err := mapper.Update(original)

	suite.NoError(err)
	suite.True(original.IsDeepEqual(actual))
	suite.externalUserCognitoAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserSQLAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserOrganizationAPI.AssertNumberOfCalls(suite.T(), "Update", 0)
}

func (suite *ExternalUserSuite) TestUpdate_CallsExternalUserCognitoWithUpdater() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	// tell mocks not to error
	suite.externalUserCognitoAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserSQLAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserOrganizationAPI.On("Update", mock.Anything).Return(nil, nil)

	// build base model
	original := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).Value()

	// set up variables
	nextName := "nextName"
	nextGroup := "nextGroup"
	nextHasGroup := true

	// modify original
	original.SetFirstName(nextName)
	original.SetGroup(nextGroup)

	// invoke update
	actual, err := mapper.Update(original)

	// prepare expectation
	expected := externalUserCognitoRecordUpdater{
		username:         original.Username,
		firstName:        nextName,
		firstNameChanged: true,
		group:            nextGroup,
		groupChanged:     true,
		hasGroup:         nextHasGroup,
		hasGroupChanged:  true,
	}

	suite.NoError(err)
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "update", expected)
	suite.externalUserSQLAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserOrganizationAPI.AssertNumberOfCalls(suite.T(), "Update", 0)
	suite.Equal(actual.FirstName(), nextName)
	suite.False(actual.FirstNameChanged())

	suite.Equal(actual.Group(), nextGroup)
	suite.False(actual.GroupChanged())

	suite.Equal(actual.HasGroup(), nextHasGroup)
	suite.False(actual.HasGroupChanged())
}

// Update calls externalUserSQL with updater
func (suite *ExternalUserSuite) TestUpdate_CallsExternalUserSQLWithUpdater() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	// tell mocks not to error
	suite.externalUserCognitoAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserSQLAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserOrganizationAPI.On("Update", mock.Anything).Return(nil, nil)

	// build base model
	original := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).Value()

	// set up variables
	nextOrganization := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: 0,
	}).Value()

	// modify original
	original.SetOrganization(nextOrganization)

	// invoke update
	actual, err := mapper.Update(original)

	// prepare expectation
	expected := externalUserSQLRecordUpdater{
		id:                                original.ID,
		externalUserOrganizationID:        &nextOrganization.ID,
		externalUserOrganizationIDChanged: true,
	}

	suite.NoError(err)
	suite.externalUserCognitoAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "update", expected)
	suite.externalUserOrganizationAPI.AssertCalled(suite.T(), "Update", nextOrganization)
	suite.Equal(actual.Organization(), nextOrganization)
	suite.False(actual.OrganizationChanged())
}

func (suite *ExternalUserSuite) TestUpdate_CallsNestedMappers() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	// tell mocks not to error
	suite.externalUserCognitoAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserSQLAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserOrganizationAPI.On("Update", mock.Anything).Return(nil, nil)

	// set up variables
	nextOrganization := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: 0,
	}).Value()

	// build base model
	original := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithOrganization(nextOrganization).Value()

	// modify original
	original.Organization().SetName("test")

	// invoke update
	actual, err := mapper.Update(original)

	suite.NoError(err)
	suite.externalUserCognitoAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserSQLAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserOrganizationAPI.AssertCalled(suite.T(), "Update", nextOrganization)
	suite.False(actual.Organization().NameChanged())
}

func (suite *ExternalUserSuite) TestUpdate_RemovesOrganization() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	// tell mocks not to error
	suite.externalUserCognitoAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserSQLAPI.On("update", mock.Anything).Return(nil)
	suite.externalUserOrganizationAPI.On("Update", mock.Anything).Return(nil, nil)

	// set up variables
	organization := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: 0,
	}).Value()

	// build base model
	original := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithOrganization(organization).Value()

	// modify original
	original.SetOrganization(nil)

	// invoke update
	actual, err := mapper.Update(original)

	// prepare expectation
	expected := externalUserSQLRecordUpdater{
		id:                                original.ID,
		externalUserOrganizationID:        nil,
		externalUserOrganizationIDChanged: true,
	}

	suite.NoError(err)
	suite.externalUserCognitoAPI.AssertNumberOfCalls(suite.T(), "update", 0)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "update", expected)
	suite.externalUserOrganizationAPI.AssertNumberOfCalls(suite.T(), "Update", 0)
	suite.Nil(actual.Organization())
}

func (suite *ExternalUserSuite) TestFindAll_QueriesCognitoAndSQLAndMergesResults() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	expected := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).WithEnabled(true).Value()

	suite.externalUserSQLAPI.On("findAll", mock.Anything).Return([]*externalUserSQLRecord{
		{
			ID:                         expected.ID,
			ExternalUserOrganizationID: nil,
		},
	}, nil)

	suite.externalUserCognitoAPI.On("findAll", mock.Anything).Return([]*externalUserCognitoRecord{
		{
			username:  expected.Username,
			status:    expected.Status,
			email:     expected.Email,
			enabled:   expected.Enabled(),
			firstName: nil,
			lastName:  nil,
			group:     "",
			hasGroup:  false,
		},
	}, nil)

	actual, err := mapper.FindAll()
	suite.NoError(err)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "findAll")
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "findAll")
	suite.Len(actual, 1)
	suite.Equal(expected, actual[0])
}

func (suite *ExternalUserSuite) TestFindAll_ReturnsMultipleRecords() {
	mapper := newExternalUser(&newExternalUserInput{
		registry: suite.registry,
		logger:   zaptest.NewLogger(suite.T()),
	})

	organization := model.BuildExternalUserOrganization(&model.BuildExternalUserOrganizationInput{
		ID: 0,
	}).WithName("testOrganization").Value()

	expected0 := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test",
		Username: "test",
		Email:    "test@test.com",
		Status:   "CONFIRMED",
	}).
		WithEnabled(true).
		Value()

	expected1FirstName := "test2FirstName"
	expected1LastName := "test2LastName"
	expected1 := model.BuildExternalUser(&model.BuildExternalUserInput{
		ID:       "test2",
		Username: "test2",
		Email:    "test2@test.com",
		Status:   "CONFIRMED",
	}).
		WithEnabled(true).
		WithFirstName(expected1FirstName).
		WithLastName(expected1LastName).
		WithGroup("testGroup").
		WithOrganization(organization).
		Value()

	suite.externalUserOrganizationAPI.On("Find", mock.Anything).Return(organization, nil)
	suite.externalUserSQLAPI.On("findAll", mock.Anything).Return([]*externalUserSQLRecord{
		{
			ID:                         expected0.ID,
			ExternalUserOrganizationID: nil,
		},
		{
			ID:                         expected1.ID,
			ExternalUserOrganizationID: &expected1.Organization().ID,
		},
	}, nil)

	suite.externalUserCognitoAPI.On("findAll", mock.Anything).Return([]*externalUserCognitoRecord{
		{
			username:  expected0.Username,
			status:    expected0.Status,
			email:     expected0.Email,
			enabled:   expected0.Enabled(),
			firstName: nil,
			lastName:  nil,
			group:     "",
			hasGroup:  false,
		},
		{
			username:  expected1.Username,
			status:    expected1.Status,
			email:     expected1.Email,
			enabled:   expected1.Enabled(),
			firstName: &expected1FirstName,
			lastName:  &expected1LastName,
			group:     expected1.Group(),
			hasGroup:  true,
		},
	}, nil)

	actual, err := mapper.FindAll()
	suite.NoError(err)
	suite.externalUserSQLAPI.AssertCalled(suite.T(), "findAll")
	suite.externalUserCognitoAPI.AssertCalled(suite.T(), "findAll")
	suite.Len(actual, 2)
	suite.Equal(expected0, actual[0])
	suite.Equal(expected1, actual[1])
}

func TestExternalUserSuite(t *testing.T) {
	suite.Run(t, new(ExternalUserSuite))
}
