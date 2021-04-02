package main

import (
	"context"

	"github.com/iakud/plume"
	"github.com/iakud/plume/log"
)

type GameApp struct {
}

func (game *GameApp) Init() {

}

func (game *GameApp) Run(ctx context.Context) {
	log.Info("game run")
	<-ctx.Done()
}

func (game *GameApp) Destory() {

}

func main() {
	plume.Run(&GameApp{})
}
