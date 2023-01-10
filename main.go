package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/itiky/goPixelWorld/engine"
	"github.com/itiky/goPixelWorld/world"
	"github.com/itiky/goPixelWorld/world/materials"
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	materialsAll := []worldTypes.MaterialI{
		materials.NewSand(),
		materials.NewWater(),
		materials.NewWood(),
		materials.NewFire(),
		materials.NewGrass(),
		materials.NewSmoke(),
		materials.NewSteam(),
	}

	worldMap, err := world.NewMap(200, 200)
	if err != nil {
		log.Fatalf("creating world.Map: %v", err)
	}

	runner, err := engine.NewRunner(
		worldMap,
		engine.WithScreenSize(1500, 1100),
		engine.WithEditorUI(materialsAll...),
	)
	if err != nil {
		log.Fatalf("creating engine.Runner: %v", err)
	}

	if err := runner.Run(); err != nil {
		log.Fatalf("running engine.Runner: %v", err)
	}
}
