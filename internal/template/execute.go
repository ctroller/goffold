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

	err = importExtends(t, out)
	if err != nil {
		return err
	}

	err = copyStructure(t.Path, "layout", out)
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

func copyStructure(basePath, templatePath, outputDir string) error {
	infos, err := afero.ReadDir(TemplateFs, filepath.Join(basePath, templatePath))
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			err = copyDir(basePath, filepath.Join(templatePath, info.Name()), outputDir)
		} else {
			err = parseFile(basePath, filepath.Join(templatePath, info.Name()), outputDir)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyDir(basePath, templateDir, targetDir string) error {
	outputDir := filepath.Join(basePath, targetDir)
	err := OsFs.Mkdir(outputDir, 0755)
	if err != nil {
		return err
	}

	return copyStructure(basePath, templateDir, outputDir)
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

func parseFile(basePath, templateFile, targetDir string) error {
	content, fileName, err := parseFileContent(basePath, templateFile)

	if err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, fileName)
	slog.Info("Writing file", "target", targetFile, "origFile", templateFile, "targetPath", targetDir, "base", basePath)
	handle, err := OsFs.Create(targetFile)
	if err != nil {
		return err
	}
	defer handle.Close()

	err = afero.WriteFile(OsFs, targetFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func parseFileContent(basePath, templateFile string) ([]byte, string, error) {
	fileName := filepath.Base(templateFile)
	var content []byte
	var err error
	if strings.HasSuffix(templateFile, goTemplateSuffix) {
		fileName = strings.TrimSuffix(filepath.Base(templateFile), ".tmpl")
		content, err = parseGoTemplate(basePath, templateFile)
	} else {
		content, err = readFileContent(templateFile)
	}

	return content, fileName, err
}

func parseGoTemplate(basePath, templateFile string) ([]byte, error) {
	slog.Info("Parsing template", "file", templateFile)
	tmplData, err := afero.ReadFile(TemplateFs, filepath.Join(basePath, templateFile))
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

func importExtends(t *Template, targetDir string) error {
	for _, f := range t.Layout.Folders {
		if f.ExtendsFrom != "" {
			extendsPath := filepath.Join(TemplateDir, f.ExtendsFrom)
			target := filepath.Join(targetDir, f.Name)

			slog.Info("Extending folder", "from", extendsPath, "to", target)
			err := copyDir("", extendsPath, target)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
