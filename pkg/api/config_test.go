package api

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dummyConfigFile = `hosts:
  github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	emptyConfigFile = ""

	dummyHostsFile = ` github.com:
      user: testuser
      oauth_token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
`
	emptyHostsFile = ""
)

func TestCanLoadConfigFromFile(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "TestCanLoadConfigFromFile")
	if err != nil {
		panic(err)
	}
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestCanLoadConfigFromFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(configFile.Name())
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	err = ioutil.WriteFile(hostsFile.Name(), []byte(emptyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := ParseConfig(configFile.Name(), hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestCanLoadEmptyConfigFromFile(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "TestCanLoadEmptyConfigFromFile")
	if err != nil {
		panic(err)
	}
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestCanLoadEmptyConfigFromFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(configFile.Name())
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(emptyConfigFile), 0600)
	assert.NoError(t, err)

	err = ioutil.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := ParseConfig(configFile.Name(), hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseConfigFile(t *testing.T) {
	configFile, err := ioutil.TempFile(os.TempDir(), "TestParseConfigFile")
	if err != nil {
		panic(err)
	}
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestParseConfigFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(configFile.Name())
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(dummyConfigFile), 0600)
	assert.NoError(t, err)

	err = ioutil.WriteFile(hostsFile.Name(), []byte(emptyHostsFile), 0600)
	assert.NoError(t, err)

	config, err := parseConfigFile(configFile.Name(), hostsFile.Name())
	assert.NoError(t, err)

	user := config.GetUser("github.com")
	assert.Equal(t, "testuser", user)

	token := config.GetToken("github.com")
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}

func TestParseHostsFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile(os.TempDir(), "TestParseHostsFile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(hostsFile.Name())

	err = ioutil.WriteFile(hostsFile.Name(), []byte(dummyHostsFile), 0600)
	assert.NoError(t, err)

	hosts, err := parseHostsFile(hostsFile.Name())
	assert.NoError(t, err)

	user := hosts["github.com"].User
	assert.Equal(t, "testuser", user)

	token := hosts["github.com"].Token
	assert.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token)
}
