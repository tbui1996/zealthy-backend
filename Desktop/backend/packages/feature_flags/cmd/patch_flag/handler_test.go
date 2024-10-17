package main

import (
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/feature_flags/mocks"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type PatchFlagHandlerSuite struct {
	suite.Suite
}

func (s *PatchFlagHandlerSuite) Test__HandlePatchSuccess() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	nameParam := "hello"
	boolParam := true
	input := request.PatchFlagRequest{
		Name:      &nameParam,
		IsEnabled: &boolParam,
		FlagId:    0,
	}
	featureFlag := &model.FeatureFlag{
		Id:        0,
		Name:      "hello1",
		IsEnabled: false,
	}

	mockRepo.On("Find", 0).Return(featureFlag, nil)

	mockRepo.On("Save", featureFlag).Return(nil)

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Save", mock.MatchedBy(func(flag *model.FeatureFlag) bool {
		return flag.IsEnabled == *input.IsEnabled && flag.Name == *input.Name
	}))

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func (s *PatchFlagHandlerSuite) Test__HandlePatchFailNoName() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.PatchFlagRequest{
		Name:      nil,
		IsEnabled: nil,
		FlagId:    0,
	}

	result, _ := Handler(input, deps)

	mockRepo.AssertNotCalled(s.T(), "Find", mock.Anything)
	mockRepo.AssertNotCalled(s.T(), "Save", mock.Anything)
	s.Equal(http.StatusBadRequest, result.StatusCode)
}

func (s *PatchFlagHandlerSuite) Test__HandlePatchFailFind() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	nameParam := "hello"
	boolParam := true
	input := request.PatchFlagRequest{
		Name:      &nameParam,
		IsEnabled: &boolParam,
		FlagId:    1,
	}

	mockRepo.On("Find", 1).Return(nil, flagerror.New("FAKE ERROR, IGNORE", flagerror.UNKNOWN))

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Find", 1)
	mockRepo.AssertNotCalled(s.T(), "Save", mock.Anything)
	s.Nil(err)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func (s *PatchFlagHandlerSuite) Test__HandlePatchFailSave() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	nameParam := "hello"
	boolParam := true
	input := request.PatchFlagRequest{
		Name:      &nameParam,
		IsEnabled: &boolParam,
		FlagId:    1,
	}
	featureFlag := &model.FeatureFlag{
		Id:        1,
		Name:      "hello1",
		IsEnabled: false,
	}

	mockRepo.On("Find", 1).Return(featureFlag, nil)
	mockRepo.On("Save", featureFlag).Return(flagerror.New("FAKE ERROR, IGNORE", flagerror.UNKNOWN))

	result, err := Handler(input, deps)

	mockRepo.AssertCalled(s.T(), "Find", 1)
	mockRepo.AssertCalled(s.T(), "Save", mock.MatchedBy(func(flag *model.FeatureFlag) bool {
		return flag.IsEnabled == *input.IsEnabled && flag.Name == *input.Name
	}))
	s.Nil(err)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func TestCreateFlagHandler(t *testing.T) {
	suite.Run(t, new(PatchFlagHandlerSuite))
}
