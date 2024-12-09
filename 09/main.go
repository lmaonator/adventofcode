package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type blkType int

const (
	Free blkType = iota
	File
)

type Blk struct {
	Type     blkType
	Id       int
	FileSize int
}

type Disc []*Blk

func (d Disc) String() string {
	str := ""
	for _, b := range d {
		if b.Type == Free {
			str += "."
		} else {
			str += strconv.Itoa(b.Id)
		}
	}
	return str
}

func (d Disc) moveBlockLeft(blkIndex int) bool {
	for i := 0; i < blkIndex; i++ {
		if d[i].Type != Free {
			continue
		}
		d[i], d[blkIndex] = d[blkIndex], d[i]
		return true
	}
	return false
}

func (d Disc) moveFileLeft(b *Blk, blkIndex int) bool {
	freeStart := 0
	freeLen := 0
	for i := range d {
		if d[i] == b {
			return false
		}
		if d[i].Type != Free {
			freeStart = i + 1
			freeLen = 0
			continue
		}
		freeLen++
		if freeLen == b.FileSize {
			break
		}
	}
	for i := freeStart; i < freeStart+b.FileSize; i++ {
		d[i] = b
	}
	for i := blkIndex; i < blkIndex+b.FileSize; i++ {
		d[i] = &Blk{Free, 0, 0}
	}
	return true
}

func (d Disc) Compact() {
	for i := len(d) - 1; i >= 0; i-- {
		if d[i].Type == Free {
			continue
		}
		if !d.moveBlockLeft(i) {
			return
		}
	}
}

func (d Disc) Defrag() {
	last := d[len(d)-1]
	for i := len(d) - 1; i >= 0; i-- {
		if d[i] == last {
			continue
		}
		if last.Type != Free {
			d.moveFileLeft(last, i+1)
		}
		last = d[i]
	}
}

func (d Disc) Checksum() int {
	csum := 0
	for i, b := range d {
		if b.Type == Free {
			continue
		}
		csum += i * b.Id
	}
	return csum
}

func loadDisc(f string) Disc {
	input, _ := os.ReadFile(f)
	input = bytes.TrimRight(input, "\n")
	disc := Disc{}
	id := 0
	for i, b := range input {
		var blk Blk
		size, _ := strconv.Atoi(string(b))
		if i%2 == 0 {
			blk = Blk{File, id, size}
			id++
		} else {
			blk = Blk{Free, 0, 0}
		}
		blocks := slices.Repeat([]*Blk{&blk}, size)
		disc = append(disc, blocks...)
	}
	return disc
}

func main() {
	disc := loadDisc("example.txt")
	disc.Compact()
	fmt.Println("Example 1 result:", disc.Checksum())

	disc = loadDisc("example.txt")
	disc.Defrag()
	fmt.Println("Example 2 result:", disc.Checksum())

	disc = loadDisc("input.txt")
	disc.Compact()
	fmt.Println("Part 1 result:", disc.Checksum())

	disc = loadDisc("input.txt")
	disc.Defrag()
	fmt.Println("Part 2 result:", disc.Checksum())
}
