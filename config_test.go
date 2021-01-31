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

	_, err := GetConfig(td+"/config.yaml", "")
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

	_, err = GetConfig(td+"/config.yaml", "")
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

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "No LDAP servers provided", err.Error())
}

func TestErrorOnGetConfigDueToNoValidLDAPServerHostProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "",
			},
		},
		User:   BindUser{},
		Groups: AllowedGroups{},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "A valid server host must be provided", err.Error())
}

func TestErrorOnGetConfigDueToNoValidLDAPBindUserNameProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name: "",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"",
			},
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "No user name provided for group check", err.Error())
}

func TestErrorOnGetConfigDueToNoValidSecretForUserNameProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name: "/run/secrets/user",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"",
			},
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "Unable to obtain user name from: /run/secrets/user", err.Error())
}

func TestErrorOnGetConfigDueToNoValidLDAPBindPasswordProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name:     "user",
			Password: "",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"",
			},
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "No password provided for group check", err.Error())
}

func TestErrorOnGetConfigDueToNoValidSecretForPasswordProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name:     "user",
			Password: "/run/secrets/password",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"",
			},
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "Unable to obtain password from: /run/secrets/password", err.Error())
}

func TestErrorOnGetConfigDueToNoValidGroupDefinitionProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name:     "user",
			Password: "password",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"",
			},
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "No valid group definition provided", err.Error())
}

func TestErrorOnGetConfigDueToNoValidJWTKeyProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name:     "user",
			Password: "test",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"some-group",
			},
		},
		JWT: JWT{
			SigningKey: "",
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "No signing key provided for JWT", err.Error())
}

func TestErrorOnGetConfigDueToNoValidSecretJWTSigningKeyProvided(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name:     "user",
			Password: "password",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"some-group",
			},
		},
		JWT: JWT{
			SigningKey: "/run/secrets/jwt",
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	_, err = GetConfig(td+"/config.yaml", "")
	assert.Equal(t, "Unable to obtain signing key from: /run/secrets/jwt", err.Error())
}

func TestSuccessOnGetConfigWithDefaultValues(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host: "ldap.test.com",
			},
		},
		User: BindUser{
			Name:     "user",
			Password: "password",
		},
		Groups: AllowedGroups{
			Definitions: []string{
				"some-group",
			},
		},
		JWT: JWT{
			SigningKey: "test",
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	var config Config
	config, err = GetConfig(td+"/config.yaml", "")
	assert.NoError(t, err)

	assert.Equal(t, "ldaps", config.Servers[0].Protocol)
	assert.Equal(t, 636, config.Servers[0].Port)
	assert.Equal(t, "ldap.test.com", config.Servers[0].Host)
	assert.Equal(t, "user", config.Groups.ObjectClass)
	assert.Equal(t, "memberOf", config.Groups.MemberAttribute)
	assert.Equal(t, "some-group", config.Groups.Definitions[0])
	assert.Equal(t, "user", config.User.Name)
	assert.Equal(t, "password", config.User.Password)
	assert.Equal(t, "5m", config.JWT.Expiration)
	assert.Equal(t, "test", config.JWT.SigningKey)
}

func TestSuccessOnGetConfigWithCustomValues(t *testing.T) {
	setupAndTearDown := setupTestCase(t)
	defer setupAndTearDown(t)

	var data = Config{
		Servers: []Server{
			Server{
				Host:     "ldap.test.com",
				Port:     3333,
				Protocol: "ldaaaap",
			},
		},
		User: BindUser{
			Name:     td + "/user",
			Password: td + "/password",
		},
		Groups: AllowedGroups{
			ObjectClass:     "bot",
			MemberAttribute: "partOf",
			Definitions: []string{
				"some-group",
			},
		},
		JWT: JWT{
			Expiration: "15m",
			SigningKey: td + "/jwt",
		},
	}

	d, err := yaml.Marshal(&data)

	var file *os.File
	file, err = os.Create(td + "/config.yaml")
	assert.NoError(t, err)

	_, err = file.Write(d)
	assert.NoError(t, err)

	file.Close()

	file, err = os.Create(td + "/user")
	assert.NoError(t, err)

	_, err = file.WriteString("Tigger")
	assert.NoError(t, err)

	file.Close()

	file, err = os.Create(td + "/password")
	assert.NoError(t, err)

	_, err = file.WriteString("Jump1234")
	assert.NoError(t, err)

	file.Close()

	file, err = os.Create(td + "/jwt")
	assert.NoError(t, err)

	_, err = file.WriteString("test1234")
	assert.NoError(t, err)

	file.Close()

	var config Config
	config, err = GetConfig(td+"/config.yaml", td)
	assert.NoError(t, err)

	assert.Equal(t, "ldaaaap", config.Servers[0].Protocol)
	assert.Equal(t, 3333, config.Servers[0].Port)
	assert.Equal(t, "ldap.test.com", config.Servers[0].Host)
	assert.Equal(t, "bot", config.Groups.ObjectClass)
	assert.Equal(t, "partOf", config.Groups.MemberAttribute)
	assert.Equal(t, "some-group", config.Groups.Definitions[0])
	assert.Equal(t, "Tigger", config.User.Name)
	assert.Equal(t, "Jump1234", config.User.Password)
	assert.Equal(t, "15m", config.JWT.Expiration)
	assert.Equal(t, "test1234", config.JWT.SigningKey)
}

func TestSuccessOnGetConfigFromTemplate(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)
	var config Config
	config, err = GetConfig(dir+"/config.yaml", "")
	assert.NoError(t, err)

	assert.Equal(t, "ldaps", config.Servers[0].Protocol)
	assert.Equal(t, 636, config.Servers[0].Port)
	assert.Equal(t, "some-ldap-server.com", config.Servers[0].Host)
	assert.Equal(t, "user", config.Groups.ObjectClass)
	assert.Equal(t, "memberOf", config.Groups.MemberAttribute)
	assert.Equal(t, "some-group", config.Groups.Definitions[0])
	assert.Equal(t, "some-other-group", config.Groups.Definitions[1])
	assert.Equal(t, "testuser", config.User.Name)
	assert.Equal(t, "test1234", config.User.Password)
	assert.Equal(t, "30m", config.JWT.Expiration)
	assert.Equal(t, "somekey", config.JWT.SigningKey)
}
