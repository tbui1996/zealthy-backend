package mapper

// an interface internal to this package
type externalUserSQLAPI interface {
	find(id string) (*externalUserSQLRecord, error)
	findAll() ([]*externalUserSQLRecord, error)
	update(updater externalUserSQLRecordUpdater) error
}
