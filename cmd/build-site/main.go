package main

import (
	"fmt"

	"github.com/nathan-fiscaletti/md-blog/internal/datatypes"
	"github.com/nathan-fiscaletti/md-blog/internal/generator"
	"github.com/nathan-fiscaletti/md-blog/internal/parser"
)

type PostPage struct {
	Post     datatypes.Post
	Previous datatypes.Post
	Next     datatypes.Post
}

func main() {
	mainTemplateData, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	err = generator.FillTemplate("index.html", mainTemplateData)
	if err != nil {
		panic(err)
	}

	for _, postPage := range mainTemplateData.Posts {
		err = generator.FillNamedTemplate("post.html", postPage.Post.Url, postPage)
		if err != nil {
			panic(err)
		}
	}

	for _, tagPage := range mainTemplateData.Tags {
		fmt.Printf("%v\n", tagPage.Tag.Url)
		err = generator.FillNamedTemplate("tag.html", tagPage.Tag.Url, tagPage)
		if err != nil {
			panic(err)
		}
	}
}
