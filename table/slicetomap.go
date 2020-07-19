package table

import (
	"errors"
	"reflect"
)

type options struct {
	key interface{}
}

type Option func(o *options)

func WithKey(key interface{}) Option {
	return func(o *options) {
		o.key = key
	}
}

func SliceToMap(in interface{}, out interface{}, o ...Option) error {
	var opts options
	for _, option := range o {
		option(&opts)
	}
	// 检查slice类型
	inValue := reflect.ValueOf(in)
	if inValue.Kind() == reflect.Ptr {
		inValue = inValue.Elem()
	}
	if inValue.Kind() != reflect.Slice {
		return errors.New("in type error")
	}
	inType := inValue.Type()
	// 检查map类型
	outValue := reflect.ValueOf(out)
	if outValue.Kind() == reflect.Ptr {
		outValue = outValue.Elem()
	}
	if outValue.Kind() != reflect.Map {
		return errors.New("out type error")
	}
	outType := outValue.Type()
	// 检查slice和map类型匹配
	if inType.Elem() != outType.Elem() {
		return errors.New("type not match")
	}

	var keyfuncValue reflect.Value
	if key := opts.key; key == nil {
		// 检查Key方法
		method, ok := inType.Elem().MethodByName("Key")
		if !ok {
			return errors.New("key method error")
		}
		keyfuncValue = method.Func
	} else {
		// 检查keyfunc类型
		keyfuncValue = reflect.ValueOf(key)
	}

	if keyfuncValue.Kind() != reflect.Func {
		return errors.New("key type error")
	}
	keyfuncType := keyfuncValue.Type()
	// 检查keyfunc参数和返回值类型
	if keyfuncType.NumIn() != 1 || keyfuncType.In(0) != inType.Elem() {
		return errors.New("key in type error")
	}
	if keyfuncType.NumOut() != 1 || keyfuncType.Out(0) != outType.Key() {
		return errors.New("key out type error")
	}

	// slice to map
	for i := 0; i < inValue.Len(); i++ {
		objValue := inValue.Index(i)
		keyValue := keyfuncValue.Call([]reflect.Value{objValue})[0]
		// 保存到map
		outValue.SetMapIndex(keyValue, objValue)
	}
	return nil
}
