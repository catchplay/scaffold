package scaffold

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	GoScaffoldPath = "src/github.com/catchplay/go-scaffold"
)

var Gopath string

type scaffold struct {
}

func New() *scaffold {
	return &scaffold{}
}

func init() {
	Gopath = os.Getenv("GOPATH")
	if Gopath == "" {
		panic("cannot find $GOPATH environment variable")
	}
}
func (*scaffold) Generate() error {
	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	projectName := filepath.Base(currDir)
	goProjectPath := strings.TrimPrefix(currDir, filepath.Join(Gopath, "src")+string(os.PathSeparator))

	d := data{
		ProjectPath: goProjectPath,
		ProjectName: projectName,
		Quit:        "<-quit",
	}

	if err := genFromTemplate(getTemplEngines(currDir), d); err != nil {
		return err
	}

	if err := genFormStaticFle(); err != nil {
		return err
	}
	return nil
}

type data struct {
	ProjectPath string
	ProjectName string
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

func getTemplEngines(exPath string) []templateSet {
	tt := templateEngine{}
	filepath.Walk(filepath.Join(Gopath, GoScaffoldPath), tt.visit)
	return tt.Templates
}

func genFromTemplate(templateSets []templateSet, data interface{}) error {
	for _, tmpl := range templateSets {
		if err := loadTemplate(tmpl, data); err != nil {
			return err
		}
	}
	return nil
}

func unescaped(x string) interface{} { return template.HTML(x) }

func loadTemplate(tmplSet templateSet, data interface{}) error {
	tmpl := template.New(tmplSet.templateFileName)
	tmpl = tmpl.Funcs(template.FuncMap{"unescaped": unescaped})
	tmpl, err := tmpl.ParseFiles(tmplSet.templateFilePath)
	if err != nil {
		return err
	}

	path := filepath.Dir(tmplSet.genFilePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	f, err := os.Create(filepath.Join(path, filepath.Base(tmplSet.genFilePath)))
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Genarate file to %s \n", tmplSet.genFilePath)
	return tmpl.Execute(f, data)
}

func (templEngine *templateEngine) visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		log.Panicf("ParseFile panic:%+v", err)
	}

	// templ := templateSet{
	// 	templateFilePath: filepath.Join(Gopath, GoScaffoldPath, "/template/config/config.tmpl"),
	// 	templateFileName: "config.tmpl",
	// 	genFilePath:      filepath.Join(exPath, "config/config.go"),
	// }}

	if ext := filepath.Ext(path); ext == ".tmpl" {

		templateFileName := filepath.Base(path)

		genFileBaeName := strings.TrimSuffix(templateFileName, ".tmpl") + ".go"
		// fmt.Printf("templateFileName:%s\n", templateFileName)
		// fmt.Printf("genFileBaeName:%s\n", genFileBaeName)
		genFileBasePath, err := filepath.Rel(filepath.Join(Gopath, GoScaffoldPath, "template"), filepath.Join(filepath.Dir(path), genFileBaeName))
		if err != nil {
			return err
		}

		templ := templateSet{
			templateFilePath: path,
			templateFileName: templateFileName,
			genFilePath:      filepath.Join(templEngine.currDir, genFileBasePath),
		}

		//fmt.Println(path)
		templEngine.Templates = append(templEngine.Templates, templ)

	}

	return nil
}

func genFormStaticFle() error {
	//TODO:
	return nil
}
