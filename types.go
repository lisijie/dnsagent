package main

const (
	TypeA uint16 = 0x01 //指定计算机 IP 地址。

	TypeNS uint16 = 0x02 //指定用于命名区域的 DNS 名称服务器。

	TypeMD uint16 = 0x03 //指定邮件接收站(此类型已经过时了，使用MX代替)

	TypeMF uint16 = 0x04 //指定邮件中转站(此类型已经过时了，使用MX代替)

	TypeCNAME uint16 = 0x05 //指定用于别名的规范名称。

	TypeSOA uint16 = 0x06 //指定用于 DNS 区域的“起始授权机构”。

	TypeMB uint16 = 0x07 //指定邮箱域名。

	TypeMG uint16 = 0x08 //指定邮件组成员。

	TypeMR uint16 = 0x09 //指定邮件重命名域名。

	TypeNULL uint16 = 0x0A //指定空的资源记录

	TypeWKS uint16 = 0x0B //描述已知服务。

	TypePTR uint16 = 0x0C //如果查询是 IP 地址，则指定计算机名;否则指定指向其它信息的指针。

	TypeHINFO uint16 = 0x0D //指定计算机 CPU 以及操作系统类型。

	TypeMINFO uint16 = 0x0E //指定邮箱或邮件列表信息。

	TypeMX uint16 = 0x0F //指定邮件交换器。

	TypeTXT uint16 = 0x10 //指定文本信息。

	TypeAAAA uint16 = 0x1c //IPV6资源记录。

	TypeUINFO uint16 = 0x64 //指定用户信息。

	TypeUID uint16 = 0x65 //指定用户标识符。

	TypeGID uint16 = 0x66 //指定组名的组标识符。

	TypeANY uint16 = 0xFF //指定所有数据类型。
)

const (
	ClassIN uint16 = 0x01 //指定 Internet 类别。

	ClassCSNET uint16 = 0x02 //指定 CSNET 类别。(已过时)

	ClassCHAOS uint16 = 0x03 //指定 Chaos 类别。

	ClassHESIOD uint16 = 0x04 //指定 MIT Athena Hesiod 类别。
)
