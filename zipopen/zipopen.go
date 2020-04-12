package zipopen

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
)

type File struct {
	path  string
	Name  []string
	Count int
	Flag  bool
}

func (t *File) ZipOpenSetup(s string) {
	t.path = s
	t.Count = 0
	t.Flag = t.exists()
}
func (t *File) exists() bool {
	_, err := os.Stat(t.path)
	return err == nil
}
func (t *File) ZipReadList() {
	var tmp []string
	i := 0
	r, err := zip.OpenReader(t.path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	for _, f := range r.File {

		tmp = append(tmp, f.Name)
		i++
	}
	t.Count = i
	t.Name = tmp
}

func (t *File) ZipRead(i int) *bytes.Buffer {
	buf := new(bytes.Buffer)
	r, err := zip.OpenReader(t.path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	for _, f := range r.File {
		if t.Name[i] == f.Name {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(buf, rc)
			if err != nil {
				log.Fatal(err)
			}
			rc.Close()
		}
	}
	return buf

}
