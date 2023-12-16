package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	ERRCODE_INVALID_ARGS = 1
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

var modMap = map[uint8][]string{
	0b000: {"bx", "si"},
	0b001: {"bx", "di"},
	0b010: {"bp", "si"},
	0b011: {"bp", "di"},
	0b100: {"si"},
	0b101: {"di"},
	0b110: {"bp"},
	0b111: {"dx"},
}

const (
	movR2rPattern       = byte(0b10_0010)
	movImmRegMemPattern = byte(0b110_0011)
	movImmRegPattern    = byte(0b1011)

	memoryMode    = byte(0b00)
	memoryMode8d  = byte(0b01)
	memoryMode16d = byte(0b10)
	registerMode  = byte(0b11)
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
		if b>>2 == movR2rPattern {
			_, err = fOut.WriteString(movR2r(b, br))
			check(err)
		} else if b>>1 == movImmRegMemPattern {
			_, err = fOut.WriteString(movImmRegMem(b, br))
			check(err)
		} else if b>>4 == movImmRegPattern {
			_, err = fOut.WriteString(movImmReg(b, br))
			check(err)
		} else {
			fmt.Printf("skipping byte (0x%x, 0b%b)...", b, b)
		}
	}
}

func movR2r(h byte, br io.ByteReader) string {
	l, err := br.ReadByte()
	check(err)

	mode := l >> 6

	// true = Register is Destination, false = Register is source.
	d := (h & 0b00000010) == 0b00000010
	// true = Word, false = byte.
	w := h & 0b00000001

	// Registry
	reg := (l & 0b00_111_000) >> 3

	// Registry or Memory
	rm := l & 0b00_000_111

	dest := registers[reg][w]

	var src string
	if mode == registerMode {
		src = registers[rm][w]
	} else {
		src = getMemoryModeSrc(rm, mode, br)
	}

	// if d is false, then reg is source, rm is destination.
	if !d {
		src, dest = dest, src
	}

	return fmt.Sprintf("mov %s, %s\n", dest, src)
}

func getMemoryModeSrc(rm, mode byte, br io.ByteReader) string {
	parts := modMap[rm]

	var add uint16
	switch mode {
	case memoryMode:
		// Nothing to do
		if rm == 0b110 {
			panic("not ready")
		}

	case memoryMode8d:
		add = readExtra(0, br)
	case memoryMode16d:
		add = readExtra(1, br)
	}

	if add > 0 {
		parts = append(parts, fmt.Sprintf("%d", add))
	}

	return "[" + strings.Join(parts, " + ") + "]"
}

func movImmRegMem(h byte, br io.ByteReader) string {
	return "YES\n"
}

func movImmReg(h byte, br io.ByteReader) string {
	reg := h & 0b00_000_111
	w := (h & 0b0000_1000) >> 3

	dest := registers[reg][w]

	src := readExtra(w, br)

	return fmt.Sprintf("mov %s, %d\n", dest, src)
}

func readExtra(w byte, br io.ByteReader) uint16 {
	l, err := br.ReadByte()
	check(err)

	src := uint16(l)
	if w == 1 {
		l, err = br.ReadByte()
		check(err)
		src += uint16(l) << 8
	}

	return src
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
