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

const secrets = "/run/secrets"
const defaultProtocol = "ldaps"
const defaultPort = 3269
const defaultObjectClass = "user"
const defaultMemberAttribute = "memberOf"

// GetConfig reads the provided config file and
// performs some checks
func GetConfig(path string) (Config, error) {
	var config Config
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return Config{}, fmt.Errorf("Unable to read config: %s", path)
	}

	err = yaml.UnmarshalStrict(content, &config)

	if err != nil {
		return Config{}, fmt.Errorf("Unable to unmarshal config: %s", path)
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
		if strings.TrimSpace(config.User.Password) == "" {
			return config, fmt.Errorf("No password provided for group check")
		}

		if strings.TrimSpace(config.Groups.ObjectClass) == "" {
			config.Groups.ObjectClass = defaultObjectClass
		}
		if strings.TrimSpace(config.Groups.MemberAttribute) == "" {
			config.Groups.ObjectClass = defaultMemberAttribute
		}

		for _, groupDefinitions := range config.Groups.Definitions {
			if strings.TrimSpace(groupDefinitions) == "" {
				return config, fmt.Errorf("No valid group definition provided")
			}
		}
	}

	return config, nil
}
