package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func getAttribute(name string, attributes []*cognitoidentityprovider.AttributeType) *string {
	for _, attribute := range attributes {
		if *attribute.Name == name {
			return attribute.Value
		}
	}

	return nil
}

func getGroupName(userGroupOutput *cognitoidentityprovider.AdminListGroupsForUserOutput) *string {
	// each user should only ever be assigned to 1 group!!!
	var group *string

	if len(userGroupOutput.Groups) > 0 {
		group = userGroupOutput.Groups[0].GroupName
	} else {
		group = aws.String("no_group")
	}

	return group
}
