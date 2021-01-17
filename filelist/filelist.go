package filelist

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

type filelists struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Pdfpass    string    `json:"pdfpass"`
	Zippass    string    `json:"zippass`
	Tag        string    `json:"tag"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
type Data struct {
	List   []filelists
	Tmp    filelists
	dbpath string
	Renew  bool
}

const database = "filelists"
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
		name varchar,
		pdfpass varchar,
		zippass varchar,
		tag varchar,
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
		rows.Scan(&data.Id, &data.Name, &data.Pdfpass, &data.Zippass, &data.Tag, &data.Created_at, &data.Updated_at)
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
		rows.Scan(&data.Id, &data.Name, &data.Pdfpass, &data.Zippass, &data.Tag, &data.Created_at, &data.Updated_at)
		tmp = append(tmp, data)
	}
	sort.Slice(tmp, func(i, j int) bool { return t.strNameSelect(tmp[i].Name, tmp[j].Name) })
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
		rows.Scan(&data.Id, &data.Name, &data.Pdfpass, &data.Zippass, &data.Tag, &data.Created_at, &data.Updated_at)
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
	rows.Scan(&data.Id, &data.Name, &data.Pdfpass, &data.Zippass, &data.Tag, &data.Created_at, &data.Updated_at)
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
	data := filelists{}
	rows.Next()
	rows.Scan(&data.Id, &data.Name, &data.Pdfpass, &data.Zippass, &data.Tag, &data.Created_at, &data.Updated_at)
	t.Tmp = data

	return E_OK
}
func (t *Data) Add(name, pdfpass, zippass, tag string) int {
	var id int
	id = t.readid()
	time := time.Now()
	DbConnection, _ := sql.Open("sqlite3", t.dbpath)
	defer DbConnection.Close()
	cmd := "INSERT INTO " + database + " (id, name, pdfpass, zippass, tag,created_at,updated_at) VALUES ("
	cmd += strconv.Itoa(id) + ","
	cmd += "'" + name + "'" + ","
	cmd += "'" + pdfpass + "',"
	cmd += "'" + zippass + "',"
	cmd += "'" + tag + "',"
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
	t.Tmp.Pdfpass = pdfpass
	t.Tmp.Zippass = zippass
	t.Tmp.Tag = tag
	t.Tmp.Created_at = time
	t.Tmp.Updated_at = time
	return 0
}
func (t *Data) Update(id, name, pdfpass, zippass, tag string) int {
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
	// if pdfpass != "" {
	tmp += ", pdfpass=" + "'" + pdfpass + "'"
	// }
	// if zippass != "" {
	tmp += ", zippass=" + "'" + zippass + "'"
	// }
	// if tag != "" {
	tmp += ", tag=" + "'" + tag + "'"
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

func (t *Data) strNameSelect(A, B string) bool {
	var num_s, num_e int
	var str_s, str_e string

	var err error
	str_s = strings.ToLower(A)
	for z := 0; z < len(A); z++ {
		tmp := A[z:]
		num_s, err = strconv.Atoi(tmp)
		if err == nil {
			str_s = strings.ToLower(A[:z])
			break
		}
		num_s = 0
	}
	str_e = strings.ToLower(B)
	for z := 0; z < len(B); z++ {
		tmp := B[z:]
		num_e, err = strconv.Atoi(tmp)
		if err == nil {
			str_e = strings.ToLower(B[:z])
			break
		}
		num_e = 0
	}
	if str_s == str_e {
		return num_s < num_e

	}
	return A < B
}
