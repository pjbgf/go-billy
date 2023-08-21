module github.com/go-git/go-billy/v5

// go-git supports the last 3 stable Go versions.
go 1.18

replace github.com/cyphar/filepath-securejoin => github.com/pjbgf/filepath-securejoin v0.0.0-20230821001828-0ca74e6d4bf8

require (
	github.com/cyphar/filepath-securejoin v0.2.3
	github.com/onsi/gomega v1.27.10
	golang.org/x/sys v0.11.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
