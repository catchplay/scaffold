package scaffold

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	pkgErr "github.com/pkg/errors"
)

const (
	GoScaffoldPath = "src/github.com/catchplay/scaffold"
)

func init() {
	Gopath = os.Getenv("GOPATH")
	if Gopath == "" {
		panic("cannot find $GOPATH environment variable")
	}
}

var Gopath string

type scaffold struct {
}

func New() *scaffold {
	return &scaffold{}
}

func (*scaffold) Generate() error {
	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	projectName := filepath.Base(currDir)
	goProjectPath := strings.TrimPrefix(currDir, filepath.Join(Gopath, "src")+string(os.PathSeparator))

	d := data{
		ProjectPath: goProjectPath,
		ProjectName: projectName,
		Quit:        "<-quit",
	}

	if err := genFromTemplate(getTemplateSets(), d); err != nil {
		return err
	}

	if err := genFormStaticFle(d); err != nil {
		return err
	}
	return nil
}

type data struct {
	ProjectPath string //The relative path of the generated file (eg:github.com/catchplay/foo)
	ProjectName string //The current directory name of the generated file //TODO:seems like don't needs
	Quit        string
}

type templateEngine struct {
	Templates []templateSet
	currDir   string
}

type templateSet struct {
	templateFilePath string
	templateFileName string
	genFilePath      string
}

func getTemplateSets() []templateSet {
	tt := templateEngine{}
	filepath.Walk(filepath.Join(Gopath, GoScaffoldPath), tt.visit)
	return tt.Templates
}

func genFromTemplate(templateSets []templateSet, data interface{}) error {
	for _, tmpl := range templateSets {
		if err := tmplExec(tmpl, data); err != nil {
			return err
		}
	}
	return nil
}

func unescaped(x string) interface{} { return template.HTML(x) }

func tmplExec(tmplSet templateSet, data interface{}) error {
	tmpl := template.New(tmplSet.templateFileName)
	tmpl = tmpl.Funcs(template.FuncMap{"unescaped": unescaped})
	tmpl, err := tmpl.ParseFiles(tmplSet.templateFilePath)
	if err != nil {
		pkgErr.WithStack(err)
	}

	path := filepath.Dir(tmplSet.genFilePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	if err != nil {
		pkgErr.WithStack(err)
	}

	distFilePath := filepath.Join(path, filepath.Base(tmplSet.genFilePath))
	dist, err := os.Create(distFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}
	defer dist.Close()

	log.Printf("Genarate go file to %s\n", distFilePath)
	return tmpl.Execute(dist, data)
}

func (templEngine *templateEngine) visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if ext := filepath.Ext(path); ext == ".tmpl" {
		templateFileName := filepath.Base(path)

		genFileBaeName := strings.TrimSuffix(templateFileName, ".tmpl") + ".go"
		genFileBasePath, err := filepath.Rel(filepath.Join(Gopath, GoScaffoldPath, "template"), filepath.Join(filepath.Dir(path), genFileBaeName))
		if err != nil {
			return pkgErr.WithStack(err)
		}

		templ := templateSet{
			templateFilePath: path,
			templateFileName: templateFileName,
			genFilePath:      filepath.Join(templEngine.currDir, genFileBasePath),
		}

		templEngine.Templates = append(templEngine.Templates, templ)

	}

	return nil
}

func genFormStaticFle(d data) error {
	walkerFuc := func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() == true {
			src, err := os.Open(path)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer src.Close()

			sf, err := filepath.Rel(filepath.Join(Gopath, GoScaffoldPath, "static"), path)
			if err != nil {
				return pkgErr.WithStack(err)
			}

			distFilePath := filepath.Join(Gopath, "src", d.ProjectPath, sf)
			if err := os.MkdirAll(filepath.Dir(distFilePath), os.ModePerm); err != nil {
				return pkgErr.WithStack(err)
			}

			dist, err := os.Create(distFilePath)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer dist.Close()

			if _, err := io.Copy(dist, src); err != nil {
				return pkgErr.WithStack(err)
			}

			log.Printf("Genarate static file to %s \n", distFilePath)
		}

		return nil
	}

	return filepath.Walk(filepath.Join(Gopath, GoScaffoldPath, "static"), walkerFuc)
}
