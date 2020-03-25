package main

import (
	"log"
	"os"
	"strings"
)

//ReadHtmlはpathに入力したファイルパスから読み取る
//pathはテキストパスでそのテキスト値をもとに戻す

func ReadHtml(path string) string {
	var output string
	fp, err := os.Open(path)
	if err != nil {
		log.Panic(err)
		return ""
	}
	defer fp.Close()
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		output += string(buf[:n])
	}
	return output
}

//ConvertDataはstrに入力されたデータから<%%>に囲まれた文字列から
//data[文字列]の入力された値と置き換えて変換値を戻り値にする

func ConvertData(str string, data map[string]string) string {
	//func ConvertData(str string) string {
	tmp := str
	output := str
	for {
		n := strings.Index(tmp, "<%")
		if n >= 0 {
			m := strings.Index(tmp, "%>")
			if m >= 0 {
				dtmp := tmp[n+2 : m]
				output = strings.Replace(output, "<%"+dtmp+"%>", data[dtmp], 1)
				//output += tmp[n+2:m] + ","
				tmp = tmp[m+2:]
			} else {
				break
			}
		} else {
			break
		}
	}
	//output := str
	return output
}
