package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

/* interface implementation functions */
func (m *MockDatabase) Get(item *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(item)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDatabase) Create(item interface{}) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockDatabase) Update(expression string, item interface{}, key interface{}) error {
	args := m.Called(expression, item, key)
	return args.Error(0)
}

func (m *MockDatabase) Scan(tableName string) (*dynamodb.ScanOutput, error) {
	args := m.Called(tableName)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (m *MockDatabase) Query(item *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *MockDatabase) Delete(item *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(item)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}
