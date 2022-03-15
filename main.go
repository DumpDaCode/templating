package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main() {
	file, err := os.Open("sample.json")
	handleError(err)
	fileInfo, err := file.Stat()
	handleError(err)
	b := make([]byte, fileInfo.Size())
	_, err = file.Read(b)
	handleError(err)
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	handleError(err)
	templateContent, err := ioutil.ReadFile("./ycmu/result.tmpl")
	handleError(err)
	MKCLFunc := template.FuncMap{
		"hPad": func(value int) string {
			return strings.Repeat("\n", value)
		},
		"vPad": func(padString string, times int) string {
			return strings.Repeat(padString, times)
		},
		"wrap": func(data string, length int) []string {
			data = strings.TrimSpace(data)
			b := make([]string, 0)
			for len(data) > length {
				id := strings.LastIndex(data[:length], " ")
				b = append(b, data[:id]+strings.Repeat(" ", length-len(data[:id])))
				data = data[id+1:]
			}
			if len(data) != 0 {
				b = append(b, data+strings.Repeat(" ", length-len(data)))
			}
			return b
		},
		"center": func(data string, length int) string {
			times := (length - len(data)) / 2
			return strings.Repeat(" ", times) + data + strings.Repeat(" ", length-len(data)-times)
		},
		"inc": func(data int, byValue int) int {
			return data + byValue
		},
		"dec": func(data int, byValue int) int {
			return data - byValue
		},
	}
	MKCLFunc2 := template.FuncMap{
		"isNewPage": func(index int, currLineNo int, data interface{}, headLen int, maxLen int, maxStudents int) bool {
			templateContent, err := ioutil.ReadFile("./ycmu/support.tmpl")
			handleError(err)
			temp := template.New("support").Funcs(MKCLFunc)
			_, err = temp.Parse(string(templateContent))
			handleError(err)
			var p bytes.Buffer
			err = temp.Execute(&p, data)
			handleError(err)
			b := p.Bytes()
			sep := []byte{'\n'}
			noOfLine := bytes.Count(b, sep)
			return index%maxStudents == 0 || (currLineNo+noOfLine+headLen+1 > maxLen)
		},
	}
	temp := template.New("result").Funcs(MKCLFunc).Funcs(MKCLFunc2)
	_, err = temp.Parse(string(templateContent))
	handleError(err)
	var p bytes.Buffer
	err = temp.Execute(&p, data)
	handleError(err)
	file, err = os.Create("output.txt")
	handleError(err)
	_, err = file.Write(p.Bytes())
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
