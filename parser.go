package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

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
		case stage == 0 && len(line) > 8 && line[0] == '"' &&
			line[4] == '"' && line[5] == ' ' && line[6] == '=' &&
			line[7] == ' ' && line[8] == '(' &&
			line[len(line)-1] == ')':
			var tt TileType
			tt.Key = string(line[1:4])
			line = line[9 : len(line)-1]

			for i := 0; i < len(line); i++ {
				if line[i] == ',' {
					tt.Instances = append(tt.Instances, parseInstance(line[0:i]))
					line = line[i+1:]
					i = 0
				}
			}
			tt.Instances = append(tt.Instances, parseInstance(line))
			m.Types[tt.Key] = tt

		case stage == 0 && len(line) == 0:
			stage++

		case stage == 3 && len(line) == 0:
			cz = nil

		case (stage == 1 || stage == 3) && len(line) > 11 && line[0] == '(' &&
			line[len(line)-1] == '"' && line[len(line)-2] == '{' &&
			line[len(line)-3] == ' ' && line[len(line)-4] == '=' &&
			line[len(line)-5] == ' ' && line[len(line)-6] == ')':
			stage = 2

			coords := bytes.SplitN(line[1:len(line) - 6], []byte{','}, 3)
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

		case stage == 2 && len(line) != 0 && len(line) % 3 == 0:
			var row []string
			for i := 0; len(line) != 0; i, line = i + 1, line[3:] {
				key := m.Types[string(line[:3])].Key
				row = append(row, key)
			}
			cz.Map = append(cz.Map, row)

		case stage == 2 && len(line) == 2 && line[0] == '"' && line[1] == '}':
			stage++

		default:
			return nil, fmt.Errorf("Unparsed line: %q", line)
		}
	}

	if stage != 3 {
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
	defer f.Close()
	_, err = Parse(f)
	if err != nil {
		fmt.Println("Error parsing map: ", err)
		return
	}
	fmt.Println("It's all working!")

}
