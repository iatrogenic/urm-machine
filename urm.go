package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

type Program struct {
	Name         string
	Instructions []string
}

func (p Program) loc() int {
	return len(p.Instructions)
}

// Converts []string into []int, given certain assumptions
func strToIntSlice(original []string) []int {
	var new_slice []int 
	for _,n := range original {
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
	fmt.Println("Setting initial configuration")	
	tape := make(map[int]int)		
	for i, n := range init_c {
		tape[i] = n
	}
	
	// Running the instructions
	fmt.Printf("Running program \"%s\"\n", program.Name)

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
			case "Z" :
				tape[parsed_args[0]] = 0	
				flow_pointer += 1
			case "S" :
				tape[parsed_args[0]] += 1 
				flow_pointer += 1
			case "T" :
				tape[parsed_args[1]] = tape[parsed_args[0]]
				flow_pointer += 1
			case "J" :
				if tape[parsed_args[0]] == tape[parsed_args[1]] {
					flow_pointer = parsed_args[2]
				} else {
					flow_pointer += 1
				}
		}
	}

	fmt.Printf("Execution finished. R1 = %d.\n", tape[0])

}

func parseProg(filename string, name string) Program {

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("An error has occured while parsing %s", filename)
	}

	instructions := strings.Split(strings.ReplaceAll(string(contents), "\n", ""), ";")

	return Program{name, instructions}
}

func main() {
	var programs []Program

	fmt.Println("\n(un)Limited Register Machine Console\nType \"help\" for the commands.")

	// REPL
	for {
		fmt.Printf("\n\n> ")
		var command string
		fmt.Scanf("%s", &command)

		switch command {

		case "run":

			var init string
			var debug bool
			var pindex int

			fmt.Printf("Program index: ")
			fmt.Scanf("%d", &pindex)
			fmt.Printf("Initial configuration (comma separated integers): ")
			fmt.Scanf("%s", &init)
			fmt.Printf("Debug mode (true/false): ")
			fmt.Scanf("%t", &debug)

			if pindex <= len(programs) && pindex >= 0 {
				runProgram(programs[pindex], init, debug)	
			} else {
				fmt.Println("ERROR: Selected index is out of bounds.")
			}


		// Loads an URM program into memory
		case "load":
			//Implement check on program size (> 0)

			var filename, name string
			fmt.Println("Filename:")
			fmt.Scanf("%s", &filename)

			fmt.Println("Name for this program:")
			fmt.Scanf("%s", &name)

			programs = append(programs, parseProg(filename, name))
			fmt.Printf("Program %s successfully loaded.")

		case "show":
			// Chosen program
			var choice int 

			fmt.Printf("There are %d loaded programs.\n", len(programs))
			fmt.Println("i :  name \t LOC")

			for index, p := range programs {
				fmt.Printf("%d : %s \t %d", index, p.Name, p.loc())
			}

			fmt.Printf("\nWrite the program's index: ")
			fmt.Scanf("%d", &choice)

			for index, line := range programs[choice].Instructions {
				fmt.Printf("\n I[%d] : %s", index, line)
			}

		case "help":

			hlp := `
			Command List:
			run \t runs a chosen program
			load \t loads program onto memory
			help \t Console command list \n display \t Displays program's instructions"
			`
			fmt.Println(hlp)
		}

	}
}
