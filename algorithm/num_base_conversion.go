/*
	进制转换问题
	进制转10进制
		个+10位*进制+100位*进制的平方 ...
	10进制转其他进制
		除机制，等余数
		余数继续除进制，又得余数，知道余数小于进制数

		余数的反序列表位10机制数

	二进制 转 8，16进制
		3位2进制为8进制，4位2进制位16进制
		一位8进制等于3位二进制

	go strconv.Format支持进制转换
*/

package main

import (
	"bytes"
	"fmt"
	//"math"
	"strconv"
)

var (
	s2n map[byte]int64
	n2s map[int64]byte
)

const (
	data_base  = 26
	start_byte = 'A'
	start_num  = int64(1)
)

func init() {
	s2n = make(map[byte]int64, 40)
	n2s = make(map[int64]byte, 40)

	for i := 0; i < data_base; i++ {
		s2n[start_byte+byte(i)] = start_num + int64(i) // A=1,B=2...Z=26
		n2s[start_num+int64(i)] = start_byte + byte(i) // 1=A,2=B...26=Z
	}

	fmt.Printf("s2n: %+v\n", s2n)
	fmt.Printf("n2s: %+v\n", n2s)
}

func main() {

	num := strToNum("A")
	fmt.Printf("strToNum(A): %d\n", num)
	fmt.Printf("toHex26(%d): %s\n", num, toHex26(num))
	fmt.Printf("numToStr(%d): %s\n", num, numToStr(num))

	num = strToNum("AB")
	fmt.Printf("strToNum(AB): %d\n", num)
	fmt.Printf("toHex26(%d): %s\n", num, toHex26(num))
	fmt.Printf("numToStr(%d): %s\n", num, numToStr(num))

	num = strToNum("ABC")
	fmt.Printf("strToNum(ABC): %d\n", num)
	fmt.Printf("toHex26(%d): %s\n", num, toHex26(num))
	fmt.Printf("numToStr(%d): %s\n", num, numToStr(num))

	num = strToNum("BCDE")
	fmt.Printf("strToNum(BCDE): %d\n", num)
	fmt.Printf("toHex26(%d): %s\n", num, toHex26(num))
	fmt.Printf("numToStr(%d): %s\n", num, numToStr(num))

}

func strToNum(str string) int64 {
	var num int64
	var rate int64 = 1

	n := len(str)
	for i := n - 1; i >= 0; i-- {
		num += s2n[str[i]] * rate
		rate *= data_base
	}

	return num
}

func numToStr(inNum int64) string {
	result := make([]byte, 0, 32)

	for {
		if inNum <= data_base {
			result = append(result, n2s[inNum])
			break
		}

		// 除数余数
		chu := inNum / data_base
		yu := inNum - chu*data_base
		result = append(result, n2s[yu])

		inNum = chu
	}

	// 翻转余数
	n := len(result)
	for i := 0; i < n/2; i++ {
		result[i], result[n-i-1] = result[n-i-1], result[i]
	}

	return string(result)
}

func toHex26(num int64) string {
	return string(
		bytes.Map(func(r rune) rune {
			if r >= 65 {
				return r + 9
			} else {
				return r + 16
			}
		},
			[]byte(strconv.FormatInt(num, 26)),
		),
	)
}
