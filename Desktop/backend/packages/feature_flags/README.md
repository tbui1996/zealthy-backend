## API Reference

#### Create a flag

```http
  POST /feature_flags
```

_Parameters_:
| Parameter | Type | Description |
| :-------- | :------- | :------------------------- |
| `name` | `string` | **Required**. Human readable name for the flag. (must be unique) |
| `key` | `string` | **Required**. Used to identify the flag in clietnts. |

_Returns_: 201, 409, 400

#### Get flags

```http
  GET /feature_flags
```

_Parameters_: none

_Returns_: 200
| Property: | Type | Description |
| :-------- | :------- | :-------------------------------- |
| `id` | `in`t | the primary key for the feature flag (referenced as admin)|
| `key` | `string` | the user created key that is referenced in clients |
| `name` | `string` | human readable identifier for the flag |
| `isEnabled` | `bool` | whether or not the flag is turned on |
| `updatedAt` | `string` | iso datetime of when the flag was last updated |
| `createdAt` | `string` | iso datetime of when the flag was created |
| `updatedBy` | `string` | user id of the last user who updated the flag |
| `createdAt` | `string` | user id of the user who created teh flag |

example:

```
{
    result: [
        {
            id: 1,
            key: 'testFlag',
            name: 'Test Flag',
            isEnabled: true,
            updatedAt: 'ISO DATE',
            createdAt: 'ISO DATE',
            updatedBy: 'user-id',
            createdBy: 'user-id'
        }
    ]
}
```

#### Update Flag (partial update):

```http
  PATCH /feature_flags/{id}
```

_Parameters_:
| Property: | Type | Description |
| :-------- | :------- | :-------------------------------- |
| `id` | `int` (route) | the id of the flag being updated|
| `name` | `string` (body)| human readable identifier for the flag |
| `isEnabled` | `bool` (body)| whether or not the flag is turned on |

_Returns_: 200, 400

_note_: at least one body parameter must not be sent

#### Delete Flag (soft delete):

Sets `isDeleted` to true, will not be returned from Get flags endpoint anymore

```http
  DELETE /feature_flags/{id}
```

_Parameters_:
| Property: | Type | Description |
| :-------- | :------- | :-------------------------------- |
| `id` | `int` (route) | the id of the flag being deleted|

_Returns_: 200

#### Evaluate Flags:

Used by clients to get back a dictionary of flagKeys mapped to a boolean that indicates whether or not that feature should be on.

_Note: this endpoint exists in both the internal and external api gateways_

```http
  GET /feature_flags/evaluate
```

_Parameters_: none

_Returns_: 200

| Key:        | Type     | Value     | Type    | Description                                 |
| :---------- | :------- | :-------- | :------ | :------------------------------------------ |
| `[flagKey]` | `string` | isEnabled | boolean | a key value mapping of flagKey to isEnabled |

example:

```
{
    result: {
        testFlag: true,
        realFlag: false
    }
}
```

## Usage Examples

### Sonar Backend
*note*: this requires env vars set up in the lambda for a read only db connection to our rds instance + vpc configuration

```
	db, err := dao.OpenConnectionWithTablePrefix(dao.FeatureFlags)

	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error while creating db: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

	evaluator := featureflags.NewEvaluatorWithDB(db)

	results, err := evaluator.Evaluate()

	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error while evaluating: %s", err.Error()))
		return exception.ErrorMessageApiGatewayV2(http.StatusInternalServerError, err.Error())
	}

    if results.FlagForDefault("testFlag", false) {
        config.Logger.Debug("Test flag is on!")
    }
```

### Web
```
  const { data: evaluatedFlags } = useEvaluateFeatureFlags({
    defaultFlagValues: {
      testFlag: false
    }
  });

  if (evaluatedFlags.testFlag) {
      console.log("Test flag is on!")
  }
```

### Loop

```
  const { data: flags } = useEvaluateFeatureFlags({
    defaultFlagValues: { testFlag: false },
  });

    if (evaluatedFlags.testFlag) {
      console.log("Test flag is on!")
  }
```