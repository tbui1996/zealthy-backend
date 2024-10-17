package session

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/response"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
	"gorm.io/gorm"
)

type RDSChatSessionRepository struct {
	DB                    *gorm.DB
	ChatMessageRepository iface.ChatMessageRepository
	DynamoDB              *dynamodb.DynamoDB
	EventBridge           eventbridgeiface.EventBridgeAPI
}

func NewRDSChatSessionRepositoryWithSession(session *session.Session) (*RDSChatSessionRepository, error) {
	newDynamoDB := dynamodb.New(session)

	chatMessageRepository := NewDynamoDBChatMessageRepositoryWithDB(newDynamoDB)

	gormDB, err := dao.OpenConnectionWithTablePrefix(dao.Chat)

	if err != nil {
		return nil, err
	}

	return &RDSChatSessionRepository{
		DB:                    gormDB,
		ChatMessageRepository: chatMessageRepository,
		DynamoDB:              newDynamoDB,
		EventBridge:           eventbridge.New(session),
	}, nil
}

func (repo *RDSChatSessionRepository) AssignPending(sessionId int, internalUserId string) (iface.ChatSession, error) {

	var possibleSession model.SessionStatus
	tx := repo.DB.Begin()

	result := tx.First(&possibleSession, sessionId)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if possibleSession.Status != model.PENDING {
		errMsg := fmt.Errorf("chat session has already been assigned")
		return nil, errMsg
	}

	result = tx.Model(&model.SessionStatus{}).Where("session_id = ?", sessionId).Update("status", model.OPEN)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	result = tx.Create(&model.SessionUser{UserId: internalUserId, SessionId: sessionId})

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	tx.Commit()

	dto := &model.ChatSessionDTO{
		ID:             strconv.Itoa(sessionId),
		InternalUserID: &internalUserId,
	}

	return &RDSChatSession{
		pending:               false,
		db:                    repo.DB,
		chatMessageRepository: repo.ChatMessageRepository,
		dbChatMessages:        repo.DynamoDB,
		dto:                   dto,
		eventBridge:           repo.EventBridge,
	}, nil
}

func (repo *RDSChatSessionRepository) Create(entity *request.ChatSessionCreateRequest) (iface.ChatSession, error) {
	// TODO: this is actually never called. We are NOT creating sessions from sonar-web
	return repo.createSession(model.OPEN, model.CIRCULATOR, nil, nil, entity.UserID, *entity.InternalUserID)
}

func (repo *RDSChatSessionRepository) CreatePending(entity *model.PendingChatSessionCreate) (iface.ChatSession, error) {
	// TODO: Topic is deprecated, remove ASAP
	return repo.createSession(model.PENDING, *entity.Description, entity.Topic, entity.Patient, entity.UserID)
}

// This method currently assumes that the Internal user is always first, and the external user is the second.
func (repo *RDSChatSessionRepository) createSession(status model.ChatStatus, description model.ChatType, topic *string, patientReq *model.Patient, users ...string) (*RDSChatSession, error) {
	var newPatient *model.Patient
	newSession := model.Session{Created: time.Now(), ChatType: description}
	allUsers := make([]model.SessionUser, len(users))
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newSession).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.SessionStatus{SessionId: newSession.ID, Status: status}).Error; err != nil {
			return err
		}
		if patientReq != nil {
			newPatient = patientReq
			if err := tx.Where(model.Patient{InsuranceID: patientReq.InsuranceID}).FirstOrCreate(&newPatient).Error; err != nil {
				return err
			}

			if err := tx.Create(&model.SessionPatient{SessionId: newSession.ID, PatientId: newPatient.ID}).Error; err != nil {
				return err
			}
		}
		if topic != nil && patientReq == nil {
			if err := tx.Create(&model.SessionDescriptor{SessionId: newSession.ID, Name: model.TOPIC, Value: *topic}).Error; err != nil {
				return err
			}
		}
		for index, value := range users {
			allUsers[index] = model.SessionUser{SessionId: newSession.ID, UserId: value}
		}
		if err := tx.Create(&allUsers).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	dto := &model.ChatSessionDTO{ID: strconv.Itoa(newSession.ID), CreatedTimestamp: newSession.Created.Unix(), UserID: allUsers[0].UserId, ChatOpen: true}
	if newPatient != nil {
		dto.Patient = *newPatient
	}
	if topic != nil {
		dto.Topic = *topic
	}

	pending := true
	if status != model.PENDING {
		pending = false
		dto.InternalUserID = &allUsers[1].UserId
	}
	return &RDSChatSession{
		pending:               pending,
		db:                    repo.DB,
		dto:                   dto,
		chatMessageRepository: repo.ChatMessageRepository,
		dbChatMessages:        repo.DynamoDB,
		eventBridge:           repo.EventBridge,
	}, nil
}

