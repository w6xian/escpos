```
package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/w6xian/escpos"
	"github.com/w6xian/printer"
)

func findDefaultPrinter() string {
	p, err := printer.Default()
	if err != nil {
		return ""
	}
	return p
}

func selectPrinter() (string, error) {
	printerId := findDefaultPrinter()
	n, err := strconv.Atoi(printerId)
	if err != nil {
		// must be a printer name
		return printerId, nil
	}
	printers, err := printer.ReadNames()
	if err != nil {
		return "", err
	}
	if n < 0 {
		return "", fmt.Errorf("printer index (%d) cannot be negative", n)
	}
	if n >= len(printers) {
		return "", fmt.Errorf("printer index (%d) is too large, there are only %d printers", n, len(printers))
	}
	return printers[n], nil
}

func main() {
	printerName, err := selectPrinter()
	if err != nil {
		os.Exit(0)
	}
	p, err := printer.Open(printerName)
	if err != nil {
		os.Exit(0)
	}
	defer p.Close()

	err = p.StartRawDocument("")
	if err != nil {
		os.Exit(0)
	}
	defer p.EndDocument()

	err = p.StartPage()
	if err != nil {
		os.Exit(0)
	}

	pr := escpos.New(
		escpos.DeviceType(escpos.PAPER_58),
		escpos.Printer(p),
	)

	pr.Begin()
	pr.Title("米粒工厂测试账套")
	pr.SubTitle("销售小票")
	pr.Content(func(p *escpos.Escpos) {
		p.InLine("单号:", "D12345")
		p.InLine("单据时间:", "2023-01-01")
		p.InLine("打印时间:", "2023-01-01")
		p.InLine("经手人:", "张三/18875028965")
	})
	pr.FillAround("订单详情", escpos.FillWith("-"))

	t := escpos.Table(
		escpos.Header(
			escpos.ColumnHeader("品名", 6),
			escpos.ColumnHeader("数量/重量", 9),
			escpos.ColumnHeader("单价", 8),
			escpos.ColumnHeader("小计", 9),
		),
		escpos.Row(
			escpos.ColumnData("啤酒", escpos.FULL_LINE_WIDTH, escpos.Position(escpos.POSITION_LEFT)),

			escpos.ColumnData("1", 9, escpos.Position(escpos.POSITION_CENTER)),
			escpos.ColumnData("50", 8, escpos.Position(escpos.POSITION_CENTER)),
			escpos.ColumnData("50", 9, escpos.Position(escpos.POSITION_RIGHT)),
		),
		escpos.Row(
			escpos.ColumnData("水", escpos.FULL_LINE_WIDTH, escpos.Position(escpos.POSITION_LEFT)),
			escpos.ColumnData("1", 9, escpos.Position(escpos.POSITION_CENTER)),
			escpos.ColumnData("50", 8, escpos.Position(escpos.POSITION_CENTER)),
			escpos.ColumnData("50", 9, escpos.Position(escpos.POSITION_RIGHT)),
		),
	)
	pr.PrintTable(t)
	pr.FillAround("", escpos.FillWith("-"))
	pr.Content(func(p *escpos.Escpos) {
		p.InLine("数量:", "2")
		p.InLine("订单金额:", "100")
		p.InLine("优惠金额:", "0")
		p.InLine("应付金额:", "100")
	})
	pr.Println("本文旨在向消费者提供一份消费者保护声明范本，以保障消费者的权益和提升购物体验。 以下是消费者保护声明的主要内容： 我们承诺提供准确、清晰的产品和服务信息，包括但不限于产品功能、规格、性能、制造商信息、有效期和售后服务等。 我们保证我们所提供的信息真实可靠，并且会尽力及时更新和修正。 在购买产品或使用服务前，请仔细阅读产品说明书和服务条款，并咨询我们的客户服务团队以获取更多信息")
	pr.Println("电话：18875028965")
	pr.Println("地址：北京市海淀区")
	pr.Feed()
	pr.FontAlign(escpos.AlignCenter)
	pr.QRCode("https://www.baidu.com", true, 8, 2)
	pr.Feed()
	pr.Println("扫描以上二维码，即可查看订单详情")
	pr.Println("谢谢惠顾")
	pr.Println("积分手机号：18875028965")
	pr.Println("本次积分：1000")

	pr.FeedN(3)
	pr.End()
	p.EndPage()
}
```