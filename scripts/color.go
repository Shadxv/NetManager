package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	var isRunning = true

	for isRunning {
		var input string
		fmt.Print("Hex color: ")
		_, err := fmt.Scanln(&input)
		if err != nil {continue}

		PrintColored(input)
	}
}

func PrintColored(hexColor string) {
	r, g, b := hexToRGB(hexColor)
	xtermColor := rgbToXterm(r, g, b)
	coloredMessage := fmt.Sprintf("\033[38;5;%dmxterm-256colour: %d\033[0m", xtermColor, xtermColor)
	fmt.Println(coloredMessage)
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255 // default to white if invalid
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)

	return int(r), int(g), int(b)
}

// Remove that later and predifine xterm-256colour in model
func rgbToXterm(r, g, b int) int {
	// Check for grayscale
	if r == g && g == b {
		return 232 + (r/10)*10
	}

	// Determine nearest color in the 6x6x6 color cube
	// The 256 color palette has 16 base colors, 6x6x6 colors, and 24 grayscale colors

	// Find closest color
	closest := 0
	minDistance := math.MaxFloat64
	for i := 0; i < 256; i++ {
		color := xtermColorToRGB(i)
		distance := colorDistance(r, g, b, color.r, color.g, color.b)
		if distance < minDistance {
			minDistance = distance
			closest = i
		}
	}

	return closest
}

// Convert xterm-256 color index to RGB
func xtermColorToRGB(index int) (color struct{ r, g, b int }) {
	if index >= 232 {
		// Grayscale colors
		gray := (index - 232) * 10
		color.r, color.g, color.b = gray, gray, gray
	} else if index < 16 {
		// Base colors
		color.r = ((index & 1) * 255) / 1
		color.g = (((index >> 1) & 1) * 255) / 1
		color.b = (((index >> 2) & 1) * 255) / 1
	} else {
		// Colors in 6x6x6 cube
		index -= 16
		color.r = ((index / 36) % 6) * 51
		color.g = ((index / 6) % 6) * 51
		color.b = (index % 6) * 51
	}
	return
}

// Calculate Euclidean distance between two colors
func colorDistance(r1, g1, b1, r2, g2, b2 int) float64 {
	return math.Sqrt(float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)))
}
