## Running as a Lambda function
Running as a function implies the results will be sent via email and that requires additional configuration and permissions.

### configuration
Configuration needs to be stored in AWS S3 from where the function will download it when executed. See [README](../README.md) for examples of the report, filters, and targets configuration files.
Place report.yml and the optional filters.yml and targets.yml files in the same directory in an S3 bucket.

### permissions

#### IAM policy 
Create an IAM policy called awsInspectorReporterLambda with the following policy:
  
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "ReadAccountAlias",
                "Effect": "Allow",
                "Action": "iam:ListAccountAliases",
                "Resource": "*"
            },
            {
                "Sid": "EmailResults",
                "Effect": "Allow",
                "Action": "ses:SendRawEmail",
                "Resource": "arn:aws:ses:us-east-1:012345678901:identity/<identity>"
            }
        ]
    }

To give the function permission to download the configuration from S3, either add the following statement to the policy:

    {
        "Sid": "DownloadConfigFromS3",
        "Effect": "Allow",
        "Action": "s3:GetObject",
        "Resource": "arn:aws:s3:::my-bucket/config/*"
    }

Or add the following to the S3 Bucket policy (after creating the role below):

    {
        "Sid": "AIRConfigDownload",
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::012345678901:role/awsInspectorReporterLambda"
        },
        "Action": "s3:PutObject",
        "Resource": "arn:aws:s3:::my-bucket/config/*"
    }

#### IAM role
- Create an IAM role choosing 'AWS service' and then 'Lambda' as the service that will use the role
- Add the following policies:
    - AWSLambdaBasicExecutionRole
    - AmazonInspectorReadOnlyAccess
    - awsInspectorReporterLambda    # the custom policy created above  

### creating the function

- Choose 'Create Function' and then 'Author from scratch'  
- Name the function awsInspectorReporter and choose Runtime 'Go 1.x' 
- Permissions:
    - choose 'Use an existing role'
    - choose the role created above
- Function code: 
    - Upload 'lambda_deployment.zip' from the [latest release page](https://github.com/jonhadfield/aws-inspector-reporter/releases/)
    - Ensure the runtime 'Go 1.x' is selected
    - Set Handler as 'main'
- Environment variables
    - Add AIR_CONFIG_PATH with value as the S3 directory where the configuration is uploaded, e.g.: s3://my-bucket/config
    - Optionally, add AIR_MAX_REPORT_AGE with value being the maximum number of days a report is considered valid for  