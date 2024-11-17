# Determine which Go executable to use.
# 'richgo' is a wrapper around the 'go' executable that prints more colorful
# information and test summary. We use richgo if it is avaliable.
GO := `if command -v richgo >/dev/null 2>&1; then echo richgo; else echo go; fi`
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