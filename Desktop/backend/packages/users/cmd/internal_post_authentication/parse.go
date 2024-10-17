package main

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

type ParsedEvent struct {
	OktaSonarGroups map[string][]string
}

func parse(event events.CognitoEventUserPoolsPostAuthentication, environment string, logger *zap.Logger) ParsedEvent {
	rawCustomGroups, rawCustomGroupsExist := event.Request.UserAttributes["custom:groups"]

	if !rawCustomGroupsExist {
		return ParsedEvent{
			OktaSonarGroups: nil,
		}
	}

	// get list of groups that the user belongs to in okta
	var groups []string
	err := json.Unmarshal([]byte(rawCustomGroups), &groups)
	if err != nil {
		logger.Error("parsing custom:groups from okta" + err.Error())
		return ParsedEvent{
			OktaSonarGroups: nil,
		}
	}
	logger.Debug("custom:groups " + strings.Join(groups, ", "))

	// filter groups down to 'sonar' groups
	// groups in okta prefixed with 'internals_' are assumed to be sonar groups
	// the expectation is that for a given user, they will have only one group per environment (e.g. dev-yourname, dev, prod, etc.)
	re := regexp.MustCompile(`^(?P<group>internals_[[:word:]]+)(.)(?P<environment>.*)`)

	oktaSonarGroups := make(map[string][]string)
	for _, group := range groups {
		match := re.FindStringSubmatch(group)
		names := re.SubexpNames()

		paramsMap := make(map[string]string)
		for i, m := range match {
			name := names[i]
			paramsMap[name] = m
		}

		env := paramsMap["environment"]
		grp := paramsMap["group"]

		if env == environment && grp != "" {
			if oktaSonarGroups[env] == nil {
				oktaSonarGroups[env] = []string{}
			}
			oktaSonarGroups[env] = append(oktaSonarGroups[env], grp)
		}
	}

	return ParsedEvent{
		OktaSonarGroups: oktaSonarGroups,
	}
}
