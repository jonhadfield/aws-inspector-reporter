To enable AIR to report for multiple accounts the provided AWS credentials must be able to assume a role with Inspector permissions in each of the target accounts. 

Example Trust Relationship policy document:

```
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::012345678901:role/<assuming-role-or-user>"   
      }
    }
  ]
}
```

Example document requiring an external id:

```
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::012345678901:role/<assuming-role-or-user>"
                             
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "<external id>"
        }
      }
    }
  ]
}
```