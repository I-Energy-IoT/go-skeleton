package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed template/*
var templateFS embed.FS

type TemplateData struct {
	ModuleName  string
	ProjectName string
}

var projectName string

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new REST API project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			log.Fatal("Please provide a project name using --name flag.")
		}

		if err := generateProject(projectName); err != nil {
			log.Fatalf("failed to generate project: %v", err)
		}
		fmt.Println("Project created:", projectName)
	},
}

func init() {
	newCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project (required)")
	newCmd.MarkFlagRequired("name")
}

func Execute() {
	rootCmd := &cobra.Command{Use: "ienergy-template"}
	rootCmd.AddCommand(newCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func generateProject(name string) error {
	data := TemplateData{
		ModuleName:  name,
		ProjectName: name,
	}

	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, "template/")
		if relPath == "" {
			return nil
		}

		outputPath := filepath.Join(name, relPath)

		if d.IsDir() {
			return os.MkdirAll(outputPath, os.ModePerm)
		}

		if strings.HasSuffix(relPath, ".tmpl") {
			relPath = strings.TrimSuffix(relPath, ".tmpl")
			outputPath = filepath.Join(name, relPath)
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := strings.ReplaceAll(string(content), "{{.name}}", name)

		tpl, err := template.New(path).Parse(contentStr)
		if err != nil {
			return err
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()

		return tpl.Execute(f, data)
	})

	if err != nil {
		return err
	}

	templateDirPath := filepath.Join(name, "template")
	return os.RemoveAll(templateDirPath)
}
