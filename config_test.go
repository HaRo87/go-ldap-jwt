package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

var td string

func setupTestCase(t *testing.T) func(t *testing.T) {
	td, _ = ioutil.TempDir("", "config")
	return func(t *testing.T) {
		td = ""
	}
}

func TestErrorOnGetConfigDueToFileDoesNotExist(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	_, err := GetConfig(td + "/config.yaml")
	assert.Equal(t, "Unable to read config: "+td+"/config.yaml", err.Error())
}

func TestErrorOnGetConfigDueToUnmarshallingFails(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = "wrong: entry"

	file, err := os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.WriteString(data)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td + "/config.yaml")
	assert.Equal(t, "Unable to unmarshal config: "+td+"/config.yaml", err.Error())
}

func TestErrorOnGetConfigDueToNoLDAPServersProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{},
		User:    BindUser{},
		Groups:  AllowedGroups{},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td + "/config.yaml")
	assert.Equal(t, "No LDAP servers provided", err.Error())
}
