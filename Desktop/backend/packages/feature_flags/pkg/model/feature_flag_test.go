package model

import (
	"fmt"
	"github.com/circulohealth/sonar-backend/packages/common/dao"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/suite"
)

type FeatureFlagModelSuite struct {
	suite.Suite
}

func (s *FeatureFlagModelSuite) Test__IsNew() {
	model := &FeatureFlag{}
	now := time.Now()

	s.True(model.IsNew())

	model.CreatedAt = &now

	s.False(model.IsNew())
	s.Equal(fmt.Sprintf("%sflags", dao.FeatureFlags), model.TableName())
}

func (s *FeatureFlagModelSuite) Test__NewFeatureFlagWithUserId() {
	model := NewFeatureFlagWithUserId(FeatureFlag{
		Key:  "TestFlag",
		Name: "TestName",
	}, aws.String("test-id"))

	s.Equal(model.Key, "TestFlag")
	s.Equal(model.Name, "TestName")
	s.Equal(model.IsEnabled, false)
	s.Equal(*model.CreatedBy, "test-id")
	s.Equal(*model.UpdatedBy, "test-id")
	s.Nil(model.CreatedAt)
	s.Nil(model.UpdatedAt)
}

func TestFeatureFlagModel(t *testing.T) {
	suite.Run(t, new(FeatureFlagModelSuite))
}
