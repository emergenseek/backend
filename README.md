## EmergenSeek Backend
Trello Board: https://trello.com/b/FgW72yJ5/emergenseek

## Setup
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
      - To give the user access to Lambda, DynamoDB, and SNS, give the user the `AWSLambdaDynamoDBExecutionRole` policy
        - Additional policies may be added later on
7. Download the AWS CLI
    - Reference: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
    - In a terminal, run `aws` to verify that the tool installed successfully
8. Configure the AWS CLI:
    - Run `aws configure`
    - Paste the `Access Key ID` that was generated for your user
    - Paste the `Secret Access Key` that was generated for your user
    - Set your default region name (recommended: `us-east-2`)
    - Set the output format (recommended: `table`)
9. Test that everything is working correctly:
   - Coming soon

## Lambda Function Breakdown
Below is a table detailing the function and purpose of each of the Lambda functions contained within this repository. Function packages are prefixed with ES (EmergenSeek) followed by their operation/responsibility.

|Function Name          |Purpose                               |
|-----------------------|--------------------------------------|
|`ESCreateUser`         |Create a new EmergenSeek user         |
|`ECSendSMSNotification`|Send SMS notifcation on S.O.S. trigger|
|...                    |...                                   |

## Project Structure
  - `/common`
    - Common helper functions which are shared by the Lambda functions in `/functions`
  - `/functions`
    - Contains Lambda functions that are automatically built, tested, and deployed using AWS CodeBuild, CodeDeploy, and CodePipeline
  - `/hooks`
    - [Git hooks](https://git-scm.com/docs/githooks) for ensureing `goreportcard` passes with a decent grade. Make sure to run `git config core.hooksPath githooks` to make sure Git can find these hooks. The default directory is `.git/hooks/`
#### Helpful References
 - https://stackoverflow.com/questions/48619686/project-structure-for-go-program-to-run-on-aws-lambda
 - https://outcrawl.com/go-url-shortener-lambda
 - https://artem.krylysov.com/blog/2018/01/18/porting-go-web-applications-to-aws-lambda/
 - https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda

