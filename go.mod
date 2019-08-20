module github.com/jaypipes/gdt

go 1.12

replace github.com/jaypipes/gdt/interfaces => ./interfaces

replace github.com/jaypipes/gdt/errors => ./errors

replace github.com/jaypipes/gdt/testcase => ./testcase

require (
	github.com/ghodss/yaml v1.0.0
	github.com/stretchr/testify v1.4.0
)
