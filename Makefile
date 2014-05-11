PROG_NAME=elwyn
DEFAULT_INSTALL=/usr/local/bin
LESS_FILES=$(wildcard web/static/css/*.less)
CSS_FILES=$(LESS_FILES:.less=.css)

default: deps fmt less
	@go build -o bin/$(PROG_NAME)

run: default
	@bin/elwyn -d ./ -l ./logs

deps:
	@go list -f "{{ range .Deps }}{{ . }} {{ end }}" ./ | tr ' ' '\n' | awk '!/^.\//' | xargs go get

fmt:
	@go fmt

less: $(CSS_FILES)

clean:
	@rm -rf logs bin

todo:
	@grep -nri "todo"

%.css: %.less
	@lessc $< > $@

install:
	service $(PROG_NAME) stop
	@if test "$(PREFIX)" = "" ; then \
		cp bin/$(PROG_NAME) $(DEFAULT_INSTALL)/$(PROG_NAME) ; \
	else \
		cp bin/$(PROG_NAME) $(PREFIX)/$(PROG_NAME); \
	fi
	@mkdir -p /etc/$(PROG_NAME)/web
	@cp -r ./web/static /etc/$(PROG_NAME)/web/
	@cp -r ./web/views /etc/$(PROG_NAME)/web/
	@cp ./config.json /etc/$(PROG_NAME)/config.json
	@cp deploy/nginx/elwyn.conf /etc/nginx/sites-enabled/elwyn
	@service nginx reload
	service $(PROG_NAME) start
