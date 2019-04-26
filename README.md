[![Go Report Card](https://goreportcard.com/badge/emergenseek/backend)](https://goreportcard.com/report/emergenseek/backend)

## EmergenSeek Backend
Trello Board: https://trello.com/b/O864ONXb/golang-backend

## Setup
Note: AWS CLI tools may require you to have [Python](https://www.python.org/downloads/release/python-371/) installed
1. Download, install, and setup Go 1.11+: http://howistart.org/posts/go/1/
2. Install dep: https://github.com/golang/dep (Requires GOBIN environment variable)
3. Within the `src` directory of `$GOPATH`, clone this repository using: `git clone https://github.com/emergenseek`
4. Change directories into this repository (`cd backend`) and install dependencies: `dep ensure`
5. Create an AWS account and enable billing
6. Create a new IAM user:
      - Configure the `Access type`
        - Enable `Programmatic access` and take note of the `access key ID` and `secret access key`
        - If you would like to create keys and accounts for multiple users, repeat step 4 and provide each account with a `Console password`
        - Configure user permissions; select `Attach existing polices`
      - For simplicity, give the new account the `AdministratorAccess` policy
        - Policies may be configured later on
7. Download the AWS CLI
    - Reference: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
    - In a terminal, run `aws --version` to verify that the tool installed successfully
8. Configure the AWS CLI:
    - Run `aws configure`
    - Paste the `Access Key ID` that was generated for your user
    - Paste the `Secret Access Key` that was generated for your user
    - Set your default region name (recommended: `us-east-2`)
    - Set the output format (recommended: `table`)

## Testing Locally 
1. Install Docker (hub.docker.com account required)
   - On Windows, if you do not have access to Hyper-V (i.e. Windows 10 Pro), install Docker Toolbox
2. Install aws-sam-local using one of the methods listed at https://aws.amazon.com/serverless/sam/
3. For convenience, install `make`
    - Windows: `choco install make`

Every subdirectory within `/functions` will have:
  - A Makefile with a `sam` rule, this will invoke the test. Docker must be running. To run the rule, change into the function directy install make and invoke `make sam`
  - A main.go file containing the source code for each Lambda function handler (`func Handler`)
  - If the function accepts [non-idempotent](https://developer.mozilla.org/en-US/docs/Glossary/idempotent) HTTP methods, an event.json file used for locally testing using cURL, wget, or HTTPie. If the accepted method is anything other than POST, the JSON's filename will be prefixed with the request method.
  - A template.yaml for the SAM CLI tool to provision resources, for local testing
  - Optionally: a request.go file to encapsulate a Request struct, typcially used for validating the JSON request body

Example invocation of `ESSendSMSNotification`
```bash
curl -X POST http://127.0.0.1:3000/sms -d @event.json --header "Content-Type: application/json"
Successfully sent SMS to contacts of user: e78e0c86-f9ba-4375-9554-6dc1426f5605
```

## Lambda Function Breakdown
Below is a table detailing the function and purpose of each of the Lambda functions contained within this repository. Function packages are prefixed with ES (EmergenSeek) followed by their operation/responsibility.

|Function Name             |Purpose                                                                                              |
|--------------------------|-----------------------------------------------------------------------------------------------------|
|`ESCreateUser`            |Create a new EmergenSeek user                                                                        |
|`ESSendSMSNotification`   |Send SMS notifcations on S.O.S. trigger via Twilio                                                   |
|`ESSendEmergencyVoiceCall`|Send voice calls on S.O.S. trigger via Twilio                                                        |
|`ESManageProfile`         |Retrieve *and* update settings for an EmergenSeek user                                               |
|`ESPollLocation`          |Send periodic location information to a user's contacts via Twilio                                   |
|`ESManageSettings`        |Retrieve *and* update settings for an EmergenSeek user                                               |
|`ESServiceLocator`        |Retrieve hospitals and pharmacies near a user's location via the Google Places API                   |
|`ESUpdateTier`            |Update the alert tier of a user's contacts                                                           |
|`ESGetEmergencyNumbers`   |Retrieve emergency phone numbers (i.e. 911 number for a given country, etc), given a user's location |
|`ESAddContact`            |Associate additional contacts to a user                                                              |
|`ESGetLockScreenInfo`     |Retrieve information necessary for first responders                                                  |

## Project Structure
  - `/common`
    - Common helper functions and packages which are shared by the Lambda functions in `/functions`
  - `/functions`
    - Contains Lambda functions that are automatically built, tested, and deployed using AWS CodeBuild, CodeDeploy, and CodePipeline
  - `/githooks`
    - [Git hooks](https://git-scm.com/docs/githooks) for ensuring `goreportcard` passes with a decent grade. Make sure to run `git config core.hooksPath githooks` to make sure Git can find these hooks. The default directory is `.git/hooks/`

#### Helpful References
 - https://stackoverflow.com/questions/48619686/project-structure-for-go-program-to-run-on-aws-lambda
 - https://outcrawl.com/go-url-shortener-lambda
 - https://artem.krylysov.com/blog/2018/01/18/porting-go-web-applications-to-aws-lambda/
 - https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda
 - https://blog.rowanudell.com/aws-sam-local-error-with-golang/
 - https://github.com/awslabs/aws-sam-cli/issues/437
 - https://medium.com/a-man-with-no-server/running-go-aws-lambdas-locally-with-sls-framework-and-sam-af3d648d49cb

#### Acknowledgements 
 - Data source for ESGetEmergencyNumbers: https://github.com/Alex0x47/EmergencyNumbers
