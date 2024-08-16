package parser

import (
	"errors"
	"html/template"
	"os"
	"path"
	"regexp"
	"strings"

	md "github.com/gomarkdown/markdown"
	mdparser "github.com/gomarkdown/markdown/parser"
	striphtml "github.com/grokify/html-strip-tags-go"
	"github.com/nathan-fiscaletti/md-blog/internal/datatypes"
	"github.com/nathan-fiscaletti/md-blog/internal/util"
	stripmd "github.com/writeas/go-strip-markdown"
)

type metaDataParser func(post *datatypes.Post, val string) error

var metaDataParsers map[string]metaDataParser = map[string]metaDataParser{
	"title": func(post *datatypes.Post, val string) error {
		post.Title = val
		return nil
	},

	"author": func(post *datatypes.Post, val string) error {
		post.Author.Name = val
		return nil
	},

	"author_bio": func(post *datatypes.Post, val string) error {
		post.Author.Bio = val
		return nil
	},

	"image": func(post *datatypes.Post, val string) error {
		post.Image = val
		return nil
	},

	"author_avatar": func(post *datatypes.Post, val string) error {
		post.Author.Avatar = val
		return nil
	},

	"tags": func(post *datatypes.Post, val string) error {
		tagValues := strings.Split(strings.Trim(val, " \r\n\t"), ",")

		tags := []datatypes.Tag{}
		for _, tagValue := range tagValues {
			tagValue = strings.Trim(tagValue, " \r\n\t")
			tags = append(tags, datatypes.Tag{
				Name: strings.ToUpper(tagValue),
				Url:  "tag-" + util.UrlSafe(strings.ToLower(tagValue)) + ".html",
			})
		}

		post.Tags = tags
		return nil
	},
}

type markdownPost struct {
	name     string
	location string
}

func (p markdownPost) Read(main datatypes.Main) (datatypes.Post, error) {
	var zeroValue datatypes.Post

	file, err := os.ReadFile(p.location)
	if err != nil {
		return zeroValue, err
	}

	fileStat, err := os.Stat(p.location)
	if err != nil {
		return zeroValue, err
	}

	fileModifiedDate := fileStat.ModTime().Format(main.Site.DateFormat)

	data := string(file)
	lines := strings.Split(data, "\n")

	post := datatypes.Post{
		Date:   fileModifiedDate,
		Author: main.Author,
	}

	for _, line := range lines {
		patternKeys := ``
		for key := range metaDataParsers {
			if patternKeys != `` {
				patternKeys += `|`
			}
			patternKeys += key
		}

		pattern := regexp.MustCompile(`!!(?P<key>` + patternKeys + `)\s(?P<val>.*)`)
		matches := pattern.FindAllStringSubmatch(line, -1)

		for _, match := range matches {
			matchcount := len(match)

			if matchcount > 0 && matchcount < 3 {
				return zeroValue, errors.New(line)
			}

			if matchcount > 0 {
				if _, exists := metaDataParsers[match[1]]; exists {
					metaDataParsers[match[1]](&post, strings.Trim(match[2], " \r\n\t"))
					data = strings.Replace(data, line, "", 1)
				}
			}
		}
	}

	markdownParser := mdparser.NewWithExtensions(
		mdparser.CommonExtensions | mdparser.Tables |
			mdparser.Footnotes | mdparser.Titleblock |
			mdparser.AutoHeadingIDs | mdparser.SuperSubscript |
			mdparser.LaxHTMLBlocks,
	)

	content := md.ToHTML([]byte(data), markdownParser, nil)

	post.Content = template.HTML(string(content))

	previewContent := striphtml.StripTags(data)
	previewContent = stripmd.Strip(previewContent)
	if len(previewContent) > 425 {
		previewContent = previewContent[0:425] + " ..."
	}

	post.Preview = previewContent
	post.Url = util.UrlSafe(post.Title) + ".html"

	return post, nil
}

type markdownPosts struct {
	dir   string
	posts []markdownPost
}

func (p *markdownPosts) Posts(main datatypes.Main) ([]datatypes.PostPage, error) {
	out := []datatypes.PostPage{}

	for idx, post := range p.posts {
		var nextPost *datatypes.Post
		var previousPost *datatypes.Post

		if idx > 0 {
			prev, err := p.posts[idx-1].Read(main)
			if err != nil {
				return nil, err
			}

			previousPost = &prev
		}

		if idx < len(p.posts)-1 {
			next, err := p.posts[idx+1].Read(main)
			if err != nil {
				return nil, err
			}

			nextPost = &next
		}

		parsedPost, err := post.Read(main)
		if err != nil {
			return nil, err
		}

		out = append(out, datatypes.PostPage{
			Site:     main.Site,
			Post:     parsedPost,
			Next:     nextPost,
			Previous: previousPost,
		})
	}

	return out, nil
}

// newMarkdownPosts creates a new MarkdownPosts object by reading all of the
// markdown files in the specified directory and creating a MarkdownPost object
// for each file.
func newMarkdownPosts(dir string) (*markdownPosts, error) {
	out := &markdownPosts{
		dir:   dir,
		posts: []markdownPost{},
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			out.posts = append(out.posts, markdownPost{
				name:     strings.TrimSuffix(file.Name(), path.Ext(file.Name())),
				location: path.Join(dir, file.Name()),
			})
		}
	}

	return out, nil
}
