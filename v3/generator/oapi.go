//go:generate oapi-codegen -generate types -package v3 -o ../types.gen.go -templates templates source.json
//go:generate oapi-codegen -generate client -package v3 -o ../client.gen.go -templates templates source.json

package main
