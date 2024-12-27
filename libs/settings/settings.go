package settings

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const (
	AuthMethodDBToken  = "dbtoken"
	AuthMethodHTPasswd = "htpasswd"
	AuthMethodLDAP     = "ldap"
)

type Static struct {
	Location string
	Root     string
	SPA      string
}
type Resource struct {
	Location string
	Method   string
	Query    string
	Template string
	Params   map[string]interface{}
}

type Database struct {
	Driver     string
	Location   string
	Migrations string
}

type Authentication struct {
	Method string
	Params map[string]string
}

type proxy struct {
	Location string
	Backends []string
}

type Settings struct {
	Port           uint16
	Static         Static
	Scripts        string
	Templates      string
	LogLevel       string
	TLSCert        string
	TLSKey         string
	workingDir     string
	Database       Database
	Resources      []Resource `toml:"resource"`
	Proxy          proxy
	Authentication []Authentication
	Secret         string
	config_file    string
}

func (s *Settings) validate() error {
	for idx, r := range s.Resources {
		switch r.Method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
			continue
		default:
			return errors.New("Invalid HTTP Method '" + r.Method + "' found on #" + strconv.Itoa(idx) + " resource")
		}
	}
	for idx, r := range s.Authentication {
		switch r.Method {
		case AuthMethodDBToken, AuthMethodLDAP, AuthMethodHTPasswd:
			continue
		default:
			return errors.New("Invalid Authentication Method '" + r.Method + "' found on #" + strconv.Itoa(idx) + " authentication")
		}
	}
	return nil
}

func envToMap() map[string]string {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		envMap[split_v[0]] = strings.Join(split_v[1:], "=")
	}
	return envMap
}

func (s *Settings) Display() {
	fmt.Println("Server settings")
	fmt.Println("===============")
	fmt.Println("WorkingDir:", s.workingDir)
	fmt.Println("config file:", s.config_file)
	fmt.Println("Port:", s.Port)
	fmt.Println("Scripts:", s.Scripts)
	fmt.Println("Templates:", s.Templates)
	fmt.Println("TLS Certificate:", s.TLSCert)
	fmt.Println("TLS Key:", s.TLSKey)
	fmt.Println("")
	fmt.Println("Static settings")
	fmt.Println("===============")
	fmt.Println("Location:", s.Static.Location)
	fmt.Println("Root:", s.Static.Root)
	fmt.Println("SPA:", s.Static.SPA)
	fmt.Println("")
	fmt.Println("Database settings")
	fmt.Println("===============")
	fmt.Println("Driver:", s.Database.Driver)
	fmt.Println("Migrations:", s.Database.Migrations)
	fmt.Println("")
	if len(s.Authentication) > 0 {
		fmt.Println("Authentication settings")
		fmt.Println("===============")
		for idx, a := range s.Authentication {
			fmt.Printf("Method %d: %s", idx, a.Method)
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func New(full_file_path string) Settings {
	pwd := filepath.Dir(full_file_path)
	config_file := filepath.Base(full_file_path)

	return Settings{
		Port:      8080,
		Scripts:   filepath.Join(pwd, "scripts"),
		Templates: filepath.Join(pwd, "templates"),
		LogLevel:  "info",
		Secret:    uuid.New().String(),
		Static: Static{
			Root: filepath.Join(pwd, "statics"),
		},
		workingDir:  pwd,
		config_file: config_file,
	}
}

func Parse(full_file_path string) (Settings, error) {
	var config bytes.Buffer
	s := New(full_file_path)
	envMap := envToMap()
	tmpl, err := template.ParseFiles(s.config_file)
	if err != nil {
		return s, err
	}

	if err := tmpl.Execute(&config, envMap); err != nil {
		return s, err
	}

	if _, err := toml.Decode(config.String(), &s); err != nil {
		return s, err
	}

	if err := s.validate(); err != nil {
		return s, err
	}
	return s, nil
}
