package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/itiky/goPixelWorld/engine"
	"github.com/itiky/goPixelWorld/monitor"
	"github.com/itiky/goPixelWorld/world"
	"github.com/itiky/goPixelWorld/world/materials"
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	monitorKeeper, err := monitor.NewKeeper(10 * time.Second)
	if err != nil {
		log.Fatalf("monitor.NewKeeper: %v", err)
	}

	materialsAll := []worldTypes.MaterialI{
		materials.SandM,
		materials.WaterM,
		materials.WoodM,
		materials.FireM,
		materials.GrassM,
		materials.SmokeM,
		materials.SteamM,
		materials.MetalM,
		materials.RockM,
		materials.GravitonM,
	}

	worldMap, err := world.NewMap(
		world.WithWidth(250),
		world.WithHeight(250),
		//world.WithMonitor(monitorKeeper),
	)
	if err != nil {
		log.Fatalf("creating world.Map: %v", err)
	}

	runner, err := engine.NewRunner(
		worldMap,
		engine.WithScreenSize(1500, 1100),
		engine.WithEditorUI(materialsAll...),
		//engine.WithMonitor(monitorKeeper),
	)
	if err != nil {
		log.Fatalf("creating engine.Runner: %v", err)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	monitorKeeper.Start(ctx)

	if err := runner.Run(); err != nil {
		log.Fatalf("running engine.Runner: %v", err)
	}
}
