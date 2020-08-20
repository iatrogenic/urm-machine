// To-Do:
// Write an actual parser
// Polish the debug mode

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Program struct {
	Instructions []string
}

func (p Program) loc() int {
	return len(p.Instructions)
}

func (p Program) rho() int {
	l_reg := 0
	for _, instruction := range p.Instructions {
		inst_type := string(instruction[0])
		parsed_args := strToIntSlice(strings.Split(instruction[2:len(instruction)-1], ","))
		switch inst_type {
		case "Z", "S":
			if parsed_args[0] > l_reg {
				l_reg = parsed_args[0]
			}

		case "T", "J":
			if (parsed_args[1] > l_reg) || (parsed_args[0] > l_reg) {
				if parsed_args[1] >= parsed_args[0] {
					l_reg = parsed_args[1]
				} else {
					l_reg = parsed_args[0]
				}
			}
		}
	}
	return l_reg
}

// Converts []string into []int, given certain assumptions
func strToIntSlice(original []string) []int {
	var new_slice []int
	for _, n := range original {
		content, _ := strconv.Atoi(n)
		new_slice = append(new_slice, content)
	}
	return new_slice
}

func runProgram(program Program, init string, debug bool) {
	var loc, flow_pointer int
	loc = program.loc() - 1
	flow_pointer = 0

	// Parse initial configuration
	init_c := strToIntSlice(strings.Split(init, ","))

	// Preparing the tape
	tape := make(map[int]int)
	for i, n := range init_c {
		tape[i] = n
	}
	for i := len(tape); i <= program.rho(); i++ {
		tape[i] = 0
	}

	// Running the instructions
	for flow_pointer <= loc {
		if debug == true {
			fmt.Println(tape)
			fmt.Println(flow_pointer)
			fmt.Scanln()
		}
		instruction := program.Instructions[flow_pointer]
		inst_type := string(instruction[0])
		parsed_args := strToIntSlice(strings.Split(instruction[2:len(instruction)-1], ","))
		switch inst_type {
		case "Z":
			tape[parsed_args[0]] = 0
			flow_pointer += 1
		case "S":
			tape[parsed_args[0]] += 1
			flow_pointer += 1
		case "T":
			tape[parsed_args[1]] = tape[parsed_args[0]]
			flow_pointer += 1
		case "J":
			if tape[parsed_args[0]] == tape[parsed_args[1]] {
				flow_pointer = parsed_args[2]
			} else {
				flow_pointer += 1
			}
		}
	}

	fmt.Printf("Execution finished: R1 = %d.\n", tape[0])

}

func parseProg(filename string) Program {

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("An error has occured while parsing %s", filename)
		os.Exit(1)
	}

	instructions := strings.Split(strings.ReplaceAll(string(contents), "\n", ""), ";")

	return Program{instructions}
}

func main() {
	// var programs []Program

	var init_ptr = flag.String("init", "", "The initial configuration of the URM")
	var debug_ptr = flag.Bool("debug", false, "Debug mode")
	flag.Parse()
	config := *init_ptr
	debug := *debug_ptr

	if debug {
		fmt.Println("Debug mode is on.")
	}
	if flag.NArg() != 1 {
		fmt.Println("Incompatible number of arguments.")
		os.Exit(1)
	}

	var filename string
	filename = flag.Args()[0]

	program := parseProg(filename)
	runProgram(program, config, debug)
}