// GetEntity TODO: this function is only called in unit tests, delete?
func (repo *RDSChatSessionRepository) GetEntity(sessionId string) (*RDSChatSession, error) {
	var joinResponse response.ChatSessionJoinResponses
	tx := repo.DB.Session(&gorm.Session{PrepareStmt: true})
	result := tx.Model(&model.Session{}).
		Select(model.JoinSelect).
		Joins(model.JoinStatus).
		Joins(model.JoinUsers).
		Joins(model.JoinDescriptors).
		Joins(model.JoinPatientsSession).
		Joins(model.JoinPatients).
		Joins(model.JoinMessages).
		Joins(model.JoinRead).
		Where(model.WhereSession, sql.Named("id", sessionId)).
		Order(model.OrderSent).
		Scan(&joinResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	chatSessionDto, pending := joinResponse.NormalizeSingleResultToDTO()

	return &RDSChatSession{
		db:                    repo.DB,
		dto:                   &chatSessionDto,
		chatMessageRepository: repo.ChatMessageRepository,
		dbChatMessages:        repo.DynamoDB,
		pending:               pending,
		eventBridge:           repo.EventBridge,
	}, nil
}

func (repo *RDSChatSessionRepository) GetEntities(userID string, chatType string) ([]iface.ChatSession, error) {
	var joinResponse response.ChatSessionJoinResponses
	tx := repo.DB.Session(&gorm.Session{PrepareStmt: true})
	result := tx.Model(&model.Session{}).
		Select(model.JoinSelectWithNotes).
		Joins(model.JoinStatus).
		Joins(model.JoinUsers).
		Joins(model.JoinDescriptors).
		Joins(model.JoinPatientsSession).
		Joins(model.JoinPatients).
		Joins(model.JoinNotes).
		Joins(model.JoinMessages).
		Joins(model.JoinRead).
		Where(model.WhereRole, sql.Named("chat_type", chatType), sql.Named("id", userID)).
		Order(model.OrderCreatedSent).
		Scan(&joinResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	chatSessions := joinResponse.NormalizeMultipleResultsToDTO()

	l := make([]iface.ChatSession, len(chatSessions))
	for i := 0; i < len(chatSessions); i++ {
		sess := &RDSChatSession{
			db:                    repo.DB,
			dto:                   &chatSessions[i].ChatSessionDTO,
			chatMessageRepository: repo.ChatMessageRepository,
			dbChatMessages:        repo.DynamoDB,
			pending:               chatSessions[i].Pending,
			eventBridge:           repo.EventBridge,
		}

		l[i] = sess
	}

	return l, nil
}

func (repo *RDSChatSessionRepository) GetEntitiesByExternalID(userId string) ([]iface.ChatSession, error) {
	var joinResponse response.ChatSessionJoinResponses
	tx := repo.DB.Session(&gorm.Session{PrepareStmt: true})
	result := tx.Model(&model.Session{}).
		Select(model.JoinSelect).
		Joins(model.JoinStatus).
		Joins(model.JoinUsers).
		Joins(model.JoinDescriptors).
		Joins(model.JoinPatientsSession).
		Joins(model.JoinPatients).
		Joins(model.JoinMessages).
		Joins(model.JoinRead).
		Where(model.WhereIdByUser, sql.Named("id", userId)).
		Order(model.OrderCreatedSent).
		Scan(&joinResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	chatSessions := joinResponse.NormalizeMultipleResultsToDTO()

	l := make([]iface.ChatSession, 0, len(chatSessions))
	for i := 0; i < len(chatSessions); i++ {
		sess := &RDSChatSession{
			db:                    repo.DB,
			dto:                   &chatSessions[i].ChatSessionDTO,
			chatMessageRepository: repo.ChatMessageRepository,
			dbChatMessages:        repo.DynamoDB,
			pending:               chatSessions[i].Pending,
			eventBridge:           repo.EventBridge,
		}

		l = append(l, sess)
	}

	return l, nil
}

func (repo *RDSChatSessionRepository) GetEntityWithStatus(sessionId string) (iface.ChatSession, *model.ChatStatus, error) {
	var status model.SessionStatus
	result := repo.DB.First(&status, "session_id = @id", sql.Named("id", sessionId))

	if result.Error != nil {
		return nil, nil, result.Error
	}

	pending := false
	if status.Status == model.PENDING {
		pending = true
	}

	dto := model.ChatSessionDTO{
		ID: strconv.Itoa(status.SessionId),
	}

	sess := RDSChatSession{
		db:                    repo.DB,
		dto:                   &dto,
		chatMessageRepository: repo.ChatMessageRepository,
		dbChatMessages:        repo.DynamoDB,
		pending:               pending,
		eventBridge:           repo.EventBridge,
	}

	return &sess, &status.Status, nil
}

func (repo *RDSChatSessionRepository) GetEntityWithUsers(sessionId string) (iface.ChatSession, error) {
	var users []model.SessionUser
	var status model.SessionStatus
	var s model.Session
	const maxAllowedUsers = 2

	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("session_id = @id", sql.Named("id", sessionId)).Find(&users).Error; err != nil {
			return err
		}

		if err := tx.Where("session_id = @id", sql.Named("id", sessionId)).Take(&status).Error; err != nil {
			return err
		}

		if err := tx.Where("id = @id", sql.Named("id", sessionId)).Take(&s).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	dto := model.ChatSessionDTO{
		ID:       sessionId,
		ChatType: s.ChatType,
	}

	if len(users) == maxAllowedUsers {
		if strings.Contains(users[0].UserId, "Okta") {
			dto.InternalUserID = &users[0].UserId
			dto.UserID = users[1].UserId
		} else {
			dto.InternalUserID = &users[1].UserId
			dto.UserID = users[0].UserId
		}
	} else {
		dto.UserID = users[0].UserId
	}

	pending := false
	if status.Status == model.PENDING {
		pending = true
	}

	sess := RDSChatSession{
		db:                    repo.DB,
		dto:                   &dto,
		chatMessageRepository: repo.ChatMessageRepository,
		dbChatMessages:        repo.DynamoDB,
		pending:               pending,
		eventBridge:           repo.EventBridge,
	}

	return &sess, nil
}

func (repo *RDSChatSessionRepository) HandleReadReceipt(session iface.ChatSession, readTime int64, request request.ReadReceiptRequest) error {
	tx := repo.DB.Begin()

	result := tx.Exec(
		model.LastReadUpsert,
		sql.Named("read", readTime),
		sql.Named("id", session.ID()),
		sql.Named("user", request.UserID),
	)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to store timestamp for session %s (%s)", request.SessionID, result.Error)
	}

	tx.Commit()

	return nil
}
