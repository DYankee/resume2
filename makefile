# Makefile

.PHONY: dev build generate css templ

# Generate templ files, build CSS, then run
dev: templ css
	go run main.go

# Generate Go code from .templ files
templ:
	templ generate

# Build Tailwind CSS
css:
	npx tailwindcss -i static/css/input.css -o static/css/output.css --minify

# Watch mode (run in separate terminals)
watch-templ:
	templ generate --watch

watch-css:
	npx tailwindcss -i static/css/input.css -o static/css/output.css --watch

build: templ css
	go build -o portfolio main.go