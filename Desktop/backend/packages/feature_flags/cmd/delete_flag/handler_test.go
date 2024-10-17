package main

import (
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"net/http"
	"testing"
	"time"

	"github.com/circulohealth/sonar-backend/packages/feature_flags/mocks"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

type DeleteFlagHandlerSuite struct {
	suite.Suite
}

func (s *DeleteFlagHandlerSuite) Test__HandleSoftDeleteSuccess() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteFlagRequest{
		FlagId: 1,
	}
	featureFlag := &model.FeatureFlag{
		Id:        1,
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}

	mockRepo.On("Find", 1).Return(featureFlag, nil)

	mockRepo.On("Delete", featureFlag).Return(nil)

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Delete", mock.Anything)

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func (s *DeleteFlagHandlerSuite) Test__HandleSoftDeleteFailId() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteFlagRequest{
		FlagId: 0,
	}

	result, _ := Handler(input, deps)

	mockRepo.AssertNotCalled(s.T(), "Find", mock.Anything)
	mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything)
	s.Equal(http.StatusBadRequest, result.StatusCode)
}

func (s *DeleteFlagHandlerSuite) Test__HandleSoftDeleteFailFind() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteFlagRequest{
		FlagId: 1,
	}
	featureFlag := &model.FeatureFlag{
		Id:        1,
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}

	mockRepo.On("Find", 1).Return(nil, flagerror.New("FAKE ERROR, IGNORE", flagerror.NOT_FOUND))
	mockRepo.On("Delete", featureFlag).Return(nil)

	result, _ := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Find", mock.Anything)
	mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything)

	s.Equal(http.StatusNotFound, result.StatusCode)
}

func (s *DeleteFlagHandlerSuite) Test__HandleSoftDeleteFailDelete() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.DeleteFlagRequest{
		FlagId: 1,
	}
	featureFlag := &model.FeatureFlag{
		Id:        1,
		DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}

	mockRepo.On("Find", 1).Return(featureFlag, nil)
	err := flagerror.New("FAKE ERROR, IGNORE", flagerror.UNKNOWN)
	mockRepo.On("Delete", featureFlag).Return(err)

	result, _ := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Delete", mock.Anything)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func TestDeleteFlagHandler(t *testing.T) {
	suite.Run(t, new(DeleteFlagHandlerSuite))
}
