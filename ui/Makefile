.PHONY: clean build

SEMANTIC_UI="https://github.com/Semantic-Org/Semantic-UI-CSS/archive/2.2.6.tar.gz"

build: static/semantic

static/semantic:
	curl -Ls $(SEMANTIC_UI) | tar -zxf- -C static
	mv static/Semantic-UI-CSS-* static/semantic

clean:
	rm -rf static/semantic
