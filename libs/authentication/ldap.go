package authentication

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"net/http"
	"strconv"
	"time"
)

type adLDAP struct {
	server  string
	basedn  string
	filter  string
	timeout time.Duration
}

func (l *adLDAP) setup(params map[string]string) error {
	var found bool
	var swap string

	if l.server, found = params["server"]; !found {
		return errors.New("server key not found on ldap authentication method")
	}

	if l.basedn, found = params["basedn"]; !found {
		return errors.New("basedn key not found on ldap authentication method")
	}

	if l.filter, found = params["filter"]; !found {
		return errors.New("filter key not found on ldap authentication method")
	}

	if swap, found = params["timeout"]; found {
		if t, err := strconv.Atoi(swap); err != nil {
			return errors.New("timeout must be an integer, found instead: " + swap)
		} else {
			l.timeout = time.Duration(t) * time.Second
		}
	} else {
		l.timeout = ldap.DefaultTimeout
	}

	return nil
}

func (l *adLDAP) validate(c *gin.Context) (int, error) {
	var bindUsername, bindPassword string
	var found bool
	var conn *ldap.Conn
	var err error

	if bindUsername, bindPassword, found = c.Request.BasicAuth(); !found {
		return http.StatusNetworkAuthenticationRequired, nil
	}
	fmt.Println(bindUsername, bindPassword)
	fmt.Println(l.basedn)
	fmt.Println(fmt.Sprintf(l.filter, bindUsername))

	ldap.DefaultTimeout = l.timeout
	if conn, err = ldap.DialURL(l.server); err != nil {
		return http.StatusInternalServerError, err
	}

	// if err = conn.Bind(bindUsername, bindPassword); err != nil {
	// 	return http.StatusUnauthorized, fmt.Errorf("Search Error: %s", err)
	// }
	// searchReq := ldap.NewSearchRequest(
	// 	l.basedn,
	// 	ldap.ScopeBaseObject,
	// 	ldap.NeverDerefAliases,
	// 	0,
	// 	0,
	// 	false,
	// 	fmt.Sprintf(l.filter, bindUsername),
	// 	[]string{},
	// 	nil,
	// )
	// result, err := conn.Search(searchReq)
	// if err != nil {
	// 	return http.StatusUnauthorized, fmt.Errorf("Search Error: %s", err)
	// }

	controls := []ldap.Control{}
	controls = append(controls, ldap.NewControlBeheraPasswordPolicy())
	bindRequest := ldap.NewSimpleBindRequest("cn=admin,dc=example,dc=com", "password", controls)

	result, err := conn.SimpleBind(bindRequest)
	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("Search Error: %s", err)
	}

	fmt.Println(result)
	conn.Debug.Printf("foobar")

	return http.StatusUnauthorized, nil
}
