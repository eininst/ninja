package ninja

import (
	"github.com/facebookgo/inject"
	"sync"
)

var graph inject.Graph
var lazyObjects []any
var mux = &sync.Mutex{}

func LazyProvide(objects ...any) {
	mux.Lock()
	defer mux.Unlock()
	for _, obj := range objects {
		lazyObjects = append(lazyObjects, obj)
	}
}

func Provide(objects ...any) {
	for _, obj := range objects {
		err := graph.Provide(&inject.Object{Value: obj})
		if err != nil {
			panic(err)
		}
	}
}

func Populate(objects ...any) {
	mux.Lock()
	defer mux.Unlock()

	Provide(objects...)

	if err := graph.Populate(); err != nil {
		panic(err)
	}

	Provide(lazyObjects...)
	if err := graph.Populate(); err != nil {
		panic(err)
	}
	lazyObjects = []any{}
}

func Objects() []*inject.Object {
	return graph.Objects()
}

type Init interface {
	Init()
}

func Install(i Init, objects ...any) {
	Provide(i)
	Populate(objects...)
	i.Init()
}
