package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	images map[string]*ebiten.Image
	others others
	mc     mc
	idlist map[int]string
	menu   menus
}

func (g *Game) Update() error {
	// fmt.Println(g.menu.menu)
	if g.menu.menu < g.menu.limit/2 {
		if g.mc.lives > 0 {
			g.generateothers()
			g.updateothersmov()
			g.updatemcmov()
			g.updatebulletmov()
			g.checkCollisions()
		} else {
			return ebiten.Termination
		}
		g.updatemenu()

	} else {
		g.updatemenu()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.menu.menu < g.menu.limit/2 {
		if g.mc.lives > 0 {
			g.drawothers(screen)
			g.drawmc(screen)
			g.drawbullets(screen)
			g.drawighud(screen)
			// fmt.Println(g.mc.lives)

		} else {
			// fmt.Println(g.mc.score)
		}
	} else {
		g.drawpause(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	return 720, 1280
}

func main() {
	ebiten.SetWindowSize(720, 1280)
	ebiten.SetWindowTitle("spacewar")
	game := &Game{images: returnimages("./assets"),
		others: others{others: map[int]other{}, size: size{x: 32, y: 32, scale: 3}, bullets: bullets{bullets: map[int]bullet{}, size: size{x: 32, y: 32, scale: 3}, recharge: reload{state: 0, limit: 200}}, generationcyc: reload{state: 0, limit: 70}},
		mc:     mc{pos: position{x: 100, y: 200}, lives: 3, score: 0, size: size{x: 32, y: 32, scale: 3}, binds: keybinds{up: ebiten.KeyW, right: ebiten.KeyD, left: ebiten.KeyA, down: ebiten.KeyS, attack: ebiten.KeyShiftLeft}, bullets: bullets{bullets: map[int]bullet{}, size: size{x: 32, y: 32, scale: 3}, recharge: reload{state: 0, limit: 7}}},
		idlist: map[int]string{},
		menu:   menus{menu: 0, limit: 12},
	}

	err := ebiten.RunGame(game)
	if err != nil {
		if err == ebiten.Termination {
			fmt.Println(game.mc.score)
			return
		}
		log.Println(err)
	}
}
