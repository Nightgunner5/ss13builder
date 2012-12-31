package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"io"
	"os"
	"sort"
	"strconv"
)

type sortInstances []Instance

func (s sortInstances) Len() int {
	return len(s)
}

func (s sortInstances) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortInstances) Less(i, j int) bool {
	return s[i].Layer() < s[j].Layer()
}

func parseInstance(b []byte) (ins Instance) {
	if i := bytes.IndexByte(b, '{'); i != -1 {
		ins.Path = string(b[:i])
		ins.Extra = string(b[i:])
	} else {
		ins.Path = string(b)
	}
	return
}

func Parse(r io.Reader) (*Map, error) {
	stage := 0
	in := bufio.NewReader(r)
	m := new(Map)
	m.Types = make(map[string]TileType)

	var cz *ZLevel
	var keyLen int

	for {
		line, err := in.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = line[:len(line)-1]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		switch {
		case stage == 0 && len(line) > 8 && line[0] == '"' && line[len(line)-1] == ')':
			for i := 1; i < len(line) && line[i] != '"'; i++ {
				keyLen++
			}
			stage++
			fallthrough
		case stage == 1 && len(line) > 8 && line[0] == '"' && line[len(line)-1] == ')':
			if !(line[keyLen+1] == '"' && line[keyLen+2] == ' ' && line[keyLen+3] == '=' && line[keyLen+4] == ' ' && line[keyLen+5] == '(') {
				return nil, fmt.Errorf("Unparseable line: %s")
			}
			var tt TileType
			tt.Key = string(line[1 : keyLen+1])
			line = line[keyLen+6 : len(line)-1]

			inExtra := false
			for i := 0; i < len(line); i++ {
				if !inExtra && line[i] == ',' {
					tt.Instances = append(tt.Instances, parseInstance(line[0:i]))
					line = line[i+1:]
					i = 0
				}
				if line[i] == '{' {
					inExtra = true
				}
				if line[i] == '}' {
					inExtra = false
				}
			}
			tt.Instances = append(tt.Instances, parseInstance(line))
			sort.Sort(sortInstances(tt.Instances))
			m.Types[tt.Key] = tt

		case stage == 1 && len(line) == 0:
			stage++

		case stage == 4 && len(line) == 0:
			cz = nil

		case (stage == 2 || stage == 4) && len(line) > 11 && line[0] == '(' &&
			line[len(line)-1] == '"' && line[len(line)-2] == '{' &&
			line[len(line)-3] == ' ' && line[len(line)-4] == '=' &&
			line[len(line)-5] == ' ' && line[len(line)-6] == ')':
			stage = 3

			coords := bytes.SplitN(line[1:len(line)-6], []byte{','}, 3)
			x, err := strconv.ParseUint(string(coords[0]), 10, 32)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseUint(string(coords[1]), 10, 32)
			if err != nil {
				return nil, err
			}
			z, err := strconv.ParseUint(string(coords[2]), 10, 32)
			if err != nil {
				return nil, err
			}

			m.ZLevels = append(m.ZLevels, ZLevel{
				Start: Coord{uint32(x), uint32(y), uint32(z)},
			})
			cz = &m.ZLevels[len(m.ZLevels)-1]

		case stage == 3 && len(line) == 2 && line[0] == '"' && line[1] == '}':
			stage++

		case stage == 3 && len(line)%keyLen == 0:
			var row []string
			for i := 0; len(line) != 0; i, line = i+1, line[keyLen:] {
				key := m.Types[string(line[:keyLen])].Key
				row = append(row, key)
			}
			cz.Map = append(cz.Map, row)

		default:
			return nil, fmt.Errorf("Unparsed line:%d: %q", stage, line)
		}
	}

	if stage != 4 {
		return nil, fmt.Errorf("Unexpected EOF: %d", stage)
	}

	return m, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "nameofmap.dmm")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening map: ", err)
		return
	}
	m, err := Parse(f)
	f.Close()
	if err != nil {
		fmt.Println("Error parsing map: ", err)
		return
	}

	images, err := Render(m)
	if err != nil {
		fmt.Println("Error rendering map: ", err)
		return
	}

	for i, img := range images {
		f, err := os.Create(fmt.Sprintf("%s-zl%d.png", os.Args[1], i+1))
		if err != nil {
			fmt.Println("Error saving zlevels: ", err)
			return
		}
		png.Encode(f, img)
		f.Close()
	}
	fmt.Println("It all worked!")
}
