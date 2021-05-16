package cpu

import (
	"testing"
)

// TestCycle exercises the clock() function
func TestCycle(t *testing.T) {

	testAddrModeFunc := func(c *MOS6502) uint8 {
		return 0
	}

	testOpFunc := func(c *MOS6502) func() uint8 {
		return func() uint8 {
			c.A = 0xDE
			return 0
		}
	}

	b := DevBus{}
	c := MOS6502{Bus: &b}

	instruction := Instruction{
		name: "test",
		op: testOpFunc(&c),
		addrMode: "testAddrMode",
		cycles: 3,
	}

	c.opLookup = []Instruction{instruction}
	c.addrModeLookup = map[string]func(*MOS6502)uint8{"testAddrMode": testAddrModeFunc}

	if c.A != 0 {
		t.Errorf("Accumulator was not initialized to 0")
	}

	if c.cycles != 0 {
		t.Errorf("Cycles not initialized to 0")
	}

	expectedCycles := instruction.cycles - 1

	for i := 1; i <= int(instruction.cycles - 1); i++ {
		c.clock()
		if c.cycles != expectedCycles {
			t.Errorf("Got cycles=%d, expected %d", c.cycles, expectedCycles)
		}
		expectedCycles--
	}

	if c.A != 0xDE {
		t.Errorf("Accumulator value not correctly set after op: got A=%#02x, expected %#02x", c.A, 0xDE)
	}
}

func TestCycleExtraCycles(t *testing.T) {

	testAddrModeFunc := func(c *MOS6502)  uint8 {
		return 1
	}

	testOpFunc := func(c *MOS6502) func() uint8 {
		return func() uint8 {
			return 1
		}
	}

	b := DevBus{}
	c := MOS6502{Bus: &b}

	instruction := Instruction{
		name: "test",
		op: testOpFunc(&c),
		addrMode: "testAddrMode",
		cycles: 3,
	}

	c.opLookup = []Instruction{instruction}
	c.addrModeLookup = map[string]func(*MOS6502)uint8{"testAddrMode": testAddrModeFunc}

	if c.cycles != 0 {
		t.Errorf("Cycles not initialized to 0")
	}

	expectedCycles := instruction.cycles
	c.clock()
	
	if c.cycles != expectedCycles {
		t.Errorf("Got cycles=%d, expected %d", c.cycles, expectedCycles)
	}
}

