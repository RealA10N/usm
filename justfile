# Determine which Go executable to use.
# 'richgo' is a wrapper around the 'go' executable that prints more colorful
# information and test summary. We use richgo if it is avaliable.
GO := `if command -v richgo >/dev/null 2>&1; then echo richgo; else echo go; fi`

build:
	{{GO}} build

test:
	{{GO}} test ./...

install:
	{{GO}} install github.com/dave/courtney@ccf8e7a919f4e25cb9bd482bb2bcbf4a647e2b85
	{{GO}} install github.com/kyoh86/richgo@98af5f3a762dabdd7f3c30a122a7950fc3cdb4f1

cover:
	courtney -v

fmt:
	{{GO}} mod tidy
	{{GO}} fmt ./...
	mdformat .