package pagebuilder

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	// "github.com/russross/blackfriday/v2"
	"github.com/shurcooL/github_flavored_markdown"
)

type Builder struct {
	template string // The path to the template file
	outputPath string // Where the final site will be
	contentPath string // Where the content is
	ProjectTitle string
}

type Page struct {
	Title string
	Body  string
	Link string
	Toc string
}

type NavLink struct {
	Title string
	Link string
}

func NewBuilder(template string, contentPath string, outputPath string) *Builder {
	return &Builder{
		template: template,
		outputPath: outputPath,
		contentPath: contentPath,
	}
}

func (b *Builder) Build() {
	os.RemoveAll(b.outputPath)
	paths := b.scanFiles()
	pages := b.parsePages(paths)
	nav := b.makeNavigation(pages)
	tmpl := template.Must(template.ParseFiles(b.template))

	var wg sync.WaitGroup
	for _, page := range pages {
		wg.Add(1)
		go func(p Page, n []NavLink, t *template.Template) {
			b.outputPage(p, n, t)
			wg.Done()
		}(page, nav, tmpl)
	}
	wg.Wait()
}

type PageData struct {
	Navigation []NavLink
	Title string
	ProjectTitle string
	Body string
	Toc string
}

func (b *Builder) outputPage(p Page, navigation []NavLink, t *template.Template ) {
	outputDir := b.outputPath + "/" + p.Link

	os.MkdirAll(outputDir, 0755)

	output, err := os.Create(outputDir + "/" + "index.html")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// Execute the template
	pagedata := PageData{
		Navigation: navigation,
		Title: p.Title,
		ProjectTitle: b.ProjectTitle,
		Body: p.Body,
		Toc: p.Toc,
	}

	err = t.Execute(output, pagedata)
	if err != nil {
		panic(err)
	}
}

func (b *Builder) makeNavigation(pages []Page) []NavLink {
	// Build navigation
	links := []NavLink{}
	
	// Make sure the home page is first
	links = append(links, NavLink{
		Title: "Home",
		Link: "/",
	})

	for _, page := range pages {
		// If home, continue
		if page.Title == "Home" {
			continue
		}

		links = append(links, NavLink{
			Title: page.Title,
			Link: page.Link,
		})
	}

	return links
}

func (b *Builder) scanFiles() []string {
	files := []string{}

	walkDirectory(b.contentPath, &files)
	
	return files
}

func (b *Builder) parsePages(paths []string) []Page {
	pages := []Page{}

	var pwaitgroup sync.WaitGroup // Heavy operations here so we'll do it concurrently
	for _, path := range paths {
		pwaitgroup.Add(1)
		go func (pth string)  {
			p := Page{
				Title: makeTitleFromPath(pth),
				Link: makeLinkFromPath(pth),
				Body: parseMarkdown(pth),
				Toc: MakeToc(pth),
			}

			pages = append(pages, p)
			pwaitgroup.Done()
		}(path)
	}
	pwaitgroup.Wait()

	return pages
}

// Utility functions
func basename(path string) string {
	return filepath.Base(path)[:len(filepath.Base(path)) - len(filepath.Ext(path))]
}

func ucfirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func parseMarkdown(path string) string {
	md, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	parsed := github_flavored_markdown.Markdown(md)
	// parsed := blackfriday.Run(md)
	return string(parsed)
}

func makeTitleFromPath(p string) string {
	if basename(p) == "index" {
		return "Home"
	}

	return ucfirst(basename(p))
}

func makeLinkFromPath(p string) string {
	if basename(p) == "index" {
		return "/"
	} 

	return "/" + basename(p)
}

/*
Get all the markdown files in the directory
*/
func walkDirectory(current string, paths *[]string) {
	// Get all the files in the directory
	files, err := os.ReadDir(current)
	if err != nil {
		panic(err)
	}
	
	for _, file := range files {
		switch {
		case file.IsDir():
			walkDirectory(current + "/" + file.Name(), paths)
			break
		case filepath.Ext(file.Name()) != ".md":
			break
		default:
			*paths = append(*paths, current + "/" + file.Name())
		}
	}
}