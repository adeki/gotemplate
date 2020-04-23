# Add in your Makefile
# you have to install https://github.com/cespare/reflex
.PHONY: watch
watch:
	reflex -r '\.tmpl\z' -- sh -c 'go generate tools/statik/main.go'
