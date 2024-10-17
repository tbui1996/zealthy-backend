# sonar-backend [![coverage report](https://gitlab.com/circulohealth/sonar/backend/badges/develop/coverage.svg)](https://gitlab.com/circulohealth/sonar/backend/-/commits/develop)

Sonar Backend/Platform

## Development Environment

### Prerequisites

**_Requirements_**

- [docker](https://docs.docker.com/get-docker/)
  - [Must have daemon running and add user to docker group (on linux)](https://docs.docker.com/engine/install/linux-postinstall/)
  - Windows: use [wsl2](https://www.omgubuntu.co.uk/how-to-install-wsl2-on-windows-10), wsl1 does not support iptables properly for docker
  - [Tips here](https://bce.berkeley.edu/enabling-virtualization-in-your-pc-bios.html) for enabling virtualization in BIOS
- [bazel](https://bazel.build/)
  - On windows (wsl2), you may have to install python manually. Suggestion: `pyenv install 3.8.10` with `python-is-python3` to resolve path
- [go](https://golang.org/)
- [task](https://taskfile.dev/#/)
- [terraform CLI](https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/aws-get-started)

**_AWS configuration_**

- [amazon-ecr-credential-helper](https://github.com/awslabs/amazon-ecr-credential-helper)
  - You've done the following: https://github.com/awslabs/amazon-ecr-credential-helper#docker (the credHelpers route not the credStore route)
- On your personal AWS account [create an IAM user](https://console.aws.amazon.com/iamv2/home?#/home) and give them the `AdministratorAccess` permissions, save the access credentials for this user to use in the next step
  - note: currently Tony is handing out personal aws accounts
- [aws-cli-v2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
  - do [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) using the credentials from previous step

**_IDE_**

- (Optionally) [VSCode](https://code.visualstudio.com/)
  - if you like VSCode, see the [VSCode](#VSCode) section

### Steps to Setup Environment

1. Create a .env file in the root directory with `AWS_ACCOUNT_ID`. This is referenced by the Taskfile.
2. Follow [Setting up Gitlab Remote State](#terraform-configuration-to-gitlab-remote-state)
   - Alternative: [Setting up s3 Remote State](#terraform-configuration-s3-state) If you would like to use an s3 bucket as your backend instead of GitLab.
3. Follow [Setting up Okta Dev Account For Sonar Developers](packages/users/README.md#setting-up-okta-dev-account-for-sonar-developers)
4. Read and follow [Using the Taskfile like a pro](#using-the-taskfile-like-a-pro)
5. Read and follow [Deploy Everything](#deploy-everything)
6. Read and follow [Database Migrations](#database-migrations)
7. Read and follow the [README in packages/users](packages/users/README.md#assigning-to-cognito-user-group)

### Endpoints available after deployment

- api.[name].circulo.dev
- ws-sonar.[name].circulo.dev
- cms.[name].circulo.dev

### Considerations

- There is a root `tf` folder with global infrastructure
- There is also a `tf` folder in each package/service specific to those packages
  and imported as modules through the root `tf` folder
- Each environment has its own terraform workspace (your developer workspace is unique
  to you)

# Using the Taskfile like a pro

The Taskfile is a useful helper to run terraform commands, aws commands, or really anything else you would like in your terminal.
It allows you to not have to think about the order in which you have to run things. Here are some of the most common commands that need to be run.

### Updating the image-variables

The image variables are needed to deploy infrastructure properly. Most task commands take this into consideration
as a dependency. Can be used to update `./tf/.account.auto.tfvars` and `./tf/.image.auto.tfvars` manually, or to see
the current commit hash, and your aws account id .

```shell
$ task image-variables
```

### Initializing infrastructure

This should only need to be done once, unless you add a new module or a new ECR repository under the root `main.tf`
folder. The purpose of this command is to initialize terraform workspace and infrastructure in preparation for planning and
deploying.

```shell
$ task init
```

### Update Provider Versions

Update the version number under the following files:

- tf > main.tf
- packages > companulo > tf > main.tf
- tf > modules > ecr-repository > main.tf

Example:

```json
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.63.0"
    }
  }
}
```

Then, go to your terminal and run the following:

```shell
$ task init-tf-upgrade
```

Finally, commit the changes to the `.terraform.lock.hcl` file.

### Logging into Docker

May need to do this in order to push resources ECR / ECS repositories if you do not have the
[amazon-ecr-credential-helper](https://github.com/awslabs/amazon-ecr-credential-helper) setup.
Run this before any deployment steps. If you do not run this, and do not have Should only need to be run once per day.

```shell
$ task docker-login
```

### Planning terraform infrastructure changes

Before deploying any AWS changes through terraform, it is often a good idea to view the plan for the changes to make sure
there are no errors.

```shell
$ task plan-all
```

### Deploying Terraform infrastructure

Once the plan is successful, you are most likely to ready to deploy your infra to AWS. There are a few different ways to do this.
You can deploy everything at once, or deploy on an individual project basis.

#### Deploy Everything

All aws infra can be deployed through one command. This command will show you the plan and then ask you to confirm it.
If it errors that is ok. Sometimes infrastructure is dependent on each other and takes more time to complete than terraform
allows for once given a "done" signal. Just run the command again and see if it completes successfully. Try running 3-4 times (if necessary)
before asking for help.

**You may want to comment out the elastic module in `packages.tf` before deploying.
It can take up to 30 min to complete and unless you are modifying that code you most likely will not need
the infrastructure**

```shell
$ task deploy
```

#### Deploy Everything Auto Approve

This can avoid having to confirm the plan, do not use this unless you know the infrastructure changes well.

```shell
$ task deploy-auto-approve
```

#### Deploy only Terraform resources

This can be used if changes are not made to image tag, but everything else still needs to be deployed.

```shell
$ task deploy-tf
```

#### Deploy Router only

This can be used to deploy router infra only, use this if changes have happened only in router.

```shell
$ task router-deploy
```

#### Deploy Forms only

This can be used to deploy forms infra only, use this if changes have happened only in forms.

```shell
$ task forms-deploy
```

#### Deploy CMS only

This can be used to deploy CMS infra only, use this if changes have happened only in CMS.

```shell
$ task cms-deploy
```

#### Deploy Search only

This can be used to deploy search infra only, use this if changes have happened only in search.

```shell
$ task search-deploy
```

#### Deploy Pearls only

This can be used to deploy pearls infra only, use this if changes have happened only in pearls.

```shell
$ task pearls-deploy
```

#### Deploy Support only

This can be used to deploy support infra only, use this if changes have happened only in support.

```shell
$ task support-deploy
```

#### Deploy Companion only

- This can be used to push companulo resources to S3 bucket.

```shell
$ task companion-push
```

- This can be used to deploy tf infra for companulo

```shell
$ task companion-tf-infra
```

# Terraform Configuration to Gitlab Remote State

If you have not already, upgrade your terraform CLI version to be 1.0. The upgrade guides can be found
[here](https://www.terraform.io/upgrade-guides/1-0.html). If you are on version 0.15.5 you do not necessarily _need_ to
upgrade, it is just recommended for stability purposes. The commands shown below will work with the latest version (0.15.5, 1.0\*)
of the Terraform CLI.

- Log into GitLab
- Go to profile dropdown and click `Edit profile`
- Go to `Access Tokens` section
- Create a new access token.
  - Give token a name: i.e `dev-name-terraform-remote-state-token`
  - Do **NOT** give token an expiration date
  - Select `api` scope
  - Click `Create Personal Access Token`
  - Copy Access token for future use
- Create an AWS Route53 Hosted Zone
  - A Hosted Zone is container that holds information about how you want to route traffic for a domain
  - Go to Route53 in your AWS Console and select the option to create one
  - Enter `[your-name].circulo.dev` under Domain Name
  - Select the `Public hosted zone` type
  - Select `Create hosted zone`
  - Expand the `Hosted zone details` section and copy the `Name servers` available
    (i.e.
    ns-1662.awsdns-15.co.uk
    ns-613.awsdns-12.net
    ns-1391.awsdns-45.org
    ns-213.awsdns-26.com)
  - Provide your `Name servers` to Tony for now (09/09/21). Eventually everyone will have access to the dev account and will be able to add them on their own.
- Create a `.env` file at the root of the sonar-backend project
  - You may already have one, if that is the case add to that file
  - Follow [.env file](#.env-file)
- Request `Maintainer` role permissions on this repo if you don't have that already. ([Documentation](https://docs.gitlab.com/ee/user/infrastructure/iac/terraform_state.html#permissions-for-using-terraform))
- Run the terraform [command](#command) (cd into `tf` folder before running command)
- Run `task plan-all` to see that the migration was successfully.
  - Can also deploy if you feel like you are in state to do so.

> If there seem to be issues after the migration pushing new architecture, can always try `terraform refresh` and see if that helps.

### .env file

```dotenv
AWS_ACCOUNT_ID=<MY_AWS_ACCOUNT_ID>
TF_VAR_environment=dev-<DEV_NAME>
TF_VAR_db_password=<ANY_DB_PASSWORD>
TF_VAR_hosted_zone_id=<HOSTED_ZONE_ID>
TF_VAR_domain_name=<DEV_NAME>.circulo.dev
```

Variables that TF needs to read should be prefixed with `TF_VAR`, otherwise AWS won't be able to read your Account ID.

- AWS_ACCOUNT_ID: The ID found in the dropdown of `workflow-dev-<DEV_NAME>`
- environment: The environment as represented by dev-your-name-here
- db_password: Used to create RDS cluster, can be any value. If you have already created one in the old remote state make sure to use the same one.
- hosted_zone_id: Hosted zone ID as created in Route53
- domain_name: The domain name of your hosted zone
- DEV_NAME: Replace with your name (assuming you have a unique name to the Workflow Circle)

Once you have created the `.env` file, you can migrate your terraform state.

### Command

```shell
terraform init -backend-config="address=https://gitlab.com/api/v4/projects/28022043/terraform/state/dev-<DEV_NAME>" \
-backend-config="lock_address=https://gitlab.com/api/v4/projects/28022043/terraform/state/dev-<DEV_NAME>/lock" \
-backend-config="unlock_address=https://gitlab.com/api/v4/projects/28022043/terraform/state/dev-<DEV_NAME>/lock" \
-backend-config="username=<GITLAB_USERNAME>" \
-backend-config="password=<GITLAB_ACCESS_TOKEN>" \
-backend-config="lock_method=POST" \
-backend-config="unlock_method=DELETE" \
-backend-config="retry_wait_min=5"
```

- DEV_NAME: your (the person reading this) name
- GITLAB_USERNAME: The username of your gitlab account found in the account dropdown
- GITLAB_ACCESS_TOKEN: The access token you created above

> Note: If you are migrating from GitHub, add `-reconfigure` to the terraform command. i.e. `terraform init -reconfigure -backend-config=...`.
> If the CLI output asks you a few yes / no questions. Answer **yes** to all questions.

# Terraform Configuration s3 State

- Create an s3 bucket and a DynamoDB table to store your state. You can use [sonar-milu-tfstate-s3](https://gitlab.com/circulohealth/sonar/terraform-backend/sonar-milu-tfstate-s3) as reference to create them.
- Create a file called `override.tf` in the tf directory and add the following:

  ```json
  terraform {
    backend "s3" {
      bucket         = "<S3-BUCKET-NAME>"
      key            = "states/terraform.tfstate"
      dynamodb_table = "<DYNAMO-DB-NAME>"
      region         = "us-east-2"
    }
  }
  ```

  This file is git ignored and processed after loading configuration by default. The rest of the team is currently using GitLab to store their state so we need to override the backend type "http" to "s3" without affecting the others

- Run `terraform init -reconfigure`

# Database Migrations

This repo does migrations with [golang-migrate](https://github.com/golang-migrate/migrate).
Allows for migrations to be written in plain sql, and will fail hard if things are incorrect.
Eventually migrations will be ran through CI/CD.

## How to add migrations

- Migrations will go into the `~/db/migrations/` folder.
- In order to create a migrations run the following command: `migrate create -ext sql name_of_migration` which will generate a timestamp and up/down .sql files.
- The naming for migrations is as follows: `YYYYMMDDHHMM_action_schema_table_name.up.sql`. Note that the `HHMM` must be on a 24hr clock to preserve order of migrations.
- Migrations must have both an `.up.sql` && a `.down.sql` file.
  What are `.up` & `.down` files you ask? They provide a simple way to wipe the state of the DB and recreate it easily.
- Migrations are written in plain `sql` and ran against a remote DB environment.
- Migrations must be reviewed before they can be run

## How to add other team's migrations

Sometimes you will need to add a new feature to Sonar that connects to another team's database (eg: Sonar -> Doppler)
In live environments, these databases will usually be deployed and you can connect to them via basic auth and using .env variables.
However, when testing locally, you will not have access to such databases since they are behind a VPC in a different AWS account.

In order to test your changes locally:

- Grab the schema from that team's project, eg: [doppler schemas for retool](https://gitlab.com/circulohealth/applications/doppler/data-retool-idd/-/tree/main/pg-schemas)
- Then create the necessary up and down migration files in db/external_migrations/

You can now task migrate-up and task migrate-down normally, as there is a conditional check in the taskfile command. Files in db/external_migrations/ will only be ran in your local dev environment, and will NOT be deployed in live environments. These migrations also create tables in a database named `external` in your "local" RDS cluster, so make sure to point to `external` as the database name wherever necessary.

## Running migratrions through golang-migrate CLI

### Download the CLI

#### All platforms

```bash
$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz
```

#### MacOS

```bash
$ brew install golang-migrate
```

#### Linux (\*.deb package)

```bash
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
```

### Run Migrations

Running migrations using task is easy. It will automatically pull the DB endpoint necessary to run migrations against.

Running migrations assumes the following things:

- You have configured AWS CLI and Terraform correctly
- You have downloaded task && golang-migrate
- You have deployed all necessary architecture using terraform
- You have [jq](https://stedolan.github.io/jq/) installed

#### Populate the DB:

```bash
$ task migrate-up
```

#### Teardown the DB:

```bash
$ task migrate-down

```

#### Working with a dirty database

If you ran an unsuccessful migration, then a dirty bit will be set in the `schema_migrations` table.

##### Fixing local migrations

In order to fix it locally, take the following actions:

1. Identify and fix any issues with your SQL.
2. Identify the last successful version and set that in your `.env`, e.g. `DB_VERSION=202109081622`
3. Run the following command which will update the `schema_migrations` table to reflect the last known good migration as well as reset the dirty bit to `false`.

```bash
task migrate-force
```

4. Run your corrected SQL again.

```bash
task migrate-up
```

##### Fixing deployed migrations

In order to fix migrations in a live environment, take the following actions:

1. Identify and fix any issues with the SQL that failed to run during migrations.
2. Identify the last successful version (typically the version previous to the failed version)
3. On Gitlab in the repository, click `CI/CD` -> `Pipelines` -> `Run pipeline` or click [here](https://gitlab.com/circulohealth/sonar/backend/-/pipelines/new) to create a new pipeline
4. Select the environment where migrations are failing and set up the `DB_VERSION` variable to the value found in step 2.
5. **DOUBLE CHECK:** Is the correct environment set? Is the correct `DB_VERSION` being passed?
6. Click the `Run pipeline` button and the new migrations should be run.

##### Fixing deployed migrations

In order to fix migrations on a deployment, take the following actions:

1. Identify and fix any issues with your SQL.
2. Identify the last successful version.
3. Create a new pipeline by clicking `CI/CD` -> `Pipelines` -> `Run pipeline` or clicking [here](https://gitlab.com/circulohealth/sonar/backend/-/pipelines/new).
4. Select the appropriate environment and add a variable `DB_VERSION` with the value equal to the last good migration.
5. Click the `Run pipeline` button.

# GORM

The backend uses an Object-Relational-Mapper (ORM) to interact with the database. Check the [GORM DOCS](https://gorm.io/docs/) for more information.

# VSCode

Assuming you want to use VSCode, for the best experience:

- Install the Go extension
  - After installed, open the command palette `ctrl + shift + p` or `cmd + shift + p`
  - Type in `Go: Install/Update tools`
  - Select all and click `OK`
- Use a workspace (VSCode does not play well with monorepos, the solution is to use a [multi-root workspace](https://code.visualstudio.com/docs/editor/workspaces#_multiroot-workspaces))
  - In a fresh window, just add each package (e.g. `common`, `forms`, etc) individually using the option `File -> Add Folder to Workspace...`
  - Save the workspace (`File -> Save Workspace As...`)
  - Use the workspace the next time you want to work on `sonar-backend`

# Testing

## Unit Testing

Unit testing defines the specification for what is guaranteed by a unit (i.e. function). Each test has an input state and validates the output state after executing the unit/function. In this sense, the unit tests for a unit define what the original developers guarantee when you use it. In a sense, it's a form of self documentation. Unit testing is required for every merge request that goes into the backend, and where possible and timely, attempt to add unit tests to legacy code that was implemented before unit tests. This is a system of compounding benefits, as we write more unit tests, we can develop faster, ensure increased stability, and identify regressions or breaking changes.

The current testing stack uses `go test` and the [testify](https://github.com/stretchr/testify) framework. In order to invoke tests, you can call `task <package>-test` or `task test-all`. Most IDE's will also provide testing capabalities. For example, on the sidebar of VSCode, you can open testing and it will identify our tests. From there you can run then after you make changes.

### Reading

- [testing in go dependency injection](https://ieftimov.com/post/testing-in-go-dependency-injection/)

### Basics

1. Test files are defined within a package, next to the files that it tests. The name of the file is the same as the file or unit that it's testing but ends with `_test.go`.
2. Tests are grouped into test suites. A Test Suite could be for a function, and each test could be for different requirements of that function.

### Best Practices

When creating a lambda function, start with 3 files:

1. `connector.go`: This file is where the `main` function should live. Any resources needed (e.g. Environment Variables, DB Connections, Logger, etc.) should be set up here and passed into the `Handler` function defined in `handler.go`. Additionally, because no buisness logic should be placed here, add the following to the top of the file to ignore it for code coverage purposes:

```golang
//go:build !test
// +build !test

package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func connector(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
  // connector should set up a HandleInput struct that houses all dependencies to carry out the request
  Handler(HandleInput{})
}

func main() {
	lambda.Start(connector)
}
```

2. `handler.go`: This is where all business logic should occur.
3. `handler_test.go`: This is where all business logic inside of `handler.go` should live.

#### Test names as a contract

Following the idea that tests define the spec for a unit, test names should correspond to the low level requirements of a unit so that other developers can browser not only the function but also the tests to figure out what the unit/function will do based on certain inputs. There should be unit tests for each guarantee that a function makes, including (or especially) edge cases. Tests are grouped into suites, and each test in the suite is for different functionality.

An example test name could be,

`func (s *suite.Suite) TestFunctionName_ShouldDoSomething() {...}`

#### Dependency Injection and Decoupling

Dependency Injection just means that if your unit/function uses something, it should be passed in the arguments. For example, say our function wants to put an item into a dynamodb table. Within the function, we could create a dynamodb session and call put item within the function, but then when we get to testing, our unit test isn't testing just the unit, it's also testing the dynamodb function (and we'd have to have a live dynamodb table!). Instead, if we pass in a function to the unit, then when we get to testing, we can `mock` that function. This means creating a fake function that is passed to the unit, which we can guarantee acts a certain way based on inputs and outputs. That way, we're only testing the unit and not additionally the functions that the unit is calling. A similar concept is decoupling.

#### Creating a mock

Do not create mocks manually. [Mockery](https://github.com/vektra/mockery) was built to generate mocks, so we should use it unless there is some weird outlier.

##### Third party mocks

When creating a mock for an interface and implementation for which we do not own (e.g. AWS), assuming you already have a local copy to work from, add the respective command to the [Taskfile.yml](Taskfile.yml) in the `auto-gen-third-party-mocks` task. Reference the existing commands for examples. You should only need to copy and existing line and update the directory and name:

```sh
mockery --case=underscore --dir $GOPATH/<PATH_TO_FILE_CONTAINING_THE_INTERFACE> --name <NAME_OF_THE_INTERFACE>
```

##### Sonar mocks

When creating a mock for a Sonar owned interface, if a block for the package does not already exist, create one. For example:

```sh
  authorizer-auto-gen-mocks:
    dir: packages/common/authorizer
    cmds:
      - mockery --all --case=underscore
```

If you created a new block, be sure to add it to the `auto-gen-mocks` task deps.

Once you're block is in place, simply run `task auto-gen-mocks` to recreate any new or updated mocks.

# Terraform Practices

Terraform has a simple abstraction system, called "modules". These accept inputs and outputs, and the side effects of the module are resources. In general, it's a best practice to keep our terraform code "flat". For us, this means not nesting modules within modules, but instead wherever a module calls another module, consider bringing that to the root terraform module and passing it's outputs as inputs to the module that would have called it.

### Packages as Modules

We have a variety of "packages" that are implemented in terraform as "modules" (i.e. "users", "support", "router", etc). These modules should define the infrastructure that they need to execute business logic in their domain, but not be concerned with things such as how that functionality is accessed externally or internally. For example, this would indicate that API Gateways and routes should be defined at the root terraform module and not in the package module.

### Shared resources

Shared resources, which are resources that are shared between 2 or more "packages" (a BAZEL term, i.e. "users", "support", "router", etc), should be defined at the root terraform module and passed down as input to packages that need it. Similarly, if you ever have a package module that references another package modules outputs, consider bringing the resources to the root terraform module and passing it down to both.

### Helpful Guidelines

- [How to create a new lambda](packages/forms/README.md#how-to-create-a-new-lambda)
- [How to test GORM with go-sqlmock](packages/forms/README.md#how-to-test-gorm-with-go-sqlmock)

# Backwards Compatibility

There now exists a way to introduce breaking changes to the loop (or anything else that consumes our API's)

The way is moving a route to a version 2 api gateway. There is an example of how a v2...vN api gateway can be created in `loop_http_v2.tf`

Essentially you would use `"./modules/api-gateway"` and provide the correct args in order to correctly create a new api gateway.

## Creating a new API Gateway

i.e

```hcl
module "my_api_gateway_v2" {
  source         = "./modules/api-gateway"
  aws_account_id = data.aws_caller_identity.current.id
  aws_region     = data.aws_region.current.name
  deployment_triggers = [ ... A LIST OF ROUTES FROM GATEWAY MODULE THAT REDEPLOY RESOURCE ... ]
  domain_name         = needs to be domain id of a route53 configuration module
  environment         = var.environment
  gateway_name        = new gateway name
  acm_certificate_arn = needs to be certificate arn of route53 configuration module
  api_mapping_key     = "v2", "v3", "v4" ... "vN"
}
```

each new api gateway also will need an authorizer. This can also be created using a module and will vary based on your needs.

```hcl
module "my_http_authorizer" {
  source                    = "./modules/http-authorizer"
  api_gateway_execution_arn = module.my_api_gateway_v2.api_gateway_execution_arn
  apigateway_id             = module.my_api_gateway_v2.api_gateway_id
  authorizer_name           = NAME YOU WANT TO GIVE AUTHORIZER
  credentials_role_arn      = module.global.<CREDENTIALS_ROLE_ARN>
  identity_sources          = ["$request.header.Authorization"]
  lambda_function_name      = module.global.<LAMBDA_FUNCTION_NAME>
  lambda_invoke_arn         = module.global.<LAMBDA_INVOKE_ARN>
  statement_id              = i.e -> "AllowExecutionFromAPIGatewayForMyHTTPAuthorizerV2"
}
```

## Adding to existing v2 api gateway

To add a new route + integration to existing v2 api gateway, you need to use the gateway route module, as well as outputs from both the gateway module and authorizer module

```hcl
module "my_v2_route" {
  source        = "./modules/gateway-route"
  route_key     = local.v2_http_route.route
  requires_auth = true

  // module router variables
  api_id     = module.my_api_gateway_v2.api_gateway_id
  source_arn = "${module.my_api_gateway_v2.api_gateway_execution_arn}/*/*"

  // module global variables
  authorizer_id = module.my_http_authorizer.<AUTHORIZER_ID>

  lambda_function = module.<MY_MODULE>.<LAMBDA_FUNCTION>
}
```

## Finding last stable lambda version and testing

What you need to do here is go to the ECR console in the environment of your choosing, and cross-reference the commit hashes with what is in the environment you are in.
You then need to update

### Locally

Finding the last local commit hash is the hardest one. It can be done on the command line if you prepend your commits with the story number.

```shell
$ git reflog | grep -i MY_BRANCH_NAME
```

Then you can cross-reference any `HEAD` states where you moved from develop to your branch.

THen set this variable in your .env file

```dotenv
TF_VAR_last_local_commit=<MY_LAST_COMMIT>
```

And then, you need to un-comment the line that says `uncomment the line below to test backwards compatibility locally`
And comment out the line above `uncomment the line below (and comment out line above) only when steps to introduce backwards compatibility have been taken`

You can then test your backwards compatibility locally

### Development

The last stable commit is most likely the one that was committed before your branch in the history.

Once you think you know it, go to CI pipeline and update `last_stable_develop_commit`

You can now deploy and test against develop

### Test

The last stable commit is most likely the one that was committed before last develop squash in the history.

Once you think you know it, go to CI pipeline and update `last_stable_test_commit`

You can now deploy and test against test

### Prod

The last stable commit is most likely the one that was committed before your branch in the history.

Once you think you know it, go to CI pipeline and update `last_stable_production_commit`

You can now deploy to prod

## Cleanup

There will come a time when the old image versions will delete. we should only need them for one release otherwise something is very wrong.
So once release is over and your changes have been merged, go into develop and set `last_production_commit = var.image_version` and redeploy!
