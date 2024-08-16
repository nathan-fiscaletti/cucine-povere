package datatypes

type Tag struct {
	Name string
	Url  string
}

type TagPage struct {
	Tag   Tag
	Site  Site
	Posts []Post
}
