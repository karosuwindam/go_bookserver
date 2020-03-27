package main

import (
	// ビルド時のみ使用する
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
}
type Setupdate struct {
	Serverdata serverdata `json:"serverdata"`
	Zippath    string     `json:"zippath"`
	Pdfpath    string     `json:"pdfpath"`
	Uploadpath string     `json:"uploadpath"`
}

// // コネクションプールを作成
var DbConnection *sql.DB

//サーバの各設定
var ServersetUp Setupdate

const config_json = "config/setup.json"

func main() {
	raw, err := ioutil.ReadFile(config_json)
	if err != nil {
		ServersetUp.Serverdata.Serverip = ""
		ServersetUp.Serverdata.Serverport = "8080"
		ServersetUp.Serverdata.Booknamedb = "test1.db"
		ServersetUp.Serverdata.Filelistdb = "test1.db"
		ServersetUp.Uploadpath = "upload/"
		ServersetUp.Zippath = "upload/zip/"
		ServersetUp.Pdfpath = "upload/pdf/"
		bytes, _ := json.Marshal(ServersetUp)

		fp, err := os.Create(config_json)
		if err != nil {
			panic(err)
		}
		fp.Write(bytes)
		fp.Close()
		fmt.Println(config_json + " creat OK")
		// fmt.Println(string(bytes))
	} else {
		var fc Setupdate

		json.Unmarshal(raw, &fc)
		// fmt.Println(string(raw))
		fmt.Println(config_json + " read OK")
		ServersetUp = fc
	}

	webserversetup("output.log")
	webserverstart()
	return
}
