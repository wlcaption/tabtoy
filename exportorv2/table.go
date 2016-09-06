package exportorv2

import (
	"bytes"
	"fmt"
	"os"

	"github.com/davyxu/tabtoy/util"
)

type Table struct {
	list []*Record
	buf  bytes.Buffer

	needSpliter bool
}

func (self *Table) Add(r *Record) {
	self.list = append(self.list, r)
}

func (self *Table) Print(rootName string) bool {

	self.Printf("# Generated by github.com/davyxu/tabtoy\n")

	// 遍历每一行
	for _, r := range self.list {

		self.Printf("%s { ", rootName)

		// 遍历每一列
		for _, cell := range r.cells {

			self.PrintSpliter(" ")

			if cell.IsRepeated {

				self.Printf("%s : [ ", cell.Name)

				for _, rv := range cell.ValueList {
					self.PrintSpliter(", ")
					self.Printf("%s", valueWrapper(cell.Type, rv))
					self.needSpliter = true
				}

				self.Printf(" ]")

			} else {

				self.Printf("%s : %s", cell.Name, valueWrapper(cell.Type, cell.Value))
				self.needSpliter = true
			}

		}

		self.Printf(" }\n")

		self.needSpliter = false
	}

	return true

}

func (self *Table) PrintSpliter(spliter string) {
	if self.needSpliter {

		self.Printf("%s", spliter)

		self.needSpliter = false
	}
}

func valueWrapper(t FieldType, v string) string {

	switch t {
	case FieldType_String:
		return util.StringEscape(v)
	}

	return v
}

func (self *Table) Printf(format string, args ...interface{}) {
	self.buf.WriteString(fmt.Sprintf(format, args...))
}

func (self *Table) WriteToFile(filename string) bool {

	// 创建输出文件
	file, err := os.Create(filename)
	if err != nil {
		log.Errorln(err.Error())
		return false
	}

	// 写入文件头

	file.WriteString(self.buf.String())

	file.Close()

	return true
}
