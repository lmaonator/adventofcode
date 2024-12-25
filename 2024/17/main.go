package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type CPU struct {
	ra, rb, rc   int // Registers
	ip           int // Instruction Pointer
	instructions []int
	output       string
}

func (cpu *CPU) ComboOperand(operand int) int {
	if operand >= 0 && operand <= 3 {
		return operand
	}
	switch operand {
	case 4:
		return cpu.ra
	case 5:
		return cpu.rb
	case 6:
		return cpu.rc
	case 7:
		panic("Reserved combo operand 7")
	default:
		panic(fmt.Sprint("Invalid combo operand", operand))
	}
}

func (cpu *CPU) Adv(combo int) {
	cpu.ra = int(math.Trunc(float64(cpu.ra) / math.Pow(2, float64(cpu.ComboOperand(combo)))))
}

func (cpu *CPU) Bxl(literal int) {
	cpu.rb ^= literal
}

func (cpu *CPU) Bst(combo int) {
	cpu.rb = cpu.ComboOperand(combo) % 8
}

func (cpu *CPU) Jnz(literal int) bool {
	if cpu.ra == 0 {
		return false
	}
	cpu.ip = literal
	return true
}

func (cpu *CPU) Bxc(ignored int) {
	cpu.rb ^= cpu.rc
}

func (cpu *CPU) Out(combo int) {
	if len(cpu.output) > 0 {
		cpu.output += ","
	}
	cpu.output += fmt.Sprint(cpu.ComboOperand(combo) % 8)
}

func (cpu *CPU) Bdv(combo int) {
	cpu.rb = int(math.Trunc(float64(cpu.ra) / math.Pow(2, float64(cpu.ComboOperand(combo)))))
}

func (cpu *CPU) Cdv(combo int) {
	cpu.rc = int(math.Trunc(float64(cpu.ra) / math.Pow(2, float64(cpu.ComboOperand(combo)))))
}

func (cpu *CPU) Run() {
	for {
		if cpu.ip < 0 || cpu.ip+1 >= len(cpu.instructions) {
			return
		}
		opcode, operand := cpu.instructions[cpu.ip], cpu.instructions[cpu.ip+1]
		incIp := true
		switch opcode {
		case 0:
			cpu.Adv(operand)
		case 1:
			cpu.Bxl(operand)
		case 2:
			cpu.Bst(operand)
		case 3:
			if cpu.Jnz(operand) {
				incIp = false
			}
		case 4:
			cpu.Bxc(operand)
		case 5:
			cpu.Out(operand)
		case 6:
			cpu.Bdv(operand)
		case 7:
			cpu.Cdv(operand)
		default:
			panic(fmt.Sprint("Invalid opcode", opcode))
		}
		if incIp {
			cpu.ip += 2
		}
	}
}

func (cpu *CPU) LoadInstructions(filename string) string {
	d, _ := os.ReadFile(filename)
	split := strings.Split(string(d), "\n\n")

outer:
	for i, line := range strings.Split(split[0], "\n") {
		s := strings.Split(line, ": ")
		v, _ := strconv.Atoi(s[1])
		switch i {
		case 0:
			cpu.ra = v
		case 1:
			cpu.rb = v
		case 2:
			cpu.rc = v
		default:
			break outer
		}
	}

	progStr := strings.Split(strings.TrimSpace(split[1]), ": ")[1]
	prog := strings.Split(progStr, ",")
	cpu.instructions = make([]int, 0, len(prog))
	for _, inst := range prog {
		v, _ := strconv.Atoi(inst)
		cpu.instructions = append(cpu.instructions, v)
	}
	return progStr
}

func NewCPU(filename string) CPU {
	cpu := CPU{}
	cpu.LoadInstructions(filename)
	return cpu
}

func part2(f string) int {
	cpu := CPU{}
	wanted := cpu.LoadInstructions(f)
	i := 1
	for {
		c := cpu
		c.ra = i
		c.Run()
		if wanted == c.output {
			break
		}
		if len(c.output) > len(wanted) {
			panic("unable to find solution")
		}
		if strings.HasSuffix(wanted, c.output) {
			i <<= 3
			continue
		}
		i++
	}
	return i
}

func main() {
	cpu := NewCPU("example.txt")
	cpu.Run()
	fmt.Println("Example 1:", cpu.output)
	result := part2("example2.txt")
	fmt.Println("Example 2:", result)

	cpu = NewCPU("input.txt")
	cpu.Run()
	fmt.Println("Part 1:", cpu.output)
	result = part2("input.txt")
	fmt.Println("Part 2:", result)
}
