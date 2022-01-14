package models

import "time"

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int32  `json:"threads,omitempty"`
}

type Thread struct {
	Id      int64     `json:"id,omitempty"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message"`
	Votes   int32     `json:"votes,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

type Post struct {
	Id       int64     `json:"id,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int32     `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Path     int64     `json:"-"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}

type ForumRep interface {
	CreateForum(newForum Forum) (Forum, error)
	GetForumBySlug(slug string) (Forum, error)
	CreateThread(newThread Thread) (Thread, error)
	GetThreads(slug, limit, since, desc string) ([]Thread, error)
	CreatePosts(thread Thread, newPost []Post) ([]Post, error)
	GetThreadBySlugOrId(slugOrId string) (Thread, error)
	Vote(thread Thread, vote Vote) (Thread, error)
	GetPosts(thread Thread, limit, since, sort, desc string) ([]Post, error)
}
