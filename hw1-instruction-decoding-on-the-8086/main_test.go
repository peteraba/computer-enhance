package main

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_deassemble(t *testing.T) {
	tests := []struct {
		name     string
		in       []byte
		expected string
	}{
		{
			name:     "empty",
			in:       []byte{},
			expected: "bits 16\n",
		},
		{
			name:     "default",
			in:       []byte{0x89, 0xd9},
			expected: "bits 16\nmov cx, bx\n",
		},
		{
			name: "complex",
			in:   []byte{0x89, 0xd9, 0x88, 0xe5, 0x89, 0xda, 0x89, 0xde, 0x89, 0xfb, 0x88, 0xc8, 0x88, 0xed, 0x89, 0xc3, 0x89, 0xf3, 0x89, 0xfc, 0x89, 0xc5},
			expected: `
bits 16
mov cx, bx
mov ch, ah
mov dx, bx
mov si, bx
mov bx, di
mov al, cl
mov ch, ch
mov bx, ax
mov bx, si
mov sp, di
mov bp, ax
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader(tt.in))
			writer := new(bytes.Buffer)

			deassemble(reader, writer)

			if tt.expected != writer.String() {
				t.Errorf("deassemble() = %v, want %v", writer.String(), tt.expected)
			}
		})
	}
}
