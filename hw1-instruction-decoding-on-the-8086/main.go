package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func check(e error, allowedErrors ...error) {
	for _, allowedError := range allowedErrors {
		if e == allowedError {
			return
		}
	}

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
	movPattern = byte(0b10_0010)

	regMovPattern = byte(0b11)
)

func main() {
	inputFileName, outputFileName := getArgs()

	fIn, err := os.Open(inputFileName)
	check(err)
	defer fIn.Close()

	fOut, err := os.Create(outputFileName)
	check(err)
	defer fOut.Close()

	br := bufio.NewReader(fIn)

	deassemble(br, fOut)
}

func deassemble(br io.ByteReader, fOut io.StringWriter) {
	_, err := fOut.WriteString("bits 16\n")
	check(err)

	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			break
		}
		check(err)

		// process mov
		if b>>2 == movPattern {
			b2, err := br.ReadByte()
			check(err, io.EOF)

			_, err = fOut.WriteString(doMov(b, b2))
			check(err)
		} else {
			fmt.Printf("skipping byte (0x%x, 0b%b)...", b, b)
		}
	}
}
func doMov(h, l byte) string {
	if l>>6 != regMovPattern {
		fmt.Printf("The first byte is okay. Byte: 0x%x, 0b%b\n", h, h)
		fmt.Printf("The second byte does not start with the pattern 11. Byte: 0x%x, 0b%b\n", l, l)
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

	return fmt.Sprintf("mov %s, %s\n", dest, src)
}

func getArgs() (string, string) {
	if len(os.Args) < 2 {
		fmt.Println("input argument is missing")
		os.Exit(ERRCODE_INVALID_ARGS)
	}
	if len(os.Args) < 3 {
		fmt.Println("output argument is missing")
		os.Exit(ERRCODE_INVALID_ARGS)
	}

	return os.Args[1], os.Args[2]
}
