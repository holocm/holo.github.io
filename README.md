# Contents for <https://holocm.org>

This website is maintained at <https://github.com/holocm/holocm.github.io>.
Please submit issues and pull requests over there.

As suggested by the repository name, the URL <http://holocm.github.io> also
works, but <https://holocm.org> is the canonical site location.

## Build process

`index.html` and `site.css` are handwritten. All other HTML files are produced
by a bespoke static site generator at `build/main.go`. The main Makefile target
will run that program. Build requirements include:

- a Go compiler
- a reasonably recent Perl interpreter

All outputs produced by the static site generator are to be committed into Git,
since the deployment process just does a `git checkout` into the document root
of the webserver.
