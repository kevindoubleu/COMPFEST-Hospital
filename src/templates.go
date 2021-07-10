package src

import (
	"log"
	"text/template"
)

var tpl *template.Template

func init() {
	log.Println("initializing templates")

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	tpl = template.New("").Funcs(funcMap)
	tpl.ParseGlob("templates/components/*.gohtml")
	tpl.ParseGlob("templates/sections/*.gohtml")
	tpl.ParseGlob("templates/sections/index/*.gohtml")
	tpl.ParseGlob("templates/sections/appointments/*.gohtml")
	tpl.ParseGlob("templates/sections/admin-appointments/*.gohtml")
	tpl.ParseGlob("templates/sections/admin-patients/*.gohtml")
	tpl.ParseGlob("templates/pages/*.gohtml")

	log.Println("initialized templates")
}