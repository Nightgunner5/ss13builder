package main

type Map struct {
	Types   map[string]TileType
	ZLevels []ZLevel
}

type Coord struct {
	X, Y, Z uint32
}

type ZLevel struct {
	Start Coord
	Map   [][]string
}

type Instance struct {
	Path  string
	Extra string
}

func (i Instance) String() string {
	return i.Path + i.Extra
}

type TileType struct {
	Key       string
	Instances []Instance
}

func (tt TileType) String() string {
	var buf []byte
	buf = append(buf, '"')
	buf = append(buf, tt.Key...)
	buf = append(buf, '"', ' ', '=', ' ', '(')

	first := true
	for _, i := range tt.Instances {
		if first {
			first = false
		} else {
			buf = append(buf, ',')
		}
		buf = append(buf, i.String()...)
	}
	buf = append(buf, '"', '\r', '\n')
	return string(buf)
}
