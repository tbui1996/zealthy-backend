package mapper

type externalUserCognitoAPI interface {
	find(id string) (*externalUserCognitoRecord, error)
	findAll() ([]*externalUserCognitoRecord, error)
	update(record externalUserCognitoRecordUpdater) error
}
