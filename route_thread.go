/*
 * @Author: catcatproer
 * @Date: 2023-01-30 20:11:17
 * @LastEditors: catcatproer
 * @LastEditTime: 2023-02-01 20:44:27
 * @FilePath: \go_web\src\chit_chat\route_thread.go
 * @Description:
 *
 * Copyright (c) 2023 by catcatproer, All Rights Reserved.
 */
package main

import (
	"chit_chat/data"
	"fmt"
	"net/http"
)

// GET /threads/new
// show the new thread form page
func newThread(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "new.thread")
	}
}

// POST /sigup
// Create new thread
func createThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {

		err = request.ParseForm()

		if err != nil {
			danger(err, "Cannot  parse form")
		}
		user, err := sess.User()

		if err != nil {
			danger(err, "Cannot get user from session")
		}

		topic := request.PostFormValue("topic")
		if _, err := user.CreateThread(topic); err != nil {
			danger(err, "Cannot create thread")
		}

		http.Redirect(writer, request, "/", 302)
	}
}

// GET /thread/read
// Show the details of the thread, including  the posts and the form to write a post
func readThread(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByUUID(uuid)

	if err != nil {
		error_message(writer, request, "Gannot read thread")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, &thread, "layout", "public.navbar", "public.thread")
		} else {
			generateHTML(writer, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}

// POST /thread/post
// Create the post
func postThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}

		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")
		thread, err := data.ThreadByUUID(uuid)

		if err != nil {
			error_message(writer, request, "Cannot reaxd thread")
		}

		if _, err := user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot create post")
		}

		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(writer, request, url, 302)
	}
}
