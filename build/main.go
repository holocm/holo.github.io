/*******************************************************************************
*
* Copyright 2021 Stefan Majewsky <majewsky@gmx.net>
*
* This program is free software: you can redistribute it and/or modify it under
* the terms of the GNU General Public License as published by the Free Software
* Foundation, either version 3 of the License, or (at your option) any later
* version.
*
* This program is distributed in the hope that it will be useful, but WITHOUT ANY
* WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
* A PARTICULAR PURPOSE. See the GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License along with
* this program. If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/golang-commonmark/markdown"
)

var allManpages = findManpages()

func main() {
	//load page template
	tmplBytes, err := ioutil.ReadFile("build/template.html.gotpl")
	must(err)
	tmpl, err := template.New("build/template.html.gotpl").Parse(string(tmplBytes))
	must(err)

	//render HTML from Markdown
	paths, err := filepath.Glob("build/pages/*.md")
	must(err)
	for _, path := range paths {
		loadPageFromMarkdown(path).Render(tmpl)
	}

	//render HTML from manpages
	for _, manpage := range allManpages {
		page := loadPageFromManpageJSON(manpage)
		page.AllManpages = allManpages
		page.Render(tmpl)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////
// pages in general

type page struct {
	OutPath string
	Title   string
	Content template.HTML
	//if present, causes the manpage subnavigation to render
	AllManpages []manpageInfo
}

func (p page) Render(tmpl *template.Template) {
	if strings.Contains(p.OutPath, "/") {
		must(os.MkdirAll(filepath.Dir(p.OutPath), 0777))
	}
	out, err := os.Create(p.OutPath)
	must(err)
	must(tmpl.Execute(out, p))
	must(out.Close())
}

////////////////////////////////////////////////////////////////////////////////
// pages from Markdown

var pageTitles = map[string]string{
	"example": "Example - Holo",
	"install": "Installation - Holo",
}

func loadPageFromMarkdown(path string) page {
	mdBytes, err := ioutil.ReadFile(path)
	must(err)

	md := markdown.New(
		markdown.HTML(true),
		markdown.Typographer(false),
	)
	out := strings.TrimSpace(md.RenderToString([]byte(mdBytes)))
	out = strings.TrimPrefix(out, "<p>")
	out = strings.TrimSuffix(out, "</p>")

	baseName := strings.TrimSuffix(filepath.Base(path), ".md")
	title := pageTitles[baseName]
	if title == "" {
		log.Fatalf("no title found for %s.html (add one to the `pageTitles` variable in build/main.go!)", baseName)
	}

	return page{
		OutPath: baseName + ".html",
		Title:   title,
		Content: template.HTML(out),
	}
}

////////////////////////////////////////////////////////////////////////////////
// pages from manpage (Perl POD -> JSON AST -> HTML)

type manpageInfo struct {
	OutPath string
	Title   string
	Section string
}

func findManpages() (result []manpageInfo) {
	paths, err := filepath.Glob("build/man/*.json")
	must(err)

	for _, path := range paths {
		baseName := strings.TrimSuffix(filepath.Base(path), ".json")
		baseNameFields := strings.Split(baseName, ".")
		result = append(result, manpageInfo{
			OutPath: fmt.Sprintf("man/%s.html", baseName),
			Title:   baseNameFields[0],
			Section: baseNameFields[1],
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})
	return result
}

func loadPageFromManpageJSON(info manpageInfo) page {
	jsonBytes, err := ioutil.ReadFile(fmt.Sprintf("build/man/%s.%s.json", info.Title, info.Section))
	must(err)
	var doc docNode
	must(json.Unmarshal(jsonBytes, &doc))

	return page{
		OutPath: info.OutPath,
		Title:   fmt.Sprintf("%s(%s)", info.Title, info.Section),
		Content: template.HTML(doc.ToHTML()),
	}
}

type docNode struct {
	Type     string            `json:"name"`
	Attrs    map[string]string `json:"attrs"`
	Children []docNode         `json:"children"`
}

//UnmarshalJSON implements the json.Unmarshaler interface.
func (n *docNode) UnmarshalJSON(buf []byte) error {
	//Nodes can have plain strings as children. When that happens, we parse the
	//plain string into a dummy node with the "__TEXT__" type.
	var s string
	err := json.Unmarshal(buf, &s)
	if err == nil {
		*n = docNode{"__TEXT__", map[string]string{"text": s}, nil}
		return nil
	}

	var data struct {
		Type     string            `json:"name"`
		Attrs    map[string]string `json:"attrs"`
		Children []docNode         `json:"children"`
	}
	err = json.Unmarshal(buf, &data)
	*n = data
	return err
}

func (n docNode) ChildrenToHTML() string {
	var result []string
	for _, child := range n.Children {
		result = append(result, child.ToHTML())
	}
	return strings.Join(result, "")
}

func (n docNode) ToHTML() string {
	switch n.Type {
	case "Document":
		return n.ChildrenToHTML()
	case "head1":
		return fmt.Sprintf(`<h1>%s</h1>`, n.ChildrenToHTML())
	case "head2":
		return fmt.Sprintf(`<h2>%s</h2>`, n.ChildrenToHTML())
	case "Para":
		return fmt.Sprintf(`<p>%s</p>`, n.ChildrenToHTML())
	case "Verbatim":
		return fmt.Sprintf(`<pre><code>%s</code></pre>`, removeRedundantIndentation(n.ChildrenToHTML()))
	case "over-text":
		return renderDescriptionList(n.Children)
	case "B":
		return fmt.Sprintf(`<strong>%s</strong>`, n.ChildrenToHTML())
	case "I":
		return fmt.Sprintf(`<em>%s</em>`, n.ChildrenToHTML())
	case "C", "F":
		return fmt.Sprintf(`<code>%s</code>`, n.ChildrenToHTML())
	case "L":
		href := template.HTMLEscapeString(linkAttrsToHref(n.Attrs))
		return fmt.Sprintf(`<a href=%q>%s</a>`, href, n.ChildrenToHTML())
	case "__TEXT__":
		return template.HTMLEscapeString(n.Attrs["text"])
	default:
		//unknown node type -> dump contents as red text to make it stand out
		contents := struct {
			Type    string            `json:"type"`
			Attrs   map[string]string `json:"attrs"`
			Content string            `json:"content"`
		}{n.Type, n.Attrs, template.HTMLEscapeString(n.ChildrenToHTML())}
		contentsJSON, _ := json.Marshal(contents)
		return fmt.Sprintf(`<span style="color:red">UNKNOWN %s</span>`, string(contentsJSON))
	}
}

//Renders the contents of a POD node of type "over-text" into the innerHTML of
//a <dl>. Each "item-text" item becomes a <dt>, and all other contents are
//pooled into <dd becomes a <dt>, and all other contents are pooled into <dd>.
func renderDescriptionList(children []docNode) string {
	var ddContents []string
	ddFinalize := func() string {
		if len(ddContents) == 0 {
			return ""
		}
		result := strings.Join(ddContents, "")
		ddContents = nil
		return fmt.Sprintf(`<dd>%s</dd>`, result)
	}

	var result []string
	for _, child := range children {
		if child.Type == "item-text" {
			result = append(result,
				ddFinalize(),
				fmt.Sprintf(`<dt>%s</dt>`, child.ChildrenToHTML()),
			)
		} else {
			ddContents = append(ddContents, child.ToHTML())
		}
	}

	result = append(result, ddFinalize())
	return fmt.Sprintf(`<dl>%s</dl>`, strings.Join(result, ""))
}

var manpageLinkRx = regexp.MustCompile(`^(.*)\((.*)\)$`)

//Converts attributes of a POD node of type "L" into a href attribute for <a>.
func linkAttrsToHref(attrs map[string]string) string {
	if attrs["type"] == "url" {
		return attrs["to"]
	}

	if attrs["type"] == "man" {
		match := manpageLinkRx.FindStringSubmatch(attrs["to"])
		if match != nil {
			title := match[1]
			section := match[2]
			//if this is one of our manpages, link inside this site
			for _, manpage := range allManpages {
				if manpage.Title == title && manpage.Section == section {
					return "/" + manpage.OutPath
				}
			}
			//if it's a foreign manpage, link to an outside manpage hoster
			return fmt.Sprintf("https://www.man7.org/linux/man-pages/man%s/%s.%[1]s.html", section, title)
		}
	}

	panic(fmt.Sprintf("cannot convert link attributes into href: %#v", attrs))
}

//Takes a codeblock and removes leading spaces that are shared by all nonempty lines.
func removeRedundantIndentation(in string) string {
	lines := strings.Split(in, "\n")
	isNonEmptyLine := make(map[int]bool)

	//count shared indentation
	commonIndent := math.MaxInt16
	for lineIdx, line := range lines {
		indent := -1
		for idx, char := range line {
			if !unicode.IsSpace(char) {
				indent = idx
				break
			}
		}
		if indent == -1 {
			//ignore lines with no non-printable chars
			continue
		}
		isNonEmptyLine[lineIdx] = true
		if commonIndent > indent {
			commonIndent = indent
		}
	}

	//remove shared indentation from nonempty lines
	for lineIdx, line := range lines {
		if isNonEmptyLine[lineIdx] {
			lines[lineIdx] = line[commonIndent:]
		}
	}

	return strings.Join(lines, "\n")
}
