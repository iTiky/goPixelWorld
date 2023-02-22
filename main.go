package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
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

	worldMap, err := world.NewMap(
		world.WithWidth(250),
		world.WithHeight(250),
		//world.WithMonitor(monitorKeeper),
	)
	if err != nil {
		log.Fatalf("creating world.Map: %v", err)
	}

	imageData, err := parseImage()
	if err != nil {
		log.Fatalf("parsing image: %v", err)
	}
	if imageData != nil {
		if err := worldMap.SetImageData(imageData); err != nil {
			log.Fatalf("setting image data: %v", err)
		}
	}

	materialsAll := []worldTypes.MaterialI{
		materials.NewSand(),
		materials.NewWater(),
		materials.NewWood(),
		materials.NewFire(),
		materials.NewGrass(),
		materials.NewSmoke(),
		materials.NewSteam(),
		materials.NewMetal(),
		materials.NewRock(),
		materials.NewGraviton(),
		materials.NewAntiGraviton(),
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

func parseImage() (image.Image, error) {
	var filePath string

	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	if filePath == "" {
		return nil, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	imageData, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	return imageData, nil
}
