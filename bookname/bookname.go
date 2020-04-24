package bookname

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

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
type Data struct {
	List   []booknames
	Tmp    booknames
	dbpath string
	Renew  bool
}

const database = "booknames"

const (
	E_OK  = 0
	E_ERR = -1
)

func (t *Data) New(s string) {
	var tmp []booknames
	if (len(t.List) == 0) || (t.Renew) {
		t.dbpath = s
		t.List = tmp
		t.Renew = false
	}
}
func (t *Data) CreatDb() int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := `CREATE TABLE IF NOT EXISTS ` + database + `(
`
	cmd += `id integer PRIMARY KEY,    
		name varchar,
		title varchar,
		writer varchar,
		burand varchar,
		booktype varchar,
		ext varchar,
		created_at datatime,
		updated_at datatime
		)`
	_, err := DbConnection.Exec(cmd)

	if err != nil {
		// Fatalln は便利
		// エラーが発生した場合、以降の処理を実行しない
		log.Fatalln(err)
	}
	return 0
}
func (t *Data) ReadAll() int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	rows, err := DbConnection.Query(cmd)
	tmp := []booknames{}
	if err != nil {
		return -1
	}
	defer rows.Close()
	data := booknames{}
	for rows.Next() {
		rows.Scan(&data.Id, &data.Name, &data.Title, &data.Writer, &data.Brand, &data.Booktype, &data.Ext, &data.Created_at, &data.Updated_at)
		tmp = append(tmp, data)
	}
	t.List = tmp
	return 0
}

func (t *Data) Read(s string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	tmp := []booknames{}

	if s != "" {
		cmd += " " + "where " + s
	} else {
	}
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return E_ERR
	}
	defer rows.Close()
	data := booknames{}
	for rows.Next() {
		rows.Scan(&data.Id, &data.Name, &data.Title, &data.Writer, &data.Brand, &data.Booktype, &data.Ext, &data.Created_at, &data.Updated_at)
		tmp = append(tmp, data)
	}
	t.List = tmp
	return E_OK
}
func (t *Data) ReadId(id string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	cmd += " where id=" + id
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return E_ERR
	}
	defer rows.Close()
	data := booknames{}
	rows.Next()
	rows.Scan(&data.Id, &data.Name, &data.Title, &data.Writer, &data.Brand, &data.Booktype, &data.Ext, &data.Created_at, &data.Updated_at)
	t.Tmp = data

	return E_OK
}
func (t *Data) ReadName(name string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	cmd += " where name='" + name
	cmd += "'"
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return E_ERR
	}
	defer rows.Close()
	data := booknames{}
	rows.Next()
	rows.Scan(&data.Id, &data.Name, &data.Title, &data.Writer, &data.Brand, &data.Booktype, &data.Ext, &data.Created_at, &data.Updated_at)
	t.Tmp = data

	return E_OK
}

func (t *Data) Add(name, title, writer, brand, booktype, ext string) int {
	var id int
	id = t.readid()
	time := time.Now()
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	//burand,brand
	cmd := "INSERT INTO " + database + " (id, title,name,writer,burand,booktype,ext,created_at,updated_at) VALUES ("
	cmd += strconv.Itoa(id) + ","
	cmd += "'" + title + "'" + ","
	cmd += "'" + name + "'" + ","
	cmd += "'" + writer + "',"
	cmd += "'" + brand + "',"
	cmd += "'" + booktype + "',"
	cmd += "'" + ext + "',"
	cmd += "'" + time.Format("2006-01-02 15:04:05.999999999") + "',"
	cmd += "'" + time.Format("2006-01-02 15:04:05.999999999") + "'"

	cmd += ");"

	// fmt.Println(cmd)
	_, err := DbConnection.Exec(cmd)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		log.Fatalln(err)
	}
	t.Tmp.Id = id
	t.Tmp.Name = name
	t.Tmp.Title = title
	t.Tmp.Writer = writer
	t.Tmp.Brand = brand
	t.Tmp.Booktype = booktype
	t.Tmp.Ext = ext
	t.Tmp.Created_at = time
	t.Tmp.Updated_at = time
	return 0
}
func (t *Data) Update(id, name, title, writer, brand, booktype, ext string) int {
	var tmp string
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()

	if id == "" {
		return E_ERR
	}
	tmp = "id=" + id
	// if name != "" {
	tmp += ", name=" + "'" + name + "'"
	// }
	// if title != "" {
	tmp += ", title=" + "'" + title + "'"
	// }
	// if writer != "" {
	tmp += ", writer=" + "'" + writer + "'"
	// }
	// if brand != "" {
	tmp += ", burand=" + "'" + brand + "'"
	// }
	// if booktype != "" {
	tmp += ", booktype=" + "'" + booktype + "'"
	// }
	// if ext != "" {
	tmp += ", ext=" + "'" + ext + "'"
	// }
	tmp += ", updated_at=" + "'" + time.Now().Format("2006-01-02 15:04:05.999999999") + "'"

	cmd := "update " + database + " set " + tmp + " where id=" + id
	_, err := DbConnection.Exec(cmd)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		return E_ERR
		// log.Fatalln(err)
	}
	return E_OK
}
func (t *Data) Delete(id string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "delete from " + database + " where id="
	cmd += id
	_, err := DbConnection.Exec(cmd)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		return E_ERR
		// log.Fatalln(err)
	}
	return E_OK

}
func (t *Data) readid() int {
	id := 1
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "select max(id) from " + database
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return id
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {

	} else {
		id++
	}
	return id
}

func (t *Data) JsonOutList() string {
	bytes, _ := json.Marshal(t.List)
	return string(bytes)
}

func (t *Data) JsonOutTmp() string {
	bytes, _ := json.Marshal(t.Tmp)
	return string(bytes)
}
