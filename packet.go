package main

// 二进制数据包读写

type Packet struct {
	buf  []byte
	size uint
	pos  uint
}

func NewPacket(buf []byte) *Packet {
	return &Packet{buf: buf, size: uint(len(buf)), pos: 0}
}

func (p *Packet) ReadUint16() uint16 {
	i := uint16(p.buf[p.pos])<<8 | uint16(p.buf[p.pos+1])
	p.pos += 2
	return i
}

func (p *Packet) ReadByte() byte {
	b := p.buf[p.pos]
	p.pos += 1
	return b
}

func (p *Packet) ReadBytes(n uint) []byte {
	if n > p.size-p.pos {
		n = p.size - p.pos
	}
	b := p.buf[p.pos : p.pos+n]
	p.pos += n
	return b
}

func (p *Packet) WriteByte(b byte) bool {
	p.buf = append(p.buf, b)
	p.size++
	return true
}

func (p *Packet) WriteUint16(n uint16) bool {
	p.buf = append(p.buf, byte(n>>8&0xFF), byte(n&0xFF))
	p.size += 2
	return true
}

func (p *Packet) WriteUint(n uint) bool {
	p.buf = append(p.buf, byte(n>>24&0xFF), byte(n>>16&0xFF), byte(n>>8&0xFF), byte(n&0xFF))
	p.size += 4
	return true
}

func (p *Packet) WriteBytes(bs []byte) bool {
	p.buf = append(p.buf, bs...)
	p.size += uint(len(bs))
	return true
}

func (p *Packet) WriteString(s string) bool {
	return p.WriteBytes([]byte(s))
}

func (p *Packet) Seek(n uint) {
	if n < p.size {
		p.pos = n
	}
}

func (p *Packet) Bytes() []byte {
	return p.buf
}
