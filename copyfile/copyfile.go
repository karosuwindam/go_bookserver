package copyfile

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type filelists struct {
	Id         int       `json:"id"`
	Zippass    string    `json:"zippass`
	Filesize   int       `json:"filesize`
	Copyflag   int       `json:"copyflag"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
type Data struct {
	List   []filelists
	Tmp    filelists
	dbpath string
	Renew  bool
}

const database = "copyfile"
const (
	E_OK  = 0
	E_ERR = -1
)

func (t *Data) New(s string) {
	var tmp []filelists
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
		zippass varchar,
		filesize integer,
		copyflag integer,
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
	if err != nil {
		return -1
	}
	defer rows.Close()
	tmp := []filelists{}
	data := filelists{}
	for rows.Next() {
		rows.Scan(&data.Id, &data.Zippass, &data.Filesize, &data.Copyflag, &data.Created_at, &data.Updated_at)
		tmp = append(tmp, data)
	}
	t.List = tmp
	return 0
}

func (t *Data) Read(s string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	tmp := []filelists{}

	if s != "" {
		cmd += " " + "where " + s
	} else {
	}
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return E_ERR
	}
	defer rows.Close()
	data := filelists{}
	for rows.Next() {
		rows.Scan(&data.Id, &data.Zippass, &data.Filesize, &data.Copyflag, &data.Created_at, &data.Updated_at)
		tmp = append(tmp, data)
	}
	t.List = tmp
	return E_OK
}
func (t *Data) ReadRand() int {
	tmp := []filelists{}
	t.ReadAll()
	max := len(t.List) - 1
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		t.ReadId(strconv.Itoa(rand.Intn(max)))
		if t.Tmp.Id != 0 {
			break
		}
	}
	// t.ReadId("1")
	tmp = append(tmp, t.Tmp)
	for i := 0; i < 3; i++ {
		t.ReadId(strconv.Itoa(rand.Intn(max)))
		if t.Tmp.Id != 0 {
			break
		}
	}
	// t.ReadId("2")
	tmp = append(tmp, t.Tmp)
	for i := 0; i < 3; i++ {
		t.ReadId(strconv.Itoa(rand.Intn(max)))
		if t.Tmp.Id != 0 {
			break
		}
	}
	// t.ReadId("3")
	tmp = append(tmp, t.Tmp)
	t.List = tmp
	return E_OK
}
func (t *Data) ReadTime(s string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	tmp := []filelists{}
	nowtime := time.Now()

	switch s {
	case "today":
		cmd += " " + "where " + "updated_at "
		cmd += "between '" + nowtime.Format("2006-01-02") + "' and '"
		cmd += nowtime.Add(24*time.Hour).Format("2006-01-02") + "'"
	case "toweek":
		cmd += " " + "where " + "updated_at "
		cmd += "between '" + nowtime.Add(-24*time.Hour*7).Format("2006-01-02") + "' and '"
		cmd += nowtime.Add(24*time.Hour).Format("2006-01-02") + "'"
	case "tomonth":
		cmd += " " + "where " + "updated_at "
		cmd += "between '" + nowtime.Add(-24*time.Hour*30).Format("2006-01-02") + "' and '"
		cmd += nowtime.Add(24*time.Hour).Format("2006-01-02") + "'"
	default:

	}
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return E_ERR
	}
	defer rows.Close()
	data := filelists{}
	for rows.Next() {
		rows.Scan(&data.Id, &data.Zippass, &data.Filesize, &data.Copyflag, &data.Created_at, &data.Updated_at)
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
	data := filelists{}
	rows.Next()
	rows.Scan(&data.Id, &data.Zippass, &data.Filesize, &data.Copyflag, &data.Created_at, &data.Updated_at)
	t.Tmp = data

	return E_OK
}
func (t *Data) ReadName(name string) int {
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "SELECT * FROM " + database
	cmd += " where Zippass='" + name
	cmd += "'"
	rows, err := DbConnection.Query(cmd)
	if err != nil {
		return E_ERR
	}
	defer rows.Close()
	data := filelists{}
	rows.Next()
	rows.Scan(&data.Id, &data.Zippass, &data.Filesize, &data.Copyflag, &data.Created_at, &data.Updated_at)
	t.Tmp = data

	return E_OK
}
func (t *Data) Add(zippass string, copyflag, filesize int) int {
	var id int
	id = t.readid()
	time := time.Now()
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "INSERT INTO " + database + " (id,  zippass,filesize, copyflag,created_at,updated_at) VALUES ("
	cmd += strconv.Itoa(id) + ","
	cmd += "'" + zippass + "',"
	cmd += strconv.Itoa(filesize) + ","
	cmd += strconv.Itoa(copyflag) + ","
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
	t.Tmp.Zippass = zippass
	t.Tmp.Filesize = filesize
	t.Tmp.Copyflag = copyflag
	t.Tmp.Created_at = time
	t.Tmp.Updated_at = time
	return 0
}
func (t *Data) Update(id, zippass string, copyflag, filesize int) int {
	var tmp string
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()

	if id == "" {
		return E_ERR
	}
	tmp = "id=" + id
	tmp += ", zippass=" + "'" + zippass + "'"
	tmp += ", filesize=" + strconv.Itoa(filesize)
	tmp += ", copyflag=" + strconv.Itoa(copyflag)
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
