package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"simplerest/libs/authentication"
	"simplerest/libs/database"
	"simplerest/libs/settings"
)

type Server struct {
	settings settings.Settings
	db       database.Database
	auth     authentication.AuthenticationChain
}

func New(s settings.Settings) Server {
	return Server{
		settings: s,
		db:       database.New(s.Database),
		auth:     authentication.New(s.Authentication),
	}
}

func (s *Server) Initialize() error {
	if err := s.db.Initialize(); err != nil {
		return err
	}

	if err := s.auth.Initialize(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() error {
	var err error = nil
	r := gin.Default()

	// Middlewares
	r.Use(s.auth.Chain)

	// route handlers
	for _, res := range s.settings.Resources {
		h := resourceHandler{res, s.db.DB}
		r.Handle(res.Method, res.Location, h.handle)
	}

	if s.settings.Static.Location != "" {
		r.Static(s.settings.Static.Location, s.settings.Static.Root)
	}

	if s.settings.Static.SPA != "" {
		r.NoRoute(func(c *gin.Context) {
			c.File(s.settings.Static.SPA)
		})
	}

	if s.settings.Templates != "" {
		var files []string
		err := filepath.Walk(s.settings.Templates, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}

		r.LoadHTMLFiles(files...)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.settings.Port),
		Handler: r,
	}

	if s.settings.TLSCert != "" && s.settings.TLSKey != "" {
		if _, err := os.Stat(s.settings.TLSCert); err != nil {
			return errors.New("TLS certificate file does not exist on " + s.settings.TLSCert)
		}
		if _, err := os.Stat(s.settings.TLSCert); err != nil {
			return errors.New("TLS certificate key file does not exist on " + s.settings.TLSCert)
		}
		err = server.ListenAndServeTLS(s.settings.TLSCert, s.settings.TLSKey)
	} else {
		err = server.ListenAndServe()
	}
	return err
}
