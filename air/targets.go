package air

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type target struct {
	ID             string `yaml:"id"`
	Alias          string `yaml:"alias"`
	RoleName       string `yaml:"roleName"`
	RoleExternalID string `yaml:"roleExternalId"`
}

type targetErrorsMap struct {
	target target
	errors []annotatedError
}

type targetErrorsMaps []targetErrorsMap

type targets []target

func parseTargetsFileContent(content []byte) (accounts targets, err error) {
	var accountsInstance targets
	unmarshalErr := yaml.Unmarshal(content, &accountsInstance)
	if unmarshalErr != nil {
		err = errors.WithStack(unmarshalErr)
		return
	}
	accounts = accountsInstance
	return
}

func readTargets(targetsPath string) (ret targets, err error) {
	if _, err = os.Stat(targetsPath); err == nil {
		_, openErr := os.Open(targetsPath)
		if openErr != nil {
			err = errors.WithStack(openErr)
			return
		}
		targetsFileContent, readErr := ioutil.ReadFile(targetsPath)
		if readErr != nil {
			err = errors.WithStack(readErr)
			return
		}
		ret, err = parseTargetsFileContent(targetsFileContent)
	}

	return
}
