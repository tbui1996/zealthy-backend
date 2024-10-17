# Sonar - Users

# Setting up Okta Dev Account For Sonar Developers

Okta is used as the Identity provider for Cognito on the [sonar-web](https://gitlab.com/circulohealth/sonar/web). In order to develop and have the same auth flow as production environments, each developer will need to do a few manual steps to set up an Okta developer account in order to terraform to provision the necessary resources.

Follow the steps below and you should be able to sign into `sonar-web` running locally via a bookmark on your user page in Okta.

1. [Create an Okta Developer Account](#create-an-okta-developer-account)
2. [Generate an API token](#generate-an-api-token)
3. [Set up environment variables](#set-up-environment-variables)

## Create an Okta Developer Account

Every developer will need to set up their own Okta Developer Account.

1. Visit: [https://developer.okta.com/signup/](https://developer.okta.com/signup/)
2. Click `CONTINUE WITH GOOGLE`
3. If prompted with a `Tell us more about yourself` pop-up, select a `Country/Region` and `State`
4. Select your `*@circulohealth.com` email to sign in and create an Okta developer account
5. You may need to verify your email to complete account set up

## Generate an API token

[Official docs](https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/)

1. From your developer account, select `Security` -> `API`
2. Select the `Tokens` tab
3. Click the `Create Token` button
4. Give it a name, e.g. `sonar-internal`
5. Copy the API Token, it will not be available again! You'll need it in the next step to set it in your `.env`)

## Set up environment variables

You will need to add and set the following in your `.env` replacing the values as needed:

```dotenv
TF_VAR_okta_api_token=<TOKEN>
TF_VAR_okta_org_name=<ORG_NAME>
TF_VAR_okta_base_url=okta.com
TF_VAR_okta_user_id=<USER_ID>
TF_VAR_okta_username=<EMAIL_YOU_CREATED_DEV_ACCOUNT_WITH>
```

- okta_api_token: token from the generate an api token step, e.g. 00asdfg0YOGsyiGSIlS997gYIgOgjkhgkgLlgytUYI
- okta_org_name: prefix to your okta dev account seen in the url, e.g. dev-123456
- okta_base_url: base url is the same for all devs, just copy this one
- okta_user_id: from the okta dev account admin page, click Directory -> People, select your name and check grab from the url, e.g. 04u1jfzusmgDvUc1H5d7
- okta_username: if you followed the above steps, this should just be your circulo health email e.g. you@circulohealth.com

### Assigning to Cognito User Group

After signing in using Okta, your user will be created in Cognito, but will not have a group. In order to be assigned a group,

1. Open AWS Console
2. Go to the "Cognito" dashboard
3. Click "Manage User Pools"
4. Click "internals"
5. Click "Users and groups" on the left navbar
6. Find and click on your user
7. Click "Add to group" in top chip bar
8. Select the group you want to be added to. YOU MUST BE IN ONLY 1 GROUP OTHER THAN THE OKTA GROUP
9. If previously logged in, sign out and sign back in
