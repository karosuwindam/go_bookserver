package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"./dirread"
)

const BACKHTML = "html/index.html"
const UPLOAD = "upload"

func ck_upload_data(str string) int {
	var dirfolder dirread.Dirtype
	dirfolder.Setup(UPLOAD)
	_ = dirfolder.Read("/")
	output := 0
	for _, file := range dirfolder.Data {
		if strings.ToLower(str) == strings.ToLower(file.Name[1:]) {
			output = 1
		}
	}
	// if str != "About.txt" {
	// 	output = 0
	// }
	return output
}

func uploadlist(w http.ResponseWriter, r *http.Request) {
	var dirfolder dirread.Dirtype
	output := "[\n"
	dirfolder.Setup(UPLOAD)
	_ = dirfolder.Read("/")
	count := 0
	for _, file := range dirfolder.Data {
		if count == 0 {
			output += "\t{"
			output += "\n\t\t\"name\":" + "\"" + file.Name[1:] + "\""
			output += ",\n\t\t\"size\":" + "\"" + strconv.FormatInt(file.Size, 10) + "\""
			output += "\n\t}"
		} else {
			output += ",\n\t{"
			output += "\n\t\t\"name\":" + "\"" + file.Name[1:] + "\""
			output += ",\n\t\t\"size\":" + "\"" + strconv.FormatInt(file.Size, 10) + "\""
			output += "\n\t}"
		}
		count++
	}
	output += "\n]"
	fmt.Fprintf(w, "%v", output)

}
func upload(w http.ResponseWriter, r *http.Request) {
	var data []byte = make([]byte, 1024)
	var tmplength int64 = 0
	var output string
	urldata := ""
	searchdata := ""
	// var file multipart.File
	// var fileHeader *multipart.FileHeader
	// var e error
	// var data []byte = make([]byte, 1024)
	if r.Method == "POST" {
		// file, fileHeader, e := r.FormFile("data")
		file, fileHeader, e := r.FormFile("file")
		if e != nil {
			fmt.Fprintf(w, "%v", backHtmlUpload())
			return
		}
		writefilename := fileHeader.Filename
		fullfilepass := UPLOAD + "/"
		if strings.Index(writefilename, "pdf") > 0 {
			fullfilepass += "pdf/"
		} else if strings.Index(writefilename, "zip") > 0 {
			fullfilepass += "zip/"
		}
		fullfilepass += writefilename
		fp, err := os.Create(fullfilepass)
		if err != nil {
			Logout.Out(0, "%v: create file err\n", UPLOAD+"/"+writefilename)
		}
		defer fp.Close()
		defer file.Close()
		Logout.Out(1, "update file data :%v\n", writefilename)
		dataBaseUpdate(writefilename)
		for {
			n, e := file.Read(data)
			if n == 0 {
				break
			}
			if e != nil {
				return
			}
			fp.WriteAt(data, tmplength)
			tmplength += int64(n)
		}
		fmt.Printf("POST\n")
	} else {
		url := r.URL.Path
		// fmt.Println(url)
		count := 0
		for _, str := range strings.Split(url[1:], "/") {
			if count == 1 {
				searchdata = str
			}
			if count == 2 {
				urldata = str
				fmt.Println(str)
				break
			}
			count++
		}
		fmt.Printf("GET\n")
		if searchdata == "search" {
			output = "{\"flage\":" + strconv.Itoa(ck_upload_data(urldata)) + "}"
			fmt.Fprintf(w, "%v", output)
			return
		}
	}
	fmt.Fprintf(w, "%v", backHtmlUpload())
}

func backHtmlUpload() string {
	var output string
	fp, err := os.Open(BACKHTML)
	if err != nil {
		Logout.Out(0, "File Open err:%v\n", BACKHTML)
		log.Panic(err)
		return ""
	}
	defer fp.Close()
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		output += string(buf[:n])
	}
	return output
}

func dataBaseUpdate(name string) {
	var filename, typename string
	tmp := strings.Split(name, ".")
	i := 0
	for _, str := range tmp {
		if i == (len(tmp) - 1) {
			typename = str
		} else if i == 0 {
			filename += str
		} else {
			filename += "." + str
		}
		i++
	}
	_ = machdata(filename)
	if typename == "pdf" {

	}
	zipname := filename + ".zip"
	pdfname := filename + ".pdf"
	tag := bookname_t.Tmp.Title + "," + bookname_t.Tmp.Writer + "," + bookname_t.Tmp.Booktype + "," + bookname_t.Tmp.Ext
	filelist_t.Add(filename, pdfname, zipname, tag)
}
