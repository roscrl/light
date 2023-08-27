package views

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/roscrl/light/config"
)

//go:embed assets/dist
var FrontendDistFS embed.FS

//go:embed templates/*
var tmplFS embed.FS

const (
	DirTemplates      = "templates"
	DirTemplatesSlash = DirTemplates + "/"

	PathTemplates = "core/views/" + DirTemplates
	PathViews     = "core/views"
)

type Views struct {
	env config.Environment

	templates *template.Template
	funcMap   template.FuncMap
}

func New(env config.Environment) *Views {
	funcMap := TemplateFuncs
	views := &Views{env: env, funcMap: funcMap}

	if env == config.LOCAL {
		templates := findAndParseTemplates(os.DirFS(PathTemplates), funcMap)

		views.templates = templates
		watchLocalTemplates(views)
	} else {
		tmplFS, err := fs.Sub(tmplFS, DirTemplates)
		if err != nil {
			log.Fatalf("failed to get subdirectory %s: %v", DirTemplates, err)
		}

		templates := findAndParseTemplates(tmplFS, funcMap)

		views.templates = templates
	}

	log.Println(views.templates.DefinedTemplates())

	return views
}

func (v *Views) RenderPage(w io.Writer, name string, data any) {
	tmpl := template.Must(v.templates.Clone())

	if v.env == config.LOCAL {
		tmpl = template.Must(tmpl.ParseGlob(PathTemplates + "/" + name))
	} else {
		tmpl = template.Must(tmpl.ParseFS(tmplFS, DirTemplatesSlash+name))
	}

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") {
			return
		}

		log.Printf("failed to render template %s: %v, defined templates %v", name, err, tmpl.DefinedTemplates())
		_ = tmpl.ExecuteTemplate(w, Error, err)
	}
}

func (v *Views) RenderDefaultErrorPage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	v.RenderPage(w, Error, map[string]any{})
}

func (v *Views) RenderErrorPage(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	v.RenderPage(w, Error, map[string]any{"error": msg})
}

func (v *Views) RenderTurboStream(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", TurboStreamMIME)
	v.RenderPage(w, name, data)
}

func findAndParseTemplates(filesys fs.FS, funcMap template.FuncMap) *template.Template {
	rootTemplate := template.New("")

	err := fs.WalkDir(filesys, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			strippedTemplateDirTemplatePath := strings.TrimPrefix(path, DirTemplatesSlash)

			templateContent, err := fs.ReadFile(filesys, strippedTemplateDirTemplatePath)
			if err != nil {
				return fmt.Errorf("reading template %s: %w", strippedTemplateDirTemplatePath, err)
			}

			tmpl := rootTemplate.New(strippedTemplateDirTemplatePath).Funcs(funcMap)
			_, err = tmpl.Parse(string(templateContent))
			if err != nil {
				return fmt.Errorf("parsing template %s: %w", strippedTemplateDirTemplatePath, err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return rootTemplate
}
