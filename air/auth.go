package air

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

func getAssumeRoleCreds(input getAssumeRoleCredsInput) (creds *credentials.Credentials, err error) {
	var roleArn string
	if input.RoleArn != "" {
		roleArn = input.RoleArn
	} else {
		roleArn = genRoleArn(input.AccountID, input.RoleName)
	}
	// TODO: Test without external id specified
	creds = stscreds.NewCredentials(input.Sess, roleArn, func(p *stscreds.AssumeRoleProvider) {
		p.ExternalID = &input.ExternalID
	})
	_, err = creds.Get()
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

type getAssumeRoleCredsInput struct {
	Sess       *session.Session
	AccountID  string
	RoleArn    string
	RoleName   string
	ExternalID string
}

func genRoleArn(accountID, roleName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
}
