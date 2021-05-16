package cpu

import (
	"strconv"
	"testing"
)

func TestAdc(t *testing.T) {
	tests := []struct {
		name   string
		a      uint8
		m      uint8
		result uint8
		status uint8
	}{
		{"no overflow", 10, 15, 25, 0},
		{"overflow/negative", 127, 10, 137, 1<<6 | 1<<7},
		{"zero", 0, 0, 0, 1 << 1},
		{"carry", 250, 10, 4, 1 << 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := DevBus{}
			c := MOS6502{Bus: &b}

			c.opLookup = []Instruction{
				Instruction{
					name:     "ADC",
					op:       func() uint8 { return 1 },
					addrMode: "IMM",
					cycles:   1,
				},
			}

			c.fetched = tt.m
			c.A = tt.a

			if c.Status != 0 {
				t.Errorf("status register not 0")
			}

			c.adc()

			if c.Status != tt.status {
				t.Errorf("Expected status = %s, got %s",
					strconv.FormatUint(uint64(tt.status), 2), strconv.FormatUint(uint64(c.Status), 2))
			}

			if c.A != tt.result {
				t.Errorf("Expected result = %d, got %d", tt.result, c.A)
			}
		})
	}
}

func TestSbc(t *testing.T) {
	tests := []struct {
		name   string
		a      uint8
		m      uint8
		result uint8
		status uint8
	}{
		{"no unsigned borrow", 10, 5, 5, 1 << 0},
		{"unsigned borrow", 5, 10, 251, 1 << 7},
		{"zero", 5, 5, 0, 1<<0 | 1<<1},
		{"signed overflow", 208, 112, 96, 1<<0 | 1<<6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MOS6502{Bus: &DevBus{}}

			c.opLookup = []Instruction{
				Instruction{
					name:     "SBC",
					op:       func() uint8 { return 1 },
					addrMode: "IMM",
					cycles:   1,
				},
			}

			c.fetched = tt.m
			c.A = tt.a

			// since the "borrow" flag is the complement of the carry
			// then set the carry flag to unset the "borrow" flag
			c.SetFlag(C, true)

			if c.Status != 1<<0 {
				t.Errorf("carry flag not set")
			}

			c.sbc()

			if c.Status != tt.status {
				t.Errorf("Expected status = %s, got %s",
					strconv.FormatUint(uint64(tt.status), 2), strconv.FormatUint(uint64(c.Status), 2))
			}

			if c.A != tt.result {
				t.Errorf("Expected result = %d, got %d", tt.result, c.A)
			}
		})
	}
}
