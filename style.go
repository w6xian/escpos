package escpos

import (
	"fmt"
)

type fontfamily byte

const (
	FontA fontfamily = 0
	FontB fontfamily = 1
)

func (e *Escpos) Font(family fontfamily) {
	e.WriteRaw([]byte{ESC, 0x4D, byte(family)})
}

type fontalign byte

const (
	AlignLeft   fontalign = 0
	AlignCenter fontalign = 1
	AlignRight  fontalign = 2
)

func (e *Escpos) FontAlign(align fontalign) {
	e.WriteRaw([]byte{ESC, 0x61, byte(align)})
}

func (e *Escpos) FontSize(width, height uint8) {
	if width >= 1 && width <= 8 && height >= 1 && height <= 8 {
		e.WriteRaw([]byte{GS, 0x21, ((width - 1) << 4) | (height - 1)})
	} else {
		panic(fmt.Sprintf("Wrong font size: (%d x %d)", width, height))
	}
}

func (e *Escpos) FontUnderline(on bool) {
	e.WriteRaw([]byte{ESC, 0x2D, boolToByte(on)})
}

func (e *Escpos) FontBold(on bool) {
	e.WriteRaw([]byte{ESC, 0x45, boolToByte(on)})
}

func boolToByte(b bool) byte {
	var r byte
	if b {
		r = byte(1)
	}

	return r
}

func (e *Escpos) Title(title string) {
	e.Font(FontB)
	e.FontAlign(AlignCenter)
	e.FontSize(2, 2)
	e.FontBold(true)
	e.Write(title)
	e.FeedN(2)
}

func (e *Escpos) SubTitle(sub string) {
	e.Font(FontA)
	e.FontAlign(AlignCenter)
	e.FontSize(1, 1)
	e.FontBold(false)
	e.Write(sub)
	e.FeedN(2)
}
