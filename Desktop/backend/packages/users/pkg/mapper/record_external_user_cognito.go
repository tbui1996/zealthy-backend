package mapper

// nil values indicate that no value exists
type externalUserCognitoRecord struct {
	/* Base data */
	// Cannot be modified
	username string
	status   string

	// Can be modified
	enabled bool

	/* UserAttributes */
	email     string
	firstName *string
	lastName  *string

	/* Derived */
	group string
	// Whereas some string attributes, such as FirstName, may be suitable to have empty strings
	// even when there isn't a "first name", the group needs to know if there really "isn't" a
	// group
	hasGroup bool
}

// corresponds with a "diff" of values
type externalUserCognitoRecordUpdater struct {
	// required
	username string

	enabled        bool
	enabledChanged bool

	// although this can be null if never previously set (see externalUserCognitoRecord)
	// it can't be reset to nil, we'll just set it to an empty string
	firstName        string
	firstNameChanged bool

	// although this can be null if never previously set (see externalUserCognitoRecord)
	// it can't be reset to nil, we'll just set it to an empty string
	lastName        string
	lastNameChanged bool

	group        string
	groupChanged bool

	hasGroup        bool
	hasGroupChanged bool
}
