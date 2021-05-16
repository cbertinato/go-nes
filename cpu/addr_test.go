package cpu

import (
    "testing"
)


// TestZP0 tests the zero page addressing mode
func TestZP0(t *testing.T) {
    b := DevBus{}
	c := MOS6502{Bus: &b}

    b.ram[0] = 0xFE
    b.ram[1] = 0xCA

    c.zp0()

    if c.absAddr != 0x00FE {
        t.Errorf("Expected absAddr = %#04x, got %#04x", 0x00FE, c.absAddr)
    }
}

// TestZPX tests the zero page addressing mode
func TestZPx(t *testing.T) {
    b := DevBus{}
    c := MOS6502{Bus: &b}

    b.ram[0] = 0xFE
    c.X = 0xCA

    c.zpX()

    expected := uint16(0x00FF & (0xCA + 0xFE))

    if c.absAddr != expected {
        t.Errorf("Expected absAddr = %#04x, got %#04x", expected, c.absAddr)
    }
}

// TestRel tests the relative addressing mode
func TestRel(t *testing.T) {
    tests := []struct{
        name string
        rel uint8
        expected uint16
    }{
        {"positive", 0x50, uint16(0x50)},
        {"negative", 0x90, uint16(0x0090 | 0xFF00)},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            b := DevBus{}
            c := MOS6502{Bus: &b}
            b.ram[0] = tt.rel

            c.rel()

            if c.relAddr != tt.expected {
                t.Errorf("Expected relAddr = %#04x, got %#04x", tt.expected, c.relAddr)
            }
        })
    }
}

// TestAbs tests the absolute addressing mode
func TestAbs(t *testing.T) {
    b := DevBus{}
    c := MOS6502{Bus: &b}

    b.ram[0] = 0xFE
    b.ram[1] = 0xCA

    c.abs()

    expected := uint16(0xCAFE)
    if c.absAddr != expected {
        t.Errorf("Expected absAddr = %#04x, got %#04x", expected, c.absAddr)
    }
}

// TestAbsIndex tests the absolute addressing mode
func TestAbsIndex(t *testing.T) {
    tests := []struct{
        name string
        idx uint8
        expected uint16
        cycles uint8
    }{
        {"same page", 0x01, uint16(0xCAFF), 0},
        {"different page", 0x05, uint16(0xCB03), 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            b := DevBus{}
            c := MOS6502{Bus: &b}

            b.ram[0] = 0xFE
            b.ram[1] = 0xCA
            c.X = uint8(tt.idx)

            cycles := c.abX()

            if c.absAddr != tt.expected {
                t.Errorf("Expected absAddr = %#04x, got %#04x", tt.expected, c.absAddr)
            }

            if cycles != tt.cycles {
                t.Errorf("Expected cycles = %d, got %d", tt.cycles, cycles)
            }
        })
    }
}

// TestIndirect tests the indirect addressing mode
func TestIndirectNominal(t *testing.T) {
    b := DevBus{}
    c := MOS6502{Bus: &b}

    b.ram[0] = 0xFE
    b.ram[1] = 0xCA
    b.ram[0xCAFE] = 0xEF
    b.ram[0xCAFF] = 0xBE

    expected := uint16(0xBEEF)

    c.ind()

    if c.absAddr != expected {
        t.Errorf("Expected absAddr = %#04x, got %#04x", expected, c.absAddr)
    }
}

func TestIndirectBug(t *testing.T) {
    b := DevBus{}
    c := MOS6502{Bus: &b}

    b.ram[0] = 0xFF
    b.ram[1] = 0xCA
    b.ram[0xCAFF] = 0xEF
    b.ram[0xCA00] = 0xBE

    expected := uint16(0xBEEF)

    c.ind()

    if c.absAddr != expected {
        t.Errorf("Expected absAddr = %#04x, got %#04x", expected, c.absAddr)
    }
}

func TestIndexedIndirect(t *testing.T) {
    b := DevBus{}
    c := MOS6502{Bus: &b}

    c.PC = 0x0100
    b.ram[0x0100] = 0x10
    c.X = 0xFA

    b.ram[0x0A] = 0xFE
    b.ram[0x0B] = 0xCA

    c.izX()

    expected := uint16(0xCAFE)

    if c.absAddr != expected {
        t.Errorf("Expected absAddr = %#04x, got %#04x", expected, c.absAddr)
    }
}

func TestIndirectIndex(t * testing.T) {
    b := DevBus{}
    c := MOS6502{Bus: &b}

    c.PC = 0x0100
    b.ram[0x0100] = 0x10
    b.ram[0x0010] = 0xFE
    b.ram[0x0011] = 0xCA
    c.Y = 0x10

    expected := uint16(0xCB0E)

    c.izY()

    if c.absAddr != expected {
        t.Errorf("Expected absAddr = %#04x, got %#04x", expected, c.absAddr)
    }
}