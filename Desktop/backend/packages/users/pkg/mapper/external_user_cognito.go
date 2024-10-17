package mapper

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"go.uber.org/zap"
)

const email = "email"

var requiredAttributes = []string{email}

type newExternalUserCognitoInput struct {
	logger     *zap.Logger
	idp        cognitoidentityprovideriface.CognitoIdentityProviderAPI
	userPoolId string
}

func newExternalUserCognito(input *newExternalUserCognitoInput) *externalUserCognito {
	return &externalUserCognito{
		logger:     input.logger.With(zap.String("mapper", "externalUserCognito")),
		idp:        input.idp,
		userPoolId: input.userPoolId,
	}
}

// Implements externalUserCognitoAPI
type externalUserCognito struct {
	idp cognitoidentityprovideriface.CognitoIdentityProviderAPI

	userPoolId string

	logger *zap.Logger
}

func validateUserType(output *cognitoidentityprovider.UserType) error {
	if output.Enabled == nil {
		return errors.New("Expected Enabled to exist")
	}

	if output.UserStatus == nil {
		return errors.New("Expected UserStatus to exist")
	}

	if output.Username == nil {
		return errors.New("Expected Username to exist")
	}

	// search each required attribute and verify that it exists and is not null
	for _, requiredAttribute := range requiredAttributes {
		found := false
		for _, userAttribute := range output.Attributes {
			if userAttribute.Name != nil && requiredAttribute == *userAttribute.Name && userAttribute.Value != nil {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("Expected to find %s in UserAttributes", requiredAttribute)
		}
	}

	return nil
}

// does not validate whether required attributes exist
func parseUserTypeToRecord(record *externalUserCognitoRecord, output *cognitoidentityprovider.UserType) {
	record.status = *output.UserStatus
	record.username = *output.Username
	record.enabled = *output.Enabled

	for _, attribute := range output.Attributes {
		if attribute.Name != nil {
			switch *attribute.Name {
			case email:
				record.email = *attribute.Value
			case "given_name":
				record.firstName = attribute.Value
			case "family_name":
				record.lastName = attribute.Value
			}
		}
	}
}

// returns group, hasGroup, error
// does not find the group record if there are no groups or more than one group
func (m *externalUserCognito) findCognitoGroupRecord(id string) (string, bool, error) {
	userGroupOutput, err := m.idp.AdminListGroupsForUser(&cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: &m.userPoolId,
		Username:   &id,
	})

	if err != nil {
		return "", false, err
	}

	// each user should only ever be assigned to 1 group!!!
	if len(userGroupOutput.Groups) == 0 {
		return "", false, nil
	}

	if len(userGroupOutput.Groups) > 1 {
		m.logger.Debug(fmt.Sprintf("Expected user %s to only be assigned to 1 group %+v", id, userGroupOutput.Groups))
		return "", false, nil
	}

	if userGroupOutput.Groups[0].GroupName == nil {
		m.logger.Debug(fmt.Sprintf("Expected group name to exist for user %s", id))
		return "", false, errors.New("expected group name to exist")
	}

	return *userGroupOutput.Groups[0].GroupName, true, nil
}

func (m *externalUserCognito) findAll() ([]*externalUserCognitoRecord, error) {
	m.logger.Debug("findAll")
	output, err := m.idp.ListUsers(&cognitoidentityprovider.ListUsersInput{
		// AttributesToGet: []*string{},
		// Filter:          new(string),
		// Limit:           new(int64),
		// PaginationToken: new(string)
		UserPoolId: &m.userPoolId,
	})

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(output.Users))
	// no need to use concurrent mechanism like channel since
	// we know the exact length of the slice (no race condition on writing)
	records := make([]*externalUserCognitoRecord, len(output.Users))
	errs := make([]error, len(output.Users))
	for i, user := range output.Users {
		go func(index int, closuredUser *cognitoidentityprovider.UserType) {
			defer wg.Done()

			record, err := m.buildCompleteRecordFromUserType(closuredUser)

			if err != nil {
				errs[index] = err
				return
			}

			records[index] = record
		}(i, user)
	}

	wg.Wait()

	var concanatedErr error
	for _, err := range errs {
		if err != nil {
			if concanatedErr == nil {
				concanatedErr = err
			} else {
				concanatedErr = fmt.Errorf("%w; %s", concanatedErr, err.Error())
			}
		}
	}

	if concanatedErr != nil {
		return nil, concanatedErr
	}

	return records, nil
}

