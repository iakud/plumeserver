package main

import (
	"github.com/iakud/plume"
)

type GameApp struct {
}

func (game *GameApp) Init() {

}

func (game *GameApp) Destory() {

}

func main() {
	plume.Run(&GameApp{})
}
