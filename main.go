package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed template/*
var templateFS embed.FS

type TemplateData struct {
	Name string
}

var projectName string

const banner = `
    ╔═══════════════════════════════════════════════════════════════════╗
    ║                                                                   ║
    ║   ██████╗  ██████╗        ███████╗██╗  ██╗███████╗██╗     ███████╗║
    ║  ██╔════╝ ██╔═══██╗       ██╔════╝██║ ██╔╝██╔════╝██║     ██╔════╝║
    ║  ██║  ███╗██║   ██║ █████╗███████╗█████╔╝ █████╗  ██║     █████╗  ║
    ║  ██║   ██║██║   ██║ ╚════╝╚════██║██╔═██╗ ██╔══╝  ██║     ██╔══╝  ║
    ║  ╚██████╔╝╚██████╔╝       ███████║██║  ██╗███████╗███████╗███████╗║
    ║   ╚═════╝  ╚═════╝        ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝║
    ║                                                                   ║
    ║            🚀 Modern Go REST API Generator with UberFX            ║
    ║                                                                   ║
    ╚═══════════════════════════════════════════════════════════════════╝
`

var rootCmd = &cobra.Command{
	Use:   "go-skeleton",
	Short: "A CLI tool to generate Go REST API project structure",
	Long: banner + `
A powerful CLI tool to generate production-ready Go REST API projects with:
• Clean Architecture structure
• UberFX dependency injection
• Structured logging with Zap
• JWT authentication
• GORM with PostgreSQL
• Docker configuration
• Comprehensive middleware
• API documentation with Swagger`,
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new REST API project",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if projectName == "" {
			return fmt.Errorf("please provide a project name using --name flag")
		}

		if err := validateProjectName(projectName); err != nil {
			return err
		}

		if err := generateProject(projectName); err != nil {
			return fmt.Errorf("failed to generate project: %w", err)
		}

		printSuccessMessage(projectName)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		versionBanner := `
    ╔═══════════════════════════════════════════════════════════════╗
    ║                                                               ║
    ║   🔧 Go-Skeleton v1.0.7                                      ║
    ║   🏗️  Modern Go REST API Generator                            ║
    ║   📅 Built with ❤️  using Go 1.24                            ║
    ║                                                               ║
    ║   🚀 Features:                                                ║
    ║   • Clean Architecture                                        ║
    ║   • UberFX Dependency Injection                              ║
    ║   • Structured Logging                                        ║
    ║   • JWT Authentication                                        ║
    ║   • PostgreSQL with GORM                                      ║
    ║                                                               ║
    ╚═══════════════════════════════════════════════════════════════╝
`
		fmt.Print(versionBanner)
	},
}

func printSuccessMessage(projectName string) {
	successBanner := fmt.Sprintf(`
    ╔══════════════════════════════════════════════════════════════════╗
    ║                                                                  ║
    ║                 ✅ PROJECT CREATED SUCCESSFULLY! ✅               ║
    ║                                                                  ║
    ║  📁 Project: %-51s ║
    ║                                                                  ║
    ║  🎯 Next steps:                                                  ║
    ║     1. cd %s                                          ║
    ║     2. go mod tidy                                               ║
    ║     3. cp .env.example .env                                      ║
    ║     4. go run cmd/app/main.go                                    ║
    ║                                                                  ║
    ║  🚀 Happy coding with Go + UberFX!                              ║
    ║                                                                  ║
    ╚══════════════════════════════════════════════════════════════════╝
`, projectName, projectName)

	fmt.Print(successBanner)
}

func validateProjectName(name string) error {
	if strings.ContainsAny(name, " /\\") {
		return fmt.Errorf("project name cannot contain spaces or path separators")
	}
	return nil
}

func generateProject(name string) error {
	data := TemplateData{
		Name: name,
	}

	// Create project directory first
	if err := os.MkdirAll(name, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, "template/")
		if relPath == "" {
			return nil
		}

		if d.IsDir() {
			outputPath := filepath.Join(name, relPath)
			return os.MkdirAll(outputPath, os.ModePerm)
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Create output path
		outputPath := filepath.Join(name, relPath)
		if strings.HasSuffix(relPath, ".tmpl") {
			relPath = strings.TrimSuffix(relPath, ".tmpl")
			outputPath = filepath.Join(name, relPath)
		}

		// Create output directory if it doesn't exist
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
		}

		// Parse and execute template
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Clean up template directory
	templateDirPath := filepath.Join(name, "template")
	return os.RemoveAll(templateDirPath)
}

func init() {
	newCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	newCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
