package services

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"

	"github.com/sirupsen/logrus"

	"face-recognition-api/internal/models"
)

// ImageProcessor handles image processing operations like drawing circles and encoding
type ImageProcessor struct {
	logger *logrus.Logger
}

// CircleOptions defines options for drawing circles on images
type CircleOptions struct {
	Color     color.RGBA
	LineWidth int
}

// NewImageProcessor creates a new image processor instance
func NewImageProcessor(logger *logrus.Logger) *ImageProcessor {
	return &ImageProcessor{
		logger: logger,
	}
}

// DrawFaceCircles draws circles around detected faces and returns base64 encoded image
func (ip *ImageProcessor) DrawFaceCircles(img image.Image, faces []models.Face, opts CircleOptions) (string, error) {
	// Create a new RGBA image from the original
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Draw circles around detected faces
	for _, face := range faces {
		centerX := face.X + face.Width/2
		centerY := face.Y + face.Height/2
		radius := int(math.Max(float64(face.Width), float64(face.Height)) / 2)

		ip.drawCircle(rgba, centerX, centerY, radius, opts.Color, opts.LineWidth)
	}

	// Encode to base64
	encoded, err := ip.encodeToBase64(rgba)
	if err != nil {
		return "", err
	}

	ip.logger.WithFields(logrus.Fields{
		"faces_processed": len(faces),
		"circle_color":    opts.Color,
		"line_width":      opts.LineWidth,
	}).Info("Face circles drawn successfully")

	return encoded, nil
}

// drawCircle draws a circle using Bresenham's circle algorithm with line width
func (ip *ImageProcessor) drawCircle(img *image.RGBA, centerX, centerY, radius int, col color.RGBA, lineWidth int) {
	// Bresenham's circle algorithm with line width
	for w := 0; w < lineWidth; w++ {
		r := radius + w - lineWidth/2
		if r <= 0 {
			continue
		}

		x := 0
		y := r
		d := 3 - 2*r

		for x <= y {
			// Draw 8 symmetric points
			ip.setPixelSafe(img, centerX+x, centerY+y, col)
			ip.setPixelSafe(img, centerX-x, centerY+y, col)
			ip.setPixelSafe(img, centerX+x, centerY-y, col)
			ip.setPixelSafe(img, centerX-x, centerY-y, col)
			ip.setPixelSafe(img, centerX+y, centerY+x, col)
			ip.setPixelSafe(img, centerX-y, centerY+x, col)
			ip.setPixelSafe(img, centerX+y, centerY-x, col)
			ip.setPixelSafe(img, centerX-y, centerY-x, col)

			if d < 0 {
				d = d + 4*x + 6
			} else {
				d = d + 4*(x-y) + 10
				y--
			}
			x++
		}
	}
}

// setPixelSafe safely sets a pixel within image bounds
func (ip *ImageProcessor) setPixelSafe(img *image.RGBA, x, y int, col color.RGBA) {
	bounds := img.Bounds()
	if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
		img.Set(x, y, col)
	}
}

// encodeToBase64 encodes an image to base64 with data URL prefix
func (ip *ImageProcessor) encodeToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer

	// Encode as JPEG with high quality
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return "", err
	}

	// Convert to base64 with data URL prefix
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/jpeg;base64," + encoded, nil
}

// ParseColor converts color name to RGBA color
func (ip *ImageProcessor) ParseColor(colorName string) color.RGBA {
	switch colorName {
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "yellow":
		return color.RGBA{255, 255, 0, 255}
	case "white":
		return color.RGBA{255, 255, 255, 255}
	case "black":
		return color.RGBA{0, 0, 0, 255}
	case "orange":
		return color.RGBA{255, 165, 0, 255}
	case "purple":
		return color.RGBA{128, 0, 128, 255}
	case "pink":
		return color.RGBA{255, 192, 203, 255}
	case "cyan":
		return color.RGBA{0, 255, 255, 255}
	default:
		return color.RGBA{255, 0, 0, 255} // Default to red
	}
}