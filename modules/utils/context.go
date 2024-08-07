package utils

import (
	"context"
	"fmt"

	app "github.com/venusforest2013/config/application"
)

func Module(ctx context.Context, name string) (app.Module, error) {
	var m app.Module
	v := ctx.Value(app.Key)
	if v != nil {
		m = v.(*app.Application).Module(name)
	}
	if m == nil {
		return nil, fmt.Errorf("app.module '%s' not found", name)
	}
	return m, nil
}

//
//type Reusable interface {
//	Reset()
//}
//
//type Pool struct {
//	rt   reflect.Type
//	pool *sync.Pool
//}
//
//func NewPool(rt reflect.Type) *Pool {
//	if rt.Kind() == reflect.Ptr {
//		rt = rt.Elem()
//	}
//	return &Pool{
//		rt: rt,
//		pool: &sync.Pool{
//			New: func() interface{} {
//				return reflect.New(rt).Pointer()
//			},
//		},
//	}
//}
//
//func (p *Pool) Get() interface{} {
//	return p.pool.Get()
//}
//
//func (p *Pool) Put(reusable Reusable) {
//	reusable.Reset()
//	p.pool.Put(reusable)
//}
