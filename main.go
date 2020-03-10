package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type binData struct {
	S  uint8
	DX int8
	DY int8
}

type Mouse struct {
	S     uint8
	DX    int64
	DY    int64
	X     int64
	Y     int64
	Left  bool
	Right bool
	Name  string
}

func (m Mouse) String() string {
	return fmt.Sprintf("S:%08b, DX:%4v, DY:%4v, X:%6v, Y:%6v  Name:%20v", m.S, m.DX, m.DY, m.X, m.Y, m.Name)
}

func Follow(rc chan Mouse, dev string) {
	fp, err := os.Open(dev)
	m := Mouse{Name: dev}

	if err != nil {
		panic(err)
	}

	defer fp.Close()

	for {
		thing := binData{}
		err := binary.Read(fp, binary.LittleEndian, &thing)

		if err == io.EOF {
			break
		}
		m.DX = int64(thing.DX)
		m.DY = int64(thing.DY)
		m.X += m.DX
		m.Y += m.DY
		m.S = thing.S
		rc <- m
	}
}

func main() {
	ch := make(chan Mouse)

	ms, _ := filepath.Glob("/dev/input/mouse*")
	go Follow(ch, "/dev/input/mice")

	for _, path := range ms {
		go Follow(ch, path)
	}

	all := make(map[string]Mouse)
	var keys []string
	var ok bool

	for m := range ch {
		fmt.Printf("\r\033[%vA", len(all))

		if _, ok = all[m.Name]; !ok {
			keys = append(keys, m.Name)
			sort.Strings(keys)
		}

		all[m.Name] = m

		for _, key := range keys {
			val, _ := all[key]
			fmt.Printf("\n-%25v  %v    ", key, val)
		}
	}
}
