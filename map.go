package main

import (
	"strings"
)

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

func (i Instance) ParseExtra() map[string]string {
	if i.Extra == "" {
		return nil
	}
	m := make(map[string]string)
	extra := strings.Split(i.Extra[1:len(i.Extra)-1], "; ")
	for _, prop := range extra {
		pieces := strings.SplitN(prop, " = ", 2)
		m[pieces[0]] = m[pieces[1]]
	}
	return m
}

func (i Instance) SpritePath() string {
	extra := i.ParseExtra()
	dir, icon_state := extra["dir"], extra["icon_state"]
	if dir == "" {
		dir = "0"
	}
	if icon_state == "" {
		icon_state = "\"default\""
	}
	icon_state = icon_state[1 : len(icon_state)-1]
	return "tiles" + i.Path + "/" + icon_state + "_" + dir + ".png"
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
