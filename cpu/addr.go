package cpu

// Addressing mode functions
// -------------------------
// The 6502 can address from 0x0000 to 0xFFFF. The hi byte refers to the page and the lo byte to
// the offset into that page.  Therefore, there can be up to 256 pages and 256 bytes per page. Most
// addressing mode functions place the address of the data to be fetched into the absAddr variable. 
// Some instructions don't require fetched data as the source is implied by the instruction.
// There are 12 addressing modes available on the 6502:
//     - implied (IMP)	
//     - immediate (IMM)
//     - zero page addressing (ZP0)
//     - indexed zero page addressing (ZPX, ZPY)
//     - relative (REL)
//     - absolute (ABS)
//     - index absolute addressing (ABX, ABY)
//     - absolute indirect addressing (IND)
//     - indexed indirect addressing (IZX)
//     - indirect indexed addressing (IZY)

// IMM (immediate): The next byte is to be used as a value.
func (c *MOS6502) imm() uint8 {
	c.fetched = c.A
	return 0
}

// IMP (implied): In the implied addressing mode, the address containing the operand is implicitly stated
// in the operation code of the instruction.
func (c *MOS6502) imp() uint8 {
	c.absAddr = c.PC
	c.PC++
	return 0
}

// ZP0 (zero page addressing): Fetches only the second byte of the instruction and assumes that the high byte
// address byte is 0.
func (c *MOS6502) zp0() uint8 {
	// read 8 bits and convert to a 16-bit address
	c.absAddr = uint16(c.read(c.PC) & 0x00FF)
	c.PC++
	return 0
}

// ZPX (indexed zero page): Essentially the same as zero page addressing, the X register is used as the index. The 
// second byte is added to the contents of the X register. Due to the "zero page" nature of this mode, a carry is 
// not added to the high 8 bits of the address and crossing of page boundaries does not occur.
func (c *MOS6502) zpX() uint8 {
	// read 8 bits and convert to a 16-bit address (take just the low byte just in case)
	c.absAddr = uint16((c.read(c.PC) + c.X) & 0x00FF)
	c.PC++
	return 0
}

// ZPY
func (c *MOS6502) zpY() uint8 {
	// read 8 bits and convert to a 16-bit address (take just the low byte just in case)
	c.absAddr = uint16((c.read(c.PC) + c.Y) & 0x00FF)
	c.PC++
	return 0
}

// REL (relative): Used only with branch instructions. The second byte of the instruction is the offset added to 
// the program counter when the counter is set at the next instruction. The offset must be between -128 to +127 bytes.
func (c *MOS6502) rel() uint8 {
	c.relAddr = uint16(c.read(c.PC))
	c.PC++
	
	// relAddr < 128d => relAddr & 0x80 == 0
	if c.relAddr & 0x80 != 0 {
		// it's a negative number -> flip the msb on the least significant byte
		c.relAddr |= 0xFF00
	}
	return 0
}

// ABS (absolute): The second byte of the instruction specifies the 8 low bits of the address and the third byte the 
// 8 high bits. Thus, the absolute addressing mode allows access to the entire 64k bytes of memory.
func (c *MOS6502) abs() uint8 {
	lo := uint16(c.read(c.PC))
	c.PC++
	hi := uint16(c.read(c.PC))
	c.PC++
	c.absAddr = hi << 8 | lo

	return 0
}

// ABX: Same as absolute addressing but the X register is added to the address. If the resulting address is in the next
// page, then add a clock cycle.
func (c *MOS6502) abX() uint8 {
	lo := uint16(c.read(c.PC))
	c.PC++
	hi := uint16(c.read(c.PC))
	c.PC++
	c.absAddr = (hi << 8 | lo)
	c.absAddr += uint16(c.X)

	if (c.absAddr & 0xFF00) != (hi << 8) {
		return 1
	}

	return 0
}

func (c *MOS6502) abY() uint8 {
	lo := uint16(c.read(c.PC))
	c.PC++
	hi := uint16(c.read(c.PC))
	c.PC++
	c.absAddr = (hi << 8 | lo)
	c.absAddr += uint16(c.Y)

	if (c.absAddr & 0xFF00) != (hi << 8) {
		return 1
	}

	return 0
}

// IND (absolute indirect addressing): The second byte of the instruction contains the 8 low bits of a 
// memory location. The third byte contains the 8 high bits. The contents of the address are 
// the low 8 bits of the effective address. The high 8 bits are at the next memory location.
//
// There is a bug in the chip for this addressing mode. If the low byte of the supplied address 
// is 0xFF, then to read the high byte of the actual address we need to cross a page boundary. 
// This doesnt actually work on the chip as designed, instead it wraps back around in the same 
// page, yielding an invalid actual address.
func (c *MOS6502) ind() uint8 {
	lo := uint16(c.read(c.PC))
	c.PC++
	hi := uint16(c.read(c.PC))
	c.PC++

	ptr := (hi << 8 | lo)

	if lo == 0x00FF { // buggy behavior
		c.absAddr = uint16(uint16(c.read(ptr & 0xFF00)) << 8 | uint16(c.read(ptr)))
	} else { // normal behavior
		c.absAddr = uint16(uint16(c.read(ptr + 1)) << 8 | uint16(c.read(ptr)))
	}
	
	return 0
}

// IZX (indexed indirect addressing): The second byte of the instruction is added to the 
// contents of the X register, discarding the carry. The result of this addition points 
// to a memory location on page zero whose contents is the low 8 bits of the effective 
// address. The next memory location contains the high 8 bits.
func (c *MOS6502) izX() uint8 {
	val := uint16(c.read(c.PC))
	c.PC++
	ptr := (val + uint16(c.X)) & 0x00FF

	lo := uint16(c.read(ptr))
	hi := uint16(c.read(ptr + 1))

	c.absAddr = (hi << 8) | lo

	return 0
}

// IZY (indirect indexed addressing): The second byte of the instruction points to a 
// memory location in the zero page. The contents of this memory location are added 
// to the Y register. The result is the low 8 bits of the effective address. The 
// contents of the carry are added to the contents of the contents of the next 
// memory location. If the carry causes a change in page, then a cycle is added. The 
// result is the high 8 bits of the effective address. 
func (c *MOS6502) izY() uint8 {
	val := uint16(c.read(c.PC)) & 0x00FF
	lo := uint16(c.read(val))
	hi := uint16(c.read(val + 1))

	c.absAddr = (hi << 8) | lo
	c.absAddr += uint16(c.Y)

	if c.absAddr & 0xFF00 != (hi << 8) {
		return 1
	}
	return 0
}