/*
PROJECT: bf2c
DESCRIPTION: An interpreter that converts brainfuck code to C code
AUTHOR: Vahin Sharma
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func compileCommand(command rune, tabs int) string {
	compiledCode := "\n"
	for i := 0; i < tabs; i++ {
		compiledCode += "\t"
	}

	switch command {
	case '+':
		compiledCode += "cells[currentCell] += 1;"

	case '-':
		compiledCode += "if(cells[currentCell] != 0){cells[currentCell] -= 1;}"

	case '<':
		compiledCode += "if(currentCell != 0){currentCell -= 1;}"

	case '>':
		compiledCode += "if(currentCell != 29999){currentCell += 1;}"

	case ',':
		compiledCode += "cells[currentCell] = (int)getchar();"

	case '.':
		compiledCode += "printf(\"%c\", cells[currentCell]);"

	default:
		compiledCode += "ic" // ic == invalid command
	}

	return compiledCode
}

func main() {
	inputFile := flag.String("f", "", "Input file")
	outputFile := flag.String("o", "output.c", "Output file")

	flag.Parse()

	if *inputFile == "" {
		fmt.Printf("FLAGS:\n")
		fmt.Printf("-f=<file> Input file\n")
		fmt.Printf("-o=<file> Output file (default value: output.c)\n")
		os.Exit(1)
	}

	code, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("ERROR: Could not read file `%s`\nReason: %s\n", *inputFile, err)
		os.Exit(1)
	}

	var compiledCode string = `// Generated by bf2c
#include <stdio.h>

int main() {
	int cells[30000];
	int currentCell = 15000;
	int c;

	for (int i = 0; i < 30000; i++) {
		cells[i] = 0;
	}
`
	var loopCode string
	// Code that are contained in loops will be in loopCode
	// e.g. ++[.----]>>
	// loopCode == ".----"
	var isInLoop bool

	for _, c := range code {
		switch c {
		case '[':
			isInLoop = true

		case ']':
			isInLoop = false
			compiledCode += "\n\twhile(cells[currentCell] != 0) {"
			for _, ic := range loopCode {
				commandToAdd := compileCommand(rune(ic), 2)
				if commandToAdd != "\n\t\tic" {
					compiledCode += commandToAdd
				} else {
					fmt.Printf("WARNING: Invalid command `%c`, skipping...\n", c)
				}
			}
			compiledCode += "\n\t}"
			loopCode = ""

		case '\n', '\t':
			continue

		default:
			if isInLoop {
				loopCode += string(c)
			} else {
				commandToAdd := compileCommand(rune(c), 1)
				if commandToAdd != "\n\tic" {
					compiledCode += commandToAdd
				} else {
					fmt.Printf("WARNING: Invalid command `%c`, skipping...\n", c)
				}
			}
		}
	}

	compiledCode += "\n}"

	err = ioutil.WriteFile(*outputFile, []byte(compiledCode), 0o644)
	if err != nil {
		fmt.Printf("ERROR: Could not write file `%s`\nReason: %s\n", *outputFile, err)
		os.Exit(2)
	}
}
