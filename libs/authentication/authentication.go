package authentication

import (
	"errors"
	"net/http"
	"simplerest/libs/settings"
	"strings"

	"github.com/gin-gonic/gin"
)

type authenticationMethod interface {
	setup(map[string]string) error
	validate(c *gin.Context) (int, error)
}

type AuthenticationChain struct {
	auth    []settings.Authentication
	methods []authenticationMethod
}

func authenticationFactory(m string) authenticationMethod {
	var method authenticationMethod
	switch m {
	case settings.AuthMethodDBToken:
		method = &dbToken{}
	case settings.AuthMethodLDAP:
		method = &adLDAP{}
	default:
		method = nil
	}
	return method
}

func New(s []settings.Authentication) AuthenticationChain {
	return AuthenticationChain{auth: s}
}
func (a *AuthenticationChain) Initialize() error {
	for _, m := range a.auth {
		method := authenticationFactory(m.Method)
		if method == nil {
			return errors.New("Authentication method not implemented yet: " + m.Method)
		}
		if err := method.setup(m.Params); err != nil {
			return err
		}
		a.methods = append(a.methods, method)
	}
	return nil
}

func abort(c *gin.Context, status int, message ...string) {
	if len(message) == 0 {
		message = append(message, "unauthorized")
	}
	c.AbortWithStatusJSON(status, gin.H{
		"success": false,
		"message": strings.Join(message, " "),
	})
}
func askBasicAuth(c *gin.Context, status int, message ...string) {
	c.Writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	abort(c, status, message...)
}

func (a *AuthenticationChain) Chain(c *gin.Context) {
	status := http.StatusUnauthorized
	askcreds := false
	var err error = nil
	if len(a.methods) <= 0 {
		c.Next()
		return
	}

	for _, m := range a.methods {
		if status, err = m.validate(c); err != nil {
			break
		} else if status == http.StatusAccepted {
			break
		} else if status == http.StatusNetworkAuthenticationRequired {
			askcreds = true
		}
	}

	if err != nil {
		abort(c, http.StatusInternalServerError, err.Error())
	} else if status == http.StatusAccepted {
		c.Next()
	} else if askcreds {
		askBasicAuth(c, http.StatusUnauthorized, "credentials required")
	} else {
		abort(c, http.StatusUnauthorized, "unauthorized")
	}

}
