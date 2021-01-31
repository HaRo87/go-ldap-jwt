package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// Config represents the configuration for go-ldap-jwt
type Config struct {
	Servers []Server      `yaml:"servers"`
	User    BindUser      `yaml:"user"`
	Groups  AllowedGroups `yaml:"groups"`
	JWT     JWT           `yaml:"jwt"`
}

// Server represents the server config portion for
// LDAP servers to be used
type Server struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// BindUser represents the user which will be used to
// check group membership
type BindUser struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

// AllowedGroups represents the groups which will be
// used to limit access
type AllowedGroups struct {
	ObjectClass     string   `yaml:"objectclass"`
	MemberAttribute string   `yaml:"memberattribute"`
	Definitions     []string `yaml:"definitions"`
}

// JWT represents the JWT configuration
type JWT struct {
	Expiration string `yaml:"expire"`
	SigningKey string `yaml:"signingkey"`
}

const defaultProtocol = "ldaps"
const defaultPort = 636
const defaultObjectClass = "user"
const defaultMemberAttribute = "memberOf"
const secretsLocation = "/run/secrets/"
const defaultExpiration = "5m"

// GetConfig reads the provided config file and
// performs some checks
func GetConfig(path, secrets string) (Config, error) {
	var config Config
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return Config{}, fmt.Errorf("Unable to read config: %s", path)
	}

	err = yaml.UnmarshalStrict(content, &config)

	if err != nil {
		return Config{}, fmt.Errorf("Unable to unmarshal config: %s", path)
	}

	if strings.TrimSpace(secrets) == "" {
		secrets = secretsLocation
	}

	if len(config.Servers) == 0 {
		return config, fmt.Errorf("No LDAP servers provided")
	}

	for i, server := range config.Servers {
		if strings.TrimSpace(server.Protocol) == "" {
			config.Servers[i].Protocol = defaultProtocol
		}
		if strings.TrimSpace(server.Host) == "" {
			return config, fmt.Errorf("A valid server host must be provided")
		}
		if server.Port <= 0 {
			config.Servers[i].Port = defaultPort
		}
	}

	if len(config.Groups.Definitions) > 0 {
		if strings.TrimSpace(config.User.Name) == "" {
			return config, fmt.Errorf("No user name provided for group check")
		}
		if strings.Contains(config.User.Name, secrets) {
			var name []byte
			name, err = ioutil.ReadFile(config.User.Name)
			if err != nil {
				return config, fmt.Errorf("Unable to obtain user name from: %s", config.User.Name)
			}
			config.User.Name = string(name)
		}

		if strings.TrimSpace(config.User.Password) == "" {
			return config, fmt.Errorf("No password provided for group check")
		}
		if strings.Contains(config.User.Password, secrets) {
			var password []byte
			password, err = ioutil.ReadFile(config.User.Password)
			if err != nil {
				return config, fmt.Errorf("Unable to obtain password from: %s", config.User.Password)
			}
			config.User.Password = string(password)
		}

		if strings.TrimSpace(config.Groups.ObjectClass) == "" {
			config.Groups.ObjectClass = defaultObjectClass
		}
		if strings.TrimSpace(config.Groups.MemberAttribute) == "" {
			config.Groups.MemberAttribute = defaultMemberAttribute
		}

		for _, groupDefinitions := range config.Groups.Definitions {
			if strings.TrimSpace(groupDefinitions) == "" {
				return config, fmt.Errorf("No valid group definition provided")
			}
		}

		if strings.TrimSpace(config.JWT.Expiration) == "" {
			config.JWT.Expiration = defaultExpiration
		}

		if strings.TrimSpace(config.JWT.SigningKey) == "" {
			return config, fmt.Errorf("No signing key provided for JWT")
		}
		if strings.Contains(config.JWT.SigningKey, secrets) {
			var signingKey []byte
			signingKey, err = ioutil.ReadFile(config.JWT.SigningKey)
			if err != nil {
				return config, fmt.Errorf("Unable to obtain signing key from: %s", config.JWT.SigningKey)
			}
			config.JWT.SigningKey = string(signingKey)
		}
	}

	return config, nil
}
