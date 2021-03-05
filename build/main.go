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
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-commonmark/markdown"
)

func main() {
	//load page template
	tmplBytes, err := ioutil.ReadFile("build/template.html.gotpl")
	must(err)
	tmpl, err := template.New("build/template.html.gotpl").Parse(string(tmplBytes))
	must(err)

	//render pages from Markdown
	paths, err := filepath.Glob("build/pages/*.md")
	must(err)
	for _, path := range paths {
		loadPageFromMarkdown(path).Render(tmpl)
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
