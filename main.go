package main

import (
	// ビルド時のみ使用する
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

// const serverip = ""
// const serverport = "8080"
// const databass = "test1.db"
// const databass = "development.sqlite3"

type serverdata struct {
	Serverip   string `json:"serverip"`
	Serverport string `json:"serverport"`
	Booknamedb string `json:"booknamedb"`
	Filelistdb string `json:"filelistdb"`
	Copyfiledb string `json:"copyfiledb"`
	TmpPass    string `json:"tmppass"`
}
type Setupdate struct {
	Serverdata serverdata `json:"serverdata"`
	Zippath    string     `json:"zippath"`
	Pdfpath    string     `json:"pdfpath"`
	Publicpath string     `json:"pulicpath`
	Uploadpath string     `json:"uploadpath"`
}

// // コネクションプールを作成
var DbConnection *sql.DB

//サーバの各設定
var ServersetUp Setupdate

const config_json = "config/setup.json"

func main() {
	var buf bytes.Buffer
	raw, err := ioutil.ReadFile(config_json)
	if err != nil {
		ServersetUp.Serverdata.Serverip = ""
		ServersetUp.Serverdata.Serverport = "8080"
		ServersetUp.Serverdata.Booknamedb = "test1.db"
		ServersetUp.Serverdata.Filelistdb = "test1.db"
		ServersetUp.Serverdata.Copyfiledb = "test1.db"
		ServersetUp.Serverdata.TmpPass = "tmp"
		ServersetUp.Uploadpath = "upload/"
		ServersetUp.Zippath = "upload/zip/"
		ServersetUp.Pdfpath = "upload/pdf/"
		bytes, _ := json.Marshal(ServersetUp)

		fp, err := os.Create(config_json)
		if err != nil {
			panic(err)
		}
		json.Indent(&buf, bytes, "", "  ")
		fp.Write(buf.Bytes())
		fp.Close()
		fmt.Println(config_json + " creat OK")
		// fmt.Println(string(bytes))
	} else {
		var fc Setupdate

		json.Unmarshal(raw, &fc)
		// fmt.Println(string(raw))
		fmt.Println(config_json + " read OK")
		ServersetUp = fc
		// bytes, _ := json.Marshal(ServersetUp)
		// json.Indent(&buf, bytes, "", "  ")
		// fmt.Println(buf.String())
	}
	// ServersetUp.Publicpath = "Public/"
	Tmp, _ := exec.Command("which", "pdfimages").Output()
	if len(Tmp) == 0 {
		fmt.Println("err not install pdfimages", "run sudo apt install poppler-utils")
		return
	}
	if f, err := os.Stat(ServersetUp.Serverdata.TmpPass); os.IsNotExist(err) || !f.IsDir() {
		fmt.Printf("%vディレクトリは存在しません！\n", ServersetUp.Serverdata.TmpPass)
		return
	}
	if f, err := os.Stat(ServersetUp.Uploadpath); os.IsNotExist(err) || !f.IsDir() {
		fmt.Printf("%vディレクトリは存在しません！\n", ServersetUp.Uploadpath)
		return
	}
	if f, err := os.Stat(ServersetUp.Publicpath); os.IsNotExist(err) || !f.IsDir() {
		fmt.Printf("%vディレクトリは存在しません！\n", ServersetUp.Publicpath)
		return
	}
	if ServersetUp.Serverdata.TmpPass[len(ServersetUp.Serverdata.TmpPass)-1:] == "/" {
		ServersetUp.Serverdata.TmpPass = ServersetUp.Serverdata.TmpPass[:len(ServersetUp.Serverdata.TmpPass)-1]
	}
	if ServersetUp.Zippath[len(ServersetUp.Zippath)-1:] != "/" {
		ServersetUp.Zippath += "/"
	}
	if ServersetUp.Pdfpath[len(ServersetUp.Pdfpath)-1:] != "/" {
		ServersetUp.Pdfpath += "/"
	}
	if ServersetUp.Uploadpath[len(ServersetUp.Uploadpath)-1:] != "/" {
		ServersetUp.Uploadpath += "/"
	}
	if ServersetUp.Publicpath[len(ServersetUp.Publicpath)-1:] != "/" {
		ServersetUp.Publicpath += "/"
	}
	webserversetup("output.log")
	// tmp := FolderDataSetup()
	// data, _ := json.Marshal(tmp.CheckData())
	// fmt.Println(string(data))
	// return
	webserverstart()
	return
}
