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

func (p *packet) WriteByte(b byte) bool {
	p.buf = append(p.buf, b)
	p.size++
	return true
}

func (p *packet) WriteUint16(n uint16) bool {
	p.buf = append(p.buf, byte(n>>8&0xFF), byte(n&0xFF))
	p.size += 2
	return true
}

func (p *packet) WriteBytes(bs []byte) bool {
	p.buf = append(p.buf, bs...)
	p.size += len(bs)
	return true
}

func (p *packet) WriteString(s string) bool {
	return p.WriteBytes([]byte(s))
}

func (p *packet) Seek(n uint) {
	if n < p.size {
		p.pos = n
	}
}

func (p *packet) Bytes() {
	return p.buf
}
