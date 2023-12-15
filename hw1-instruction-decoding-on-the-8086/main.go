package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	ERRCODE_INVALID_ARGS    = 1
	ERRCODE_INVALID_CONTENT = 2
)

var registers = map[uint8][2]string{
	0b000: {"al", "ax"},
	0b001: {"cl", "cx"},
	0b010: {"dl", "dx"},
	0b011: {"bl", "bx"},
	0b100: {"ah", "sp"},
	0b101: {"ch", "bp"},
	0b110: {"dh", "si"},
	0b111: {"bh", "di"},
}

const (
	movPattern    = byte(0b10_0010)
	regMovPattern = byte(0b11)
)

func main() {
	inputFile, outputFile := getArgs()

	source := readFile(inputFile)

	var content []byte
	for i := 0; i < len(source); i += 2 {
		h := source[i]
		if h>>2 != movPattern {
			fmt.Println("The first byte does not start with the pattern 100010.")
			os.Exit(ERRCODE_INVALID_CONTENT)
		}
		l := source[i+1]
		if l>>6 != regMovPattern {
			fmt.Println("The second byte does not start with the pattern 11.")
			os.Exit(ERRCODE_INVALID_CONTENT)
		}

		// true = Register is Destination, false = Register is source.
		d := (h & 0b00000010) == 0b00000010
		// true = Word, false = byte.
		w := h & 0b00000001

		// Registry
		reg := (l & 0b00_111_000) >> 3

		// Registry or Memory
		rm := l & 0b00_000_111

		dest := registers[reg][w]
		src := registers[rm][w]

		// if d is false, then reg is source, rm is destination.
		if !d {
			src, dest = dest, src
		}

		result := fmt.Sprintf("mov %s, %s\n", dest, src)

		content = append(content, []byte(result)...)
	}

	writeFile(inputFile, outputFile, content)
}

func getArgs() (string, string) {
	if len(os.Args) < 2 {
		fmt.Println("input input argument is missing")
		os.Exit(ERRCODE_INVALID_ARGS)
	}
	if len(os.Args) < 3 {
		fmt.Println("output input argument is missing")
		os.Exit(ERRCODE_INVALID_ARGS)
	}

	return os.Args[1], os.Args[2]
}

func readFile(inputFile string) []byte {
	dat, err := os.ReadFile(inputFile)
	check(err)

	if len(dat) == 0 {
		panic("file is empty: " + inputFile)
	}

	if len(dat)%2 != 0 {
		panic(fmt.Sprintf("file does not content is expected to be even length but got odd length. input: %s, length: %d", inputFile, len(dat)))
	}

	return dat
}

func writeFile(inputFile, outputFile string, content []byte) {
	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()

	_ = inputFile

	_, err = f.WriteString("bits 16\n")
	check(err)

	_, err = f.Write(content)
	check(err)
}
