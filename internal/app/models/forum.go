package models

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int32  `json:"threads,omitempty"`
}

type ForumRep interface {
	Create(newForum Forum) (Forum, error)
	GetForumBySlug(slug string) (Forum, error)
}
