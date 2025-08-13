package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameConfig struct {
	Title         string `json:"title"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	FPS           int    `json:"fps"`
	Background    string `json:"background"`
	EntityConfigs []struct {
		Type    string `json:"type"`
		Sprites []string `json:"sprites"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Velocity int    `json:"velocity"`
	} `json:"entities"`
}

type Game struct {
	cfg      GameConfig
	win      *pixelgl.Window
	sprites  map[string]*pixel.Sprite
	entities []*Entity
}

type Entity struct {
	sprite *pixel.Sprite
	x      int
	y      int
	vel    int
}

func (g *Game) Run() {
	rand.Seed(time.Now().UnixNano())
	g.setupWindow()
	g.loadSprites()
	g.loadEntities()
	g.runLoop()
}

func (g *Game) setupWindow() {
	cfg := g.cfg
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  cfg.Title,
		Bounds: pixel.R(cfg.Width, cfg.Height),
		VSync: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.win = win
}

func (g *Game) loadSprites() {
	sprites := make(map[string]*pixel.Sprite)
	for _, spriteCfg := range g.cfg.EntityConfigs {
		for _, spritePath := range spriteCfg.Sprites {
			pic, err := pixelgl.Load.picture(spritePath)
			if err != nil {
				log.Fatal(err)
			}
			sprites[spritePath] = pixel.NewSprite(pic, pixel.R(0, 0, pic.Bounds().Max.X, pic.Bounds().Max.Y))
		}
	}
	g.sprites = sprites
}

func (g *Game) loadEntities() {
	entities := make([]*Entity, len(g.cfg.EntityConfigs))
	for i, entityCfg := range g.cfg.EntityConfigs {
		sprite, ok := g.sprites[entityCfg.Sprites[0]]
		if !ok {
			log.Fatal("sprite not found")
		}
		entities[i] = &Entity{
			sprite: sprite,
			x:      entityCfg.X,
			y:      entityCfg.Y,
			vel:    entityCfg.Velocity,
		}
	}
	g.entities = entities
}

func (g *Game) runLoop() {
	for !g.win.Closed() {
		g.update()
		g.draw()
	}
	g.win.Destroy()
}

func (g *Game) update() {
	for _, entity := range g.entities {
		entity.x += entity.vel
		if entity.x > g.cfg.Width || entity.x < 0 {
			entity.vel = -entity.vel
		}
	}
}

func (g *Game) draw() {
	g.win.Clear(g.cfg.Background)
	for _, entity := range g.entities {
		entity.sprite.Draw(g.win, pixel.IM.Moved(pixel.V(entity.x, entity.y)))
	}
	g.win.Update()
}

func main() {
	jsonCfg := `
	{
		"title": "Auto Game",
		"width": 800,
		"height": 600,
		"fps": 60,
		"background": "#ffffff",
		"entities": [
			{
				"type": "player",
				"sprites": ["player.Sprite.png"],
				"x": 400,
				"y": 300,
				"velocity": 5
			},
			{
				"type": "enemy",
				"sprites": ["enemy.Sprite.png"],
				"x": 700,
				"y": 500,
				"velocity": -3
			}
		]
	}
	`
	var cfg GameConfig
	err := json.Unmarshal([]byte(jsonCfg), &cfg)
	if err != nil {
		log.Fatal(err)
	}
	game := &Game{cfg: cfg}
	game.Run()
}