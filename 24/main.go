package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Operation func(a, b bool) bool

type Gate struct {
	name         string
	op           string
	in1, in2     *Gate
	valid, value bool
}

func (g Gate) CanExecute() bool {
	return g.in1 != nil && g.in2 != nil && g.in1.valid && g.in2.valid
}

func (g *Gate) Execute() bool {
	if !g.CanExecute() {
		return false
	}
	switch g.op {
	case "AND":
		g.value = g.in1.value && g.in2.value
	case "OR":
		g.value = g.in1.value || g.in2.value
	case "XOR":
		g.value = g.in1.value != g.in2.value
	default:
		panic(fmt.Sprint("invalid operation:", g.op))
	}
	g.valid = true
	return true
}

type Device struct {
	gates       map[string]*Gate
	outputGates []*Gate
	x, y, z     int
}

func (d *Device) Run() {
	queue := make([]*Gate, 0, len(d.gates))
	for _, g := range d.gates {
		queue = append(queue, g)
	}
outer:
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if !curr.valid && !curr.Execute() {
			queue = append(queue, curr)
		}
		for _, g := range d.outputGates {
			if !g.valid {
				continue outer
			}
		}
		break
	}
	d.setResult()
}

func (d *Device) setResult() {
	d.z = 0
	for _, g := range d.outputGates {
		d.z <<= 1
		if g.value {
			d.z |= 1
		}
	}
}

func load(f string) *Device {
	d, _ := os.ReadFile(f)
	lines := strings.Split(strings.TrimSpace(string(d)), "\n\n")
	device := Device{
		gates: map[string]*Gate{},
	}
	var ix, iy int
	for _, line := range strings.Split(lines[0], "\n") {
		s := strings.Split(line, ": ")
		v := s[1] == "1"
		g := &Gate{
			name:  s[0],
			valid: true,
			value: v,
		}
		device.gates[g.name] = g

		if s[0][0] == 'x' {
			if v {
				device.x |= 1 << ix
			}
			ix++

		}
		if s[0][0] == 'y' {
			if v {
				device.y |= 1 << iy
			}
			iy++
		}
	}
	for _, line := range strings.Split(lines[1], "\n") {
		s := strings.Split(line, " ")
		in1, op, in2, name := s[0], s[1], s[2], s[4]
		if _, ok := device.gates[in1]; !ok {
			device.gates[in1] = &Gate{name: in1}
		}
		if _, ok := device.gates[in2]; !ok {
			device.gates[in2] = &Gate{name: in2}
		}
		g, ok := device.gates[name]
		if !ok {
			g = &Gate{name: name}
			device.gates[name] = g
		}
		g.op = op
		g.in1 = device.gates[in1]
		g.in2 = device.gates[in2]
	}
	for name, gate := range device.gates {
		if name[0] == 'z' {
			device.outputGates = append(device.outputGates, gate)
		}
	}
	slices.SortFunc(device.outputGates, func(a, b *Gate) int { return strings.Compare(b.name, a.name) })
	return &device
}

func part1(d *Device) int {
	d.Run()
	return d.z
}

func part2(d *Device) string {
	wrong := []*Gate{}

	for i, g := range d.outputGates {
		// gates are sorted DESC in load
		// first output is just x0 XOR y0
		if i == len(d.outputGates)-1 {
			if g.op != "XOR" {
				wrong = append(wrong, g)
			}
			continue
		}
		// last gate is just carry out
		if i == 0 {
			if g.op != "OR" {
				wrong = append(wrong, g)
			}
			continue
		}

		// output has to be XOR
		if g.op != "XOR" {
			wrong = append(wrong, g)
			continue
		}

		in, cout := g.in1, g.in2
		// swap them to match names
		if g.in1.op == "OR" || g.in2.op == "XOR" {
			in, cout = g.in2, g.in1
		}

		// in has to be xor
		if in.op != "XOR" {
			wrong = append(wrong, in)
		} else {
			// and has to have x and y inputs
			n1, n2 := in.in1.name[0], in.in2.name[0]
			if n2 == 'y' && n1 != 'x' {
				wrong = append(wrong, in.in1)
			}
			if n1 == 'x' && n2 != 'y' {
				wrong = append(wrong, in.in2)
			}
			if n1 != 'y' && n1 != 'x' && n2 != 'y' && n2 != 'x' {
				wrong = append(wrong, in.in1, in.in2)
			}
		}

		if i == len(d.outputGates)-2 && cout.op == "AND" {
			// carry out is just X AND Y for the second gate, valid
			continue
		} else if cout.op != "OR" {
			// carry out has to be OR
			wrong = append(wrong, cout)
		} else {
			// carry out has to be connected to 2 AND gates
			if cout.in1.op != "AND" {
				wrong = append(wrong, cout.in1)
			}
			if cout.in2.op != "AND" {
				wrong = append(wrong, cout.in2)
			}
		}
	}

	slices.SortFunc(wrong, func(a, b *Gate) int { return strings.Compare(a.name, b.name) })
	r := ""
	for i, g := range wrong {
		if i > 0 {
			r += ","
		}
		r += g.name
	}
	return r
}

func main() {
	d := load("example.txt")
	fmt.Println("Example 1:", part1(d))

	d = load("input.txt")
	fmt.Println("Part 1:", part1(d))
	fmt.Println("Part 2:", part2(d))
}
