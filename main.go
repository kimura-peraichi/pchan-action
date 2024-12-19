package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const (
	screenWidth  = 640
	screenHeight = 480
	playerSpeed  = 2
	gravity      = 0.5
	jumpStrength = -10
)

type Game struct {
	playerX      float64
	playerY      float64
	playerVelY   float64
	onGround     bool
	backgroundX  float64
}

func (g *Game) Update() error {
	// プレイヤーの入力処理
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += playerSpeed
		g.backgroundX -= playerSpeed / 2 // 背景スクロール
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= playerSpeed
		g.backgroundX += playerSpeed / 2
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.onGround {
		g.playerVelY = jumpStrength
		g.onGround = false
	}

	// 重力とジャンプ処理
	g.playerVelY += gravity
	g.playerY += g.playerVelY

	// 地面との衝突処理
	if g.playerY >= screenHeight-50 { // 地面の高さを50と仮定
		g.playerY = screenHeight - 50
		g.playerVelY = 0
		g.onGround = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景の描画
	screen.Fill(color.RGBA{R: 135, G: 206, B: 250, A: 255}) // 空の色

	// プレイヤーの描画
	playerRect := ebiten.NewImage(20, 40)
	playerRect.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(playerRect, op)

	// デバッグ情報の表示
	ebitenutil.DebugPrint(screen, "Use Arrow Keys to Move, Space to Jump")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		playerX: screenWidth / 2,
		playerY: screenHeight - 50,
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Side Scrolling Action Game")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

