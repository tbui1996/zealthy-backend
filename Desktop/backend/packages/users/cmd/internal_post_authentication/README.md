The Post Authentication lambda is called each time an internal user signs in. When authenticating against Okta, any groups that the user belongs to are passed along. This lambda's job is to keep the user's Okta groups in sync with a mirrored Cognito group.

This is Post Authentication in the sense that the OAuth handshake is done between Cognito and Okta, but we can still deny access to a user by returning an error.

Access is denied in the following cases:

- User is not assigned the client app (does not belong to a group that has the client assigned)
- User belongs to more than one Okta group for a given environment e.g.:
  - User CAN belong to `internals_program_manager.dev` and `internals_program_manager.prod` (one group per environment)
  - User CAN NOT belong to `internals_program_admin.prod` and `internals_program_manager.prod` (we cannot make a determination at this time without implementing a ranking system)
- User is assigned a group in Okta that does not exist in Cognito
  - This should not be possible since everything is provisioned together with terraform, but in the event someone manually generates a group in Okta and assigns a user

https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-post-authentication.html
