package main

import (
	"context"

	"github.com/iakud/plume"
)

type GameApp struct {
}

func (game *GameApp) Init() {

}

func (game *GameApp) Run(ctx context.Context) {
	<-ctx.Done()
}

func (game *GameApp) Destory() {

}

func main() {
	plume.Run(&GameApp{})
}
