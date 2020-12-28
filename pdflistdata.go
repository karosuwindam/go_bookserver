package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"./dirread"
	"./zipopen"
)

type folderckdata struct {
	Name string `json:name`
	Flag int    `json:flag`
}

type zippdfdata struct {
	Name    string `json:"name"`
	Pdfpass string `json:"pdfpass"`
	Zippass string `json:"zippass`
	Tag     string `json:"tag"`
}

type jsonFolderCkeck struct {
	Pdf  folderckdata `json:pdf`
	Zip  folderckdata `json:zip`
	Jpg  folderckdata `json:jpg`
	Data zippdfdata   `json:Data`
}

type FolderCKeck struct {
	pdfdata   []string
	zipdata   []string
	zipindata []string
	jpgdata   []string
}

func FolderDataSetup() FolderCKeck {
	var tmp FolderCKeck
	tmp.readPdfDir()
	tmp.readZipDir()
	tmp.readJpgDir()
	return tmp
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
type filelists struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Pdfpass    string    `json:"pdfpass"`
	Zippass    string    `json:"zippass`
	Tag        string    `json:"tag"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

//フォルダ内のファイル名取得_pdf
func (t *FolderCKeck) readPdfDir() {
	tmp := []string{}
	var t_dir dirread.Dirtype
	t_dir.Setup(ServersetUp.Pdfpath)
	t_dir.Read("/")

	for _, s := range t_dir.Data {
		if !s.Folder {
			tmp = append(tmp, s.Name[1:])
		}
	}
	t.pdfdata = tmp
}

//フォルダ内のファイル名取得_zip
func (t *FolderCKeck) readZipDir() {
	tmp := []string{}
	tmp_in := []string{}
	var t_dir dirread.Dirtype
	t_dir.Setup(ServersetUp.Zippath)
	t_dir.Read("/")

	for _, s := range t_dir.Data {
		if !s.Folder {
			var ziptmp zipopen.File
			ziptmp.ZipOpenSetup(ServersetUp.Zippath + s.Name)
			ziptmp.ZipReadList()
			tmp_in = append(tmp_in, ziptmp.Name[0][:strings.Index(ziptmp.Name[0], "-000")])
			tmp = append(tmp, s.Name[1:])
		}
	}
	t.zipdata = tmp
	t.zipindata = tmp_in
}

//フォルダ内のファイル名取得_jpg
func (t *FolderCKeck) readJpgDir() {
	tmp := []string{}
	var t_dir dirread.Dirtype
	t_dir.Setup("./html/jpg")
	t_dir.Read("/")

	for _, s := range t_dir.Data {
		if !s.Folder {
			tmp = append(tmp, s.Name[1:])
		}
	}
	t.jpgdata = tmp
}

func (t *FolderCKeck) machdataZip(str string) bool {
	var fc []filelists
	filelist_t.Read("zippass='" + str + "'")
	tmp := filelist_t.JsonOutList()
	json.Unmarshal([]byte(tmp), &fc)
	if len(fc) == 0 {
		return false
	}
	// fmt.Println(str)
	return str == fc[0].Zippass
}
func (t *FolderCKeck) machdata(str string) (string, string) {

	output := "{}"
	data := map[string]string{}
	kan := ""
	num1 := strings.Index(str, ".pdf")

	data["keyword"] = str[:num1]
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
		kan = "0" + strconv.Itoa(num)
	}
	return output, kan
}

func (t *FolderCKeck) CheckData() []jsonFolderCkeck {
	output := []jsonFolderCkeck{}
	zip_tmp := t.zipdata
	zipin_tmp := t.zipindata
	jpg_tmp := t.jpgdata
	for _, s := range t.pdfdata {
		// fmt.Println(s)
		var ary jsonFolderCkeck
		ary.Pdf.Name = s
		var fc booknames
		tmp, kan := t.machdata(s)
		if tmp != "{}" {
			ary.Pdf.Flag = 1
			json.Unmarshal([]byte(tmp), &fc)
			//zipck
			zip_tmp_t := []string{}
			zipin_tmp_t := []string{}
			for i, ss := range zip_tmp {
				// fmt.Println(i)
				ck := "[" + fc.Writer + "]" + fc.Title
				if ss == ck+".zip" {
					ary.Zip.Name = ss
					if t.machdataZip(ss) {
						ary.Zip.Flag = 1
					}
					for j, sss := range zip_tmp[i+1:] {
						zip_tmp_t = append(zip_tmp_t, sss)
						zipin_tmp_t = append(zipin_tmp_t, zipin_tmp[i+1+j])
					}
					ary.Data.Name = s[:strings.Index(s, ".pdf")]
					ary.Data.Zippass = ss
					ary.Data.Pdfpass = s
					ary.Data.Tag = fc.Title + "," + fc.Writer + "," + fc.Brand + "," + fc.Ext

					break
				} else if ss == ck+kan+".zip" {
					ary.Zip.Name = ss
					if t.machdataZip(ss) {
						ary.Zip.Flag = 1
					}
					for j, sss := range zip_tmp[i+1:] {
						zip_tmp_t = append(zip_tmp_t, sss)
						zipin_tmp_t = append(zipin_tmp_t, zipin_tmp[i+1+j])
					}
					ary.Data.Name = s[:strings.Index(s, ".pdf")]
					ary.Data.Zippass = ss
					ary.Data.Pdfpass = s
					ary.Data.Tag = fc.Title + kan + "," + fc.Writer + "," + fc.Brand + "," + fc.Ext

					break
				} else if ss == ck+"0"+kan+".zip" {
					ary.Zip.Name = ss
					if t.machdataZip(ss) {
						ary.Zip.Flag = 1
					}
					for j, sss := range zip_tmp[i+1:] {
						zip_tmp_t = append(zip_tmp_t, sss)
						zipin_tmp_t = append(zipin_tmp_t, zipin_tmp[i+1+j])
					}
					ary.Data.Name = s[:strings.Index(s, ".pdf")]
					ary.Data.Zippass = ss
					ary.Data.Pdfpass = s
					ary.Data.Tag = fc.Title + "0" + kan + "," + fc.Writer + "," + fc.Brand + "," + fc.Booktype + "," + fc.Ext

					break
				}
				zip_tmp_t = append(zip_tmp_t, ss)
				zipin_tmp_t = append(zipin_tmp_t, zipin_tmp[i])
			}
			zip_tmp = zip_tmp_t
			zipin_tmp = zipin_tmp_t
			//jpgck
			ary.Jpg.Name = s[:strings.Index(s, ".pdf")] + ".jpg"
			// jpg_tmp_t := []string{}
			for _, ss := range jpg_tmp {
				if ss == ary.Jpg.Name {
					ary.Jpg.Flag = 1
					// for _, sss := range jpg_tmp[:i+1] {
					// 	jpg_tmp_t = append(jpg_tmp_t, sss)
					// }
					break
				}
				// jpg_tmp_t = append(jpg_tmp_t, ss)
			}
			// jpg_tmp = jpg_tmp_t

			// fc.Name
			// fmt.Println(fc)
		} else {
			zip_tmp_t := []string{}
			zipin_tmp_t := []string{}

			for i, ss := range zipin_tmp {
				if (ss + ".pdf") == ary.Pdf.Name {
					ary.Zip.Name = zip_tmp[i]
					for j, sss := range zipin_tmp[i+1:] {
						zip_tmp_t = append(zip_tmp_t, zip_tmp[i+1+j])
						zipin_tmp_t = append(zipin_tmp_t, sss)
					}
					break
				}
				zip_tmp_t = append(zip_tmp_t, zip_tmp[i])
				zipin_tmp_t = append(zipin_tmp_t, ss)

			}
			zip_tmp = zip_tmp_t
			zipin_tmp = zipin_tmp_t
			//jpgck
			ary.Jpg.Name = s[:strings.Index(s, ".pdf")] + ".jpg"
			// jpg_tmp_t := []string{}
			for _, ss := range jpg_tmp {
				if ss == ary.Jpg.Name {
					ary.Jpg.Flag = 1
					// for _, sss := range jpg_tmp[:i+1] {
					// 	jpg_tmp_t = append(jpg_tmp_t, sss)
					// }
					break
				}
				// jpg_tmp_t = append(jpg_tmp_t, ss)
			}
			// jpg_tmp = jpg_tmp_t
		}
		output = append(output, ary)
	}
	return output
}
