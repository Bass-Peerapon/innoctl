/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package service

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

// serviceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "generate service",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		Generate(serviceName, serviceType)
	},
}

func init() {
	ServiceCmd.Flags().StringVarP(&serviceName, "service_name", "n", "", "service name ")
	ServiceCmd.Flags().StringVarP(&serviceType, "service_type", "t", "", "service type (eg. postgres, mongodb) ")
	if err := ServiceCmd.MarkFlagRequired("service_name"); err != nil {
		fmt.Println(err)
	}
	if serviceType == "" {
		serviceType = "postgres"
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Service struct {
	Name           string
	Version        string // ex. v1
	CamelCase      string
	LowerCamelCase string
	defultPath     string
	serviceType    string
}

func NewService(name string, serviceType string) *Service {
	//{service name}
	//{service name}/{version}
	splitName := strings.Split(name, "/")
	var version string
	switch len(splitName) {
	case 2:
		version = splitName[1]
		fallthrough
	case 1:
		name = splitName[0]
	}
	defultPath := DEFULT_PATH + name
	if version != "" {
		defultPath = DEFULT_PATH + name + "/" + version
	}
	return &Service{
		Name:           name,
		Version:        version,
		CamelCase:      strcase.ToCamel(name),
		LowerCamelCase: strcase.ToLowerCamel(name),
		defultPath:     defultPath,
		serviceType:    serviceType,
	}
}

func (s Service) GetDefultPath() string {
	return s.defultPath
}

const (
	REPO        = "repository"
	USECASE     = "usecase"
	HTTP        = "http"
	HANDLER     = "handler"
	VALIDATOR   = "validator"
	DEFULT_PATH = "./service/"
	SER_NAME_GO = ".go"
)

const (
	FILE_MAIN = "./main.go"
)

var (
	GOPATH       = os.Getenv("GOPATH")
	SERVICE_PATH = "/service/"
)

var tmpRepoAdapter = `
package {{.Name}}
type {{ .CamelCase }}Repository interface {

}
`

var tmpUsecaseAdapter = `
package {{.Name}}
type {{ .CamelCase }}Usecase interface {

}
`

var tmpHttpAdapter = `
package {{.Name}} 
type {{ .CamelCase }}Handler interface {

}
`

var tmpRepo = `
package repository
import (
	"git.innovasive.co.th/backend/psql"
)
type {{.LowerCamelCase}}Repository struct {
	client             *psql.Client
}

func New{{.CamelCase}}Repository(client *psql.Client) {{.Name}}.{{.CamelCase}}Repository {
	return &{{.LowerCamelCase}}Repository{
		client:             client,	
	}
}
`

var tmpMongoRepo = `
package repository
import (
	"go.mongodb.org/mongo-driver/mongo"
	)
type {{.LowerCamelCase}}Repository struct {
	client 	*mongo.Client
	dbName 	string
}

func New{{.CamelCase}}Repository(client *mongo.Client, dbName string) {{.Name}}.{{.CamelCase}}Repository {
	return &{{.LowerCamelCase}}Repository{
		client: client,
		dbName: dbName,
	}
}
`

var tmpUsecase = `
package usecase
type {{.LowerCamelCase}}Usecase struct {
	{{.LowerCamelCase}}Repo {{.Name}}.{{.CamelCase}}Repository
}

func New{{.CamelCase}}Usecase({{.LowerCamelCase}}Repo {{.Name}}.{{.CamelCase}}Repository) {{.Name}}.{{.CamelCase}}Usecase {
	return &{{.LowerCamelCase}}Usecase{
		{{.LowerCamelCase}}Repo : {{.LowerCamelCase}}Repo,	
	}
}
`

var tmpHttp = `
package http
type {{.LowerCamelCase}}Handler struct {
	{{.LowerCamelCase}}Us {{.Name}}.{{.CamelCase}}Usecase
}

func New{{.CamelCase}}Handler({{.LowerCamelCase}}Us {{.Name}}.{{.CamelCase}}Usecase) {{.Name}}.{{.CamelCase}}Handler {
	return &{{.LowerCamelCase}}Handler{
		{{.LowerCamelCase}}Us:   {{.LowerCamelCase}}Us,	
	}
}
`

var tmpVal = `
package validator
type Validation struct{

}
`

var (
	serviceName string
	serviceType string
)

func GenerateWithChannel(serviceName string, serviceType string, c chan struct{}) error {
	serviceType = strings.ToLower(serviceType)
	s := NewService(serviceName, serviceType)

	if err := s.generateServiceDir(); err != nil {
		return err
	}
	c <- struct{}{}

	if err := s.generateReposiroryAdapter(); err != nil {
		return err
	}
	c <- struct{}{}

	if err := s.generateUsecaseAdapter(); err != nil {
		return err
	}
	c <- struct{}{}

	if err := s.generateHandlerAdapter(); err != nil {
		return err
	}
	c <- struct{}{}

	if err := s.generateHandler(); err != nil {
		return err
	}
	c <- struct{}{}

	if err := s.generateUsecase(); err != nil {
		return err
	}
	c <- struct{}{}

	switch serviceType {
	case "postgres":
		if err := s.generateReposirory(); err != nil {
			return err
		}

	case "mongodb":
		if err := s.generateMongoReposirory(); err != nil {
			return err
		}

	}
	c <- struct{}{}

	if err := s.generateValidator(); err != nil {
		return err
	}
	c <- struct{}{}

	close(c)
	return nil
}

func Generate(serviceName string, serviceType string) error {
	serviceType = strings.ToLower(serviceType)
	s := NewService(serviceName, serviceType)

	if err := s.generateServiceDir(); err != nil {
		return err
	}
	if err := s.generateReposiroryAdapter(); err != nil {
		return err
	}
	if err := s.generateUsecaseAdapter(); err != nil {
		return err
	}
	if err := s.generateHandlerAdapter(); err != nil {
		return err
	}
	if err := s.generateHandler(); err != nil {
		return err
	}
	if err := s.generateUsecase(); err != nil {
		return err
	}

	switch serviceType {
	case "postgres":
		if err := s.generateReposirory(); err != nil {
			return err
		}

	case "mongodb":
		if err := s.generateMongoReposirory(); err != nil {
			return err
		}

	}

	if err := s.generateValidator(); err != nil {
		return err
	}

	return nil
}

func (s Service) generateServiceDir() error {
	if err := os.MkdirAll(s.GetDefultPath(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

// ./service/{service_name}/repository/{serivce_name}_repository.go
func (s Service) generateReposirory() error {
	dir := "./" + s.GetDefultPath() + "/" + REPO
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + REPO + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("repository").Parse(tmpRepo)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}
	return nil
}

func (s Service) generateMongoReposirory() error {
	dir := "./" + s.GetDefultPath() + "/" + REPO
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + REPO + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("repository").Parse(tmpMongoRepo)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}
	return nil
}

// ./service/{service_name}/usecase/{serivce_name}_usecase.go
func (s Service) generateUsecase() error {
	dir := "./" + s.GetDefultPath() + "/" + USECASE
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + USECASE + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("usecase").Parse(tmpUsecase)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}
	return nil
}

// ./service/{service_name}/http/{serivce_name}_handler.go
func (s Service) generateHandler() error {
	dir := "./" + s.GetDefultPath() + "/" + HTTP
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + HTTP + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("http").Parse(tmpHttp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}
	return nil
}

// ./service/{service_name}/repository.go
func (s Service) generateReposiroryAdapter() error {
	fn := s.GetDefultPath() + "/" + REPO + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("repository").Parse(tmpRepoAdapter)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}
	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}

	return nil
}

// ./service/{service_name}/usecase.go
func (s Service) generateUsecaseAdapter() error {
	fn := s.GetDefultPath() + "/" + USECASE + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("usecase").Parse(tmpUsecaseAdapter)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}

	return nil
}

// ./service/{service_name}/http.go
func (s Service) generateHandlerAdapter() error {
	fn := s.GetDefultPath() + "/" + HTTP + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("http").Parse(tmpHttpAdapter)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}

	return nil
}

func (s Service) generateValidator() error {
	dir := "./" + s.GetDefultPath() + "/" + VALIDATOR
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + VALIDATOR + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("validator").Parse(tmpVal)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	if err != nil {
		return err
	}
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}

	return nil
}
