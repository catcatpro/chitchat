/*
 * @Author: catcatproer
 * @Date: 2023-01-30 20:11:17
 * @LastEditors: catcatproer
 * @LastEditTime: 2023-02-01 21:28:26
 * @FilePath: \go_web\src\chit_chat\route_auth.go
 * @Description:
 *
 * Copyright (c) 2023 by catcatproer, All Rights Reserved.
 */
package main

import (
	"chit_chat/data"
	"log"
	"net/http"
)

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	// t := parseTemplateFiles("login.layout", "public.navbar", "login")
	// err := t.Execute(writer, nil)
	// fmt.Println(err)
	generateHTML(writer, nil, "login.layout", "public.navbar", "login")

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")

	if err != http.ErrNoCookie {
		// warning(err, "Falled to get Cookie")
		session := data.Session{Uuid: cookie.Value}

		session.DeleteByUUID()
	}

	http.Redirect(writer, request, "/", 302)
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup_account
// Create this user account
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}

	user := data.User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}

	http.Redirect(writer, request, "/login", 302)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Fatal(err)
	}
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, _ := user.CreateSession()
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)

	}
}
