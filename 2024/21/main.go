package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var kpMove = map[byte]map[byte]string{
	'0': {
		'A': ">A",
		'0': "A",
		'1': "^<A",
		'2': "^A",
		'3': "^>A",
		'4': "^<^A",
		'5': "^^A",
		'6': "^^>A",
		'7': "^^^<A",
		'8': "^^^A",
		'9': "^^^>A",
	},
	'1': {
		'A': ">>vA",
		'0': ">vA",
		'1': "A",
		'2': ">A",
		'3': ">>A",
		'4': "^A",
		'5': "^>A",
		'6': "^>>A",
		'7': "^^A",
		'8': "^^>A",
		'9': "^^>>A",
	},
	'2': {
		'A': "v>A",
		'0': "vA",
		'1': "<A",
		'2': "A",
		'3': ">A",
		'4': "<^A",
		'5': "^A",
		'6': "^>A",
		'7': "<^^A",
		'8': "^^A",
		'9': "^^>A",
	},
	'3': {
		'A': "vA",
		'0': "<vA",
		'1': "<<A",
		'2': "<A",
		'3': "A",
		'4': "<<^A",
		'5': "<^A",
		'6': "^A",
		'7': "<<^^A",
		'8': "<^^A",
		'9': "^^A",
	},
	'4': {
		'A': ">>vvA",
		'0': ">vvA",
		'1': "vA",
		'2': "v>A",
		'3': "v>>A",
		'4': "A",
		'5': ">A",
		'6': ">>A",
		'7': "^A",
		'8': "^>A",
		'9': "^>>A",
	},
	'5': {
		'A': "vv>A",
		'0': "vvA",
		'1': "<vA",
		'2': "vA",
		'3': "v>A",
		'4': "<A",
		'5': "A",
		'6': ">A",
		'7': "<^A",
		'8': "^A",
		'9': "^>A",
	},
	'6': {
		'A': "vvA",
		'0': "<vvA",
		'1': "<<vA",
		'2': "<vA",
		'3': "vA",
		'4': "<<A",
		'5': "<A",
		'6': "A",
		'7': "<<^A",
		'8': "<^A",
		'9': "^A",
	},
	'7': {
		'A': ">>vvvA",
		'0': ">vvvA",
		'1': "vvA",
		'2': "vv>A",
		'3': "vv>>A",
		'4': "vA",
		'5': "v>A",
		'6': "v>>A",
		'7': "A",
		'8': ">A",
		'9': ">>A",
	},
	'8': {
		'A': "vvv>A",
		'0': "vvvA",
		'1': "<vvA",
		'2': "vvA",
		'3': "vv>A",
		'4': "<vA",
		'5': "vA",
		'6': "v>A",
		'7': "<A",
		'8': "A",
		'9': ">A",
	},
	'9': {
		'A': "vvvA",
		'0': "<vvvA",
		'1': "<<vvA",
		'2': "<vvA",
		'3': "vvA",
		'4': "<<vA",
		'5': "<vA",
		'6': "vA",
		'7': "<<A",
		'8': "<A",
		'9': "A",
	},
	'A': {
		'A': "A",
		'0': "<A",
		'1': "^<<A",
		'2': "<^A",
		'3': "^A",
		'4': "^^<<A",
		'5': "<^^A",
		'6': "^^A",
		'7': "^^^<<A",
		'8': "<^^^A",
		'9': "^^^A",
		'^': "<A",
		'<': "v<<A",
		'v': "<vA",
		'>': "vA",
	},
	'^': {
		'A': ">A",
		'^': "A",
		'<': "v<A",
		'v': "vA",
		'>': "v>A",
	},
	'<': {
		'A': ">>^A",
		'^': ">^A",
		'<': "A",
		'v': ">A",
		'>': ">>A",
	},
	'v': {
		'A': "^>A",
		'^': "^A",
		'<': "<A",
		'v': "A",
		'>': ">A",
	},
	'>': {
		'A': "^A",
		'^': "<^A",
		'<': "<<A",
		'v': "<A",
		'>': "A",
	},
}

func load(f string) []string {
	c, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(c)), "\n")
	return lines
}

type cacheKey struct {
	code   string
	robots int
}

var cache = map[cacheKey]int{}

func getSequenceLen(code string, robots int) int {
	if robots == 0 {
		return len(code)
	}
	key := cacheKey{code, robots}
	if v, ok := cache[key]; ok {
		return v
	}
	length := 0
	next := byte('A')
	for _, curr := range code {
		move := kpMove[next][byte(curr)]
		length += getSequenceLen(move, robots-1)
		next = byte(curr)
	}
	cache[key] = length
	return length
}

func solve(codes []string, robots int) int {
	r := 0
	for _, code := range codes {
		length := getSequenceLen(code, robots)
		n, _ := strconv.Atoi(code[:len(code)-1])
		r += length * n
	}
	return r
}

func main() {
	codes := load("example.txt")
	fmt.Println("Example 1:", solve(codes, 3))

	codes = load("input.txt")
	fmt.Println("Part 1:", solve(codes, 3))
	fmt.Println("Part 2:", solve(codes, 26))
}
