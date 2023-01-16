package authentication

import (
	"errors"
	"net/http"
	"simplerest/libs/database"
	"simplerest/libs/settings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type dbToken struct {
	db          *sqlx.DB
	tokenHeader string
	query       string
	migrations  string
}

func (d *dbToken) setup(params map[string]string) error {
	var driver, location string
	if r, found := params["driver"]; !found {
		return errors.New("driver key not found on dbtoken authentication method")
	} else {
		driver = r
	}

	if l, found := params["location"]; !found {
		return errors.New("location key not found on dbtoken authentication method")
	} else {
		location = l
	}

	if q, found := params["query"]; !found {
		return errors.New("query key not found on dbtoken authentication method")
	} else {
		d.query = q
	}

	if h, found := params["header"]; !found {
		d.tokenHeader = "X-Auth-Token"
	} else {
		d.tokenHeader = h
	}

	d.migrations, _ = params["migrations"]
	db := database.New(settings.Database{
		Driver:     driver,
		Location:   location,
		Migrations: d.migrations,
	})

	if err := db.Initialize(); err != nil {
		return err
	}

	d.db = db.DB
	return nil
}

func (d *dbToken) validate(c *gin.Context) (int, error) {
	var username []string
	token := c.GetHeader(d.tokenHeader)

	if err := d.db.Select(&username, d.query, token); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(username) > 0 {
		c.Set("username", username[0])
		return http.StatusAccepted, nil
	}

	return http.StatusUnauthorized, nil
}
