package main

import (
	"fmt"
	"net/http"
	"os"
)

// const zippass = "upload/zip/"
// const pdfpass = "upload/pdf/"

func download(w http.ResponseWriter, r *http.Request) {
	// downloaddata := []string{"test.zip", "test.txt.tar.gz"}
	data := map[string]string{}
	urldata := urlAnalysis(r.URL.Path)
	var filepass string
	fmt.Println(urldata)
	data["type"] = urldata[1]
	data["id"] = urldata[2]

	filelist_t.ReadId(data["id"])
	if filelist_t.Tmp.Id == 0 {
		Logout.Out(1, "download err URL:%v\n", r.URL.Path)
		fmt.Fprintf(w, "%s", "download err")
		return
	}

	//ダウンロードパス
	if data["type"] == "zip" {
		filepass = ServersetUp.Zippath + filelist_t.Tmp.Zippass
	} else if data["type"] == "pdf" {
		filepass = ServersetUp.Pdfpath + filelist_t.Tmp.Pdfpass
	} else {
		Logout.Out(1, "download err URL:%v\n", r.URL.Path)
		fmt.Fprintf(w, "%s", "download err")
		return
	}
	file, err1 := os.Open(filepass)

	if err1 != nil {
		fmt.Println(err1)
	}

	defer file.Close()
	buf := make([]byte, 1024)
	var buffer []byte
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			// Readエラー処理
			break
		}
		buffer = append(buffer, buf[:n]...)
	}
	Logout.Out(1, "download URL:%v,Name:%v\n", r.URL.Path, filelist_t.Tmp.Name)
	if data["type"] == "zip" {
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+filelist_t.Tmp.Zippass)
		// コンテントタイプ
		w.Header().Set("Content-Type", "application/zip")
	} else if data["type"] == "pdf" {
		// ファイル名
		w.Header().Set("Content-Disposition", "attachment; filename="+filelist_t.Tmp.Pdfpass)
		// コンテントタイプ
		w.Header().Set("Content-Type", "application/pdf")
	}
	// ファイルの長さ
	w.Header().Set("Content-Length", string(len(buffer)))
	// bodyに書き込み
	w.Write(buffer)
}
