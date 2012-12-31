package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"io/ioutil"
)

var (
	ErrInvalidMagic = errors.New("dmi: invalid magic for PNG")
)

func ReadInfo(r io.Reader) (string, error) {
	magic := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	buf := make([]byte, len(magic))
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	if !bytes.Equal(magic, buf) {
		return "", ErrInvalidMagic
	}

	four := make([]byte, 4)
	for {
		if _, err := io.ReadFull(r, four); err != nil {
			return "", err
		}
		buf = make([]byte, uint(four[0])<<24|uint(four[1])<<16|uint(four[2])<<8|uint(four[3]))
		if _, err := io.ReadFull(r, four); err != nil {
			return "", err
		}
		if _, err := io.ReadFull(r, buf); err != nil {
			return "", err
		}
		switch string(four) {
		case "zTXt":
			i := bytes.IndexByte(buf, 0)
			name := string(buf[:i])
			if name != "Description" {
				break
			}
			z, err := zlib.NewReader(bytes.NewReader(buf[i+2:]))
			if err != nil {
				return "", err
			}
			b, err := ioutil.ReadAll(z)
			if err != nil {
				return "", err
			}
			z.Close()

			return string(b), nil
		}
		if _, err := io.ReadFull(r, four); err != nil {
			return "", err
		}
	}
	panic("unreachable")
}
