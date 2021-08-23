package main

import (
	"context"

	"github.com/iakud/plume"
	"github.com/iakud/plume/log"
)

type GameApp struct {
}

func (game *GameApp) Init() {
	log.Info("game init")
}

func (game *GameApp) Run(ctx context.Context) {
	log.Info("game run")
	<-ctx.Done()
}

func (game *GameApp) Shutdown() {
	log.Info("game shutdown")
}

func main() {
	services := plume.WithServices(&GameApp{})
	plume.Run(services)
}