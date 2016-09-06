package printer

import (
	"fmt"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
	"github.com/davyxu/tabtoy/util"
)

type luaWriter struct {
	printerfile
}

func (self *luaWriter) RepeatedMessageBegin(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, msgCount int, indent int) {

	if indent == 1 {
		self.printer.WriteString(fmt.Sprintf("%s = {\n", fd.Name()))
	} else {
		self.printer.WriteString(fmt.Sprintf("%s = {", fd.Name()))
	}
}

// Value是消息的字段
func (self *luaWriter) WriteMessage(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, indent int) {

	if indent == 1 || fd.IsRepeated() {
		self.printer.WriteString("{")

	} else {
		self.printer.WriteString(fmt.Sprintf("%s = {", fd.Name()))
	}

	rawWriteMessage(&self.printer, self, msg, indent)

	self.printer.WriteString("}")

}

func (self *luaWriter) RepeatedMessageEnd(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, msgCount int, indent int) {
	self.printer.WriteString("}")
}

func (self *luaWriter) RepeatedValueBegin(fd *pbmeta.FieldDescriptor) {

}

// 普通值
func (self *luaWriter) WriteValue(fd *pbmeta.FieldDescriptor, value string, indent int) {

	var finalValue string
	switch fd.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_STRING,
		pbprotos.FieldDescriptorProto_TYPE_ENUM:
		finalValue = util.StringEscape(value)
	case pbprotos.FieldDescriptorProto_TYPE_INT64,
		pbprotos.FieldDescriptorProto_TYPE_UINT64:
		finalValue = fmt.Sprintf("\"%s\"", value)
	default:
		finalValue = value
	}

	self.printer.WriteString(fmt.Sprintf("%s = %s", fd.Name(), finalValue))
}

func (self *luaWriter) RepeatedValueEnd(fd *pbmeta.FieldDescriptor) {

}

func (self *luaWriter) WriteFieldSpliter() {

	self.printer.WriteString(", ")
}

// msg类型=XXFile
func (self *luaWriter) PrintMessage(msg *data.DynamicMessage) bool {

	self.printer.WriteString("local data = {\n\n")

	rawWriteMessage(&self.printer, self, msg, 0)

	self.printer.WriteString("\n\n}\n")

	/*

		data.ActorByID = {}
		for _, rec in pairs( data.Actor ) do

			data.ActorByID[rec.ID] = rec

		end

	*/

	// 输出lua索引
	fdset, lineFieldName := findMapperField(msg)

	if fdset == nil {
		return false
	}

	for _, fd := range fdset {

		mapperVarName := fmt.Sprintf("data.%sBy%s", lineFieldName, fd.Name())

		self.printer.WriteString("\n-- " + fd.Name() + "\n")
		self.printer.WriteString(mapperVarName + " = {}\n")
		self.printer.WriteString("for _, rec in pairs(data." + lineFieldName + ") do\n")
		self.printer.WriteString("\t" + mapperVarName + "[rec." + fd.Name() + "] = rec\n")
		self.printer.WriteString("end\n")
	}

	self.printer.WriteString("\nreturn data")

	return true
}

func findMapperField(msg *data.DynamicMessage) (fdset []*pbmeta.FieldDescriptor, lineFieldName string) {

	var lineMsgDesc *pbmeta.Descriptor
	// 找到行描述符
	for i := 0; i < msg.Desc.FieldCount(); i++ {
		fd := msg.Desc.Field(i)

		if fd.IsRepeated() {
			lineMsgDesc = fd.MessageDesc()
			lineFieldName = fd.Name()
			break
		}
	}

	// 在结构中寻找需要导出的lua字段
	for i := 0; i < lineMsgDesc.FieldCount(); i++ {
		fd := lineMsgDesc.Field(i)

		meta, ok := data.GetFieldMeta(fd)

		if !ok {
			return nil, ""
		}

		if meta == nil {
			continue
		}

		if !meta.LuaMapper {
			continue
		}

		fdset = append(fdset, fd)
	}

	return
}

func NewLuaWriter() IPrinter {

	self := &luaWriter{}
	self.printer.WriteString("-- Generated by github.com/davyxu/tabtoy\n")

	return self
}
