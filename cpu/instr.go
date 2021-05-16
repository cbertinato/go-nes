package cpu

// Legend
// ------
// *  add 1 to cycles if page boundery is crossed

// ** add 1 to cycles if branch occurs on same page
//    add 2 to cycles if branch occurs to different page

//    Legend to Flags:  + .... modified
// 					    - .... not modified
// 					    1 .... set
// 					    0 .... cleared
// 					   M6 .... memory bit 6
// 					   M7 .... memory bit 7

// ADC  Add Memory to Accumulator with Carry
// -----------------------------------------
//      A + M + C -> A, C                N Z C I D V
//                                       + + + - - +
//
//      addressing    assembler    opc  bytes  cyles
//      --------------------------------------------
//      immediate     ADC #oper     69    2     2
//      zeropage      ADC oper      65    2     3
//      zeropage,X    ADC oper,X    75    2     4
//      absolute      ADC oper      6D    3     4
//      absolute,X    ADC oper,X    7D    3     4*
//      absolute,Y    ADC oper,Y    79    3     4*
//      (indirect,X)  ADC (oper,X)  61    2     6
//      (indirect),Y  ADC (oper),Y  71    2     5*

// Logic for first signifcant bits of accumulator (A) + memory (M) = result (R)
//
// A  M  R  V  A^R & ~(A^M)
// ------------------------
// 0  0  0  0  0
// 0  0  1  1  1
// 0  1  0  0  0
// 0  1  1  0  0
// 1  0  0  0  0
// 1  0  1  0  0
// 1  1  0  1  1
// 1  1  1  0  0
func (c *MOS6502) adc() uint8 {
	c.fetch()

	accum := uint16(c.A)
	mem := uint16(c.fetched)
	carry := uint16(c.GetFlag(C))
	res := accum + mem + carry
	// ANDing with 0x80 extracts the sign bit
	v := ((accum ^ res) &^ (accum ^ mem)) & 0x80

	// set flags
	c.SetFlag(Z, res&0x00ff == 0)
	c.SetFlag(C, res > 0x00ff)
	c.SetFlag(N, res&0x80 != 0)
	c.SetFlag(V, v != 0)

	c.A = uint8(res & 0x00FF)

	// this operation could potentially get an extra cycle
	return 1
}

// SBC subtract memory from accumulator with borrow
// ------------------------------------------------
// The 6502 has a SBC operation (subtract with carry) that subtracts two numbers
// and also subtracts the borrow bit. If the (unsigned) operation results in a
// borrow (is negative), then the borrow bit is set. However, there is no explicit
// borrow flag - instead the complement of the carry flag is used. If the carry
// flag is 1, then borrow is 0, and if the carry flag is 0, then borrow is 1.
// This behavior may seem backwards, but note that both for addition and subtraction,
// if the carry flag is set, the output is one more than if the carry flag is clear.
//
// A = A - M - B
// A = A - M - B + 256		Add 256, which doesn't change the 8-bit value.
// A = A - M - (1-C) + 256	Replace B with the inverted carry flag.
// A = A + (255-M) + C		Simple algebra.
// A = A + M ^ 0x00ff + C	255 - M is the same as flipping the bits

func (c *MOS6502) sbc() uint8 {
	c.fetch()

	accum := uint16(c.A)
	mem := uint16(c.fetched) ^ 0x00ff // invert the value in memory
	carry := uint16(c.GetFlag(C))
	res := accum + mem + carry
	v := ((accum ^ res) &^ (accum ^ mem)) & 0x80

	// set flags
	c.SetFlag(Z, res&0x00ff == 0)
	c.SetFlag(C, res > 0x00ff)
	c.SetFlag(N, res&0x80 != 0)
	c.SetFlag(V, v != 0)

	c.A = uint8(res & 0x00ff)

	// this operation could potentially get an extra cycle
	return 1
}
