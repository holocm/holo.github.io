all: $(patsubst %.md,%.html,$(filter-out README.md,$(wildcard *.md)))
.PHONY: all

%.html: %.md holo.html5.template Makefile
	 < $< pandoc -f markdown -t html5 --standalone --template=holo.html5.template \
	| sed -r -e 's,<pre><code>,<pre>,g' \
	         -e 's,</code></pre>,</pre>,g' \
	         -e 's/<(h.) id="[^"]*">/<\1>/g' \
	         -e 's,<(/?)code>,<\1tt>,g' \
	         -e 's,</a> <a,</a>\n<a,g' \
	 > $@

clean:
	rm -f -- *.html
.PHONY: clean

deploy: all
	git branch -D tmp || true
	git checkout -b tmp
	echo '!*.html' >> .gitignore
	git add .
	git commit -m 'Automatic deployment'
	git checkout gh-pages # ensure it exists as a local branch
	git checkout tmp # switch back
	git merge --no-ff --no-edit -s ours gh-pages
	git checkout gh-pages
	git merge tmp
	git branch -d tmp
	git push origin gh-pages
.PHONY: deploy

.DELETE_ON_ERROR:
SHELL = bash -o pipefail
