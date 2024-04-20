package main

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"log"

	"gopkg.in/yaml.v3"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	rhtml "github.com/yuin/goldmark/renderer/html"
)

type PostMetadata struct {
	Title  string   `yaml:"title"`
	Date   string   `yaml:"date"`
	Author string   `yaml:"author"`
	Tags   []string `yaml:"tags"`
	Slug   string   `yaml:"slug"`
}

type Post struct {
	PostMetadata
	Content template.HTML
}

type DynamicTemplateData map[string]interface{}

func parseMetadata(content string) (PostMetadata, string) {
	var metadata PostMetadata
	splitContent := strings.SplitN(content, "---", 3)
	if len(splitContent) < 3 {
		log.Fatalf("error: invalid markdown file")
	}

	err := yaml.Unmarshal([]byte(splitContent[1]), &metadata)
	// convert title to url friendly format
	metadata.Title = strings.ReplaceAll(metadata.Title, " ", "-")

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return metadata, splitContent[2] // Return the rest of the content without front matter
}

// ConvertMarkdownToHTML takes a path to a Markdown file and converts it to HTML, returning the metadata and the HTML content
func ConvertMarkdownToHTML(filePath string, style_name string) Post {
	input, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Parse the metadata
	metadata, content := parseMetadata(string(input))

	// Create a new markdown parser with Chroma highlighting
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithRendererOptions(
			rhtml.WithXHTML(),
			rhtml.WithUnsafe(),
		),
	)

	var sb strings.Builder
	err = markdown.Convert([]byte(content), &sb)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return Post{metadata, template.HTML(sb.String())}
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	pageDir := rootDir // Only parse files in the 'pages' directory
	cleanPageDir := filepath.Clean(pageDir)
	root := template.New("")

	err := filepath.Walk(cleanPageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			name := strings.TrimPrefix(path, cleanPageDir+"/")
			// remove the rootdir/ from the name
			name = strings.TrimPrefix(name, rootDir+"/")
			tmpl := root.New(name).Funcs(funcMap)
			_, err = tmpl.Parse(string(contents))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return root, err
}

func getPosts(rootDir string) ([]Post, error) {
	posts := []Post{}
	// get all the markdown files in the content directory
	contentDir := filepath.Join(rootDir, "posts")
	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			// convert the markdown file to HTML
			post := ConvertMarkdownToHTML(path, "atom-one-dark")

			posts = append(posts, post)
		}
		return nil
	})
	return posts, err
}

// handles the building of the site
func Build(rootDir string, outputDir string) {
	data := DynamicTemplateData{
		"posts": []Post{},
	}
	// delete the output directory if it exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		os.RemoveAll(outputDir)
	}

	// parse markdown files in the content directory
	posts, err := getPosts(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	data["posts"] = posts
	// parse templates in the views directory
	templates, err := findAndParseTemplates(rootDir, template.FuncMap{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// create the output directory
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// copy the static directory to the output directory
	staticDir := filepath.Join(rootDir, "static")
	err = filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(staticDir, path)
			if err != nil {
				return err
			}
			outputPath := filepath.Join(outputDir, "static", relPath)
			// create the directory structure in the output directory
			os.MkdirAll(filepath.Dir(outputPath), 0755)
			// copy the file
			input, err := os.Open(path)
			if err != nil {
				return err
			}
			output, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(output, input)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// generate the static pages in the pages dir. The pages directory structure should be preserved in the output directory.
	// Each page may include templates which should be executed with the appropriate data.
	err = filepath.Walk(filepath.Join(rootDir, "pages"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			relPath, err := filepath.Rel(filepath.Join(rootDir, "pages"), path)
			if err != nil {
				return err
			}
			// evaluate the template
			outputPath := filepath.Join(outputDir, relPath)
			// create the directory structure in the output directory
			os.MkdirAll(filepath.Dir(outputPath), 0755)
			// create the file
			output, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			tmpl := templates.Lookup("..\\..\\views\\pages\\" + relPath)
			if tmpl == nil {
				log.Fatalf("Template not found for page %s", relPath)
			}
			if tmpl != nil {
				err = tmpl.Execute(output, data)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// walk the posts directory and generate the posts
	err = filepath.Walk(filepath.Join(rootDir, "posts"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			// convert the markdown file to HTML
			post := ConvertMarkdownToHTML(path, "atom-one-dark")

			// create the output path
			outputPath := filepath.Join(outputDir, "posts", post.Slug+".html")
			// create the directory structure in the output directory
			os.MkdirAll(filepath.Dir(outputPath), 0755)
			// create the file
			output, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			// execute the template
			tmpl := templates.Lookup("..\\..\\views\\templates\\post.html")
			if tmpl == nil {
				log.Fatal("Template not found for post", post.Title)
			}
			if tmpl != nil {
				err = tmpl.Execute(output, post)
				if err != nil {
					return err
				}
			}
		} else if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// resolve and copy the file
			relPath, err := filepath.Rel(filepath.Join(rootDir, "posts"), path)
			if err != nil {
				return err
			}
			outputPath := filepath.Join(outputDir, "posts", relPath)
			// create the directory structure in the output directory
			os.MkdirAll(filepath.Dir(outputPath), 0755)

			// create the file
			output, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			tmpl := templates.Lookup("..\\..\\views\\posts\\" + relPath)
			if tmpl == nil {
				log.Fatal("Template not found for post file", relPath)
			}
			if tmpl != nil {
				err = tmpl.Execute(output, data)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func main() {
	Build("../../views", "../../output")
}
