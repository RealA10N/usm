# Determine which Go executable to use.
# 'richgo' is a wrapper around the 'go' executable that prints more colorful
# information and test summary. We use richgo if it is avaliable.
GO := `if command -v richgo >/dev/null 2>&1; then echo richgo; else echo go; fi`
PY := "python3"
COVERPROFILE := "coverage.out"

build:
	{{GO}} build

test:
	{{GO}} test ./...

cover:
	courtney -o={{COVERPROFILE}} | richgo testfilter

fmt:
	{{GO}} mod tidy
	{{GO}} fmt ./...
	mdformat .

setup:
	{{GO}} get -u github.com/kyoh86/richgo
	{{GO}} get -u github.com/dave/courtney 
	{{GO}} get -t -v ./...
	{{PY}} -m pip install --upgrade pip
	{{PY}} -m pip install mdformat
