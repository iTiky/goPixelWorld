package world

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/itiky/goPixelWorld/world/materials"
)

func (m *Map) SetImageData(imageData image.Image) error {
	imageData = m.resizeImage(200, 200, imageData)

	imageBounds := imageData.Bounds()
	imageWidth := imageBounds.Max.X - imageBounds.Min.X
	imageHeight := imageBounds.Max.Y - imageBounds.Min.Y

	m.initGrid(imageWidth+2, imageHeight+2)

	mapXOffset := 1 - imageBounds.Min.X
	mapYOffset := 1 - imageBounds.Min.Y

	for imageX := imageBounds.Min.X; imageX < imageBounds.Max.X; imageX++ {
		for imageY := imageBounds.Min.Y; imageY < imageBounds.Max.Y; imageY++ {
			material := materials.FindClosestMaterialByColor(imageData.At(imageX, imageY))
			if material == nil {
				continue
			}

			tile := m.getTile(imageX+mapXOffset, imageY+mapYOffset)
			m.createParticle(tile, material)
		}
	}

	return nil
}

func (m *Map) resizeImage(widthMax, heightMax int, srcImage image.Image) image.Image {
	srcImageWidth := srcImage.Bounds().Max.X - srcImage.Bounds().Min.X
	srcImageHeight := srcImage.Bounds().Max.Y - srcImage.Bounds().Min.Y
	if srcImageWidth <= widthMax && srcImageHeight <= heightMax {
		return srcImage
	}

	widthScaleK := float64(widthMax) / float64(srcImageWidth)
	heightScaleK := float64(heightMax) / float64(srcImageHeight)

	dstImageWidth := int(float64(srcImageWidth) * widthScaleK)
	dstImageHeight := int(float64(srcImageHeight) * heightScaleK)

	dstImage := image.NewRGBA(image.Rect(0, 0, dstImageWidth, dstImageHeight))
	draw.NearestNeighbor.Scale(dstImage, dstImage.Rect, srcImage, srcImage.Bounds(), draw.Over, nil)

	return dstImage
}
