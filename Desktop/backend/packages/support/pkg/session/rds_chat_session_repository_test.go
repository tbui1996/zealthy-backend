package session

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type RDSChatSessionRepositoryTestSuite struct {
	suite.Suite
}

var Join = fmt.Sprintf(
	"SELECT %s FROM \"sessions\" %s %s %s %s %s %s %s",
	model.JoinSelect,
	model.JoinStatus,
	model.JoinUsers,
	model.JoinDescriptors,
	model.JoinPatientsSession,
	model.JoinPatients,
	model.JoinMessages,
	model.JoinRead,
)

var JoinWithNotes = fmt.Sprintf(
	"SELECT %s FROM \"sessions\" %s %s %s %s %s %s %s %s",
	model.JoinSelectWithNotes,
	model.JoinStatus,
	model.JoinUsers,
	model.JoinDescriptors,
	model.JoinPatientsSession,
	model.JoinPatients,
	model.JoinNotes,
	model.JoinMessages,
	model.JoinRead,
)

func (suite *RDSChatSessionRepositoryTestSuite) TestAssignPending_ShouldNotUpdateAndInsert() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	sessionID := 1
	internalUser := "olive"

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	row := sqlmock.
		NewRows([]string{"session_id", "status"}).
		AddRow(sessionID, model.OPEN)

	theMock.ExpectBegin()
	theMock.
		ExpectQuery("SELECT * FROM \"session_statuses\" WHERE \"session_statuses\".\"session_id\" = $1 ORDER BY \"session_statuses\".\"session_id\" LIMIT 1").
		WithArgs(sessionID).
		WillReturnRows(row)

	sess, _ := repo.AssignPending(sessionID, internalUser)
	suite.Nil(sess)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestAssignPending_ShouldUpdateAndInsert() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	sessionID := 1
	internalUser := "olive"

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	row := sqlmock.
		NewRows([]string{"id"}).
		AddRow(sessionID)

	row2 := sqlmock.
		NewRows([]string{"session_id", "status"}).
		AddRow(sessionID, model.PENDING)

	theMock.ExpectBegin()
	theMock.
		ExpectQuery("SELECT * FROM \"session_statuses\" WHERE \"session_statuses\".\"session_id\" = $1 ORDER BY \"session_statuses\".\"session_id\" LIMIT 1").
		WithArgs(sessionID).
		WillReturnRows(row2)
	theMock.
		ExpectExec("UPDATE \"session_statuses\" SET \"status\"=$1 WHERE session_id = $2").
		WithArgs("OPEN", sessionID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	theMock.
		ExpectQuery("INSERT INTO \"session_users\" (\"user_id\",\"session_id\") VALUES ($1,$2) RETURNING \"id\"").
		WithArgs(internalUser, sessionID).
		WillReturnRows(row)
	theMock.ExpectCommit()

	sess, _ := repo.AssignPending(sessionID, internalUser)
	dataRef := sess.(*RDSChatSession)

	suite.Equal("1", dataRef.dto.ID)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestCreate_ShouldCreateSession() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	externalUser := "Olive"
	internalUser := "Okta_Circulo"

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	sessionRow := sqlmock.
		NewRows([]string{"id"}).
		AddRow(1)

	usersRow := sqlmock.
		NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)

	created := time.Now()

	theMock.ExpectBegin()
	theMock.ExpectQuery("INSERT INTO \"sessions\" (\"created\",\"chat_type\") VALUES ($1,$2) RETURNING \"id\"").
		WithArgs(AnyTime{}, model.CIRCULATOR).
		WillReturnRows(sessionRow)
	theMock.ExpectExec("INSERT INTO \"session_statuses\" (\"session_id\",\"status\") VALUES ($1,$2)").
		WithArgs(1, "OPEN").
		WillReturnResult(sqlmock.NewResult(1, 1))
	theMock.ExpectQuery("INSERT INTO \"session_users\" (\"user_id\",\"session_id\") VALUES ($1,$2),($3,$4) RETURNING \"id\"").
		WithArgs("Olive", 1, "Okta_Circulo", 1).
		WillReturnRows(usersRow)
	theMock.ExpectCommit()

	req := &request.ChatSessionCreateRequest{
		Topic:          "A new chat",
		InternalUserID: &internalUser,
		UserID:         externalUser,
		Created:        created.Unix(),
	}

	sess, _ := repo.Create(req)

	suite.NotNil(sess)
	suite.Equal("Olive", sess.UserID())

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestCreate_ShouldCreatePendingSession() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	user := "Olive"

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	created := time.Now()

	patient := model.Patient{
		Name:        "John",
		LastName:    "Smith",
		Address:     "123 test Apt. 1",
		InsuranceID: "1234567891324",
		Birthday:    time.Now(),
		ProviderId:  "1",
	}

	theMock.ExpectBegin()
	theMock.ExpectQuery("INSERT INTO \"sessions\" (\"created\",\"chat_type\") VALUES ($1,$2) RETURNING \"id\"").
		WithArgs(AnyTime{}, model.CIRCULATOR).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	theMock.ExpectExec("INSERT INTO \"session_statuses\" (\"session_id\",\"status\") VALUES ($1,$2)").
		WithArgs(1, "PENDING").
		WillReturnResult(sqlmock.NewResult(1, 1))
	theMock.ExpectQuery(`SELECT * FROM "patients" WHERE "patients"."insurance_id" = $1 ORDER BY "patients"."id" LIMIT 1`).
		WithArgs(patient.InsuranceID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	theMock.ExpectQuery("INSERT INTO \"patients\" (\"name\",\"last_name\",\"address\",\"insurance_id\",\"birthday\",\"provider_id\") VALUES ($1,$2,$3,$4,$5,$6) RETURNING \"id\"").
		WithArgs(patient.Name, patient.LastName, patient.Address, patient.InsuranceID, patient.Birthday, patient.ProviderId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	theMock.ExpectExec("INSERT INTO \"session_patients\" (\"session_id\",\"patient_id\") VALUES ($1,$2)").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	theMock.ExpectQuery("INSERT INTO \"session_users\" (\"user_id\",\"session_id\") VALUES ($1,$2) RETURNING \"id\"").
		WithArgs("Olive", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	theMock.ExpectCommit()

	t := model.CIRCULATOR
	req := &model.PendingChatSessionCreate{
		UserID:      user,
		Created:     created.Unix(),
		Description: &t,
		Patient:     &patient,
	}

	sess, _ := repo.CreatePending(req)

	suite.NotNil(sess)
	suite.Equal(user, sess.UserID())

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetEntity_ShouldReturnAResponseWhenCirculoCreatesChat() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	created := time.Now()

	row := sqlmock.
		NewRows([]string{"id", "chat_type", "created", "status", "user_id", "last_read", "value", "last_message", "last_sent"}).
		AddRow(1, model.GENERAL, created, "OPEN", "Olive", 1, "Jenny Jones", "World", 2).
		AddRow(1, model.GENERAL, created, "OPEN", "Okta_Circulo", 1, "Jenny Jones", "Hello", 1)

	query := fmt.Sprintf(
		"%s WHERE %s ORDER BY %s",
		Join,
		strings.Replace(model.WhereSession, "@id", "$1", 1),
		model.OrderSent,
	)

	theMock.ExpectPrepare(query)
	theMock.ExpectQuery(query).WithArgs("1").WillReturnRows(row)

	sess, err := repo.GetEntity("1")

	deRef := *sess

	suite.Nil(err)
	suite.Equal("World", deRef.dto.LastMessagePreview)
	suite.Equal("Okta_Circulo", *deRef.dto.InternalUserID)
	suite.Equal("Olive", deRef.dto.UserID)
	suite.Equal(model.GENERAL, deRef.Type())

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetEntity_ShouldReturnAResponseWhenExternalEntityCreatesChat() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	row := sqlmock.
		NewRows([]string{"id", "chat_type", "created", "status", "user_id", "last_read", "value", "last_message", "last_sent"}).
		AddRow(1, model.CIRCULATOR, time.Now(), "PENDING", "Olive", 1, "Jenny Jones", "Hello", 2)

	query := fmt.Sprintf(
		"%s WHERE %s ORDER BY %s",
		Join,
		strings.Replace(model.WhereSession, "@id", "$1", 1),
		model.OrderSent,
	)

	theMock.ExpectPrepare(query)
	theMock.ExpectQuery(query).WithArgs("1").WillReturnRows(row)

	sess, err := repo.GetEntity("1")

	deRef := *sess

	suite.Nil(err)
	suite.Equal("Hello", deRef.dto.LastMessagePreview)
	suite.Equal("Olive", deRef.UserID())
	suite.Equal("", deRef.InternalUserID())
	suite.Equal(model.CIRCULATOR, deRef.Type())

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetEntity_ResponseIsValidWhenLastAttributesAreNil() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	row := sqlmock.
		NewRows([]string{"id", "chat_type", "created", "status", "user_id", "last_read", "value", "last_message", "last_sent"}).
		AddRow(1, model.GENERAL, time.Now(), "PENDING", "Olive", nil, "Jenny Jones", nil, nil)

	query := fmt.Sprintf(
		"%s WHERE %s ORDER BY %s",
		Join,
		strings.Replace(model.WhereSession, "@id", "$1", 1),
		model.OrderSent,
	)

	theMock.ExpectPrepare(query)
	theMock.ExpectQuery(query).WithArgs("1").WillReturnRows(row)

	sess, err := repo.GetEntity("1")

	deRef := *sess

	suite.Nil(err)
	suite.Equal("", deRef.dto.LastMessagePreview)
	suite.Equal("Olive", deRef.UserID())
	suite.Equal("", deRef.InternalUserID())
	suite.Equal(model.GENERAL, deRef.Type())

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetEntities_ReturnNormalizedResultsFromMultipleRows() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	created := time.Now()

	row := sqlmock.
		NewRows([]string{"id", "chat_type", "created", "status", "user_id", "last_read", "name", "value", "last_message", "last_sent"}).
		AddRow(3, model.CIRCULATOR, created, "PENDING", "Olive", 1, "TOPIC", "Jenny Jones", "Hello", 1).
		AddRow(2, model.CIRCULATOR, created, "OPEN", "Okta_user", 1, "TOPIC", "Test User", "How bad?", 2).
		AddRow(2, model.CIRCULATOR, created, "OPEN", "Person", 1, "TOPIC", "Test User", "Not feeling good :(", 1).
		AddRow(1, model.CIRCULATOR, created, "CLOSED", "Subject", 1, "TOPIC", "Steve Austin", "Stone Cold", 2).
		AddRow(1, model.CIRCULATOR, created, "CLOSED", "Okta_title", 1, "TOPIC", "Steve Austin", "How are you?", 1)

	query := fmt.Sprintf(
		"%s WHERE %s ORDER BY %s",
		JoinWithNotes,
		strings.Replace(strings.Replace(model.WhereRole, "@chat_type", "$1", 1), "@id", "$2", 1),
		model.OrderCreatedSent,
	)

	theMock.ExpectPrepare(query)
	theMock.ExpectQuery(query).WillReturnRows(row)

	sess, err := repo.GetEntities("1", model.ChatTypeToString(model.CIRCULATOR))
	suite.Nil(err)
	suite.Equal(3, len(sess))

	for _, v := range sess {
		suite.NotNil(v.(*RDSChatSession).dto)
		suite.Equal(model.CIRCULATOR, v.Type())
	}

	suite.Equal("3", sess[0].(*RDSChatSession).dto.ID)
	suite.Equal("2", sess[1].(*RDSChatSession).dto.ID)
	suite.Equal("1", sess[2].(*RDSChatSession).dto.ID)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetEntitiesByExternalId_ReturnNormalizedResultsFromMultipleRows() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	created := time.Now()

	row := sqlmock.NewRows([]string{"id", "chat_type", "created", "status", "user_id", "last_read", "name", "value", "last_message", "last_sent"}).
		AddRow(2, model.CIRCULATOR, created, "OPEN", "Okta_user", 1, "TOPIC", "Test User", "bleep bloop", 2).
		AddRow(2, model.CIRCULATOR, created, "OPEN", "Person", 1, "TOPIC", "Test User", "bloop bleep", 1).
		AddRow(2, model.CIRCULATOR, created, "OPEN", "Okta_user", 1, "STARRED", "true", "bleep bloop", 2).
		AddRow(2, model.CIRCULATOR, created, "OPEN", "Person", 1, "STARRED", "true", "bloop bleep", 1).
		AddRow(1, model.CIRCULATOR, created, "CLOSED", "Person", 1, "TOPIC", "Hello World", "World", 2).
		AddRow(1, model.CIRCULATOR, created, "CLOSED", "Okta_user2", 1, "TOPIC", "Hello World", "Hello", 1).
		AddRow(1, model.CIRCULATOR, created, "CLOSED", "Person", 1, "STARRED", "false", "World", 2).
		AddRow(1, model.CIRCULATOR, created, "CLOSED", "Okta_user2", 1, "STARRED", "false", "Hello", 1)

	query := fmt.Sprintf(
		"%s WHERE %s ORDER BY %s",
		Join,
		strings.Replace(model.WhereIdByUser, "@id", "$1", 1),
		model.OrderCreatedSent,
	)

	theMock.ExpectPrepare(query)
	theMock.ExpectQuery(query).WithArgs("Person").WillReturnRows(row)

	sess, err := repo.GetEntitiesByExternalID("Person")

	suite.Nil(err)
	suite.Equal(2, len(sess))
	suite.True(sess[0].ToResponseDTO().Starred)
	suite.False(sess[1].ToResponseDTO().Starred)

	for _, v := range sess {
		suite.Equal(model.CIRCULATOR, v.Type())
		suite.Equal("Person", v.(*RDSChatSession).dto.UserID)
		suite.False(v.(*RDSChatSession).pending)
	}

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetStatus_ReturnsStatus() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	row := sqlmock.NewRows([]string{"session_id", "status"}).
		AddRow(1, "OPEN")

	theMock.ExpectQuery("SELECT * FROM \"session_statuses\" WHERE session_id = $1 ORDER BY \"session_statuses\".\"session_id\" LIMIT 1").
		WithArgs("1").
		WillReturnRows(row)

	sess, status, err := repo.GetEntityWithStatus("1")

	suite.Nil(err)
	suite.NotNil(status)
	suite.Equal(model.OPEN, *status)
	suite.False(sess.(*RDSChatSession).pending)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetOtherUserInChat_ReturnsUsers() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	uRow := sqlmock.NewRows([]string{"id", "user_id", "session_id"}).
		AddRow(1, "Okta_test", 1).
		AddRow(2, "test", 1)

	theMock.ExpectBegin()

	theMock.ExpectQuery("SELECT * FROM \"session_users\" WHERE session_id = $1").
		WithArgs("1").
		WillReturnRows(uRow)

	sRow := sqlmock.NewRows([]string{"session_id", "status"}).
		AddRow(1, model.OPEN)

	theMock.ExpectQuery("SELECT * FROM \"session_statuses\" WHERE session_id = $1 LIMIT 1").
		WithArgs("1").
		WillReturnRows(sRow)

	rRow := sqlmock.NewRows([]string{"id", "created", "chat_type"}).
		AddRow(1, time.Now(), model.CIRCULATOR)

	theMock.ExpectQuery("SELECT * FROM \"sessions\" WHERE id = $1 LIMIT 1").
		WithArgs("1").
		WillReturnRows(rRow)

	theMock.ExpectCommit()

	result, err := repo.GetEntityWithUsers("1")

	suite.Equal("Okta_test", result.InternalUserID())
	suite.Equal("test", result.UserID())
	suite.Nil(err)

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RDSChatSessionRepositoryTestSuite) TestGetOtherUserInChat_ReturnsOneUser() {
	db, theMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	theDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	repo := &RDSChatSessionRepository{
		DB:                    theDB,
		ChatMessageRepository: new(MockChatMessageRepository),
	}

	row := sqlmock.NewRows([]string{"id", "user_id", "session_id"}).
		AddRow(1, "test", 1)

	theMock.ExpectBegin()

	theMock.ExpectQuery("SELECT * FROM \"session_users\" WHERE session_id = $1").
		WithArgs("1").
		WillReturnRows(row)

	sRow := sqlmock.NewRows([]string{"session_id", "status"}).
		AddRow(1, model.OPEN)

	theMock.ExpectQuery("SELECT * FROM \"session_statuses\" WHERE session_id = $1 LIMIT 1").
		WithArgs("1").
		WillReturnRows(sRow)

	rRow := sqlmock.NewRows([]string{"id", "created", "chat_type"}).
		AddRow(1, time.Now(), model.CIRCULATOR)

	theMock.ExpectQuery("SELECT * FROM \"sessions\" WHERE id = $1 LIMIT 1").
		WithArgs("1").
		WillReturnRows(rRow)

	theMock.ExpectCommit()

	result, err := repo.GetEntityWithUsers("1")

	suite.Nil(err)
	suite.NotNil(result)
	suite.Empty(result.InternalUserID())

	if err := theMock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRDSChatSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RDSChatSessionRepositoryTestSuite))
}
