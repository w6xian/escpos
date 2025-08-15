package escpos

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func (e *Escpos) SetChineseOn() {
	e.Write("\x1C&")
}

type Escpos struct {
	// destination
	opts Options
}

// reset toggles
func (e *Escpos) reset() {
	e.opts.Width = 1
	e.opts.Height = 1

	e.opts.Underline = 0
	e.opts.Emphasize = 0
	e.opts.Upsidedown = 0
	e.opts.Rotate = 0

	e.opts.Reverse = 0
	e.opts.Smooth = 0
}

// create Escpos printer
func New(opts ...Option) (e *Escpos) {
	// 默认
	opt := newOpts(opts...)
	e = &Escpos{
		opts: *opt,
	}
	return
}

// write raw bytes to printer
func (e *Escpos) WriteRaw(data []byte) (n int, err error) {
	if len(data) > 0 {
		e.opts.Io.Write(data)
	}

	return 0, nil
}

// read raw bytes from printer
func (e *Escpos) ReadRaw(data []byte) (n int, err error) {
	return e.opts.Io.Read(data)
}

// write a string to the printer
func (e *Escpos) Write(data string) (int, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(data)), simplifiedchinese.GB18030.NewEncoder())
	bs, _ := io.ReadAll(reader)
	return e.WriteRaw(bs)
}

func (e *Escpos) Print(content string) (int, error) {

	return e.Write(content)
}

/**
 * 打印文字并换行
 * @param  {string}  content  文字内容
 */
func (e *Escpos) Println(content string) {
	e.Print(content + EOL)
}

// 开钱箱
func (e *Escpos) OpenDrawer() {
	e.WriteRaw([]byte{ESC, 0x70, byte(0), byte(10), byte(10)})
}

// init/reset printer settings
func (e *Escpos) Begin() {
	e.reset()
	e.Write("\x1B@")
}

// end output
func (e *Escpos) End() {
	e.Write("\xFA")
}

// send cut
func (e *Escpos) Cut() {
	e.Write("\x1DVA0")
}

// send cut minus one point (partial cut)
func (e *Escpos) CutPartial() {
	e.WriteRaw([]byte{GS, 0x56, 1})
}

// send cash
func (e *Escpos) Cash() {
	e.Write("\x1B\x70\x00\x0A\xFF")
}

// send linefeed
func (e *Escpos) Linefeed() {
	e.Write(EOL)
}

// send N formfeeds
func (e *Escpos) FormfeedN(n int) {
	e.Write(fmt.Sprintf("\x1Bd%c", n))
}

// send formfeed
func (e *Escpos) Formfeed() {
	e.FormfeedN(1)
}

// Feed skip one line of paper
func (e *Escpos) Feed() {
	e.WriteRaw([]byte{LF})
}

// Feed skip n lines of paper
func (e *Escpos) FeedN(n byte) {
	e.WriteRaw([]byte{ESC, 0x64, n})
}

// SelfTest start self test of printer
func (e *Escpos) SelfTest() {
	e.WriteRaw([]byte{GS, 0x28, 0x41, 0x02, 0x00, 0x00, 0x02})
}

// set font
func (e *Escpos) SetFont(font string) {
	f := 0

	switch font {
	case "A":
		f = 0
	case "B":
		f = 1
	case "C":
		f = 2
	default:
		f = 0
	}

	e.Write(fmt.Sprintf("\x1BM%c", f))
}

func (e *Escpos) SendFontSize() {
	e.Write(fmt.Sprintf("\x1D!%c", ((e.opts.Width)<<4)|(e.opts.Height)))
}
func (e *Escpos) SetFontStyle(style uint8) {
	e.Write(string([]byte{ESC, 0x21, byte(style)}))
}

func (e *Escpos) SetLetterSpace(n int) {
	e.Write(string([]byte{ESC, 0x20, byte(n)}))
}

// set font size
func (e *Escpos) SetFontSize(width, height uint8) {
	if width < 8 && height < 8 {
		e.opts.Width = width
		e.opts.Height = height
		e.SendFontSize()
	}
}

