/*
 * @Author: catcatproer
 * @Date: 2023-01-30 20:11:17
 * @LastEditors: catcatproer
 * @LastEditTime: 2023-01-31 22:57:04
 * @FilePath: \go_web\src\chit_chat\route_main.go
 * @Description:
 *
 * Copyright (c) 2023 by catcatproer, All Rights Reserved.
 */
package main

import (
	"chit_chat/data"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func index(w http.ResponseWriter, r *http.Request) {
	// files := []string{
	// 	"templates/layout.html",
	// 	"templates/navbar.html",
	// 	"templates/index.html",
	// }
	// template := template.Must(template.ParseFiles(files...))
	// if threads, err := data.Threads(); err == nil {
	// 	_, err := session(w, r)
	// 	template.ExecuteTemplate(w, "layout", threads)
	// }
	threads, err := data.Threads()
	if err != nil {
		error_message(w, r, "Cannot get threads")

		// templates.ExecuteTemplate(w, "layout", threads)
	} else {
		_, err := session(w, r)

		if err != nil {
			generateHTML(w, threads, "layout", "public.navbar", "index")
		} else {
			generateHTML(w, threads, "layout", "private.navbar", "index")
		}
	}
}

func err(writer http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	_, err := session(writer, r)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")

	}
}

func generateHTML(w http.ResponseWriter, data interface{}, fn ...string) {
	var files []string
	for _, file := range fn {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Fatal(err)
	}

}
