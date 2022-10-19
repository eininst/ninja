package ninja

import (
	"github.com/eininst/flog"
	"github.com/facebookgo/inject"
	"sync"
)

var graph inject.Graph
var lazyObjects = []any{}
var mux = &sync.Mutex{}

func Provide(objects ...any) {
	for _, obj := range objects {
		flog.Info(obj)
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

func LazyProvide(objects ...any) {
	mux.Lock()
	defer mux.Unlock()

	for _, obj := range objects {
		lazyObjects = append(lazyObjects, obj)
	}
}

func LazyPopulate(objects ...any) {
	mux.Lock()
	defer mux.Unlock()

	Populate(objects...)
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
