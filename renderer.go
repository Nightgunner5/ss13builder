package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

func Render(m *Map) (images []*image.RGBA, err error) {
	tileImg := make(map[Instance]image.Image)
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
					if strings.HasPrefix(ins.Path, "/area") {
						continue
					}
					if ti, ok := tileImg[ins]; ok {
						draw.Draw(img, image.Rect(x*32, y*32, x*32+32, y*32+32), ti, image.ZP, draw.Over)
					} else {
						e := fmt.Errorf("No tile image: (%d, %d, %d) %q", uint32(x)+zl.Start.X, uint32(y)+zl.Start.Y, zl.Start.Z, ins)
						f, err := os.Open("tiles" + ins.String() + ".png")
						if err != nil {
							//return nil, e
							fmt.Println(e)
							continue
						}
						ti, err := png.Decode(f)
						f.Close()
						if err != nil {
							return nil, e
						}
						tileImg[ins] = ti
						draw.Draw(img, image.Rect(x*32, y*32, x*32+32, y*32+32), ti, image.ZP, draw.Over)
					}
				}
			}
		}
	}

	return
}
