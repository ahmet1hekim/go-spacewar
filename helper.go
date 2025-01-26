package main

import (
	"log"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type menus struct {
	menu  int
	limit int
}
type keybinds struct {
	up     ebiten.Key
	down   ebiten.Key
	left   ebiten.Key
	right  ebiten.Key
	attack ebiten.Key
}
type position struct {
	x int
	y int
}

type mc struct {
	pos     position
	binds   keybinds
	bullets bullets
	size    size
	lives   int
	score   int
}
type size struct {
	x     int
	y     int
	scale int
}

type other struct {
	pos position
}
type others struct {
	others        map[int]other
	bullets       bullets
	size          size
	generationcyc reload
}

type bullet struct {
	pos  position
	size size
}
type reload struct {
	state int
	limit int
}

func (r *reload) cycle() {
	r.state = (r.state + 1) % r.limit
}
func (r *menus) cycle() {
	r.menu = (r.menu + 1) % r.limit
}

type bullets struct {
	bullets  map[int]bullet
	size     size
	recharge reload
}

func returnimages(assetspath string) map[string]*ebiten.Image {

	assetlist, err := os.ReadDir(assetspath)

	if err != nil {
		log.Println(err)
	}

	images := make(map[string](*ebiten.Image))

	for _, path := range assetlist {
		images[path.Name()], _, err = ebitenutil.NewImageFromFile(assetspath + "/" + path.Name())

		if err != nil {
			log.Println(err)
		}
	}
	return images
}

func (g *Game) generate_new_id(from string) int {
	for {
		randid := rand.IntN(1000000)
		if _, ok := g.idlist[randid]; !ok {

			g.idlist[randid] = from
			return randid
		}
	}
}
func (g *Game) generateothers() {
	if g.others.generationcyc.state == 7 {
		newother := other{
			pos: position{
				x: rand.IntN(240),
				y: -rand.IntN(150)}}
		g.others.others[g.generate_new_id("other")] = newother
	}
	g.others.generationcyc.cycle()
}

func (g *Game) drawothers(screen *ebiten.Image) {
	tempop := ebiten.DrawImageOptions{}

	for key, _ := range g.others.others {
		tempop.GeoM.Translate(float64(g.others.others[key].pos.x)-float64(g.others.size.scale)*float64(g.others.size.x)/2, float64(g.others.others[key].pos.y)-float64(g.others.size.scale)*float64(g.others.size.y)/2)
		tempop.GeoM.Scale(float64(g.others.size.scale), float64(g.others.size.scale))
		screen.DrawImage(g.images["other.png"], &tempop)
		tempop.GeoM.Reset()
	}

}
func (g *Game) drawmc(screen *ebiten.Image) {
	tempop := ebiten.DrawImageOptions{}
	tempop.GeoM.Translate(float64(g.mc.pos.x)-(float64(g.mc.size.scale)*float64(g.mc.size.x)/2), float64(g.mc.pos.y)-(float64(g.mc.size.scale)*float64(g.mc.size.y)/2))
	tempop.GeoM.Scale(float64(g.mc.size.scale), float64(g.mc.size.scale))
	screen.DrawImage(g.images["player.png"], &tempop)
	// fmt.Println(float64(g.mc.pos.x)-(float64(g.mc.size.scale)*float64(g.mc.size.x)/2), float64(g.mc.pos.y)-(float64(g.mc.size.scale)*float64(g.mc.size.y)/2))
}

func (g *Game) drawbullets(screen *ebiten.Image) {
	tempop := ebiten.DrawImageOptions{}

	for key, _ := range g.others.bullets.bullets {
		tempop.GeoM.Translate(float64(g.others.bullets.bullets[key].pos.x)-float64(g.others.bullets.size.scale)*float64(g.others.size.x)/2, float64(g.others.bullets.bullets[key].pos.y)-float64(g.others.bullets.size.scale)*float64(g.others.size.y)/2)
		tempop.GeoM.Scale(float64(g.others.bullets.size.scale), float64(g.others.bullets.size.scale))
		screen.DrawImage(g.images["otherbullet.png"], &tempop)
		tempop.GeoM.Reset()
	}
	for key, _ := range g.mc.bullets.bullets {
		tempop.GeoM.Translate(float64(g.mc.bullets.bullets[key].pos.x)-float64(g.mc.bullets.size.scale)*float64(g.mc.bullets.size.x)/2, float64(g.mc.bullets.bullets[key].pos.y)-float64(g.mc.bullets.size.scale)*float64(g.mc.bullets.size.y)/2)
		tempop.GeoM.Scale(float64(g.mc.bullets.size.scale), float64(g.mc.bullets.size.scale))
		screen.DrawImage(g.images["bullet.png"], &tempop)
		tempop.GeoM.Reset()
	}
}

func (g *Game) updateothersmov() {
	for key, value := range g.others.others {
		value.pos.y += 1
		g.others.others[key] = value
		if g.others.bullets.recharge.state == 1 {
			bullettoadd := bullet{pos: position{x: g.others.others[key].pos.x, y: g.others.others[key].pos.y}}
			g.others.bullets.bullets[g.generate_new_id("otherbullet")] = bullettoadd
			g.others.bullets.recharge.cycle()

		} else {
			g.others.bullets.recharge.cycle()
		}

	}
}

func (g *Game) updatemcmov() {

	if ebiten.IsKeyPressed(g.mc.binds.up) {
		g.mc.pos.y -= 2
	}
	if ebiten.IsKeyPressed(g.mc.binds.down) {
		g.mc.pos.y += 2
	}
	if ebiten.IsKeyPressed(g.mc.binds.left) {
		g.mc.pos.x -= 2
	}
	if ebiten.IsKeyPressed(g.mc.binds.right) {
		g.mc.pos.x += 2
	}
	if ebiten.IsKeyPressed(g.mc.binds.attack) {
		// fmt.Println(g.mc.bullets.recharge.state)
		if g.mc.bullets.recharge.state == 1 {
			bullettoadd := bullet{pos: position{x: g.mc.pos.x, y: g.mc.pos.y}}
			g.mc.bullets.bullets[g.generate_new_id("mcbullet")] = bullettoadd
			g.mc.bullets.recharge.cycle()
		} else {
			g.mc.bullets.recharge.cycle()
		}

	}
	// fmt.Println(g.mc.pos)
}
func (g *Game) updatebulletmov() {
	for key, value := range g.others.bullets.bullets {
		value.pos.y += 2
		g.others.bullets.bullets[key] = value
	}

	for key, value := range g.mc.bullets.bullets {
		value.pos.y -= 1
		g.mc.bullets.bullets[key] = value
	}
}
func checkCollision(a, b position, aSize, bSize size, scale int) int {
	aWidth := float32(aSize.x) * 0.5
	aHeight := float32(aSize.y) * 0.5
	bWidth := bSize.x
	bHeight := bSize.y

	aLeft := float32(a.x) - aWidth/2
	aRight := float32(a.x) + aWidth/2
	aTop := float32(a.y) - aHeight/2
	aBottom := float32(a.y) + aHeight/2

	bLeft := b.x - bWidth/2
	bRight := b.x + bWidth/2
	bTop := b.y - bHeight/2
	bBottom := b.y + bHeight/2

	// Debug: Print bounding box values to check for correct scaling and positioning
	// Draw bounding boxes on screen for debugging
	// Draw bounding box for Object A (red), scaling the position for visualization only
	// vector.DrawFilledRect(screen, float32(aLeft)*float32(scale), float32(aTop)*float32(scale), float32(aWidth)*float32(scale), float32(aHeight)*float32(scale), color.RGBA{255, 0, 0, 60}, true)

	// Draw bounding box for Object B (green), scaling the position for visualization only
	// vector.DrawFilledRect(screen, float32(bLeft)*float32(scale), float32(bTop)*float32(scale), float32(bWidth)*float32(scale), float32(bHeight)*float32(scale), color.RGBA{0, 255, 0, 60}, true)

	// Check for collision (adjusted to inclusive edges)
	if aRight >= float32(bLeft) && aLeft <= float32(bRight) && aBottom >= float32(bTop) && aTop <= float32(bBottom) {
		return 1
	}
	if bLeft < 32 || bRight > 272 || bBottom > 390 {
		return 2
	}
	return 0
}

func (g *Game) checkCollisions() {
	for key, otherObj := range g.others.others {
		j := checkCollision(g.mc.pos, otherObj.pos, g.mc.size, g.others.size, g.mc.size.scale)
		if j != 0 {
			delete(g.others.others, key)
			if j == 1 {
				g.mc.lives -= 1
			}
		}
	}
	i := size{scale: 3, x: 8, y: 8}
	for bulkey, bullet := range g.mc.bullets.bullets {
		for key, otherObj := range g.others.others {
			if checkCollision(bullet.pos, otherObj.pos, i, g.others.size, g.mc.bullets.size.scale) == 1 {
				delete(g.others.others, key)
				delete(g.mc.bullets.bullets, bulkey)
				g.mc.score += 1
			}
		}
	}
	for bulkey, bullet := range g.others.bullets.bullets {

		if checkCollision(bullet.pos, g.mc.pos, i, g.others.size, g.mc.bullets.size.scale) == 1 {
			// fmt.Println("Bullet hit an object!")
			delete(g.others.bullets.bullets, bulkey)
			g.mc.lives -= 1

		}
	}
	// fmt.Println(" ")
}
func (g *Game) drawighud(screen *ebiten.Image) {
	tempop := ebiten.DrawImageOptions{}
	for i := 0; i < g.mc.lives; i++ {
		tempop.GeoM.Translate(float64(32*i), float64(32))
		tempop.GeoM.Scale(float64(2), float64(2))
		screen.DrawImage(g.images["life.png"], &tempop)
		tempop.GeoM.Reset()
	}
	for i := 0; i < 3-g.mc.lives; i++ {
		tempop.GeoM.Translate(float64(64-32*i), float64(32))
		tempop.GeoM.Scale(float64(2), float64(2))
		screen.DrawImage(g.images["lostlife.png"], &tempop)
		tempop.GeoM.Reset()
	}
	tempop.GeoM.Translate(float64(0), float64(32))
	tempop.GeoM.Scale(float64(2), float64(2))
	screen.DrawImage(g.images["frame.png"], &tempop)
	tempop.GeoM.Reset()

	str := strconv.Itoa(g.mc.score)
	for i, char := range str {
		// fmt.Println(string(char) + ".png")
		tempop.GeoM.Translate(float64(350-32*(len(str)-i)), float64(32))
		tempop.GeoM.Scale(float64(2), float64(2))
		screen.DrawImage(g.images[string(char)+".png"], &tempop)
		tempop.GeoM.Reset()
	}
}

func (g *Game) drawpause(screen *ebiten.Image) {
	tempop := ebiten.DrawImageOptions{}
	tempop.GeoM.Translate(float64(64), float64(64))
	tempop.GeoM.Scale(float64(2), float64(2))
	screen.DrawImage(g.images["paused.png"], &tempop)
	tempop.GeoM.Reset()

	str := strconv.Itoa(g.mc.score)
	for i, char := range str {
		// fmt.Println(string(char) + ".png")
		tempop.GeoM.Translate(float64(350-32*(len(str)-i)), float64(32))
		tempop.GeoM.Scale(float64(2), float64(2))
		screen.DrawImage(g.images[string(char)+".png"], &tempop)
		tempop.GeoM.Reset()
	}

}
func (g *Game) updatemenu() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.menu.cycle()
	}

}
