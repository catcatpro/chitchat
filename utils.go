/*
 * @Author: catcatproer
 * @Date: 2023-01-30 20:11:17
 * @LastEditors: catcatproer
 * @LastEditTime: 2023-01-31 22:02:43
 * @FilePath: \go_web\src\chit_chat\utils.go
 * @Description:
 *
 * Copyright (c) 2023 by catcatproer, All Rights Reserved.
 */
package main

import (
	"chit_chat/data"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var logger *log.Logger

func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}

	}
	return
}

// parse HTML templates
// pass in a list of files names, and  get a template.
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	tp := template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	t = template.Must(tp.ParseFiles(files...))
	return
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

// Convenience function to redirect to the error message page
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING")
	logger.Println(args...)
}
