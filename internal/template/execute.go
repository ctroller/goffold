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

	err := createOutputDir()
	if err != nil {
		return err
	}

	slog.Info("Prompting variables", "template", t.Name)
	tplVars, err := promptVars(t)
	if err != nil {
		return err
	}

	t.TemplateVars = tplVars

	slog.Info("Checking if template extends any folders", "template", t.Name)
	err = importExtends(t, OutputDir)
	if err != nil {
		return err
	}

	slog.Info("Copying template structure", "template", t.Name)
	err = copyStructure(t.Path, "layout", OutputDir, t)
	if err != nil {
		return err
	}

	return nil
}

func createOutputDir() error {
	err := ensureDirDoesNotExist(OutputDir)
	if err != nil {
		return err
	}

	// Create the directory if it does not exist
	err = OsFs.MkdirAll(OutputDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func copyStructure(basePath, templatePath, outputDir string, t *Template) error {
	infos, err := afero.ReadDir(TemplateFs, filepath.Join(basePath, templatePath))
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			err = copyDir(basePath, filepath.Join(templatePath, info.Name()), outputDir, t)
		} else {
			err = parseFile(basePath, filepath.Join(templatePath, info.Name()), outputDir, t)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyDir(basePath, templateDir, targetDir string, t *Template) error {
	outputDir := filepath.Join(basePath, targetDir)
	err := OsFs.Mkdir(outputDir, 0755)
	if err != nil {
		return err
	}

	return copyStructure(basePath, templateDir, outputDir, t)
}

func ensureDirDoesNotExist(output string) error {
	if output == "." {
		return nil
	}

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

func parseFile(basePath, templateFile, targetDir string, t *Template) error {
	content, fileName, err := parseFileContent(basePath, templateFile, t)

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

func parseFileContent(basePath, templateFile string, t *Template) ([]byte, string, error) {
	fileName := filepath.Base(templateFile)
	var content []byte
	var err error
	if strings.HasSuffix(templateFile, goTemplateSuffix) {
		fileName = strings.TrimSuffix(filepath.Base(templateFile), ".tmpl")
		content, err = parseGoTemplate(basePath, templateFile, t)
	} else {
		content, err = readFileContent(templateFile)
	}

	return content, fileName, err
}

func parseGoTemplate(basePath, templateFile string, t *Template) ([]byte, error) {
	slog.Info("Parsing template", "file", templateFile)
	tmplData, err := afero.ReadFile(TemplateFs, filepath.Join(basePath, templateFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	tpl := template.Must(template.New("").Parse(string(tmplData)))
	out := new(bytes.Buffer)
	if err := tpl.Execute(out, t); err != nil {
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
			err := copyDir("", extendsPath, target, t)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func promptVars(t *Template) (TemplateVars, error) {
	tplVars := TemplateVars{}

	for _, v := range t.Vars {
		value, err := v.GetValue(t)
		if err != nil {
			return nil, err
		}

		tplVars[v.Name] = value
	}

	return tplVars, nil
}
