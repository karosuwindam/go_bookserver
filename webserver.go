package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"./bookname"
	"./filelist"
	"./logoutput"
)

// const serverip = ""
// const serverport = "8080"

// const databass = "test1.db"

// const databass = "development.sqlite3"

var websetup_flag bool = false
var Logdata logoutput.Data

// Logdata.Setup("output.log")
var bookname_t bookname.Data
var filelist_t filelist.Data

var Logout logoutput.Data

func webserversetup(logname string) {
	Logout.Setup("upload.log")
	Logdata.Setup(logname)
	// bookname_t.New("development.sqlite3")
	bookname_t.New(ServersetUp.Serverdata.Booknamedb)
	filelist_t.New(ServersetUp.Serverdata.Filelistdb)
	bookname_t.CreatDb()
	filelist_t.CreatDb()
	Logdata.Out(0, "Bookname db Read %v\n", ServersetUp.Serverdata.Booknamedb)
	Logdata.Out(0, "filelist db Read %v\n", ServersetUp.Serverdata.Filelistdb)
	websetup_flag = true
}

func getserch(w http.ResponseWriter, r *http.Request) {
	output := ""
	data := map[string]string{}
	keyword_s := []string{"today", "toweek", "tomonth"}

	urldata := urlAnalysis(r.URL.Path)
	if len(urldata) > 2 {
		data["pass"] = urldata[1]
		data["keyword"] = urldata[2]
	} else {
		data["keyword"] = urldata[1]
	}
	for _, str := range keyword_s {
		if data["keyword"] == str {
			filelist_t.ReadTime(str)
			output = filelist_t.JsonOutList()
			fmt.Fprintf(w, "%s", output)
			return
		}
	}
	if data["keyword"] == "rand" {
		filelist_t.ReadRand()
		output = filelist_t.JsonOutList()
		fmt.Fprintf(w, "%s", output)
		return
	}
	if data["keyword"] == "" {
		output = "[]"
	} else {
		tmp := "name like '%" + data["keyword"] + "%'"
		if data["pass"] == "filelist" {
			tmp += " or tag like '%" + data["keyword"] + "%' order by name"
			filelist_t.Read(tmp)
			output = filelist_t.JsonOutList()
		} else {
			tmp += " or title like '%" + data["keyword"] + "%'"
			tmp += " or writer like '%" + data["keyword"] + "%'"
			tmp += " or ext like '%" + data["keyword"] + "%'"
			bookname_t.Read(tmp)
			output = bookname_t.JsonOutList()
		}
	}
	fmt.Fprintf(w, "%s", output)
}

func getlist(w http.ResponseWriter, r *http.Request) {
	output := ""
	urldata := urlAnalysis(r.URL.Path)
	if urldata[1] == "filelist" {
		filelist_t.ReadAll()
		output = filelist_t.JsonOutList()
	} else {
		bookname_t.ReadAll()
		output = bookname_t.JsonOutList()
	}

	fmt.Fprintf(w, "%s", output)
}

func urlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}
func editData(w http.ResponseWriter, r *http.Request) {
	var output string
	var tmp map[string]string
	var ret int
	typedata_tmp := []string{"name", "title", "writer", "brand", "booktype", "ext"}
	urldata := urlAnalysis(r.URL.Path)
	fmt.Println(urldata)
	url_data := map[string]string{}
	_, err := strconv.Atoi(urldata[1])
	if err != nil {
		url_data["url"] = urldata[1]
		url_data["id"] = urldata[2]
	} else {
		url_data["id"] = urldata[1]
	}
	if url_data["url"] == "filelist" {
		typedata_tmp = []string{"name", "pdfpass", "zippass", "tag"}
		tmp, ret = readFilelistId(url_data["id"])
	} else {
		tmp, ret = readBooknameId(url_data["id"])
	}
	if ret < 0 {
		output = "err input id"
		fmt.Fprintf(w, "%s", output)
		return
	}

	tmp["url"] = "/edit/"
	if url_data["url"] != "" {
		tmp["url"] += url_data["url"] + "/"
	}
	tmp["url"] += url_data["id"]
	tmp["edit"] = "<div>" + "id:" + url_data["id"] + "</div>"

	if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			// エラー処理
			output = "POST err"
		} else {
			for _, keyword := range typedata_tmp {
				tmp[keyword] = r.FormValue(keyword)
				// fmt.Println(keyword+":", r.FormValue(keyword))
			}
			if url_data["url"] == "filelist" {
				filelist_t.Update(url_data["id"], tmp["name"], tmp["pdfpass"], tmp["zippass"], tmp["tag"])
			} else {
				bookname_t.Update(url_data["id"], tmp["name"], tmp["title"], tmp["writer"], tmp["brand"], tmp["booktype"], tmp["ext"])
			}
			output = "POST"
			r.Method = "GET"
			showdata(w, r)
			return
		}
	} else {
		output = "GET"
		if url_data["url"] == "filelist" {
			tmp["inputdata"] = ConvertData(ReadHtml("html_tmp/filelist/new.html"), tmp)
			output = ConvertData(ReadHtml("html_tmp/new_.html"), tmp)
		} else {
			output = ConvertData(ReadHtml("html_tmp/new.html"), tmp)
		}
	}
	fmt.Fprintf(w, "%s", output)
}
func destory(w http.ResponseWriter, r *http.Request) {
	output := ""
	urldata := urlAnalysis(r.URL.Path)
	tmp := map[string]string{}
	_, err := strconv.Atoi(urldata[1])
	if err != nil {
		tmp["url"] = urldata[1]
		tmp["id"] = urldata[2]
	} else {
		tmp["id"] = urldata[1]
	}
	if r.Method == "POST" {
		if tmp["url"] == "filelist" {
			filelist_t.ReadId(tmp["id"])
			filelist_t.Delete(tmp["id"])
			output = filelist_t.JsonOutTmp()
		} else {
			bookname_t.ReadId(tmp["id"])
			bookname_t.Delete(tmp["id"])
			output = bookname_t.JsonOutTmp()
		}
	} else {
		output = "GET"
	}
	fmt.Fprintf(w, "%s", output)
}
func readBooknameId(str string) (map[string]string, int) {
	tmp := map[string]string{}
	if bookname_t.ReadId(str) < 0 {
		return tmp, -1
	}
	tmp["id"] = strconv.Itoa(bookname_t.Tmp.Id)
	tmp["name"] = bookname_t.Tmp.Name
	tmp["title"] = bookname_t.Tmp.Title
	tmp["title"] = bookname_t.Tmp.Title
	tmp["writer"] = bookname_t.Tmp.Writer
	tmp["brand"] = bookname_t.Tmp.Brand
	tmp["booktype"] = bookname_t.Tmp.Booktype
	tmp["ext"] = bookname_t.Tmp.Ext
	return tmp, 0
}
func readFilelistId(str string) (map[string]string, int) {
	tmp := map[string]string{}
	if filelist_t.ReadId(str) < 0 {
		return tmp, -1
	}
	tmp["id"] = strconv.Itoa(filelist_t.Tmp.Id)
	tmp["name"] = filelist_t.Tmp.Name
	tmp["pdfpass"] = filelist_t.Tmp.Pdfpass
	tmp["zippass"] = filelist_t.Tmp.Zippass
	tmp["tag"] = filelist_t.Tmp.Tag
	return tmp, 0
}
func showdata(w http.ResponseWriter, r *http.Request) {
	var output string
	data := map[string]string{}
	tmp := map[string]string{}
	err := 0
	urldata := urlAnalysis(r.URL.Path)

	if len(urldata) > 2 {
		data["pass"] = urldata[1]
		data["id"] = urldata[2]
	} else {
		data["id"] = urldata[1]
	}
	if data["pass"] == "filelist" {
		tmp, err = readFilelistId(data["id"])
	} else {
		tmp, err = readBooknameId(data["id"])
	}
	if err < 0 {
		output = "err input id"
		fmt.Fprintf(w, "%s", output)
		return
	}
	if r.Method == "POST" {
		output = "input err"
	} else {
		output = "GET"
		// output = ConvertData(ReadHtml("html_tmp/new.html"), tmp)
		if data["pass"] == "filelist" {
			tmp["show"] = ConvertData(ReadHtml("html_tmp/filelist/show.html"), tmp)
		} else {
			tmp["show"] = ConvertData(ReadHtml("html_tmp/bookname/show.html"), tmp)
		}
		output = ConvertData(ReadHtml("html_tmp/show_.html"), tmp)
	}
	fmt.Fprintf(w, "%s", output)

}
func addNew(w http.ResponseWriter, r *http.Request) {
	var output string
	tmp := map[string]string{}
	urldata := urlAnalysis(r.URL.Path)
	tmp["url"] = "/" + urldata[0]
	if urldata[1] != "" {
		tmp["url"] += "/" + urldata[1]
	}
	tmp["url"] += "/"

	var typedata_tmp []string
	if urldata[1] == "filelist" {
		typedata_tmp = []string{"name", "pdfpass", "zippass", "tag"}
	} else {
		typedata_tmp = []string{"name", "title", "writer", "brand", "booktype", "ext"}
	}
	if r.Method == "POST" {
		addnew_flag := false
		err := r.ParseForm()
		if err != nil {
			// エラー処理
			output = "POST err"
		} else {
			for _, keyword := range typedata_tmp {
				if (r.FormValue(keyword) != "") && (keyword == "name") {
					addnew_flag = true
				}
				tmp[keyword] = r.FormValue(keyword)
				// fmt.Println(keyword+":", tmp[keyword])
			}
			if addnew_flag {
				if urldata[1] == "filelist" {
					filelist_t.Add(tmp["name"], tmp["pdfpass"], tmp["zippass"], tmp["tag"])
					output = filelist_t.JsonOutTmp()
				} else {
					bookname_t.Add(tmp["name"], tmp["title"], tmp["writer"], tmp["brand"], tmp["booktype"], tmp["ext"])
					output = bookname_t.JsonOutTmp()
				}
				// output = "POST"
			} else {
				output = "input err"
			}
		}
	} else {
		output = "GET"
		// output = ConvertData(ReadHtml("html_tmp/new.html"), tmp)
		if urldata[1] == "filelist" {
			tmp["inputdata"] = ConvertData(ReadHtml("html_tmp/filelist/new.html"), tmp)
		} else {
			tmp["inputdata"] = ConvertData(ReadHtml("html_tmp/bookname/new.html"), tmp)
		}
		output = ConvertData(ReadHtml("html_tmp/new_.html"), tmp)
	}
	fmt.Fprintf(w, "%s", output)
}
func machdata(str string) (string, string) {
	output := "{}"
	data := map[string]string{}
	kan := ""

	data["keyword"] = str
	if data["keyword"] == "" {
		return "{}", kan
	}
	bookname_t.ReadName(data["keyword"])
	if bookname_t.Tmp.Id == 0 {
		for i := 1; i < 3; i++ {
			bookname_t.ReadName(data["keyword"][0 : len(data["keyword"])-i])
			if bookname_t.Tmp.Id != 0 {
				output = bookname_t.JsonOutTmp()
				kan = data["keyword"][len(data["keyword"])-i : len(data["keyword"])]
				break
			} else {

				// return "{}", kan
			}
		}
	} else {
		output = bookname_t.JsonOutTmp()
	}
	num, _ := strconv.Atoi(kan)
	if num < 10 {
		kan = "0" + kan
	}
	return output, kan
}
func machdata_filelist(str string) string {
	output := ""
	filelist_t.ReadName(str)
	if filelist_t.Tmp.Id == 0 {
		output = "{\"flag\":0}"
	} else {
		output = "{\"flag\":1}"
	}
	return output
}
func machdatahttp(w http.ResponseWriter, r *http.Request) {
	output := ""
	urldata := urlAnalysis(r.URL.Path)
	str := urldata[1]
	output, _ = machdata(str)
	if output != "{}" {
		output += "," + machdata_filelist(str)
		output = "[" + output + "]"
	}
	fmt.Fprintf(w, "%s", output)
}
func webserverstart() {
	if !websetup_flag {
		fmt.Println("web server not init")
		return
	}
	Logdata.Out(1, "web server start %v:%v\n", ServersetUp.Serverdata.Serverip, ServersetUp.Serverdata.Serverport)
	http.HandleFunc("/list/", getlist)
	http.HandleFunc("/serch/", getserch)
	http.HandleFunc("/mach/", machdatahttp)
	http.HandleFunc("/new/", addNew)
	http.HandleFunc("/edit/", editData)
	http.HandleFunc("/destory/", destory)
	http.HandleFunc("/show/", showdata)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/zip", zipdata)
	http.HandleFunc("/view/", view)
	http.HandleFunc("/download/", download)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./html"))))
	http.ListenAndServe(ServersetUp.Serverdata.Serverip+":"+ServersetUp.Serverdata.Serverport, nil)

}
