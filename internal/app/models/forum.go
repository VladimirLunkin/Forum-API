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

type ForumRep interface {
	CreateForum(newForum Forum) (Forum, error)
	GetForumBySlug(slug string) (Forum, error)
	CreateThread(newThread Thread) (Thread, error)
	GetThreads(slug, limit, since, desc string) ([]Thread, error)
}
