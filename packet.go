package main

type packet struct {
	buf  []byte
	size uint
	pos  uint
}

func NewPacket(buf []byte) *packet {
	return &packet{buf: buf, size: uint(len(buf)), pos: 0}
}

func (p *packet) ReadUint16() uint16 {
	i := uint16(p.buf[p.pos])<<8 | uint16(p.buf[p.pos+1])
	p.pos += 2
	return i
}

func (p *packet) ReadByte() byte {
	b := p.buf[p.pos]
	p.pos += 1
	return b
}

func (p *packet) ReadBytes(n uint) []byte {
	if n > p.size-p.pos {
		n = p.size - p.pos
	}
	b := p.buf[p.pos : p.pos+n]
	p.pos += n
	return b
}

func (p *packet) Seek(n uint) {
	if n < p.size {
		p.pos = n
	}
}
