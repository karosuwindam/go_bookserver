package main

import (
	// ビルド時のみ使用する
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"./bookname"
	_ "github.com/mattn/go-sqlite3"
)

// DB Path(相対パスでも大丈夫かと思うが、筆者の場合、絶対パスでないと実行できなかった)
// const dbPath = "db.sql"
const dbPath = "development.sqlite3"

// コネクションプールを作成
var DbConnection *sql.DB

// データ格納用
type Blog struct {
	id    int
	title string
}

func intosql(data Blog) {
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()
	cmd := "INSERT INTO blog (id, title) VALUES (" + strconv.Itoa(data.id) + ",'" + data.title + "')"
	fmt.Println(cmd)
	// fmt.Fprintf(cmd, "INSERT INTO blog (id, title) VALUES (%v, %v)", data.id, data.title)
	_, err := DbConnection.Exec(cmd)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		log.Fatalln(err)
	}
}

func test1() {
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// blog テーブルの作成
	cmd := `CREATE TABLE IF NOT EXISTS blog(
             id INT,    
             title STRING)`

	// cmd を実行
	// _ -> 受け取った結果に対して何もしないので、_ にする
	_, err := DbConnection.Exec(cmd)

	// エラーハンドリング(Go だと大体このパターン)
	if err != nil {
		// Fatalln は便利
		// エラーが発生した場合、以降の処理を実行しない
		log.Fatalln(err)
	}
	var tmp Blog
	tmp.id = 1
	tmp.title = "bb"
	intosql(tmp)
}

type booknames struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Title      string    `json:"title"`
	Writer     string    `json:"writer`
	Brand      string    `json:"brand"`
	Booktype   string    `json:"booktype"`
	Ext        string    `json:"ext"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func test2() {
	// var tmp booknames
	DbConnection, _ := sql.Open("sqlite3", dbPath)
	defer DbConnection.Close()
	// cmd := "SELECT id,name,title,writer,created_at FROM booknames"
	cmd := "SELECT * FROM booknames"
	fmt.Println(cmd)
	// fmt.Fprintf(cmd, "INSERT INTO blog (id, title) VALUES (%v, %v)", data.id, data.title)
	// rows, _ := DbConnection.Exec(cmd)
	// rows, err := DbConnection.Prepare(cmd)
	rows, err := DbConnection.Query(cmd)

	// content, err := ioutil.ReadFile(rows)
	if err != nil {
		// return "", err
	}
	defer rows.Close()
	data := booknames{}
	var output []booknames
	// rows.QueryRow(0).Scan(&data.id, &data.name)
	for rows.Next() {
		rows.Scan(&data.Id, &data.Title, &data.Name, &data.Writer, &data.Brand, &data.Booktype, &data.Ext, &data.Created_at, &data.Updated_at)
		// rows.Scan(data)
		output = append(output, data)
	}
	// b := bytes.NewBuffer(content)
	// // if err != nil{

	// // }
	// fmt.Println(b.String())
	// rows.Scan(tmp.id)
	// defer tmp.Close()

	// fmt.Println(rows)
	fmt.Println(output)
}

func main() {

	webserversetup("output.log")
	webserverstart()
	return
	// test1()
	// test2()
	var tmp bookname.Data
	tmp.New("test1.db")
	tmp.CreatDb()
	tmp.Add("1", "2", "3", "4", "5", "6")
	tmp.Add("11", "2", "3", "4", "5", "6")
	// tmp.Delete("6")
	// tmp.Update("1", "", "1", "1", "1", "1", "1")
	// tmp.ReadAll()
	tmp.ReadAll()
	tmp.ReadId("2")
	// fmt.Println(tmp.List)
	// fmt.Println(tmp.Tmp)
	fmt.Println(tmp.JsonOutList())
	fmt.Println(tmp.JsonOutTmp())
}
