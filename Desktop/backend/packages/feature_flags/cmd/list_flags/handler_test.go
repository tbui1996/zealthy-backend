package main

import (
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/feature_flags/mocks"
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/response"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type ListFlagsSuite struct {
	suite.Suite
}

func (s *ListFlagsSuite) Test__HandlesError() {
	mockRepo := new(mocks.FeatureFlagRepository)
	logger := zaptest.NewLogger(s.T())

	mockRepo.On("FindAll").Return(nil, flagerror.New("unknown error", flagerror.UNKNOWN))

	deps := HandlerDeps{
		repo:   mockRepo,
		logger: logger,
	}

	results, err := Handler(deps)

	s.Nil(results)
	s.NotNil(err)
}

func (s *ListFlagsSuite) Test__Success() {
	userId := "user-id"
	now := time.Now()
	expected := response.FeatureFlagResponse{
		Key:       "Key1",
		Name:      "Name1",
		IsEnabled: false,
		Id:        1,
		CreatedAt: now.String(),
		UpdatedAt: now.String(),
		CreatedBy: userId,
		UpdatedBy: userId,
	}
	mockRepo := new(mocks.FeatureFlagRepository)
	logger := zaptest.NewLogger(s.T())

	repoResults := []model.FeatureFlag{
		{
			Key:       "Key1",
			Name:      "Name1",
			IsEnabled: false,
			Id:        1,
			CreatedAt: &now,
			UpdatedAt: &now,
			CreatedBy: &userId,
			UpdatedBy: &userId,
		},
	}

	mockRepo.On("FindAll").Return(&repoResults, nil)

	deps := HandlerDeps{
		repo:   mockRepo,
		logger: logger,
	}

	results, err := Handler(deps)

	s.Equal(expected, (*results)[0])
	s.Nil(err)
}

func TestListFlagsHandler(t *testing.T) {
	suite.Run(t, new(ListFlagsSuite))
}
