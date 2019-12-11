package main

import (
	"io/ioutil"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Layer []byte

func main() {
	width, height := 25, 6
	layers := parseInput("input.txt", width, height)

	// Part 1
	fewestZeroLayer := layers[0]
	fewestZeroCount := count(layers[0], 0)
	for _, layer := range layers[1:] {
		if c := count(layer, 0); c < fewestZeroCount {
			fewestZeroCount = c
			fewestZeroLayer = layer
		}
	}
	fmt.Println("Part 1", count(fewestZeroLayer, 1) * count(fewestZeroLayer, 2))

	// Part 2
	part2File := "part2.png"
	part2 := createImage(layers, width, height, part2File)
	fmt.Println("Part 2 image saved to", part2File)
	fmt.Println(part2)

}

func createImage(layers []Layer, width, height int, filename string) string {
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	file, _ := os.Create(filename)
	img := make([][]byte, height)

	for y := 0; y < height; y++ {
		img[y] = make([]byte, width)
		for x := 0; x < width; x++ {
			for _, layer := range layers {
				if layer[y*width+x] == 0 {
					m.Set(x, y, color.Black)
					img[y][x] = '0'
					break
				}
				if layer[y*width+x] == 1 {
					m.Set(x, y, color.White)
					img[y][x] = '1'
					break
				}
				m.Set(x, y, color.Transparent)
				img[y][x] = ' '
			}
		}
	}

	png.Encode(file, m)
	imgString := ""
	for i := range img {
		imgString += string(img[i]) + "\n"
	}

	return imgString
}

func count(layer Layer, data byte) int {
	c := 0
	for i := range layer {
		if layer[i] == data {
			c++
		}
	}
	return c
}

func parseInput(filename string, width, height int) []Layer {
	data, _ := ioutil.ReadFile(filename)
	layerLen := width*height
	layers := make([]Layer, 0)
	for i := 0; i < len(data); i += layerLen {
		layer := make(Layer, layerLen)
		for j := 0; j < layerLen; j++ {
			layer[j] = data[i+j] - '0'
		}
		layers = append(layers, layer)
	}
	return layers
}
