package main

// 包头，12字节
type Header struct {
	Id      uint16 // ID
	Flags   uint16 // 标识
	Qdcount uint16 // 报文请求段中的问题记录数
	Ancount uint16 // 报文回答段中的回答记录数
	Nscount uint16 // 报文授权段中的授权记录数
	Arcount uint16 // 报文附加段中的附加记录数
}

type MsgFlags struct {
	Response           bool // QR:长度1位，值0是请求，1是应答
	Opcode             int  // 长度4位，值0是标准查询，1是反向查询，2死服务器状态查询。
	Authoritative      bool // 长度1位，授权应答(Authoritative Answer) – 这个比特位在应答的时候才有意义，指出给出应答的服务器是查询域名的授权解析服务器。
	Truncated          bool // 长度1位，截断(TrunCation) – 用来指出报文比允许的长度还要长，导致被截断。
	RecursionDesired   bool // 长度1位，期望递归(Recursion Desired) – 这个比特位被请求设置，应答的时候使用的相同的值返回。如果设置了RD，就建议域名服务器进行递归解析，递归查询的支持是可选的。
	RecursionAvailable bool // 长度1位，支持递归(Recursion Available) – 这个比特位在应答中设置或取消，用来代表服务器是否支持递归查询。
	Zero               bool // 长度3位，保留值，值为0.
	Rcode              int  // 长度4位，返回码，通常为0(没有差错)和3(名字差错)
}

// Msg contains the layout of a DNS message.
type Msg struct {
	header *Header
}

func UnpackMsg(buf []byte) *Msg {
	pkt := NewPacket(buf)
	msg := new(Msg)
	msg.header = unpackHeader(pkt)

	return msg
}

func unpackHeader(pkt *packet) *Header {
	hd := new(Header)
	hd.Id = pkt.readUint16()
	hd.Flags = pkt.readUint16()
	hd.Qdcount = pkt.readUint16()
	hd.Nscount = pkt.readUint16()
	hd.Arcount = pkt.readUint16()
	return hd
}

type packet struct {
	buf  []byte
	size uint
	pos  uint
}

func NewPacket(buf []byte) *packet {
	return &packet{buf: buf, size: uint(len(buf)), pos: 0}
}

func (p *packet) readUint16() uint16 {
	i := uint16(p.buf[p.pos])<<8 | uint16(p.buf[p.pos+1])
	p.pos += 2
	return i
}

func (p *packet) readByte() byte {
	b := p.buf[p.pos]
	p.pos += 1
	return b
}

func (p *packet) readBytes(n uint) []byte {
	if n > p.size-p.pos {
		n = p.size - p.pos
	}
	b := p.buf[p.pos : p.pos+n]
	p.pos += n
	return b
}

func (p *packet) seek(n uint) {
	if n < p.size {
		p.pos = n
	}
}
