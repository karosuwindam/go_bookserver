package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"./dirread"
	"./zipopen"
)

// const ZIPPATH = "upload/zip/"

func zipdata(w http.ResponseWriter, r *http.Request) {
	var t zipopen.File
	var t_dir dirread.Dirtype
	data := map[string]string{}
	str := r.URL.RawQuery
	data["id"] = "1"
	data["page"] = "0"
	num := 1
	page := 0

	if strings.Index(str, "&") > 0 {
		for _, tmp := range strings.Split(str, "&") {
			if strings.Index(tmp, "=") > 0 {
				tmp2 := strings.Split(tmp, "=")
				data[tmp2[0]] = tmp2[1]

			}

		}
	} else if strings.Index(str, "=") > 0 {
		tmp2 := strings.Split(str, "=")
		data[tmp2[0]] = tmp2[1]

	}
	filelist_t.ReadId(data["id"])
	num, _ = strconv.Atoi(data["id"])
	page, _ = strconv.Atoi(data["page"])
	t_dir.Setup(ServersetUp.Zippath)
	_ = t_dir.Read("/")
	if ((num - 1) >= len(t_dir.Data)) || (num == 0) {
		num = 1
	}
	filename := t_dir.Data[num-1].RootPath + t_dir.Data[num-1].Name
	if filelist_t.Tmp.Id != 0 {
		filename = ServersetUp.Zippath + filelist_t.Tmp.Zippass
		// fmt.Println(filename + " page:" + data["page"])
	}
	t.ZipOpenSetup(filename)
	t.ZipReadList()
	// page--
	if page >= t.Count {
		page = 0
	}
	fmt.Fprintf(w, "%s", t.ZipRead(page))
}
func view(w http.ResponseWriter, r *http.Request) {
	var t zipopen.File
	var t_dir dirread.Dirtype
	var datap map[string]string
	var filename string
	url := r.URL.Path
	data := map[string]string{}
	datap = data
	data["id"] = "1"
	// id := 0
	data["nowpage"] = "0"

	t_dir.Setup(ServersetUp.Zippath)
	_ = t_dir.Read("/")

	i := 0
	for _, str := range strings.Split(url[1:], "/") {
		if (i == 1) && (str != "") {
			// tmp, _ := strconv.Atoi(str)
			//			if tmp > 0 {
			//				if len(t_dir.Data) >= tmp {
			data["id"] = str
			// id = tmp - 1
			//				}
			//			}
		}
		if (i == 2) && (str != "") {
			tmp, _ := strconv.Atoi(str)
			if tmp > 0 {
				data["nowpage"] = str
			}
		}
		println(str)
		i++
	}
	filelist_t.ReadId(data["id"])
	//i = filelist_t.Tmp.Id
	//if i < 1 {
	//	Logout.Out(1, "id=%v err\n", data["id"])
	//	filelist_t.ReadId("1")
	//}
	// data["title"] = t_dir.Data[id].Name[1:]
	// filename := t_dir.Data[id].RootPath + t_dir.Data[id].Name
	if filelist_t.Tmp.Id != 0 {
		filename = ServersetUp.Zippath + filelist_t.Tmp.Zippass
		data["name"] = filelist_t.Tmp.Name
		data["title"] = filelist_t.Tmp.Zippass
		data["tag"] = filelist_t.Tmp.Tag
		fmt.Println(filename)
	} else {
		fmt.Fprintf(w, "%s", "err id")
		return
		// id = 1
		// filename = t_dir.Data[id].RootPath + t_dir.Data[id].Name
		// fmt.Println(filename)
	}
	t.ZipOpenSetup(filename)
	t.ZipReadList()
	data["pagemax"] = strconv.Itoa(t.Count)
	// output := ConvertData(ReadHtml("html/comic/view.html"), datap)
	output := ConvertData(ReadHtml("html_tmp/view.html"), datap)
	fmt.Fprintf(w, output)
	// fmt.Fprintf(w, "id=%vnowpage=%vpagemax=%v", data["id"], data["nowpage"], data["pagemax"])
}
