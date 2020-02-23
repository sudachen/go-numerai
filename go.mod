module github.com/sudachen/go-numerai

go 1.13

replace github.com/sudachen/go-ml => ./go-ml

replace github.com/sudachen/go-fp => ./go-ml/go-fp

require (
	github.com/sudachen/go-fp v0.0.0-20200221050211-7c5fbafbab76
	github.com/sudachen/go-ml v0.0.0-00010101000000-000000000000
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	gotest.tools v2.2.0+incompatible
)
