package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct { //Pass data to the template...
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

func (h *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData { //To make this work;added a Session w/manager above and then also Session the the createRender in hiSpeed.go
	td.Secure = h.Secure
	td.ServerName = h.ServerName
	td.Port = h.Port
	if h.Session.Exists(r.Context(), "userID") { //If userID exist then by default; user must be authenticated
		td.IsAuthenticated = true
	}
	return td
}

// A func to render a page
func (h *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(h.Renderer) { //Only going to use (2) two types of templates for now: standard go template and and/or a jet template...
	case "go":
		return h.GoPage(w, r, view, data)
	case "jet":
		return h.JetPage(w, r, view, variables, data)
	default: // removed a call to test err != nil and instead just returned an error if no rendering engine specified (also allows testing).

	}
	return errors.New("no rendering engine specified")
}

// GoPage renders a standard Go template...
func (h *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", h.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet template engine...
func (h *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	td = h.defaultData(td, r)

	t, err := h.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
