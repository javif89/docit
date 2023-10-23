package pagebuilder

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type TocItem struct {
	Title string
	Link string
	SubItems []TocItem
}

func MakeToc(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	headings := []string{}

	s := bufio.NewScanner(file)

	for s.Scan() {
		l := strings.Trim(s.Text(), " ")
		if strings.HasPrefix(l, "#") {
			headings = append(headings, l)
		}
	}

	// If no headings or empty page return blank
	if len(headings) == 0 {
		return ""
	}

	// get the level of the first heading
	startDepth := len(strings.Split(headings[0], " ")[0])
	toc, _ := parseToc(headings, 0, 0, startDepth)
	// printTocSimple(toc, 0)

	txt := []string{}
	for _, i := range toc {
		txt = append(txt, printItem(i))
	}

	templ, e := template.New("toc").Parse("<ul id=\"table-of-contents\">{{range .Items}} {{ . }} {{ end }}</ul>")

	if e != nil {
		panic(e)
	}

	var buf bytes.Buffer
	templ.Execute(&buf, map[string]interface{}{"Items": txt})

	// tocHtml := printToc(toc, 1)
	// os.WriteFile("./toc.html", buf.Bytes(), 0755)

	return buf.String()
}

type ItemHtmlData struct {
	Title string
	Link string
	Sublist string
}

func printItem(item TocItem) string {
	itemTemplate := `
		<li> <a href="{{ .Link }}"> {{ .Title }} </a> </li>
	`

	listTemplate := `
		{{ .Item }}
		{{ if .HasSub }}
		<ul>
			{{ range .Sub }}
			{{ . }}
			{{ end }}
		</ul>
		{{ end }}
	`

	it, err := template.New("item").Parse(itemTemplate)
	lt, e := template.New("list").Parse(listTemplate)
	if err != nil || e != nil {
		panic(e)
	}

	var itemHtml bytes.Buffer
	it.Execute(&itemHtml, map[string]interface{}{"Title": item.Title, "Link": item.Link})

	sub := []string{}

	if len(item.SubItems) > 0 {
		for _, si := range item.SubItems {
			sub = append(sub, printItem(si))
		}
	}

	var fullHtml bytes.Buffer
	lt.Execute(&fullHtml, map[string]interface{}{"Item": itemHtml.String(), "Sub": sub, "HasSub": (len(sub) > 0)})

	return fullHtml.String()
}

func printTocSimple(items []TocItem, depth int) {
	indent := strings.Repeat("*", depth+1)

	for _, i := range items {
		fmt.Println(indent,i.Title)
		if len(i.SubItems) > 0 {
			printTocSimple(i.SubItems, depth+1)
		}
	}
}

type Toc struct {
	Items []TocItem
	Sublist string
}

// Recursively parse headings based on the level above them
// depth: keeps track of where we are in the recursion
// index: our cursor for moving through the array
// level: the current level of heading we're at. ex: ### = level 3
func parseToc(headings []string, depth int, index int, level int) ([]TocItem, int) {
	toc := []TocItem{}
	
	for index < len(headings) {
		// Set up the item
		item := TocItem{Title: makeHeadingTitle(headings[index]), Link: makeHeadingLink(headings[index])}

		_, nextlvl := getNext(headings, index)

		// If the next level is above the current indentation and we're inside
		// a recursive call let the parent function handle it
		if nextlvl < level && depth > 0 {
			toc = append(toc, item)
			return toc, index
		}

		// If the next heading is higher indentation adjust the current level
		// ex: next: # current: ##
		if nextlvl < level {
			level = nextlvl
		}

		// If the next level is below the current level process the sub
		// items recursively. ex: current # next: ##
		if nextlvl > level {
			new, mov := parseToc(headings, depth+1, index+1, nextlvl)
			item.SubItems = new
			index = mov // Move the cursor to the current position after processing sub items
		}

		// Add the item and sub items to the list
		toc = append(toc, item)

		// Move cursor to the next item
		index += 1
	}

	return toc, index
}

func getlvl(h string) int {
	return len(strings.Split(h, " ")[0])
}

func getNext(headings []string, i int) (string, int) {
	if i+1 <= len(headings)-1 {
		return headings[i+1], getlvl(headings[i+1])
	}

	return "", 0
}

func makeHeadingTitle(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "#", "")
}

func makeHeadingLink(s string) string {
	return "#" + strings.ToLower(strings.Replace(strings.TrimSpace(strings.ReplaceAll(s, "#", "")), " ", "-", -1))
}