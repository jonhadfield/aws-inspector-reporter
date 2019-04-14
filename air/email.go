package air

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	gomail "gopkg.in/gomail.v2"
)

// Email has the settings to be used to connect to a mail server and the properties of the email to send
type Email struct {
	Provider   string
	Host       string
	Port       string
	Username   string
	Password   string
	Region     string
	Source     string
	Subject    string
	Recipients []string
}

func emailConfigDefined(email Email) (result bool) {
	if !reflect.DeepEqual(email, Email{}) {
		result = true
	}
	return
}
func extractEmail(input string) (output string) {
	if strings.Contains(input, "<") {
		output = getStringInBetween(input, "<", ">")
	} else {
		output = input
	}
	return
}
func getStringInBetween(str, start, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}

func validateEmailSettings(email Email) (err error) {
	supportedProviders := []string{"ses", "smtp"}
	if emailConfigDefined(email) {
		if email.Provider == "" {
			err = fmt.Errorf("email provider not specified")
			return
		}

		// TODO: Check minimum configuration (to, from, etc.)
		if email.Source == "" {
			err = fmt.Errorf("email source not specified")
			return
		}

		if !stringInSlice(email.Provider, supportedProviders) {
			err = fmt.Errorf("email provider '%s' not supported", email.Provider)
			return
		}
		emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
		// validate recipient email addresses
		for _, emailAddr := range email.Recipients {
			if !emailRegexp.MatchString(extractEmail(emailAddr)) {
				err = fmt.Errorf("invalid email address '%s'", extractEmail(emailAddr))
				return
			}
		}
		// validate source email address
		if !emailRegexp.MatchString(extractEmail(email.Source)) {
			err = fmt.Errorf("invalid email address '%s'", extractEmail(email.Source))
			return
		}

		// TODO: Check provider specific configuration
	}
	return err
}
func emailReport(sess *session.Session, reportPath string, email Email, deleteAfter bool) (err error) {
	err = validateEmailSettings(email)
	if err != nil {
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", email.Source)
	var emailSubject string
	if email.Subject != "" {
		emailSubject = email.Subject
	} else {
		emailSubject = "AWS Inspector Report"
	}

	msg.SetHeader("Subject", emailSubject)
	// TODO: Add summary to email body
	body := "attached"
	msg.SetBody("text/html", body)
	msg.Attach(reportPath)

	var emailRaw bytes.Buffer
	_, err = msg.WriteTo(&emailRaw)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	switch email.Provider {
	case "ses":
		msg.SetHeader("To", strings.Join(email.Recipients, ","))
		svc := ses.New(sess, &aws.Config{Region: ptrToStr(email.Region)})
		message := ses.RawMessage{Data: emailRaw.Bytes()}
		source := aws.String(email.Source)
		var destinations []*string
		for _, dest := range email.Recipients {
			destinations = append(destinations, ptrToStr(dest))
		}
		input := ses.SendRawEmailInput{Source: source, Destinations: destinations, RawMessage: &message}
		_, err = svc.SendRawEmail(&input)
		if err != nil {
			delErr := deleteFile(reportPath)
			if delErr != nil {
				err = errors.WithStack(delErr)
				return
			}
			return
		}

	case "smtp":
		msg.SetHeader("To", email.Recipients...)
		host := email.Host
		port, _ := strconv.Atoi(email.Port)
		dialer := gomail.NewPlainDialer(host, port, email.Username, email.Password)
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}
		dialer.TLSConfig = tlsConfig
		err = dialer.DialAndSend(msg)
		if err != nil {
			delErr := deleteFile(reportPath)
			if delErr != nil {
				err = errors.WithStack(delErr)
				return
			}
			return
		}
	}
	if deleteAfter {
		err = deleteFile(reportPath)
	}
	return err
}
