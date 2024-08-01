package template

import (
	"bytes"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ctroller/goffold/internal/dependencies"
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
	err = importExtends(t)
	if err != nil {
		return err
	}

	slog.Info("Copying template structure", "template", t.Name)
	err = copyStructure(filepath.Join(t.Path, "layout"), OutputDir, t)
	if err != nil {
		return err
	}

	resolver := dependencies.GetResolver("go")
	if resolver == nil {
		return fmt.Errorf("failed to get resolver for go dependencies")
	}

	err = installDependencies(t, resolver)
	if err != nil {
		return err
	}

	return resolver.Finisher(OutputDir)
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

func copyStructure(templatePath, outputDir string, t *Template) error {
	infos, err := afero.ReadDir(TemplateFs, templatePath)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			err = copyDir(filepath.Join(templatePath, info.Name()), filepath.Join(outputDir, info.Name()), t)
		} else {
			err = parseFile(filepath.Join(templatePath, info.Name()), outputDir, t)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func copyDir(templateDir, targetDir string, t *Template) error {
	slog.Info("Copying directory", "from", templateDir, "to", targetDir)
	err := OsFs.Mkdir(targetDir, 0755)
	if err != nil {
		return err
	}

	return copyStructure(templateDir, targetDir, t)
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
	goTemplateSuffix = ".gotmpl"
)

func parseFile(templateFile, targetDir string, t *Template) error {
	content, fileName, err := parseFileContent(templateFile, t)

	if err != nil {
		return err
	}

	targetFile := filepath.Join(targetDir, fileName)
	slog.Info("Writing file", "target", targetFile, "origFile", templateFile, "targetPath", targetDir)
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

func parseFileContent(templateFile string, t *Template) ([]byte, string, error) {
	fileName := filepath.Base(templateFile)
	var content []byte
	var err error
	if strings.HasSuffix(templateFile, goTemplateSuffix) {
		fileName = strings.TrimSuffix(filepath.Base(templateFile), goTemplateSuffix)
		content, err = parseGoTemplate(templateFile, t)
	} else {
		content, err = readFileContent(templateFile)
	}

	return content, fileName, err
}

func getTemplateFuncs(t *Template) template.FuncMap {
	return template.FuncMap{
		"output_path": func() string {
			if OutputDir == "." {
				return t.Name
			}

			return OutputDir
		},
	}
}

func parseGoTemplate(templateFile string, t *Template) ([]byte, error) {
	slog.Info("Parsing template", "file", templateFile)

	tmplData, err := afero.ReadFile(TemplateFs, templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	tpl := template.Must(template.New("").Funcs(getTemplateFuncs(t)).Parse(string(tmplData)))
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

func importExtends(t *Template) error {
	for _, f := range t.Layout.Folders {
		if f.ExtendsFrom != "" {
			extendsPath := filepath.Join(TemplateDir, f.ExtendsFrom)
			target := filepath.Join(OutputDir, f.Name)

			slog.Info("Extending folder", "from", extendsPath, "to", target)
			err := copyDir(extendsPath, target, t)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func installDependencies(t *Template, resolver *dependencies.DependencyResolver) error {
	for _, dep := range t.Dependencies {
		_, err := resolver.Resolve(OutputDir, dep)
		if err != nil {
			return err
		}
	}

	return nil
}
