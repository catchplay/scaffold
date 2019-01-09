package scaffold

import (
	"fmt"
	"html/template"
	"io"
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

func (*scaffold) Generate(path string) error {
	genAbsDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	projectName := filepath.Base(genAbsDir)
	//TODO: have to check path MUST be under the $GOPATH/src folder
	goProjectPath := strings.TrimPrefix(genAbsDir, filepath.Join(Gopath, "src")+string(os.PathSeparator))

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
	templatesFolder := filepath.Join(Gopath, GoScaffoldPath, "template/")
	//fmt.Printf("walk:%s\n", templatesFolder)
	filepath.Walk(templatesFolder, tt.visit)
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

	relateDir := filepath.Dir(tmplSet.genFilePath)
	if _, err := os.Stat(relateDir); os.IsNotExist(err) {
		os.MkdirAll(relateDir, os.ModePerm)
	}
	if err != nil {
		pkgErr.WithStack(err)
	}

	//fmt.Printf("relateDir:%s\n", relateDir)
	distFilePath := filepath.Join(relateDir, filepath.Base(tmplSet.genFilePath))
	//fmt.Printf("distFilePath:%s\n", distFilePath)
	dist, err := os.Create(distFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}
	defer dist.Close()

	fmt.Printf("Create %s\n", distFilePath)
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

	} else if mode := f.Mode(); mode.IsRegular() {
		//TODO:
		templateFileName := filepath.Base(path)
		//fmt.Printf("templateFileName:%s\n", templateFileName)

		basepath := filepath.Join(Gopath, GoScaffoldPath, "template")
		targpath := filepath.Join(filepath.Dir(path), templateFileName)

		//fmt.Printf("basepath:%s\n", basepath)
		//fmt.Printf("targpath:%s\n", targpath)

		genFileBasePath, err := filepath.Rel(basepath, targpath)
		if err != nil {
			return pkgErr.WithStack(err)
		}

		//fmt.Printf("genFileBasePath:%s\n", genFileBasePath)

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

			distFilePath, err := filepath.Rel(filepath.Join(Gopath, GoScaffoldPath, "static"), path)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			//fmt.Printf("sf:%s\n", sf)
			//	distFilePath := filepath.Join(Gopath, "src", d.ProjectPath, sf)
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

			fmt.Printf("Create %s \n", distFilePath)
		}

		return nil
	}

	return filepath.Walk(filepath.Join(Gopath, GoScaffoldPath, "static"), walkerFuc)
}
