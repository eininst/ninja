package ninja

import (
	"github.com/facebookgo/inject"
)

var graph inject.Graph

func Provide(objects ...any) {
	for _, obj := range objects {
		err := graph.Provide(&inject.Object{Value: obj})
		if err != nil {
			panic(err)
		}
	}
}

func Populate(objects ...any) {
	Provide(objects...)

	if err := graph.Populate(); err != nil {
		panic(err)
	}
}

func Objects() []*inject.Object {
	return graph.Objects()
}
