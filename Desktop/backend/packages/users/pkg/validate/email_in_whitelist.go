package validate

import (
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/circulohealth/sonar-backend/packages/common/exception"
	"github.com/circulohealth/sonar-backend/packages/users/pkg/dao/iface"
)

var EmailInWhitelist = func(email string, dao iface.EmailDomainWhitelistRepository, logger *zap.Logger) *exception.SonarError {
	emailDomain := strings.Split(email, "@")

	err := isInWhitelist(emailDomain[1], dao)
	if err == nil {
		return nil
	}

	logger.Error(err.Error())

	err = isInWhitelist(email, dao)
	if err == nil {
		return nil
	}

	logger.Error(err.Error())

	return exception.NewSonarError(http.StatusUnauthorized, err.Error())
}

func isInWhitelist(email string, dao iface.EmailDomainWhitelistRepository) error {
	res, err := dao.GetWhitelistDomain(email)

	if err != nil {
		return err
	}

	if res == nil || res.EmailDomain == "" {
		return errors.New("email domain is not valid, or is not recognized by Circulo")
	}

	return nil
}
