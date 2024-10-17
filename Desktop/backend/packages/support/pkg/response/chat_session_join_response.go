package response

import (
	"strconv"
	"strings"
	"time"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/model"
)

type ChatSessionJoinResponse struct {
	ID                 int
	ChatType           model.ChatType
	Created            time.Time
	Status             model.ChatStatus
	UserId             string
	LastRead           *int
	Name               string
	Value              string
	LastMessage        *string
	LastSent           *int
	Notes              *string
	PatientName        string
	PatientLastName    string
	PatientAddress     string
	PatientInsuranceId string
	PatientBirthday    time.Time
	PatientId          int
}

type Visited struct {
	RespIndex int
	ListIndex int
}

type ChatSessionJoinResponses []ChatSessionJoinResponse

const maxDBRowForSingleResult = 2

func (resp ChatSessionJoinResponses) NormalizeSingleResultToDTO() (model.ChatSessionDTO, bool) {
	sess := model.ChatSessionDTO{
		ID:               strconv.Itoa(resp[0].ID),
		CreatedTimestamp: int64(resp[0].Created.Unix()),
		ChatOpen:         resp[0].isOpen(),
		Topic:            resp[0].Value,
		ChatType:         resp[0].ChatType,
	}

	if resp[0].LastMessage != nil && resp[0].LastSent != nil {
		sess.LastMessageTimestamp = int64(*resp[0].LastSent)
		sess.LastMessagePreview = *resp[0].LastMessage
		sess.LastMessageSenderID = resp[0].UserId
	}

	resp[0].getLastReadMessageFromResponse(&sess, nil)
	resp[0].getUserAttributesFromResponse(&sess)

	if len(resp) == maxDBRowForSingleResult {
		resp[1].getLastReadMessageFromResponse(&sess, resp[0].LastSent)
		resp[1].getUserAttributesFromResponse(&sess)
	}

	return sess, resp[0].isPending()
}

func (resp ChatSessionJoinResponses) NormalizeMultipleResultsToDTO() []model.ChatSessionDTOWithPending {
	visited := make(map[int]Visited)
	var list []model.ChatSessionDTOWithPending
	// Sessions are guaranteed to follow one another in the list
	for i := 0; i < len(resp); i++ {
		theResp := resp[i]
		v, ok := visited[resp[i].ID]
		// initialize dto based on first row in response
		if !ok {
			dto := model.ChatSessionDTO{
				ID:               strconv.Itoa(theResp.ID),
				CreatedTimestamp: int64(theResp.Created.Unix()),
				ChatOpen:         theResp.isOpen(),
				Notes:            theResp.Notes,
				ChatType:         theResp.ChatType,
			}
			theResp.getLastReadMessageFromResponse(&dto, nil) // first call needs nil comparison passed in

			list = append(list, model.ChatSessionDTOWithPending{ChatSessionDTO: dto, Pending: theResp.isPending()})
			v = Visited{RespIndex: i, ListIndex: len(list) - 1}
			visited[resp[i].ID] = v
		}
		// update dto with the rest of the rows
		dtoWithPending := list[v.ListIndex]
		theResp.getDescriptorsFromResponse(&dtoWithPending.ChatSessionDTO)
		theResp.getPatientFromResponse(&dtoWithPending.ChatSessionDTO)
		theResp.getLastReadMessageFromResponse(&dtoWithPending.ChatSessionDTO, resp[v.RespIndex].LastSent)
		theResp.getUserAttributesFromResponse(&dtoWithPending.ChatSessionDTO)
		list[v.ListIndex] = model.ChatSessionDTOWithPending{
			ChatSessionDTO: dtoWithPending.ChatSessionDTO,
			Pending:        dtoWithPending.Pending,
		}
	}

	return list
}

func (resp ChatSessionJoinResponse) isOpen() bool {
	var open = true
	if resp.Status == model.CLOSED {
		open = false
	}
	return open
}

func (resp ChatSessionJoinResponse) isPending() bool {
	pending := false
	if resp.Status == model.PENDING {
		pending = true
	}
	return pending
}

func (resp ChatSessionJoinResponse) getUserAttributesFromResponse(dto *model.ChatSessionDTO) {
	user := resp.UserId
	lastRead := resp.LastRead

	// TODO: Not great, but concept of internal / external must remain for now
	if strings.Contains(user, "Okta") {
		dto.InternalUserID = &user

		if lastRead != nil {
			dto.InternalUserLastRead = int64(*lastRead)
		}
	} else {
		dto.UserID = user

		if lastRead != nil {
			dto.UserLastRead = int64(*lastRead)
		}
	}
}

func (resp ChatSessionJoinResponse) getLastReadMessageFromResponse(dto *model.ChatSessionDTO, compareTime *int) {
	lastMessage := resp.LastMessage
	lastSent := resp.LastSent

	if lastMessage != nil && lastSent != nil && compareTime == nil {
		dto.LastMessagePreview = *lastMessage
		dto.LastMessageTimestamp = int64(*lastSent)
		dto.LastMessageSenderID = resp.UserId
	}

	if lastMessage != nil && lastSent != nil && compareTime != nil && *compareTime < *resp.LastSent {
		dto.LastMessagePreview = *lastMessage
		dto.LastMessageTimestamp = int64(*lastSent)
		dto.LastMessageSenderID = resp.UserId
	}
}

func (resp ChatSessionJoinResponse) getDescriptorsFromResponse(dto *model.ChatSessionDTO) {
	theName := resp.Name

	if theName == "TOPIC" {
		dto.Topic = resp.Value
	}

	if theName == "STARRED" {
		b, err := strconv.ParseBool(resp.Value)
		dto.Starred = b && err == nil
	}
}

func (resp ChatSessionJoinResponse) getPatientFromResponse(dto *model.ChatSessionDTO) {
	dto.Patient = model.Patient{
		ID:          resp.PatientId,
		Name:        resp.PatientName,
		LastName:    resp.PatientLastName,
		Address:     resp.PatientAddress,
		InsuranceID: resp.PatientInsuranceId,
		Birthday:    resp.PatientBirthday,
	}

}
