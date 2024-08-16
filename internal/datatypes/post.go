package datatypes

import "html/template"

type Post struct {
	Image   string
	Title   string
	Url     string
	Tags    []Tag
	Date    string
	Content template.HTML
	Preview string
	Author  Author
}

type PostPage struct {
	Site     Site
	Post     Post
	Next     *Post
	Previous *Post
}
