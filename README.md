# AWS Inspector Reporter (AIR)

[![circleci][circleci-image]][circleci-url] [![Go Report Card][go-report-card-image]][go-report-card-url] 

## about
AIR is a tool to retrieve the latest AWS Inspector findings (latest run of each template) from your AWS accounts and presents them in an auto-filtered Excel spreadsheet.  
By specifying filters it enables you to adjust severity of specific findings, or ignore them and state the justification.  
Generated reports can be automatically emailed using AWS SES (Simple Email Service).
 

## installation
Download the latest release here: https://github.com/jonhadfield/aws-inspector-reporter/releases

#### macOS and Linux
  
Install:  
``
$ install <downloaded binary> /usr/local/bin/air
``  

## running

Type air and press enter.

## configuration

### authentication
AIR retrieves Inspector findings using the AWS API that requires a set of API credentials. See [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-quick-configuration) for instructions on how to set credentials. 

### permissions
#### basic  
In order for AIR to access the AWS Inspector findings the user or role that it runs under will need the following policy:  
``
arn:aws:iam::aws:policy/AmazonInspectorReadOnlyAccess  
``    

For AIR to be able to use the AWS Account alias (name) instead of just the AWS Account ID number, it additionally requires this permission:  
``
iam:ListAccountAliases
``  

#### sending emails with AWS SES (Simple Email Service)
AIR will additionally require the following permission for the identity resource sending the email:  
``
ses:SendRawEmail
``

### filtering
By default, AIR will report the severity stated by AWS Inspector. To override these, create a directory called config with a file called 'filters.yml' in with a list of filters to apply:  

    - title-match: <finding title to match, supporting regexp>  
      severity: <high|medium|low|informational|ignore>  
      comment: <comment to add to spreadsheet>

See [here](docs/filters.yml.example) for examples.


### email
AIR supports sending generated reports via email using AWS SES. Note: this requires the provided AWS credentials have the necessary permissions.  
To configure: create a directory called config with a file called 'report.yml' with the email settings:

    email:
      provider: ses
      region: <AWS region for SES>
      source: "<email address of sender>"
      recipients:
        - "<email recipient>"
        - "<email recipient>"
      subject: "<email subject>"
 
See [here](docs/report.yml.example) for an example.


### running against multiple-accounts
By default, AIR will retrieve findings from the AWS account that corresponds to the credentials specified.  
To run against multiple accounts you need to:  
* provide sts:AssumeRole permissions to the user AIR is run with
* create an IAM role in each target account with:
  * the policy 'AmazonInspectorReadOnlyAccess' attached
  * a trust relationship allowing the provided AWS permissions to be used to assume the role (see [here](docs/TRUST.md) for examples)

directory called 'config' with a file called 'targets.yml' that specifies a list of target account roles:
* id: the numeric account id
* alias: the account alias
* roleName: name of the role to assume
* roleExternalId _(optional)_: to match the external id specifed on trust relationship on the target role  

See [here](docs/targets.yml.example) for example.

[circleci-image]: https://circleci.com/gh/jonhadfield/aws-inspector-reporter.svg?style=svg
[circleci-url]: https://circleci.com/gh/jonhadfield/aws-inspector-reporter
[go-report-card-url]: https://goreportcard.com/report/github.com/jonhadfield/aws-inspector-reporter
[go-report-card-image]: https://goreportcard.com/badge/github.com/jonhadfield/aws-inspector-reporter