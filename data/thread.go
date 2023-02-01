/*
 * @Author: catcatproer
 * @Date: 2023-01-30 20:11:16
 * @LastEditors: catcatproer
 * @LastEditTime: 2023-02-01 21:08:17
 * @FilePath: \go_web\src\chit_chat\data\thread.go
 * @Description:
 *
 * Copyright (c) 2023 by catcatproer, All Rights Reserved.
 */
package data

import (
	"time"
)

type Thread struct {
	Id       int
	Uuid     string
	Topic    string
	UserId   int
	CreateAt time.Time
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

func Threads() (threads []Thread, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return
	}

	for rows.Next() {
		th := Thread{}
		if err = rows.Scan(&th.Id, &th.Uuid, &th.Topic, &th.UserId, &th.CreateAt); err != nil {
			return
		}

		threads = append(threads, th)
	}

	rows.Close()
	return
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1,$2,$3,$4,$5) returning id,uuid,body,user_id,thread_id,created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	//user QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}

func (thread *Thread) NumReplies() (count int) {
	rows, err := Db.Query("SELECT count(*) from posts where thread_id = $l", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}

	}
	rows.Close()
	return
}

// get posts to a thread
func (thread *Thread) Posts() (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}

	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Get a thread by the UUID
func ThreadByUUID(uuid string) (conv Thread, err error) {
	conv = Thread{}
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreateAt)
	return
}

// Get the user who staarted this thread
func (thread *Thread) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreateAt)
	return
}

// Get the user who wrote the post
func (post *Post) PostUser() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", post.UserId).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreateAt)
	return
}
