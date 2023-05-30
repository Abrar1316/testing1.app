
# Billing App

![Dashboard](https://github.com/surajtft/billing-app/blob/Feature/updatedREADME/services/webbff/public/img/Readme.png)

## Overview

The ability to budget, track, and optimize spending is becoming increasingly important for organizations that are using to the cloud services. This is especially true for Data Platform teams that support multiple Analysts and Data Scientists. In order to effectively manage costs across various cloud platforms, we developed a systematic approach to cost management to address our concerns around accountability and budgeting.

## Features

- Display cost information and analytics from service providers like AWS.

- Allow users to view cost information for different services.

- Provide visualizations and graphs to help users understand cost trends.

- Allow users to set cost budgets and get notifications when they are exceeded.

- Provide recommendations to optimize costs based on usage patterns.


# To start using Billing App

## Requirements

* Golang 
* AWScli
* PostgreSQL

## Pre-requisites

 - Clone GitHub repo `git clone https://github.com/surajtft/billing-app`

 - An AWS user with Administrator/Power user access.
Refer the below AWS documentation to create a user and generate Access Keys.
https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html

- Install AWS sdk for Go
```bash
$ go get github.com/aws/aws-sdk-go-v2/aws
```

## To Run Locally

### Setup postgres database
* Open psql postgres:- 
  ```bash
  - Create new postgres user:- 
    * CREATE ROLE dev WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION BYPASSRLS PASSWORD '12345';
  - Create new local database
    * CREATE DATABASE "billapp-dev" OWNER dev;
  - Exit psql
    * \q
  - To Login at billapp-dev database
    * psql billapp-dev
  ```
* To Initilize migration, run below command from project's root folder
  ```bash
  cd tools/migrations && go run *.go -command=migrate && cd ../..  OR  cd tools/migrations && go run . -command=migrate && cd ../..
  ```
* To Rollback migration, run below command from project's root folder
  ```bash
  cd tools/migrations && go run *.go -command=rollback && cd ../..  OR  cd tools/migrations && go run . -command=rollback && cd ../..
  ```
* To dump test data with Migrations
    ```bash
    go run tools/test_data/migrate_test_data.go
    ```
    - Login: social@tftus.com
    - Password: Tftus@1234

* To start server
  ```bash
  go run services/server/server.go
  ```
* To start Cron job 
  ```bash 
  go run tools/cronjob_aws/engine.go
  ```
### To test db
* Open psql postgres:- 

  ```bash
  - Create new postgres test user:-
    * CREATE ROLE test WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION BYPASSRLS PASSWORD '12345';
  - Create new local database
    * CREATE DATABASE "billapp-test" OWNER test;
  - Exit psql
    * \q
  - Run migration
    * cd tools/migrations && go run *.go -command=migrate -test=true && cd ../..
  ```  
## License

[MIT](https://choosealicense.com/licenses/mit/)

## Support

If you have any questions or issues with the Billing App, please contact us at info@tftus.com

