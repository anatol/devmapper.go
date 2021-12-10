module examples

go 1.17

require (
	github.com/anatol/devmapper.go v0.0.0-20211209025129-80464f6a11ee
	github.com/stretchr/testify v1.7.0
	github.com/tych0/go-losetup v0.0.0-20170407175016-fc9adea44124
	golang.org/x/sys v0.0.0-20211209171907-798191bca915
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/anatol/devmapper.go => ./..
