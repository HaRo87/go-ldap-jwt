# go-ldap-jwt

![Testing Go Code](https://github.com/HaRo87/go-ldap-jwt/workflows/Testing%20Go%20Code/badge.svg)
[![codecov](https://codecov.io/gh/HaRo87/go-ldap-jwt/branch/main/graph/badge.svg?token=MPUVSQ1TYA)](https://codecov.io/gh/HaRo87/go-ldap-jwt)

<img src="https://img.shields.io/badge/Go-1.15+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;<img src="https://img.shields.io/badge/license-mit-red?style=for-the-badge&logo=none" alt="license" />

A small Go library for creating JWTs based on LDAP user credentials 
including group checks.

## ⚡️ Getting started

### Installation

To install go-ldap-jwt, use `go get`:

```bash
go get github.com/haro87/go-ldap-jwt
```

This will make the following packages available to you:

```
github.com/haro87/go-ldap-jwt/config
```

### Usage

1. Import the `github.com/haro87/go-ldap-jwt/config` and use it to retrieve the configuration:

```go
package yours

import (
  "github.com/haro87/go-ldap-jwt/config"
)

func LoadConfig() {

  // In case you use secrets in a Docker container you can
  // provide the default secrets location via the `secrets`
  // parameter. If nothing is provided then the default:
  // `/run/secrets/` is used.
  config, err := config.GetConfig("path/to/config.yaml", "")

}
```

## ⚙️ Configuration

```yaml
# ./configs/config.yaml

# Servers config
servers:
  -
    protocol: ldaps # defaults to ldaps if not provided
    host: some-ldap-server.com
    port: 636 # defaults to 636 if not provided

# Bind user config
# can also work with reading secrets from file
user: 
  name: testuser
  password: test1234 

# Groups config
groups:
  objectclass: user # defaults to user if not provided
  memberattribute: memberOf # defaults to memberOf if not provided
  - some-group
  - some-other-group
```

## ⚠️ License

MIT &copy; [HaRo87](https://github.com/HaRo87).

