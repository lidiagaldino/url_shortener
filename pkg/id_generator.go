package pkg

import "github.com/teris-io/shortid"

type IDGenerator interface {
	Generate() (string, error)
}

type ShortIDGenerator struct{}

func (g *ShortIDGenerator) Generate() (string, error) {
	return shortid.Generate()
}
