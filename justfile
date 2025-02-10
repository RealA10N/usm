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
	{{GO}} test -coverprofile={{COVERPROFILE}} ./...

fmt:
	{{GO}} mod tidy
	{{GO}} fmt ./...
	mdformat .

# count line of code, excluding tests.
cloc:
	cloc --not-match-f='.*_test.go' . 

setup:
	{{GO}} install github.com/kyoh86/richgo@v0.3.12
	{{PY}} -m pip install --upgrade pip
	{{PY}} -m pip install mdformat
