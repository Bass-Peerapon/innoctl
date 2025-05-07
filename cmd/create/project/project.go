/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package project

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	templater "github.com/Bass-Peerapon/innoctl/cmd/create/project/template"
	"github.com/Bass-Peerapon/innoctl/utils"
	"github.com/spf13/cobra"
)

type Project struct {
	ProjectName string
	ModuleName  string
}

// create/project/projectCmd represents the create/project/project command
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create project",
	Long:  `Create a new Golang project with a specific structure`,
	Run: func(cmd *cobra.Command, args []string) {
		flagName := cmd.Flag("name").Value.String()
		if flagName != "" && doesDirectoryExistAndIsNotEmpty(flagName) {
			log.Fatalf(
				"directory '%s' already exists and is not empty. Please choose a different name",
				flagName,
			)
		}

		gopath := os.Getenv("GOPATH")
		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get current working directory: %v", err)
		}
		moduleName := fmt.Sprintf(
			"%s/%s",
			strings.ReplaceAll(currentWorkingDir, gopath+"/src/", ""),
			flagName,
		)

		project := Project{
			ProjectName: flagName,
			ModuleName:  moduleName,
		}

		projectPath := filepath.Join(currentWorkingDir, project.ProjectName)
		createDirectory(projectPath)
		fmt.Println("Created project directory")

		templates := map[string][]byte{
			"go.mod":                         templater.GoMod,
			"docker-compose.yml":             templater.DockerCompose,
			"Dockerfile-development":         templater.DockerfileDevelopment,
			"Dockerfile-production":          templater.DockerfileProduction,
			".gitignore":                     templater.GitIgnore,
			".dockerignore":                  templater.DockerIgnore,
			".env.example":                   templater.EnvExample,
			"Makefile":                       templater.Makefile,
			"README.md":                      templater.Readme,
			"sonar-project.properties":       templater.SonarProperties,
			"constants/constants.go":         templater.Constants,
			"middleware/middleware.go":       templater.Middleware,
			"middleware/openapi.go":          templater.OpenAPI,
			"middleware/tracer.go":           templater.Tracer,
			"utils/opentracing/init.go":      templater.OpenTracingInit,
			"utils/redis/client.go":          templater.RedisClient,
			"utils/pagination/pagination.go": templater.Pagination,
			"main.go":                        templater.Main,
		}

		totalSteps := len(
			templates,
		) + 4 // +4 for directories creation git init, go get, and go tidy
		currentStep := 1

		// git init
		if err := utils.GitInit(projectPath); err != nil {
			log.Fatalf("Error in git init: %v", err)
		}
		fmt.Printf("[%d/%d] Ran git init\n", currentStep, totalSteps)
		currentStep++

		for path, template := range templates {
			renderTemplateToFile(projectPath, path, string(template), project)
			fmt.Printf("[%d/%d] Created file: %s\n", currentStep, totalSteps, path)
			currentStep++
		}

		createEmptyDirectoriesWithGitkeep(projectPath, []string{
			"assets",
			"migrations/database/postgres",
			"models",
			"proto",
			"service",
		})
		fmt.Printf("[%d/%d] Created empty directories with .gitkeep\n", currentStep, totalSteps)
		currentStep++

		// go get and go mod tidy
		if err := utils.GoGetPackage(projectPath, []string{"google.golang.org/genproto@latest"}); err != nil {
			log.Fatalf("Error in go get: %v", err)
		}
		fmt.Printf("[%d/%d] Ran go get\n", currentStep, totalSteps)
		currentStep++

		if err := utils.GoTidy(projectPath); err != nil {
			log.Fatalf("Error in go mod tidy: %v", err)
		}
		fmt.Printf("[%d/%d] Ran go mod tidy\n", currentStep, totalSteps)

		fmt.Printf("Project %s created\n", project.ProjectName)
	},
}

func init() {
	ProjectCmd.Flags().StringP("name", "n", "", "Name of project to create")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// create/project/projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// create/project/projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// doesDirectoryExistAndIsNotEmpty checks if the directory exists and is not empty
func doesDirectoryExistAndIsNotEmpty(name string) bool {
	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Printf("could not read directory: %v", err)
		}
		if len(dirEntries) > 0 {
			return true
		}
	}
	return false
}

func createDirectory(path string) {
	if err := os.MkdirAll(path, 0o751); err != nil {
		log.Fatalf("Error creating directory %s: %v", path, err)
	}
}

func renderTemplateToFile(basePath, relativePath, templateStr string, data interface{}) {
	t, err := template.New(filepath.Base(relativePath)).Parse(templateStr)
	if err != nil {
		log.Fatalf("Error parsing template for %s: %v", relativePath, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Fatalf("Error executing template for %s: %v", relativePath, err)
	}
	fullPath := filepath.Join(basePath, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o751); err != nil {
		log.Fatalf("Error creating directories for %s: %v", fullPath, err)
	}
	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		log.Fatalf("Error writing file %s: %v", fullPath, err)
	}
}

func createEmptyDirectoriesWithGitkeep(basePath string, directories []string) {
	for _, dir := range directories {
		path := filepath.Join(basePath, dir)
		createDirectory(path)
		if err := os.WriteFile(filepath.Join(path, ".gitkeep"), []byte(""), 0644); err != nil {
			log.Fatalf("Error creating .gitkeep in %s: %v", path, err)
		}
	}
}
