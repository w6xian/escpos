package escpos

import (
	"math"
	"regexp"
	"strings"
)

const PAGE_WIDTH = 384
const MAX_CHAR_COUNT_EACH_LINE = 32

func isChinese(str string) bool {
	return regexp.MustCompile("^[\u4e00-\u9fa5]$").MatchString(str)
}

func isEnglish(str string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9]*$").MatchString(str)
}

/**
 * 返回字符串宽度(1个中文=2个英文字符)
 * @param str
 * @returns {number}
 */
func getStringWidth(str string) int {
	width := 0
	strs := strings.Split(str, "")
	for _, char := range strs {
		if isChinese(char) {
			// 中文字符宽度为2
			width += 2
		} else {
			// 英文字符宽度为1
			width += 1
		}
	}
	return width
}

/**
 * 同一行输出str1, str2，str1居左, str2居右
 * @param {string} str1 内容1
 * @param {string} str2 内容2
 * @param {number} fontWidth 字符宽度 1/2
 * @param {string} fillWith str1 str2之间的填充字符
 *
 */
func Inline(maxChar int, str1, str2 string, fillWith string, fontWidth int, pos int) string {
	lineWidth := maxChar / fontWidth
	str1Width := getStringWidth(str1)
	str2Width := getStringWidth(str2)

	// 需要填充的字符数量
	fillCount := lineWidth - (str1Width+str2Width)%lineWidth
	fillStr := strings.Repeat(fillWith, fillCount)
	// 内容已经超过一行了，没必要填充
	if getStringWidth(str1+fillStr+str2) > lineWidth {
		return str1 + str2
	}
	if pos == POSITION_LEFT {
		return str1 + str2 + fillStr
	} else if pos == POSITION_CENTER {
		leftCount := int(math.Round(float64(fillCount) / 2))
		// 两侧的填充字符，需要考虑左边需要填充，右边不需要填充的情况
		fillStr := strings.Repeat(fillWith, leftCount)
		return str1 + fillStr + str2 + fillStr[0:fillCount-leftCount]
	}

	return str1 + fillStr + str2
}

/**
 * 用字符填充一整行
 * @param {string} fillWith 填充字符
 * @param {number} fontWidth 字符宽度 1/2
 */
func fillLine(fillWith string, fontWidth int) string {
	lineWidth := MAX_CHAR_COUNT_EACH_LINE / fontWidth
	return strings.Repeat(fillWith, lineWidth)
}

/**
 * 文字内容居中，左右用字符填充
 * @param {string} str 文字内容
 * @param {number} fontWidth 字符宽度 1/2
 * @param {string} fillWith str1 str2之间的填充字符
 */
func fillAround(maxChar int, str string, fillWith string, fontWidth int) string {
	lineWidth := maxChar / fontWidth
	strWidth := getStringWidth(str)
	// 内容已经超过一行了，没必要填充
	if strWidth >= lineWidth {
		return str
	}
	// 需要填充的字符数量
	fillCount := lineWidth - strWidth
	// 左侧填充的字符数量
	leftCount := int(math.Round(float64(fillCount) / 2))
	// 两侧的填充字符，需要考虑左边需要填充，右边不需要填充的情况
	fillStr := strings.Repeat(fillWith, leftCount)
	return fillStr + str + fillStr[0:fillCount-leftCount]
}

func fillColumn(maxChar int, str string, fillWith string, fontWidth int, pos int) string {
	lineWidth := maxChar / fontWidth
	strWidth := getStringWidth(str)
	// 内容已经超过一行了，没必要填充
	if strWidth >= lineWidth {
		return str
	}
	// 需要填充的字符数量
	fillCount := lineWidth - strWidth
	if pos == POSITION_CENTER {
		// 左侧填充的字符数量
		leftCount := int(math.Round(float64(fillCount) / 2))
		// 两侧的填充字符，需要考虑左边需要填充，右边不需要填充的情况
		fillStr := strings.Repeat(fillWith, leftCount)
		return fillStr + str + fillStr[0:fillCount-leftCount]
	} else if pos == POSITION_RIGHT {
		// 右侧填充的字符数量
		rightCount := int(math.Round(float64(fillCount)))
		// 右侧的填充字符，需要考虑右边需要填充，左边不需要填充的情况
		fillStr := strings.Repeat(fillWith, rightCount)
		return fillStr + str
	}
	// 左侧填充的字符数量
	leftCount := int(math.Round(float64(fillCount)))
	// 左侧的填充字符，需要考虑左边需要填充，右边不需要填充的情况
	fillStr := strings.Repeat(fillWith, leftCount)
	return str + fillStr
}

// text replacement map
var textReplaceMap = map[string]string{
	// horizontal tab
	"&#9;":  "\x09",
	"&#x9;": "\x09",

	// linefeed
	"&#10;": "\n",
	"&#xA;": "\n",

	// xml stuff
	"&apos;": "'",
	"&quot;": `"`,
	"&gt;":   ">",
	"&lt;":   "<",

	// ampersand must be last to avoid double decoding
	"&amp;": "&",
}

// replace text from the above map
func textReplace(data string) string {
	for k, v := range textReplaceMap {
		data = strings.Replace(data, k, v, -1)
	}
	return data
}
