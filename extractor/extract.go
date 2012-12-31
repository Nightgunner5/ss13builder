package main

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"image/png"
	"strings"
	"strconv"
)

func Extract(basedir string, desc io.RuneReader, baseimg image.Image) {
	var (
		state   []rune
		newLine bool
	)

	left := 0

	img := func(w io.Writer, stride, count int) {
		stride *= 32
		l := left
		img := image.NewRGBA(image.Rect(0, 0, count * 32, 32))

		for i := 0; i < count; i++ {
			draw.Draw(img, image.Rect(i*32, 0, i*32+32, 32), baseimg, image.Pt(l, 0), draw.Src)
			l += stride
		}

		png.Encode(w, img)
	}

	extract := func() {
		if state[0] == '#' || (state[0] == 'v' &&
			state[1] == 'e' && state[2] == 'r') {
			state = nil
			return
		}
		s := strings.Split(string(state), "\n\t")
		props := make(map[string]string)

		for _, prop := range s {
			parts := strings.SplitN(prop, " = ", 2)
			props[parts[0]] = parts[1]
		}

		name := props["state"]
		name = name[1:len(name)-1]
		count, _ := strconv.ParseUint(props["dirs"], 10, 64)
		width, _ := strconv.ParseUint(props["frames"], 10, 64)

		_ = width

		for i := uint64(0); i < count; i++ {
			f, _ := os.Create(fmt.Sprintf("%s/%s_dir%d.png", basedir, name, i))
			img(f, int(count), int(width))
			f.Close()
			left += 32
		}
		left += int((width - 1) * count) * 32

		state = nil
	}

	for {
		r, _, err := desc.ReadRune()
		if err != nil {
			return
		}

		switch r {
		default:
			if newLine {
				extract()
			}
			state = append(state, r)

		case '\t':
			if newLine {
				state = append(state, '\n')
				newLine = false
			}
			state = append(state, r)

		case '\n':
			newLine = true
			continue
		}
		newLine = false
	}
}
