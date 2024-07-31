package template

import (
	"bytes"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

func (t *Template) Execute() error {
	slog.Info("Executing template", "template", t.Name)

	out, err := createOutputDir(t)
	if err != nil {
		return err
	}

	err = copyStructure(t, "layout", out)
	if err != nil {
		return err
	}

	return nil
}

func createOutputDir(t *Template) (string, error) {
	output := OutputDir
	if output == "" {
		path := strings.Split(t.Path, "/")
		output = path[len(path)-1]
	}

	err := ensureDirDoesNotExist(output)
	if err != nil {
		return "", fmt.Errorf("failed to ensure directory does not exist: %w", err)
	}

	// Create the directory if it does not exist
	err = OsFs.MkdirAll(output, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return output, nil
}

func copyStructure(t *Template, templatePath, outputDir string) error {
	infos, err := afero.ReadDir(TemplateFs, filepath.Join(t.Path, templatePath))
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			outputDir := filepath.Join(outputDir, info.Name())
			newTemplatePath := filepath.Join(templatePath, info.Name())
			err = OsFs.Mkdir(outputDir, 0755)
			if err != nil {
				return err
			}

			err = copyStructure(t, newTemplatePath, outputDir)
			if err != nil {
				return err
			}
		} else {
			err = parseFile(t, filepath.Join(templatePath, info.Name()), outputDir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ensureDirDoesNotExist(output string) error {
	exists, err := afero.DirExists(OsFs, output)
	if err != nil {
		return fmt.Errorf("failed to check if directory exists: %w", err)
	}
	if exists {
		return fmt.Errorf("output path '%s' already exists", output)
	}
	return nil
}

const (
	goTemplateSuffix = ".go.tmpl"
)

func parseFile(t *Template, templateFile, targetDir string) error {
	fileName := filepath.Base(templateFile)
	var content []byte
	var err error

	if strings.HasSuffix(templateFile, goTemplateSuffix) {
		fileName = strings.TrimSuffix(fileName, ".tmpl")
		content, err = parseGoTemplate(t, templateFile)
	} else {
		content, err = readFileContent(templateFile)
	}

	if err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, fileName)
	slog.Info("Writing file", "target", targetFile, "origFile", templateFile, "targetPath", targetDir)
	err = afero.WriteFile(OsFs, targetFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseGoTemplate(t *Template, templateFile string) ([]byte, error) {
	slog.Info("Parsing template", "file", templateFile)
	tmplData, err := afero.ReadFile(TemplateFs, filepath.Join(t.Path, templateFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	tpl := template.Must(template.New("").Parse(string(tmplData)))
	out := new(bytes.Buffer)
	if err := tpl.Execute(out, nil); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return out.Bytes(), nil
}

func readFileContent(templateFile string) ([]byte, error) {
	content, err := afero.ReadFile(TemplateFs, templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return content, nil
}
