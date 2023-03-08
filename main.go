package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/itiky/goPixelWorld/engine"
	"github.com/itiky/goPixelWorld/world"
	"github.com/itiky/goPixelWorld/world/materials"
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//monitorKeeper, err := monitor.NewKeeper(10 * time.Second)
	//if err != nil {
	//	log.Fatalf("monitor.NewKeeper: %v", err)
	//}

	worldMap, err := world.NewMap(
		world.WithWidth(250),
		world.WithHeight(250),
		world.WithNatureEffects(),
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
		materials.NewSand(),         // 1
		materials.NewWater(),        // 2
		materials.NewWood(),         // 3
		materials.NewGrass(),        // 4
		materials.NewFire(),         // 5
		materials.NewRock(),         // 6
		materials.NewMetal(),        // 7
		materials.NewBug(),          // 8
		materials.NewGraviton(),     // 9
		materials.NewAntiGraviton(), // 0
		materials.NewSmoke(),
		materials.NewSteam(),
	}

	runner, err := engine.NewRunner(
		worldMap,
		engine.WithScreenSize(1200, 1100),
		engine.WithEditorUI(materialsAll...),
		//engine.WithMonitor(monitorKeeper),
	)
	if err != nil {
		log.Fatalf("creating engine.Runner: %v", err)
	}

	//ctx, ctxCancel := context.WithCancel(context.Background())
	//defer ctxCancel()
	//monitorKeeper.Start(ctx)

	if err := runner.Run(); err != nil {
		log.Fatalf("running engine.Runner: %v", err)
	}
}

// parseImage parses an image by path (if the corresponding CLI args is provided).
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
