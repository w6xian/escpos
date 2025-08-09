package escpos

import "io"

type Options struct {
	DeviceType int
	Io         io.ReadWriter
	// font metrics
	Width, Height uint8

	// state toggles ESC[char]
	Underline, Emphasize, Upsidedown, Rotate uint8

	// state toggles GS[char]
	Reverse, Smooth uint8
	// paper metrics
	PaperWidth, MaxChar, LineHeight int
}

func newOpts(opts ...Option) *Options {
	opt := &Options{
		DeviceType: 80,
		Width:      1,
		Height:     1,

		Underline:  0,
		Emphasize:  0,
		Upsidedown: 0,
		Rotate:     0,

		Reverse:    0,
		Smooth:     0,
		PaperWidth: 576,
		MaxChar:    48,
		LineHeight: 24,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type Option func(*Options)

type Paper int

const (
	PAPER_58 Paper = 58
	PAPER_80 Paper = 80
)

// Server to be used for service.
func DeviceType(width Paper) Option {
	return func(o *Options) {
		if width == 80 {
			o.PaperWidth = 576
			o.MaxChar = 48
			o.LineHeight = 24
		} else {
			o.PaperWidth = 384
			o.MaxChar = 32
			o.LineHeight = 24
		}

		o.DeviceType = int(width)
	}
}

func Printer(pt io.ReadWriter) Option {
	return func(o *Options) {
		o.Io = pt
	}
}

const (
	POSITION_LEFT   = 0
	POSITION_RIGHT  = 1
	POSITION_CENTER = 2
)

type FillOptions struct {
	FillWith  string
	FontWidth int
	Position  int
	Width     int
}

type FillOption func(*FillOptions)

func newFillOptions(opts ...FillOption) *FillOptions {
	opt := &FillOptions{
		FillWith:  " ",
		FontWidth: 1,
		Position:  POSITION_LEFT,
		Width:     -1, // 没有宽度，根据内容自动调整

	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func FillWith(fillWith string) FillOption {
	return func(o *FillOptions) {
		o.FillWith = fillWith
	}
}

func FontWidth(fontWidth int) FillOption {
	return func(o *FillOptions) {
		o.FontWidth = fontWidth
	}
}

func Position(position int) FillOption {
	return func(o *FillOptions) {
		o.Position = position
	}
}

func Width(width int) FillOption {
	return func(o *FillOptions) {
		o.Width = width
	}
}

// 单行文本，不换行
// 宽度为-2时，不换行
// 宽度为-1时，根据内容自动调整宽度
const FULL_LINE_WIDTH = -2

func FullLine() FillOption {
	return func(o *FillOptions) {
		o.Width = FULL_LINE_WIDTH
	}
}
