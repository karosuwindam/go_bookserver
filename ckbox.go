package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Exists(name string) bool {
	_, err := os.Stat(ServersetUp.Publicpath + name)
	return !os.IsNotExist(err)
}
func zipExists(name string) bool {
	_, err := os.Stat(ServersetUp.Zippath + name)
	return !os.IsNotExist(err)
}
func zipFilesize(name string) int {
	f, _ := os.Open(ServersetUp.Zippath + name)
	defer f.Close()
	if fi, err := f.Stat(); err == nil {
		return int(fi.Size())
	}
	return 0
}
func cpfile(filename string, flag int) {
	srcName := ServersetUp.Zippath + filename
	dstName := ServersetUp.Publicpath + filename
	if Exists(filename) {
		if flag == 0 {
			if err := os.Remove(dstName); err != nil {
				fmt.Println(err)
			}
		}
	} else {
		if flag == 1 {
			src, err := os.Open(srcName)
			if err != nil {
				panic(err)
			}
			defer src.Close()

			dst, err := os.Create(dstName)
			if err != nil {
				panic(err)
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				panic(err)
			}
		}
	}
}
func ckbox(w http.ResponseWriter, r *http.Request) {
	var output string
	tmp := map[string]string{}
	flag := 0
	filesize := 0
	typedata_tmp := []string{"zippass", "copyflag"}
	if r.Method == "POST" {
		// addnew_flag = false
		err := r.ParseForm()
		if err != nil {
			// エラー処理
			output = "POST err"
		} else {
			for _, keyword := range typedata_tmp {
				tmp[keyword] = r.FormValue(keyword)
				fmt.Println(keyword+":", tmp[keyword])
			}
			if tmp["copyflag"] == "true" {
				flag = 1
			}
			copyfile_t.ReadName(tmp["zippass"])
			if copyfile_t.Tmp.Id != 0 {
				copyfile_t.Update(strconv.Itoa(copyfile_t.Tmp.Id), copyfile_t.Tmp.Zippass, flag, copyfile_t.Tmp.Filesize)
				Logout.Out(1, "update:%v\n", copyfile_t.JsonOutTmp())
			} else {
				//ファイル名よりファイルサイズ取得
				if zipExists(tmp["zippass"]) {
					filesize = zipFilesize(tmp["zippass"])
					//データ更新
					copyfile_t.Add(tmp["zippass"], flag, filesize)
					Logout.Out(1, "new:%v\n", copyfile_t.JsonOutTmp())
				} else {
					return
				}
			}
			//ファイルコピー処理
			cpfile(tmp["zippass"], flag)
		}
	} else {
		data := r.URL.RawQuery
		// if strings.Index(data, "&") > 0 {
		for _, data_tmp := range strings.Split(data, "&") {
			if strings.Index(data_tmp, "=") > 0 {
				data_tmp2 := strings.Split(data_tmp, "=")
				data_tmp2[1], _ = url.QueryUnescape(data_tmp2[1])
				tmp[data_tmp2[0]] = data_tmp2[1]
			}
		}
		// }
		if data == "" {
			output = "err get"
		} else {
			copyfile_t.ReadName(tmp["zippass"])
			output = strconv.Itoa(copyfile_t.Tmp.Copyflag)
		}
		fmt.Fprintf(w, "%s", output)
	}
}
