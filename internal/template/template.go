package template

import (
	"html/template"
	"log"
)

// WebRoot is the path to the web templates + static files
var WebRoot string

// CacheTemplates is the setting on whether to cache the template files or read from disk each time
var CacheTemplates bool

// template cache
var templates = make(map[string]*template.Template)

// ReadTemplate is used by the various handlers to read the template file off disk, or return
// the template from cache if we already did that.  -cache_templates=false can be passed on the
// command line to always read off disk, useful for developing
func ReadTemplate(filename string) *template.Template {
	if CacheTemplates {
		if tmpl, ok := templates[filename]; ok {
			return tmpl
		}
	}

	// web templates always have the base.tmpl that provides the overall layout, and then the requested template
	// provides all the content
	t, err := template.New(filename).ParseFiles(WebRoot+"/templates/base.tmpl", WebRoot+"/templates/"+filename)
	if err != nil {
		log.Fatal("Could not open template: " + WebRoot + "/templates/" + filename + " : " + err.Error())
	}

	// drop the template in cache for later
	templates[filename] = t

	return t
}
