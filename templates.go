package main

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
	tpl.ParseGlob("templates/pages/*.gohtml")
	tpl.ParseGlob("templates/layouts/*.gohtml")

	log.Println("initialized templates")
}