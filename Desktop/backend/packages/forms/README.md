# Local Development

### Setup commitizen globally

Start by installing globally, make sure the global packages are in your PATH,

```
yarn global add commitizen cz-conventional-changelog
```

Create a .czrc file in your home directory, with path referring to the preferred, globally installed, commitizen prompt

```
echo '{ "path": "cz-conventional-changelong" }' > ~/.czrc
```

Now instead of running git commit, run git cz.

### Setup pre-commit

[Website](https://pre-commit.com/),

Install the package, it might take a couple of minutes.

```
curl https://pre-commit.com/install-local.py | python -
```

This installs to ~/bin, so you'll have to add that to your path or move it.

Verify that it's in your path with,

```
pre-commit --version
```

Install the hooks for this repository (make sure your at the root of the repo),

```
pre-commit install
```

### Lint

1. Install [golangci-lint](https://github.com/golangci/golangci-lint)
2. Linting will run before you commit a file, otherwise simply run,

```
golangci-lint run
```

---

# Helpful Guidelines

### How to create a new lambda

1. Modify the following Terraform files under the package where your lambda is going to live (users, forms, etc)

   `lambda.tf`: Here is where we create our lambda and we assign it a role

   ```
   resource "aws_lambda_function" "[lambda-name]" {
     function_name = "sonar_service_[lambda-name]" // this is the name that is going to show in the AWS console
     role          = aws_iam_role.[role-name].arn

     ...

   }
   ```

   `iam.tf`: This is where we set our roles and policies, so our lambda has permissions

   ```
   resource "aws_iam_role" "[lambda-name]" {
     name               = "sonar_service__[lambda-name]"
     assume_role_policy = data.aws_iam_policy_document.lambda_execution_document.json
   }

   resource "aws_iam_role_policy" "sonar_service_[lambda-name]_policy" {
     name   = "sonar_service_[lambda-name]_policy"
     policy = data.aws_iam_policy_document.logging_role_policy_document.json
     role   = aws_iam_role.sonar_service_[lambda-name].id
   }
   ```

> Go back to iam.tf to specify permissions as you continue working on your lambda

    `Outputs.tf`


      ```
      output "[lambda_name]" {
        value = aws_lambda_function.sonar_service_[lambda-name]
      }
      ```

2. Create lambda directory under the cmd folder in the corresponding package (i.e. forms > cmd > [name])

3. Create a `request` file if your lambda takes in any parameters

   - location: package > requests > [request-file-name].go
   - Defines the parameters that the lambda will accept

4. Create a `response` file if your lambda returns any specific parameters

   - location: package > responses > [response-file-name].go
   - Defines the parameters that the lambda will return

5. Go to the BUILD file at the root of the directory you are working on (i.e. forms, users) and add your route in the files array

```
container_image(
name = "users_lambda_image",
base = "@lambda_go//image",
cmd = ["/lambda/connect"],
directory = "/lambda",
files = [
    ...,
    "//packages/users/cmd/[lambda-name]", <---- add your route here
],
)
```

6.  Run `task build` -> this will build a `BUILD.bazel` file in your lambda folder

    Add the following lines:

        ```
        go_binary(
            name = "receive",
            embed = [":receive_lib"],
            goarch = "amd64",     <------ add this line
            goos = "linux",       <------ add this line
            visibility = ["//visibility:public"],
        )
        ```

7.  Go to the root `tf` directory and do the following:

    - Add your route to the `locals.tf ` file
    - Add your module in the corresponding file (i.e. If you created a lambda in the `forms` package, add module in `forms.tf` file)

      ```
      module "[lambda-name]" {
          source        = "./modules/gateway-route"
          route_key     = local.forms_internal_http_routes.[name]
          requires_auth = true

          // resource variables
          api_id     = aws_apigatewayv2_api.gateway.id
          source_arn = "${aws_apigatewayv2_api.gateway.execution_arn}/ */* "

          // module global variables
          authorizer_id = module.global.internal_http_authorizer_id
          lambda_function = module.forms.[lambda-name]
      }
      ```

    - Add your module to the triggers object in the `http.tf` file

      ```
      triggers = {
        redeployment = sha1(join(",", [
        ...
        jsonencode(module.[lambda-name]),
        ]))
      }
      ```

### How to test GORM with go-sqlmock

In the following steps, we are going to replace the GORM dependency by using [sqlmock](https://github.com/DATA-DOG/go-sqlmock), a controlled replacement object that simulates the behavior of the real GORM code.

1. Create a mock database using sqlmock:

```
db, mock, _ := sqlmock.New()
```

2. Open a connection using GORM and mock database:

```
gdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
```

3. Create a table by adding column declarations:

```
sqlmock.NewRows([]string{"id", "column_name_1", "column_name_2"})
```

4. Add the rows you need to test:

```
sqlmock.AddRow("1", "rowValue1", "rowValue2")
```

5. Use `ExpectQuery` to mock a Query

```
mock.ExpectQuery(
	"SELECT(.*)").
	WithArgs(formID).
	WillReturnRows(row)
```

or use `ExpectExec` to mock an INSERT or UPDATE

```
mock.ExpectBegin()
mock.ExpectExec(sqlUpdate).WithArgs(&date, formID).WillReturnResult(sqlmock.NewResult(0, 1)) // 0 new rows, 1 row updated
mock.ExpectCommit()
```

#### Resources:

- [Comment Example](https://github.com/go-gorm/gorm/issues/1525#issuecomment-376164189)
- [Unit test for Gorm application with go-sqlmock](https://tienbm90.medium.com/unit-test-for-gorm-application-with-go-sqlmock-ecb5c369e570)
- [Example Tests](https://github.com/dche423/dbtest)
- [Go-sqlmock Examples](https://github.com/DATA-DOG/go-sqlmock/tree/master/examples)
