package air

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Target struct {
	ID             string `yaml:"id"`
	Alias          string `yaml:"alias"`
	RoleName       string `yaml:"roleName"`
	RoleExternalID string `yaml:"roleExternalId"`
}

type targetErrorsMap struct {
	target Target
	errors []annotatedError
}

type targetErrorsMaps []targetErrorsMap

type Targets []Target

func loadTargets(targetsFilePath string, debug bool) (targets Targets) {
	var err error
	targets, err = readTargets(targetsFilePath)
	if err != nil && debug {
		fmt.Println(err)
	}
	return
}

func parseTargetsFileContent(content []byte) (accounts Targets, err error) {
	var accountsInstance Targets
	unmarshalErr := yaml.Unmarshal(content, &accountsInstance)
	if unmarshalErr != nil {
		err = errors.WithStack(unmarshalErr)
		return
	}
	accounts = accountsInstance
	return
}

func readTargets(targetsPath string) (ret Targets, err error) {
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
