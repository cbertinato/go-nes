package cpu

// The 6502 is a little-endian 8-bit processor with a 16-bit address bus.

// Instruction represents a single 6502 instruction
type Instruction struct {
	name     string
	op       func() uint8
	addrMode string
	cycles   uint8
}

// MOS6502 represents the state of the CPU
type MOS6502 struct {
	Bus            *DevBus
	A              uint8         // Accumulator
	X              uint8         // X register
	Y              uint8         // Y register
	Status         uint8         // Status register
	SP             uint8         // Stack pointer
	PC             uint16        // Program counter
	fetched        uint8         //
	opLookup       []Instruction // Instruction lookup table
	cycles         uint8         // Number of remaining cycles for current instruction
	absAddr        uint16
	relAddr        uint16
	opcode         uint8
	addrModeLookup map[string]func(*MOS6502) uint8
}

// CPU is the primary interface for the 6502 emulator
// TODO: Is this interface really necessary?
type CPU interface {
	read(address uint16)
	write(address uint16, data uint8)
	reset()
	irq()
	nmi()
	fetch()
	clock()
}

// Flags
const (
	C uint8 = 1 << iota // Carry
	Z                   // Zero
	I                   // Disable interrupts
	D                   // Decimal mode (not used)
	B                   // Break
	U                   // Unused
	V                   // Overflow
	N                   // Negative
)

func (c *MOS6502) read(address uint16) uint8 {
	return c.Bus.Read(address, false)
}

func (c *MOS6502) write(address uint16, data uint8) {
	c.Bus.Write(address, data)
}

// GetFlag returns the value of the specified flag
func (c *MOS6502) GetFlag(f uint8) uint8 {
	return c.Status & f
}

// SetFlag sets the given bit in the status register if the condition evaluates to true,
// and clears the bit if the condition evaluates to false
func (c *MOS6502) SetFlag(f uint8, cond bool) {
	if cond {
		c.Status |= f
	} else {
		c.Status &^= f
	}
}

// Perform one clock cycle of computation
func (c *MOS6502) clock() {
	// When the cycle counter has reached 0, the instruction is complete and the next is ready
	// to be executed
	if c.cycles == 0 {
		c.opcode = c.read(c.PC)
		instruction := c.opLookup[c.opcode]
		c.PC++

		// Always set the unused flag to 1 (WHY?)
		c.SetFlag(U, true)

		// Set number of cycles for instruction
		c.cycles = instruction.cycles

		// Additional cycles for address mode
		addrModeCycles := c.addrModeLookup[instruction.addrMode](c)

		// Get additional instruction cycles
		addlOpCycles := instruction.op()

		c.cycles += (addrModeCycles & addlOpCycles)

		c.SetFlag(U, true)
	}

	c.cycles--
}

// Create6502 returns an instance of the CPU
func Create6502() MOS6502 {
	c := MOS6502{}

	// populate address mode lookup table
	c.addrModeLookup = map[string]func(*MOS6502) uint8{
		"IMP": (*MOS6502).imp,
		"IMM": (*MOS6502).imm,
		"ZP0": (*MOS6502).zp0,
		"ZPX": (*MOS6502).zpX,
		"ZPY": (*MOS6502).zpY,
		"REL": (*MOS6502).rel,
		"ABS": (*MOS6502).abs,
		"ABX": (*MOS6502).abX,
		"ABY": (*MOS6502).abY,
		"IND": (*MOS6502).ind,
		"IZX": (*MOS6502).izX,
		"IZY": (*MOS6502).izY,
	}

	return c
}

// Fetch retrieves data given an address and stores it in the instance variable "fetched" and
// returns it as well.
func (c *MOS6502) fetch() uint8 {
	if c.opLookup[c.opcode].addrMode != "IMM" {
		c.fetched = c.read(c.absAddr)
	}
	return c.fetched
}
