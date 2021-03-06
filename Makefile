default: build

build: FORCE $(patsubst %.pod,%.json,$(wildcard build/man/*.pod))
	go run ./build
build/man/%.json: build/man/%.pod build/man2json.pl
	perl build/man2json.pl $< > $@

import-manpages: FORCE
	sh ./build/import-man-pages.sh

.PHONY: FORCE
