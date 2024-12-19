package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	playerSpeed  = 2
	gravity      = 0.5
	jumpStrength = -10
)

type Obstacle struct {
	x, y          float64
	width, height float64
	image         *ebiten.Image
}

type Enemy struct {
	x, y          float64
	width, height float64
	velocityX     float64
	image         *ebiten.Image
}

type Coin struct {
	x, y          float64
	width, height float64
	image         *ebiten.Image
	collected     bool // コインが収集されたかどうか
}

type Game struct {
	playerX     float64
	playerY     float64
	playerVelY  float64
	onGround    bool
	playerImage *ebiten.Image
	obstacles   []*Obstacle
	enemies     []*Enemy
	coins       []*Coin
	score       int  // スコア
	hp          int  // プレイヤーのHP
	isGameOver  bool // ゲームオーバーフラグ
}

func NewObstacle(x, y float64, imagePath string) *Obstacle {
	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatalf("failed to load obstacle image: %v", err)
	}
	return &Obstacle{x: x, y: y, width: float64(img.Bounds().Dx()), height: float64(img.Bounds().Dy()), image: img}
}

func NewEnemy(x, y, velocityX float64, imagePath string) *Enemy {
	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatalf("failed to load enemy image: %v", err)
	}
	return &Enemy{x: x, y: y, width: float64(img.Bounds().Dx()), height: float64(img.Bounds().Dy()), velocityX: velocityX, image: img}
}

func NewCoin(x, y float64, imagePath string) *Coin {
	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatalf("failed to load coin image: %v", err)
	}
	return &Coin{x: x, y: y, width: float64(img.Bounds().Dx()), height: float64(img.Bounds().Dy()), image: img, collected: false}
}

func (g *Game) resetGame() {
	g.playerX = screenWidth / 2
	g.playerY = screenHeight - 50
	g.playerVelY = 0
	g.onGround = true
	g.isGameOver = false
	g.score = 0
}

func (g *Game) handleCollision() {
	g.hp -= 1 // HPを減らす
	log.Println("Player hit! HP:", g.hp)

	if g.hp <= 0 {
		g.isGameOver = true // HPが0になったらゲームオーバー
	}
}

func (g *Game) Update() error {
	if g.isGameOver {
		// ゲームオーバー時の入力処理（リスタート）
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.resetGame()
		}
		return nil
	}

	// プレイヤーの入力処理
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.onGround {
		g.playerVelY = jumpStrength
		g.onGround = false
	}

	// 重力とジャンプ処理
	g.playerVelY += gravity
	g.playerY += g.playerVelY

	// 地面との衝突処理
	if g.playerY >= screenHeight-50 {
		g.playerY = screenHeight - 50
		g.playerVelY = 0
		g.onGround = true
	}

	// 敵の動き
	for _, enemy := range g.enemies {
		enemy.x += enemy.velocityX
	}

	// 衝突判定
	for _, obstacle := range g.obstacles {
		if g.checkCollision(g.playerX, g.playerY, 20, 40, obstacle.x, obstacle.y, obstacle.width, obstacle.height) {
			g.handleCollision()
		}
	}

	for _, enemy := range g.enemies {
		if g.checkCollision(g.playerX, g.playerY, 20, 40, enemy.x, enemy.y, enemy.width, enemy.height) {
			g.handleCollision()
		}
	}

	// プレイヤーとコインの衝突判定
	for _, coin := range g.coins {
		if !coin.collected && g.checkCollision(g.playerX, g.playerY, 20, 40, coin.x, coin.y, coin.width, coin.height) {
			coin.collected = true
			g.score += 10 // スコアを増加
			log.Println("Coin collected! Score:", g.score)
		}
	}

	return nil
}

// 簡易的な衝突判定関数
func (g *Game) checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景の描画
	screen.Fill(color.RGBA{R: 135, G: 206, B: 250, A: 255}) // 空の色

	// 障害物の描画
	for _, obstacle := range g.obstacles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(obstacle.x, obstacle.y)
		screen.DrawImage(obstacle.image, op)
	}

	// 敵の描画
	for _, enemy := range g.enemies {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(enemy.x, enemy.y)
		screen.DrawImage(enemy.image, op)
	}

	// プレイヤーの描画
	if !g.isGameOver {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.playerX, g.playerY)
		screen.DrawImage(g.playerImage, op)
	}

	// コインの描画
	for _, coin := range g.coins {
		if !coin.collected { // 収集済みのコインは描画しない
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(coin.x, coin.y)
			screen.DrawImage(coin.image, op)
		}
	}

	// スコアとHPの表示
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d  HP: %d", g.score, g.hp), 10, 10)

	// ゲームオーバー画面
	if g.isGameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over!", screenWidth/2-50, screenHeight/2-20)
		ebitenutil.DebugPrintAt(screen, "Press R to Restart", screenWidth/2-70, screenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		playerX: screenWidth / 2,
		playerY: screenHeight - 50,
		obstacles: []*Obstacle{
			NewObstacle(200, screenHeight-50, "obstacle.png"),
			NewObstacle(400, screenHeight-50, "obstacle.png"),
		},
		enemies: []*Enemy{
			NewEnemy(300, screenHeight-90, -1, "enemy.png"),
		},
		score: 100,
		hp:    5,
	}

	// プレイヤー画像の読み込み
	var err error
	game.playerImage, _, err = ebitenutil.NewImageFromFile("player.png")
	if err != nil {
		log.Fatalf("failed to load player image: %v", err)
	}

	// コインの生成
	game.coins = []*Coin{
		NewCoin(100, 320, "coin.png"), // 1つ目のコイン
		NewCoin(300, 350, "coin.png"), // 2つ目のコイン
		NewCoin(500, 330, "coin.png"), // 3つ目のコイン
	}

	// ウィンドウ設定
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game with Obstacles and Enemies")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
