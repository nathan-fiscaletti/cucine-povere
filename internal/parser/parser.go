package parser

import (
	"os"
	"time"

	"github.com/nathan-fiscaletti/md-blog/internal/datatypes"
	"gopkg.in/yaml.v2"
)

func Parse() (datatypes.Main, error) {
	var main datatypes.Main

	yamldata, err := os.ReadFile("config.yml")
	if err != nil {
		return main, err
	}

	err = yaml.Unmarshal(yamldata, &main)
	if err != nil {
		return main, err
	}

	main.Site.CopyrightYear = time.Now().Year()

	parser, err := newMarkdownPosts("./posts")
	if err != nil {
		return main, err
	}

	posts, err := parser.Posts(main)
	if err != nil {
		return main, err
	}

	var tagsToPosts map[datatypes.Tag][]datatypes.Post = map[datatypes.Tag][]datatypes.Post{}
	for _, post := range posts {
		for _, tag := range post.Post.Tags {
			if _, ok := tagsToPosts[tag]; !ok {
				tagsToPosts[tag] = []datatypes.Post{}
			}

			tagsToPosts[tag] = append(tagsToPosts[tag], post.Post)
		}
	}

	tags := []datatypes.TagPage{}
	for tag, posts := range tagsToPosts {
		tags = append(tags, datatypes.TagPage{
			Tag:   tag,
			Site:  main.Site,
			Posts: posts,
		})
	}

	main.Tags = tags
	main.Posts = posts

	return main, nil
}
