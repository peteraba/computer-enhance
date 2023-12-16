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
			name: "register to register",
			in:   []byte{0x89, 0xde, 0x88, 0xc6},
			expected: `bits 16
mov si, bx
mov dh, al
`,
		},
		{
			name: "8-bit immediate-to-register",
			in:   []byte{0xb1, 0x0c, 0xb5, 0xf4},
			expected: `bits 16
mov cl, 12
mov ch, 244
`,
		},
		{
			name: "16-bit immediate-to-register",
			in:   []byte{0xb9, 0x0c, 0x00, 0xb9, 0xf4, 0xff, 0xba, 0x6c, 0x0f, 0xba, 0x94, 0xf0},
			expected: `bits 16
mov cx, 12
mov cx, 65524
mov dx, 3948
mov dx, 61588
`,
		},
		{
			name: "source address calculation",
			in:   []byte{0x8a, 0x00, 0x8b, 0x1b, 0x8b, 0x56, 0x00},
			expected: `bits 16
mov al, [bx + si]
mov bx, [bp + di]
mov dx, [bp]
`,
		},
		{
			name: "source address calculation + 8bit",
			in:   []byte{0x8a, 0x60, 0x04},
			expected: `bits 16
mov ah, [bx + si + 4]
`,
		},
		{
			name: "source address calculation + 16bit",
			in:   []byte{0x8a, 0x80, 0x87, 0x13},
			expected: `bits 16
mov al, [bx + si + 4999]
`,
		},
		{
			name: "dest address calculation",
			in:   []byte{0x89, 0x09, 0x88, 0x0a, 0x88, 0x6e, 0x00},
			expected: `bits 16
mov [bx + di], cx
mov [bp + si], cl
mov [bp], ch
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
