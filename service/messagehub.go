package service

import (
	"context"
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var messageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var ErrNoHandler = errors.New("No handler")
var ErrUnknowType = errors.New("Unknow type")

type messageHandler struct {
	handler reflect.Value
	pbType  reflect.Type

	requireContext bool
}

type MessageHub struct {
	handlerMap map[int16]*messageHandler
}

func NewMessageHub() *MessageHub {
	messageHub := &MessageHub{
		handlerMap: make(map[int16]*messageHandler),
	}
	return messageHub
}

func (this *MessageHub) Register(cmd int16, cb interface{}) {
	handler := reflect.ValueOf(cb)
	handlerType := handler.Type()
	if handlerType.Kind() != reflect.Func {
		panic("no function")
	}
	var nextArg int = 0
	var requireContext bool = false
	var pbType reflect.Type
	switch handlerType.NumIn() {
	case 3: // 3个参数
		// context必须实现接口context.Context
		argCtx := handlerType.In(nextArg)
		if argCtx.Kind() != reflect.Interface {
			panic("unknow args")
		}
		if !argCtx.Implements(contextType) {
			panic("unknow args")
		}
		requireContext = true // 需要context
		nextArg++
		fallthrough // 继续解析后2个参数
	case 2:
		// Cmd 是int16
		argCmd := handlerType.In(nextArg)
		if argCmd.Kind() != reflect.Int16 {
			panic("unknow args")
		}
		nextArg++
		// pb必须实现接口proto.Message
		argPb := handlerType.In(nextArg)
		if argPb.Kind() != reflect.Ptr {
			panic("unknow args")
		}
		if !argPb.Implements(messageType) {
			panic("unknow args")
		}
		pbType = argPb.Elem() // 保存pbType
	default:
		panic("unknow args")
	}
	this.handlerMap[cmd] = &messageHandler{handler, pbType, requireContext}
}

func (this *MessageHub) Dispatch(ctx context.Context, cmd int16, buf []byte) error {
	// 查找注册的消息
	messageHandler, ok := this.handlerMap[cmd]
	if !ok {
		return ErrNoHandler
	}

	var args []reflect.Value
	if messageHandler.requireContext {
		if ctx == nil {
			args = append(args, reflect.Zero(reflect.TypeOf((*context.Context)(nil)).Elem()))
		} else {
			args = append(args, reflect.ValueOf(ctx))
		}
	}
	argPb := reflect.New(messageHandler.pbType)
	pb, ok := argPb.Interface().(proto.Message)
	if !ok {
		return ErrUnknowType
	}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return err
	}

	args = append(args, reflect.ValueOf(cmd), argPb)
	messageHandler.handler.Call(args)
	return nil
}
