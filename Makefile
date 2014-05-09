PROG_NAME=elwyn
DEFAULT_INSTALL=/usr/local/bin

default: deps fmt
	@go build -o bin/$(PROG_NAME)

run: default
	@bin/elwyn

deps:
	@go list -f "{{ range .Deps }}{{ . }} {{ end }}" ./ | tr ' ' '\n' | awk '!/^.\//' | xargs go get

fmt:
	@go fmt

clean:
	@rm -rf log bin

todo:
	@grep -nri "todo"

install: default
	@if test "$(PREFIX)" = "" ; then \
		cp bin/$(PROG_NAME) $(DEFAULT_INSTALL)/$(PROG_NAME) ; \
	else \
		cp bin/$(PROG_NAME) $(PREFIX)/$(PROG_NAME); \
	fi
