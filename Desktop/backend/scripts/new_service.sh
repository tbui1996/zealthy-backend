#!/bin/bash

while getopts ":h" option; do
   case $option in
      h) # display Help
        cat << EOF
        positional arguments:
            service_name - the name of the new service
            lambda_name - the name of the first lambda the service will have

        example:
            new_service.sh feature_flags create_flag
EOF
         exit;;
   esac
done

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
pushd "$parent_path"

cp -rf ./new_service_template ../packages/$1

pushd ../packages/$1

mv BUILD.template BUILD
mv main.go_template main.go
mv cmd/LAMBDA_NAME/connector.go_template cmd/LAMBDA_NAME/connector.go
mv cmd/LAMBDA_NAME/handler.go_template cmd/LAMBDA_NAME/handler.go

find . -type f -exec sed -i '' -e "s/SERVICE_NAME/$1/g" {} + 
find . -type f -exec sed -i '' -e "s/\/\/go:build exclude//g" {} + 
find . -type f -exec sed -i '' -e "s/\/\/ +build exclude//g" {} + 

if [ "$2" ]; then
    mv cmd/LAMBDA_NAME cmd/$2
    find . -type f -exec sed -i '' -e "s/LAMBDA_NAME/$2/g" {} + 
fi

go mod init github.com/circulohealth/sonar-backend/packages/$1

echo "replace github.com/circulohealth/sonar-backend/packages/common => ../common/" >> go.mod

go mod tidy

popd
popd

cat << EOF
TODO
===============
1. create ecr repository named "circulo/sonar-backend-$1-lambda"
2. instantiate your new services tf module in the root tf folder's main.tf
3. write tf for api gateway route + update lambda role as needed
3. add appropriate scripts to task file (build, deploy, test, lint, etc)
4. update appropriate scripts in task file (init-ecr, build, deploy, test, lint, etc)
5. add appropriate replace to root go.mod file
6. update .gitlab-ci.yml to match the other packages
EOF