func (m *externalUserCognito) buildCompleteRecordFromUserType(u *cognitoidentityprovider.UserType) (*externalUserCognitoRecord, error) {
	record := &externalUserCognitoRecord{}
	err := validateUserType(u)

	if err != nil {
		return nil, err
	}

	parseUserTypeToRecord(record, u)

	record.group, record.hasGroup, err = m.findCognitoGroupRecord(*u.Username)

	if err != nil {
		return nil, err
	}

	return record, nil
}

// internal to package
func (m *externalUserCognito) find(username string) (*externalUserCognitoRecord, error) {
	m.logger.Debug("find")
	adminGetUserOutput, err := m.idp.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &m.userPoolId,
		Username:   &username,
	})

	if err != nil {
		m.logger.Error(fmt.Sprintf("find error: %s", err.Error()))
		return nil, err
	}

	userType := &cognitoidentityprovider.UserType{
		Enabled:    adminGetUserOutput.Enabled,
		Username:   adminGetUserOutput.Username,
		UserStatus: adminGetUserOutput.UserStatus,
		Attributes: adminGetUserOutput.UserAttributes,
	}

	m.logger.Debug("found record")
	return m.buildCompleteRecordFromUserType(userType)
}

func (m *externalUserCognito) updateAttributes(record externalUserCognitoRecordUpdater) error {
	userAttributes := make([]*cognitoidentityprovider.AttributeType, 0)
	if record.firstNameChanged {
		userAttributes = append(userAttributes, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("given_name"),
			Value: aws.String(record.firstName),
		})
	}

	if record.lastNameChanged {
		userAttributes = append(userAttributes, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("family_name"),
			Value: aws.String(record.lastName),
		})
	}

	if len(userAttributes) == 0 {
		return nil
	}

	_, err := m.idp.AdminUpdateUserAttributes(&cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: userAttributes,
		UserPoolId:     &m.userPoolId,
		Username:       &record.username,
	})

	return err
}

func (m *externalUserCognito) clearGroups(id string) error {
	groups, err := m.idp.AdminListGroupsForUser(&cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: aws.String(m.userPoolId),
		Username:   aws.String(id),
	})

	if err != nil {
		m.logger.Error(fmt.Sprintf("while clearing groups, failed to list groups for user %s due to %s", id, err.Error()))
		return err
	}

	// First, remove all groups the user currently has
	errs := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(groups.Groups))
	for _, currGroup := range groups.Groups {
		go func(closuredGroup *cognitoidentityprovider.GroupType) {
			defer wg.Done()

			if closuredGroup.GroupName == nil {
				return
			}
			_, err := m.idp.AdminRemoveUserFromGroup(&cognitoidentityprovider.AdminRemoveUserFromGroupInput{
				GroupName: aws.String(*closuredGroup.GroupName),
				Username:  aws.String(id),
			})

			if err != nil {
				errs <- err
			}
		}(currGroup)
	}

	wg.Wait()
	close(errs)

	updateErrs := make([]string, 0)
	for err := range errs {
		updateErrs = append(updateErrs, err.Error())
	}

	if len(updateErrs) > 0 {
		errMessage := strings.Join(updateErrs, ", ")
		return fmt.Errorf(errMessage)
	}

	return nil
}

// Asserts that the user only has 1 group. group must represent a valid cognito group
// Does not handle the case where the user is assigned no groups
func (m *externalUserCognito) updateGroup(id string, group string) error {
	// add user to valid group
	_, err := m.idp.AdminAddUserToGroup(&cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(group),
		Username:   aws.String(id),
		UserPoolId: aws.String(m.userPoolId),
	})

	return err
}

func (m *externalUserCognito) updateEnabled(id string, enabled bool) error {
	if enabled {
		_, err := m.idp.AdminEnableUser(&cognitoidentityprovider.AdminEnableUserInput{
			UserPoolId: aws.String(m.userPoolId),
			Username:   aws.String(id),
		})

		return err
	}

	_, err := m.idp.AdminDisableUser(&cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: aws.String(m.userPoolId),
		Username:   aws.String(id),
	})

	return err
}

func (m *externalUserCognito) update(record externalUserCognitoRecordUpdater) error {
	m.logger.Debug("update")

	err := m.updateAttributes(record)

	if err != nil {
		return err
	}

	// update enabled
	if record.enabledChanged {
		err = m.updateEnabled(record.username, record.enabled)

		if err != nil {
			return err
		}
	}

	if record.hasGroupChanged && !record.hasGroup {
		/* User previously had a group and now has no groups */
		err = m.clearGroups(record.username)

		if err != nil {
			return err
		}
	} else if record.groupChanged {
		/* We don't know if the user previously had a group, but it's irrelevant since we're changing the group */
		err = m.clearGroups(record.username)

		if err != nil {
			return err
		}

		err = m.updateGroup(record.username, record.group)

		if err != nil {
			return err
		}
	}

	return nil
}
