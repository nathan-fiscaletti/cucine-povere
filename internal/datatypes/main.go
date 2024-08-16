package datatypes

type Main struct {
	Site   Site   `yaml:"site"`
	Author Author `yaml:"author"`

	Posts []PostPage
	Tags  []TagPage
}
