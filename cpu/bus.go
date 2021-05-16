package cpu

// Bus interface
type Bus interface {
    Read(address uint16, readOnly bool) uint8
    Write(address uint16, data uint8)
}

// DevBus is a simple bus that consists only of RAM
type DevBus struct {
    ram [64 * 1024]uint8 // 64k of RAM
}

func (b *DevBus) Read(address uint16, readOnly bool) uint8 {
    if address >= 0 && address < 0xFFFF {
        return b.ram[address]
    }
    return 0
}

func (b *DevBus) Write(address uint16, data uint8) {
    if address >= 0 && address < 0xFFFF {
        b.ram[address] = data
    }
}