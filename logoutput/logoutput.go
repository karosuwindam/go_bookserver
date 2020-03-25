package logoutput

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Data struct {
	fileName string
}

func (t *Data) Setup(str string) {
	t.fileName = str
}

func (t *Data) Out(printflag int, str string, v ...interface{}) {
	// fp, err := os.Create(t.fileName)
	var output string
	fp, err := os.OpenFile(t.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer fp.Close()
	output = "[" + time.Now().String() + "]  "
	output += str
	fmt.Fprintf(fp, output, v...)
	if printflag != 0 {
		fmt.Printf(str, v...)
	}
	// fp,err := os.OpenFile(t.FileName,os.O_APPEND)
}
