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
	debug bool
}

func New(debug bool) *scaffold {
	return &scaffold{debug: debug}
}

func (s *scaffold) Generate(path string) error {
	genAbsDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	projectName := filepath.Base(genAbsDir)
	//TODO: have to check path MUST be under the $GOPATH/src folder
	goProjectPath := strings.TrimPrefix(genAbsDir, filepath.Join(Gopath, "src")+string(os.PathSeparator))

	d := data{
		AbsGenProjectPath: genAbsDir,
		ProjectPath:       goProjectPath,
		ProjectName:       projectName,
		Quit:              "<-quit",
	}

	if err := s.genFromTemplate(getTemplateSets(), d); err != nil {
		return err
	}

	if err := s.genFormStaticFle(d); err != nil {
		return err
	}
	return nil
}

type data struct {
	AbsGenProjectPath string // The Abs Gen Project Path
	ProjectPath       string //The Go import project path (eg:github.com/fooOrg/foo)
	ProjectName       string //The project name which want to generated
	Quit              string
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

func (s *scaffold) genFromTemplate(templateSets []templateSet, d data) error {
	for _, tmpl := range templateSets {
		if err := s.tmplExec(tmpl, d); err != nil {
			return err
		}
	}
	return nil
}

func unescaped(x string) interface{} { return template.HTML(x) }

func (s *scaffold) tmplExec(tmplSet templateSet, d data) error {
	tmpl := template.New(tmplSet.templateFileName)
	tmpl = tmpl.Funcs(template.FuncMap{"unescaped": unescaped})
	tmpl, err := tmpl.ParseFiles(tmplSet.templateFilePath)
	if err != nil {
		pkgErr.WithStack(err)
	}

	relateDir := filepath.Dir(tmplSet.genFilePath)

	distRelFilePath := filepath.Join(relateDir, filepath.Base(tmplSet.genFilePath))
	distAbsFilePath := filepath.Join(d.AbsGenProjectPath, distRelFilePath)

	s.debugPrintf("distRelFilePath:%s\n", distRelFilePath)
	s.debugPrintf("distAbsFilePath:%s\n", distAbsFilePath)

	if err := os.MkdirAll(filepath.Dir(distAbsFilePath), os.ModePerm); err != nil {
		return pkgErr.WithStack(err)
	}

	dist, err := os.Create(distAbsFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}
	defer dist.Close()

	fmt.Printf("Create %s\n", distRelFilePath)
	return tmpl.Execute(dist, d)
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
		templateFileName := filepath.Base(path)

		basepath := filepath.Join(Gopath, GoScaffoldPath, "template")
		targpath := filepath.Join(filepath.Dir(path), templateFileName)
		genFileBasePath, err := filepath.Rel(basepath, targpath)
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

func (s *scaffold) genFormStaticFle(d data) error {
	walkerFuc := func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() == true {
			src, err := os.Open(path)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer src.Close()

			basepath := filepath.Join(Gopath, GoScaffoldPath, "static")
			distRelFilePath, err := filepath.Rel(basepath, path)
			if err != nil {
				return pkgErr.WithStack(err)
			}

			distAbsFilePath := filepath.Join(d.AbsGenProjectPath, distRelFilePath)

			if err := os.MkdirAll(filepath.Dir(distAbsFilePath), os.ModePerm); err != nil {
				return pkgErr.WithStack(err)
			}

			dist, err := os.Create(distAbsFilePath)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer dist.Close()

			if _, err := io.Copy(dist, src); err != nil {
				return pkgErr.WithStack(err)
			}

			fmt.Printf("Create %s \n", distRelFilePath)
		}

		return nil
	}

	walkPath := filepath.Join(Gopath, GoScaffoldPath, "static")
	return filepath.Walk(walkPath, walkerFuc)
}

func (s *scaffold) debugPrintf(format string, a ...interface{}) {
	if s.debug == true {
		fmt.Printf(format, a...)
	}
}
