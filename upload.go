package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"./dirread"
)

const BACKHTML = "html/index.html"

// const UPLOAD = "upload"

func ck_upload_data(str string) int {
	var dirfolder dirread.Dirtype
	dirfolder.Setup(ServersetUp.Uploadpath)
	_ = dirfolder.Read("/")
	output := 0
	for _, file := range dirfolder.Data {
		if strings.ToLower(str) == strings.ToLower(file.Name[1:]) {
			output = 1
		}
	}
	return output
}

//upload
//upload時の動作作業
func upload(w http.ResponseWriter, r *http.Request) {
	var data []byte = make([]byte, 1024)
	var tmplength int64 = 0
	var output string
	urldata := ""
	searchdata := ""
	if r.Method == "POST" {
		file, fileHeader, e := r.FormFile("file")
		if e != nil {
			fmt.Fprintf(w, "%v", backHtmlUpload())
			return
		}
		writefilename := fileHeader.Filename
		fullfilepass := ServersetUp.Uploadpath + "/"
		if strings.Index(writefilename, "pdf") > 0 {
			fullfilepass += "pdf/"
		} else if strings.Index(writefilename, "zip") > 0 {
			fullfilepass += "zip/"
		}
		fullfilepass += writefilename
		fp, err := os.Create(fullfilepass)
		if err != nil {
			Logout.Out(0, "%v: create file err\n", ServersetUp.Uploadpath+"/"+writefilename)
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
		go uploadCHdata(writefilename)
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

//backHtmlUpload
//Upload時の戻り用HTML
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

//Upload実行時のデータベースアップロード
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
	_, kan := machdata(filename)
	if typename == "pdf" {

	}
	zipname := "[" + bookname_t.Tmp.Writer + "]" + bookname_t.Tmp.Title + kan + ".zip"
	pdfname := filename + ".pdf"
	tag := bookname_t.Tmp.Title + kan + "," + bookname_t.Tmp.Writer + "," + bookname_t.Tmp.Booktype + "," + bookname_t.Tmp.Ext
	filelist_t.Add(filename, pdfname, zipname, tag)
}

func uploadCHdata(str string) {
	var dirfolder dirread.Dirtype
	if strings.Index(str, "pdf") > 0 {
		filename := str[0 : len(str)-4]
		subcmd := "pdfimages" + " " + ServersetUp.Uploadpath + "/pdf" + "/" + str + " " + "tmp" + "/" + filename + " " + "-j"
		_, kan := machdata(filename)
		fmt.Println(subcmd)
		err := exec.Command("sh", "-c", subcmd).Run()
		if err != nil {

		} else {
			subcmd = "cp " + "tmp/" + filename + "-000.jpg " + "html/jpg/" + filename + ".jpg"
			_ = exec.Command("sh", "-c", subcmd).Run()
			dirfolder.Setup("tmp")
			_ = dirfolder.Read("/")
			zipname := "[" + bookname_t.Tmp.Writer + "]" + bookname_t.Tmp.Title + kan + ".zip"
			dest, errzip := os.Create(ServersetUp.Uploadpath + "/" + "zip" + "/" + zipname)
			if errzip != nil {
				return
			}
			zipWriter := zip.NewWriter(dest)
			defer zipWriter.Close()
			for _, file := range dirfolder.Data {
				if strings.Index(file.Name, filename) > 0 {
					_ = addToZip("tmp"+file.Name, zipWriter)
				}
			}
			subcmd = "rm -rf" + " " + "tmp/" + filename + "*"
			fmt.Println(subcmd)
			err = exec.Command("sh", "-c", subcmd).Run()

		}
	} else if strings.Index(str, "zip") > 0 {

	}

}
func addToZip(filename string, zipWriter *zip.Writer) error {
	info, _ := os.Stat(filename)
	hdr, _ := zip.FileInfoHeader(info)
	hdr.Name = filename
	for _, s := range strings.Split(filename, "/") {
		hdr.Name = s
	}
	f, err := zipWriter.CreateHeader(hdr)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	f.Write(body)
	return nil
}

/*func addToZip(filename string, zipWriter *zip.Writer) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()
	zipfilename := filename
	for _, s := range strings.Split(filename, "/") {
		zipfilename = s
	}
	writer, err := zipWriter.Create(zipfilename)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, src)
	if err != nil {
		return err
	}

	return nil
}
*/
