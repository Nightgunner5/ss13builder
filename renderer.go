package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

func Render(m *Map) (images []*image.RGBA, err error) {
	tileImg := make(map[string]image.Image)
	missingTile := image.NewUniform(color.RGBA{255, 0, 255, 255})
	for _, zl := range m.ZLevels {
		width, height := 0, len(zl.Map)*32
		for _, row := range zl.Map {
			if w := len(row) * 32; w > width {
				width = w
			}
		}
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		images = append(images, img)
		for y, row := range zl.Map {
			for x, t := range row {
				tt := m.Types[t]
				for _, ins := range tt.Instances {
					if strings.HasPrefix(ins.Path, "/area") || strings.HasPrefix(ins.Path, "/obj/landmark") {
						continue
					}
					path := ins.SpritePath()
					if ti, ok := tileImg[path]; ok {
						draw.Draw(img, image.Rect(x*32, y*32, x*32+32, y*32+32), ti, image.ZP, draw.Over)
					} else {
						e := fmt.Errorf("No tile image: (%d, %d, %d) %q", uint32(x)+zl.Start.X, uint32(y)+zl.Start.Y, zl.Start.Z, path)
						f, err := os.Open(path)
						if err != nil {
							fmt.Println(e)
							tileImg[path] = missingTile
							draw.Draw(img, image.Rect(x*32, y*32, x*32+32, y*32+32), missingTile, image.ZP, draw.Src)
							continue
						}
						ti, err := png.Decode(f)
						f.Close()
						if err != nil {
							return nil, e
						}
						tileImg[path] = ti
						draw.Draw(img, image.Rect(x*32, y*32, x*32+32, y*32+32), ti, image.ZP, draw.Over)
					}
				}
			}
		}
	}

	return
}
