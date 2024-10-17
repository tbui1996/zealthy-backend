package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func GetConnectionItemsInput(table string, idColumnName string, IDs []string) (*dynamodb.ScanInput, error) {
	var input *dynamodb.ScanInput
	if len(IDs) == 0 {
		input = &dynamodb.ScanInput{
			TableName: aws.String(table),
		}

		return input, nil
	}

	IDValues := make([]expression.OperandBuilder, 0, len(IDs))
	for _, ID := range IDs {
		IDValues = append(IDValues, expression.Value(ID))
	}

	var condition expression.ConditionBuilder
	if len(IDValues) == 1 {
		condition = expression.Name(idColumnName).In(IDValues[0])
	} else {
		condition = expression.Name(idColumnName).In(IDValues[0], IDValues[1:]...)
	}

	exp, err := expression.NewBuilder().WithFilter(condition).Build()
	if err != nil {
		return nil, err
	}
	input = &dynamodb.ScanInput{
		FilterExpression:          exp.Filter(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		TableName:                 aws.String(table),
	}

	return input, nil
}