func (e *Escpos) SetFontColor(color uint8) {
	e.WriteRaw([]byte{ESC, 0x72, byte(color)})
}

// send underline
func (e *Escpos) SetUnderline() {
	e.Write(fmt.Sprintf("\x1B-%c", e.opts.Underline))
}

// send emphasize / doublestrike
func (e *Escpos) SendEmphasize() {
	e.Write(fmt.Sprintf("\x1BG%c", e.opts.Emphasize))
}

// send upsidedown
func (e *Escpos) SendUpsidedown() {
	e.Write(fmt.Sprintf("\x1B{%c", e.opts.Upsidedown))
}

// send rotate
func (e *Escpos) SendRotate() {
	e.Write(fmt.Sprintf("\x1BR%c", e.opts.Rotate))
}

// send reverse
func (e *Escpos) SendReverse() {
	e.Write(fmt.Sprintf("\x1DB%c", e.opts.Reverse))
}

// send smooth
func (e *Escpos) SendSmooth() {
	e.Write(fmt.Sprintf("\x1Db%c", e.opts.Smooth))
}

// 光标移动到x位置
func (e *Escpos) SendMoveX(x int) {
	e.Write(string([]byte{ESC, 0x24, byte(x % 256), byte(x / 256)}))
}

// send move y
func (e *Escpos) SendMoveY(y int) {
	e.Write(string([]byte{0x1d, 0x24, byte(y % 256), byte(y / 256)}))
}

// set emphasize
func (e *Escpos) SetEmphasize(u uint8) {
	e.opts.Emphasize = u
	e.SendEmphasize()
}

// set upsidedown
func (e *Escpos) SetUpsidedown(v uint8) {
	e.opts.Upsidedown = v
	e.SendUpsidedown()
}

// set rotate
func (e *Escpos) SetRotate(v uint8) {
	e.opts.Rotate = v
	e.SendRotate()
}

// set reverse
func (e *Escpos) SetReverse(v uint8) {
	e.opts.Reverse = v
	e.SendReverse()
}

// set smooth
func (e *Escpos) SetSmooth(v uint8) {
	e.opts.Smooth = v
	e.SendSmooth()
}

// pulse (open the drawer)
func (e *Escpos) Pulse() {
	// with t=2 -- meaning 2*2msec
	e.Write("\x1Bp\x02")
}

// set alignment
func (e *Escpos) SetAlign(align string) {
	a := 0
	switch align {
	case "left":
		a = 0
	case "center":
		a = 1
	case "right":
		a = 2
	}
	e.Write(fmt.Sprintf("\x1Ba%c", a))
}

func (e *Escpos) SetMarginLeft(size uint16) {
	if size <= 47 {
		e.Write(string([]byte{0x1d, 0x4c, byte(size % 256), byte(size / 256)}))
	}
}

// set language -- ESC R
func (e *Escpos) SetLang(lang string) {
	l := 0

	switch lang {
	case "en":
		l = 0
	case "fr":
		l = 1
	case "de":
		l = 2
	case "uk":
		l = 3
	case "da":
		l = 4
	case "sv":
		l = 5
	case "it":
		l = 6
	case "es":
		l = 7
	case "ja":
		l = 8
	case "no":
		l = 9
	}
	e.Write(fmt.Sprintf("\x1BR%c", l))
}

// feed and cut based on parameters
func (e *Escpos) FeedAndCut(params map[string]string) {
	if t, ok := params["type"]; ok && t == "feed" {
		e.Formfeed()
	}

	e.Cut()
}

// Barcode sends a barcode to the printer.
func (e *Escpos) Barcode(barcode string, format int) {
	code := ""
	switch format {
	case 0:
		code = "\x00"
	case 1:
		code = "\x01"
	case 2:
		code = "\x02"
	case 3:
		code = "\x03"
	case 4:
		code = "\x04"
	case 73:
		code = "\x49"
	}

	// reset settings
	e.reset()

	// set align
	e.SetAlign("center")

	// write barcode
	if format > 69 {
		e.Write(fmt.Sprintf("\x1dk"+code+"%v%v", len(barcode), barcode))
	} else if format < 69 {
		e.Write(fmt.Sprintf("\x1dk"+code+"%v\x00", barcode))
	}
	e.Write(fmt.Sprintf("%v", barcode))
}

