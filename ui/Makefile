.PHONY: clean build

SEMANTIC_UI="https://github.com/Semantic-Org/Semantic-UI-CSS/archive/2.2.6.tar.gz"

build: static/semantic

static/semantic: | static
	curl -Ls $(SEMANTIC_UI) | tar -zxf- -C static
	mv static/Semantic-UI-CSS-* static/semantic

static:
	mkdir static

clean:
	rm -rf static/semantic
