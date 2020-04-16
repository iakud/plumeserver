package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/iakud/starry/service/pb"

	"github.com/golang/protobuf/proto"
)

const cmd1 int16 = 0x0001
const cmd2 int16 = 0x0012

type userKey struct{}

func newUserContext(ctx context.Context, u *user) context.Context {
	return context.WithValue(ctx, userKey{}, u)
}

func fromUserContext(ctx context.Context) (*user, bool) {
	u, ok := ctx.Value(userKey{}).(*user)
	return u, ok
}

type user struct {
	name string
}

func testHandler1(cmd int16, message *pb.Test) {
	fmt.Printf("testHandler1 message: id=%v, name=%v\n", message.GetId(), message.GetName())
}

func testHandler2(ctx context.Context, cmd int16, message *pb.Test) {
	u, ok := fromUserContext(ctx)
	if !ok {
		return
	}
	fmt.Printf("testHandler2 user: %v, message: id=%v, name=%v\n", u.name, message.GetId(), message.GetName())
}

func createMessage() ([]byte, error) {
	// 先构造一个message
	message := pb.Test{
		Id:   proto.Int32(101),
		Name: proto.String("上海"),
	}

	return proto.Marshal(&message)
}

func TestMessageHub(t *testing.T) {
	messageHub := NewMessageHub()
	messageHub.Register(cmd1, testHandler1)
	messageHub.Register(cmd2, testHandler2)

	// 创建一条消息
	buf, err := createMessage()
	if err != nil {
		t.Fatal(err)
	}
	// 派发消息
	if err := messageHub.Dispatch(nil, cmd1, buf); err != nil {
		t.Fatal(err)
	}
	// 传递参数
	u := &user{name: "暖暖"}
	ctx := newUserContext(context.Background(), u)
	if err := messageHub.Dispatch(ctx, cmd2, buf); err != nil {
		t.Fatal(err)
	}
}
