package main

import (
	"net/http"
	"testing"

	"github.com/circulohealth/sonar-backend/packages/feature_flags/mocks"
	flagerror "github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/data/flag_error"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/model"
	"github.com/circulohealth/sonar-backend/packages/feature_flags/pkg/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type CreateFlagHandlerSuite struct {
	suite.Suite
}

func (s *CreateFlagHandlerSuite) Test__ValidatesKeyAndName() {
	deps := HandlerDeps{
		repo:   new(mocks.FeatureFlagRepository),
		logger: zaptest.NewLogger(s.T()),
	}
	input1 := request.CreateFlagRequest{
		Key: "key",
	}

	input2 := request.CreateFlagRequest{
		Name: "key",
	}

	res1, _ := Handler(input1, deps)
	res2, _ := Handler(input2, deps)

	s.Equal(res1.StatusCode, http.StatusBadRequest)
	s.Equal(res2.StatusCode, http.StatusBadRequest)
}

func (s *CreateFlagHandlerSuite) Test__HandlesError() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreateFlagRequest{
		Key:  "key",
		Name: "name",
	}

	mockRepo.On("Save", mock.MatchedBy(func(flag *model.FeatureFlag) bool {
		return flag.IsEnabled == false &&
			flag.CreatedAt == nil &&
			flag.Name == "name" &&
			flag.Key == "key"
	})).Return(flagerror.New("unknown error", flagerror.UNKNOWN))

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusInternalServerError, result.StatusCode)
}

func (s *CreateFlagHandlerSuite) Test__HandlesDuplicate() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreateFlagRequest{
		Key:  "key",
		Name: "name",
	}

	mockRepo.On("Save", mock.MatchedBy(func(flag *model.FeatureFlag) bool {
		return flag.IsEnabled == false &&
			flag.CreatedAt == nil &&
			flag.Name == "name" &&
			flag.Key == "key"
	})).Return(flagerror.New("dupe", flagerror.KEY_CONFLICT))

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusConflict, result.StatusCode)
}

func (s *CreateFlagHandlerSuite) Test__HandlesSuccess() {
	mockRepo := new(mocks.FeatureFlagRepository)
	deps := HandlerDeps{
		repo:   mockRepo,
		logger: zaptest.NewLogger(s.T()),
	}
	input := request.CreateFlagRequest{
		Key:  "key",
		Name: "name",
	}

	mockRepo.On("Save", mock.MatchedBy(func(flag *model.FeatureFlag) bool {
		return flag.IsEnabled == false &&
			flag.CreatedAt == nil &&
			flag.Name == "name" &&
			flag.Key == "key"
	})).Return(nil)

	result, err := Handler(input, deps)

	s.Nil(err)
	s.Equal(http.StatusCreated, result.StatusCode)
}

func TestCreateFlagHandler(t *testing.T) {
	suite.Run(t, new(CreateFlagHandlerSuite))
}