// Prints a QR Code.
// code specifies the data to be printed
// model specifies the qr code model. false for model 1, true for model 2
// size specifies the size in dots. It needs to be between 1 and 16
// QRCode("https://www.baidu.com", true, 8, 2)二维码居中
func (e *Escpos) QRCode(code string, model bool, size uint8, correctionLevel uint8) (int, error) {
	if len(code) > 7089 {
		return 0, fmt.Errorf("the code is too long, it's length should be smaller than 7090")
	}
	if size < 1 {
		size = 1
	}
	if size > 16 {
		size = 16
	}
	var m byte = 49
	var err error
	// set the qr code model
	if model {
		m = 50
	}
	_, err = e.WriteRaw([]byte{GS, '(', 'k', 4, 0, 49, 65, m, 0})
	if err != nil {
		return 0, err
	}

	// set the qr code size
	_, err = e.WriteRaw([]byte{GS, '(', 'k', 3, 0, 49, 67, size})
	if err != nil {
		return 0, err
	}

	// set the qr code error correction level
	if correctionLevel < 48 {
		correctionLevel = 48
	}
	if correctionLevel > 51 {
		correctionLevel = 51
	}
	_, err = e.WriteRaw([]byte{GS, '(', 'k', 3, 0, 49, 69, size})
	if err != nil {
		return 0, err
	}

	// store the data in the buffer
	// we now write stuff to the printer, so lets save it for returning

	// pL and pH define the size of the data. Data ranges from 1 to (pL + pH*256)-3
	// 3 < pL + pH*256 < 7093
	var codeLength = len(code) + 3
	var pL, pH byte
	pH = byte(int(math.Floor(float64(codeLength) / 256)))
	pL = byte(codeLength - 256*int(pH))

	written, err := e.WriteRaw(append([]byte{GS, '(', 'k', pL, pH, 49, 80, 48}, []byte(code)...))
	if err != nil {
		return written, err
	}

	// finally print the buffer
	_, err = e.WriteRaw([]byte{GS, '(', 'k', 3, 0, 49, 81, 48})
	if err != nil {
		return written, err
	}

	return written, nil
}

// used to send graphics headers
func (e *Escpos) gSend(m byte, fn byte, data []byte) {
	l := len(data) + 2

	e.Write("\x1b(L")
	e.WriteRaw([]byte{byte(l % 256), byte(l / 256), m, fn})
	e.WriteRaw(data)
}

// ReadStatus Read the status n from the printer
func (e *Escpos) ReadStatus(n byte) (byte, error) {
	e.WriteRaw([]byte{DLE, EOT, n})
	data := make([]byte, 1)
	_, err := e.ReadRaw(data)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (e *Escpos) Content(f func(p *Escpos)) {
	e.Font(FontA)
	e.FontAlign(AlignLeft)
	e.FontSize(1, 1)
	e.FontBold(false)
	f(e)
	e.Feed()
}

func (e *Escpos) InLine(str1, str2 string, opts ...FillOption) (int, error) {
	opt := newFillOptions(opts...)
	fontWidth := opt.FontWidth
	fillWith := opt.FillWith
	position := opt.Position
	return e.Print(Inline(e.opts.MaxChar, str1, str2, fillWith, fontWidth, int(position)))
}

func (e *Escpos) FillAround(str1 string, opts ...FillOption) (int, error) {
	opt := newFillOptions(opts...)
	fontWidth := opt.FontWidth
	fillWith := opt.FillWith
	return e.Print(fillAround(e.opts.MaxChar, str1, fillWith, fontWidth))
}

func (e *Escpos) Divider(opts ...FillOption) (int, error) {
	opt := newFillOptions(opts...)
	fontWidth := opt.FontWidth
	fillWith := opt.FillWith
	return e.Print(fillAround(e.opts.MaxChar, fillWith, fillWith, fontWidth))
}
