package printer

import (
	"fmt"
	"text/template"

	"strings"

	"tabtoy/v2/i18n"
	"tabtoy/v2/model"
)

const tsTemplate = `// Generated by github.com/davyxu/tabtoy
// Version: {{.ToolVersion}}
// DO NOT EDIT!!

namespace {{.Namespace}}{{$globalIndex:=.Indexes}}{{$verticalFields:=.VerticalFields}}
{
	{{range .Enums}} {{if ne .DefinedTable "Globals"}}
	// Defined in table: {{.DefinedTable}}
	export enum {{.Name}}
	{
	{{range .Fields}}	
		{{.Comment}}
		{{.FieldDescriptor.Name}} = {{.FieldDescriptor.EnumValue}}, {{.Alias}}
	{{end}}
	}


	{{end}}{{end}}
	{{range .Classes}} {{if ne .DefinedTable "Globals"}}

	// Defined in table: {{.DefinedTable}}
	{{.TsClassHeader}}
	export class {{.Name}}
	{
	{{range .Fields}}	
		{{.Comment}}
		{{.TypeCode}} {{.Alias}}
	{{end}}
	

	} {{end}}{{end}}

}
`

type tsIndexField struct {
	TableIndex
}

func (self tsIndexField) IndexName() string {
	return self.Index.Name
}

func (self tsIndexField) RowType() string {
	return self.Row.Complex.Name
}

func (self tsIndexField) RowName() string {
	return self.Row.Name
}

func (self tsIndexField) IndexType() string {

	switch self.Index.Type {
	case model.FieldType_Int32:
		return "int"
	case model.FieldType_UInt32:
		return "uint"
	case model.FieldType_Int64:
		return "long"
	case model.FieldType_UInt64:
		return "ulong"
	case model.FieldType_String:
		return "string"
	case model.FieldType_Float:
		return "float"
	case model.FieldType_Bool:
		return "bool"
	case model.FieldType_Enum:

		return self.Index.Complex.Name
	default:
		log.Errorf("%s can not be index ", self.Index.String())
	}

	return "unknown"
}

type tsField struct {
	*model.FieldDescriptor

	IndexKeys []*model.FieldDescriptor

	parentStruct *tsStructModel
}

func (self tsField) Alias() string {

	v := self.FieldDescriptor.Meta.GetString("Alias")
	if v == "" {
		return ""
	}

	return "// " + v
}

func (self tsField) MakeIndex() bool {

	v := self.FieldDescriptor.Meta.GetBool("MakeIndex")

	return v
}

func (self tsField) GetTableName() string {
	return strings.Replace(self.FieldDescriptor.Parent.Name, "Define", "", 1)
}

func (self tsField) Comment() string {

	if self.FieldDescriptor.Comment == "" {
		return ""
	}

	// zjwps 建议修改
	return "/// <summary> \n		/// " + strings.Replace(self.FieldDescriptor.Comment, "\n", "\n		///", -1) + "\n		/// </summary>"
}

func (self tsField) SReadCode() string {

	var raw string
	var baseType string
	var readType string
	var hasReadArgs string

	switch self.Type {
	case model.FieldType_Int32:
		raw = "int[]"
		if self.IsRepeated {
			raw = "int[][]"
		}
		baseType = "Int32"
		readType = "reader"
	case model.FieldType_UInt32:
		raw = "uint[]"
		if self.IsRepeated {
			raw = "uint[][]"
		}
		baseType = "UInt32"
		readType = "reader"
	case model.FieldType_Int64:
		raw = "long[]"
		if self.IsRepeated {
			raw = "long[][]"
		}
		baseType = "Int64"
		readType = "reader"
	case model.FieldType_UInt64:
		raw = "ulong[]"
		if self.IsRepeated {
			raw = "ulong[][]"
		}
		baseType = "UInt64"
		readType = "reader"
	case model.FieldType_String:
		raw = "string[]"
		if self.IsRepeated {
			raw = "string[][]"
		}
		baseType = "String"
		readType = "reader"
	case model.FieldType_Float:
		raw = "float[]"
		if self.IsRepeated {
			raw = "float[][]"
		}
		baseType = "Float"
		readType = "reader"
	case model.FieldType_Bool:
		raw = "bool[]"
		if self.IsRepeated {
			raw = "bool[][]"
		}
		baseType = "Bool"
		readType = "reader"
	case model.FieldType_Enum:
		if self.Complex == nil {
			log.Errorln("unknown enum type ", self.Type)
			return "unknown"
		}

		raw = fmt.Sprintf("List<%s>", self.Complex.Name)
		if self.IsRepeated {
			raw = fmt.Sprintf("List<List<%s>>", self.Complex.Name)
		}
		baseType = self.Complex.Name
		readType = fmt.Sprintf("%sHelp", self.Complex.Name)
		hasReadArgs = "reader,"

	case model.FieldType_Struct:
		if self.Complex == nil {
			log.Errorln("unknown struct type ", self.Type, self.FieldDescriptor.Name, self.FieldDescriptor.Parent.Name)
			return "unknown"
		}

		raw = self.Complex.Name
		raw = fmt.Sprintf("List<%s>", self.Complex.Name)
		if self.IsRepeated {
			raw = fmt.Sprintf("List<List<%s>>", self.Complex.Name)
		}
		baseType = self.Complex.Name
		readType = self.Complex.Name
		hasReadArgs = "reader,"
	default:
		raw = "unknown"
	}

	if self.IsRepeated {
		return fmt.Sprintf("%s %ss = %s.Read%sArray2(%slen);", raw, self.Name, readType, baseType, hasReadArgs)
	}

	return fmt.Sprintf("%s %ss = %s.Read%sArray(%slen);", raw, self.Name, readType, baseType, hasReadArgs)
}

func (self tsField) ReadCode() string {

	var baseType string

	var descHandlerCode string

	var fullTypeName string

	switch self.Type {
	case model.FieldType_Int32:
		baseType = "Int32"
	case model.FieldType_UInt32:
		baseType = "UInt32"
	case model.FieldType_Int64:
		baseType = "Int64"
	case model.FieldType_UInt64:
		baseType = "UInt64"
	case model.FieldType_String:
		baseType = "String"
	case model.FieldType_Float:
		baseType = "Float"
	case model.FieldType_Bool:
		baseType = "Bool"
	case model.FieldType_Enum:

		if self.Complex == nil {
			log.Errorln("unknown enum type ", self.Type)
			return "unknown"
		}

		baseType = fmt.Sprintf("(%s)reader.ReadInt32()", self.Complex.Name)
	case model.FieldType_Struct:
		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("Struct<%s>", self.Complex.Name)
		fullTypeName = fmt.Sprintf(", \"%s%s\"", "Data.", self.Complex.Name)
	}

	if self.Type == model.FieldType_Struct {
		descHandlerCode = fmt.Sprintf("%sDeserializeHandler", self.Complex.Name)
	}

	if self.IsRepeated {
		if self.Type == model.FieldType_Enum {
			return fmt.Sprintf("ins.%s.Add (%s);", self.Name, baseType)
		}
		return fmt.Sprintf("ins.%s.Add( reader.Read%s(%s%s) );", self.Name, baseType, descHandlerCode, fullTypeName)
	} else {
		if self.Type == model.FieldType_Enum {
			// log.Infof("enum")
			return fmt.Sprintf("ins.%s = %s ;", self.Name, baseType)
		}

		return fmt.Sprintf("ins.%s = reader.Read%s(%s%s);", self.Name, baseType, descHandlerCode, fullTypeName)
	}

}

func (self tsField) IsArraySimple() bool {
	if self.Complex == nil && self.IsRepeated {
		return true
	}

	return false
}

func (self tsField) Tag() string {

	if self.parentStruct.IsCombine() {
		tag := model.MakeTag(int32(model.FieldType_Table), self.Order)

		return fmt.Sprintf("0x%x", tag)
	}

	return fmt.Sprintf("0x%x", self.FieldDescriptor.Tag())
}

func (self tsField) StructName() string {
	if self.Complex == nil {
		return "[NotComplex]"
	}

	return self.Complex.Name
}

func (self tsField) IsVerticalStruct() bool {
	if self.FieldDescriptor.Complex == nil {
		return false
	}

	return self.FieldDescriptor.Complex.File.Pragma.GetBool("Vertical")
}

func (self tsField) ArrayTypeCode() string {

	var raw string

	switch self.Type {
	case model.FieldType_Int32:
		raw = "number"
	case model.FieldType_UInt32:
		raw = "number"
	case model.FieldType_Int64:
		raw = "number"
	case model.FieldType_UInt64:
		raw = "number"
	case model.FieldType_String:
		raw = "string"
	case model.FieldType_Float:
		raw = "number"
	case model.FieldType_Bool:
		raw = "boolean"
	default:
		raw = "unknown"
	}
	return raw
}

func (self tsField) TypeCode() string {

	var raw string

	switch self.Type {
	case model.FieldType_Int32:
		raw = "number"
	case model.FieldType_UInt32:
		raw = "number"
	case model.FieldType_Int64:
		raw = "number"
	case model.FieldType_UInt64:
		raw = "number"
	case model.FieldType_String:
		raw = "string"
	case model.FieldType_Float:
		raw = "number"
	case model.FieldType_Bool:
		raw = "boolean"
	case model.FieldType_Enum:
		if self.Complex == nil {
			log.Errorln("unknown enum type ", self.Type)
			return "unknown"
		}

		raw = self.Complex.Name
	case model.FieldType_Struct:
		if self.Complex == nil {
			log.Errorln("unknown struct type ", self.Type, self.FieldDescriptor.Name, self.FieldDescriptor.Parent.Name)
			return "unknown"
		}

		raw = self.Complex.Name

		// 非repeated的结构体
		if !self.IsRepeated {
			return fmt.Sprintf("public %s : %s = new %s();", self.Name, raw, raw)
		}

	default:
		raw = "unknown"
	}

	if self.IsRepeated {
		return fmt.Sprintf("public %s : List<%s> = new List<%s>();", self.Name,raw, raw)
	}

	return fmt.Sprintf("public %s : %s = %s;",  self.Name,raw, wrapCSharpDefaultValue(self.FieldDescriptor))
}

func wrapTsDefaultValue(fd *model.FieldDescriptor) string {
	switch fd.Type {
	case model.FieldType_Enum:
		return fmt.Sprintf("%s.%s", fd.Complex.Name, fd.DefaultValue())
	case model.FieldType_String:
		return fmt.Sprintf("\"%s\"", fd.DefaultValue())
	case model.FieldType_Float:
		return fmt.Sprintf("%s", fd.DefaultValue())
	}

	return fd.DefaultValue()
}

type tsStructModel struct {
	*model.Descriptor
	Fields        []tsField
	IndexedFields []tsField // 与csharpField.IndexKeys组成树状的索引层次
}

func (self *tsStructModel) TsClassHeader() string {

	// zjwps 提供需求
	return self.File.Pragma.GetString("TsClassHeader")
}

func (self *tsStructModel) DefinedTable() string {
	return self.File.Name
}

func (self *tsStructModel) Name() string {
	return self.Descriptor.Name
}

func (self *tsStructModel) IsTableDefine() bool {
	return strings.Contains(self.Descriptor.Name, "Define")
}

func (self *tsStructModel) TableConfigName() string {
	return strings.Replace(self.Descriptor.Name, "Define", "Config", 1)
}

func (self *tsStructModel) TableDefineName() string {
	return strings.Replace(self.Descriptor.Name, "Define", "", 1)
}

func (self *tsStructModel) IsTableStruct() bool {
	return strings.Contains(self.Descriptor.Name, "Define") == false && strings.Contains(self.Descriptor.Name, "Config") == false
}

func (self *tsStructModel) IsCombine() bool {
	return self.Descriptor.Usage == model.DescriptorUsage_CombineStruct
}

type tsFileModel struct {
	Namespace   string
	ToolVersion string
	Classes     []*tsStructModel
	Enums       []*tsStructModel
	Indexes     []tsIndexField // 全局的索引

	VerticalFields []tsField

	GenSerializeCode bool
}

type tsPrinter struct {
}

func (self *tsPrinter) Run(g *Globals) *Stream {

	tpl, err := template.New("ts").Parse(tsTemplate)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	var m tsFileModel

	if g.PackageName != "" {
		m.Namespace = g.PackageName
	} else {
		m.Namespace = g.FileDescriptor.Pragma.GetString("Package")
	}

	m.ToolVersion = g.Version
	m.GenSerializeCode = g.GenCSSerailizeCode

	// combinestruct的全局索引
	for _, ti := range g.GlobalIndexes {

		// 索引也限制
		if !ti.Index.Parent.File.MatchTag(".ts") {
			continue
		}

		m.Indexes = append(m.Indexes, tsIndexField{TableIndex: ti})
	}

	// 遍历所有类型
	for _, d := range g.FileDescriptor.Descriptors {

		// 这给被限制输出
		if !d.File.MatchTag(".ts") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), d.Name)
			continue
		}

		var sm tsStructModel
		sm.Descriptor = d

		switch d.Kind {
		case model.DescriptorKind_Struct:
			m.Classes = append(m.Classes, &sm)
		case model.DescriptorKind_Enum:
			m.Enums = append(m.Enums, &sm)
		}

		// 遍历字段
		for _, fd := range d.Fields {

			// 对CombineStruct的XXDefine对应的字段
			if d.Usage == model.DescriptorUsage_CombineStruct {

				// 这个字段被限制输出
				if fd.Complex != nil && !fd.Complex.File.MatchTag(".ts") {
					continue
				}

				// 这个结构有索引才创建
				if fd.Complex != nil && len(fd.Complex.Indexes) > 0 {

					// 被索引的结构
					indexedField := tsField{FieldDescriptor: fd, parentStruct: &sm}

					// 索引字段
					for _, key := range fd.Complex.Indexes {
						indexedField.IndexKeys = append(indexedField.IndexKeys, key)
					}

					sm.IndexedFields = append(sm.IndexedFields, indexedField)
				}

				if fd.Complex != nil && fd.Complex.File.Pragma.GetBool("Vertical") {
					m.VerticalFields = append(m.VerticalFields, tsField{FieldDescriptor: fd, parentStruct: &sm})
				}

			}

			csField := tsField{FieldDescriptor: fd, parentStruct: &sm}

			sm.Fields = append(sm.Fields, csField)

		}

	}

	bf := NewStream()

	err = tpl.Execute(bf.Buffer(), &m)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	return bf
}

func (self tsField) GetArrayValue() string {
	var baseType string

	var descHandlerCode string

	var fullTypeName string

	switch self.Type {
	case model.FieldType_Int32:
		baseType = "Int32"
	case model.FieldType_UInt32:
		baseType = "UInt32"
	case model.FieldType_Int64:
		baseType = "Int64"
	case model.FieldType_UInt64:
		baseType = "UInt64"
	case model.FieldType_String:
		baseType = "String"
	case model.FieldType_Float:
		baseType = "Float"
	case model.FieldType_Bool:
		baseType = "Bool"
	case model.FieldType_Enum:

		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("(%s)reader.ReadInt32()", self.Complex.Name)
	case model.FieldType_Struct:
		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("Struct<%s>", self.Complex.Name)
		fullTypeName = fmt.Sprintf(", \"%s%s\"", "Data.", self.Complex.Name)
	}

	if self.Type == model.FieldType_Struct {
		descHandlerCode = fmt.Sprintf("%sDeserializeHandler", self.Complex.Name)
	}

	if self.IsRepeated {
		if self.Type == model.FieldType_Enum {
			return fmt.Sprintf("ins.%s.Add (%s);", self.Name, baseType)
		}
		return fmt.Sprintf("ins.%s.Add( reader.Read%s(%s%s) );", self.Name, baseType, descHandlerCode, fullTypeName)
	} else {
		if self.Type == model.FieldType_Enum {
			log.Infof("enum")
			return fmt.Sprintf("ins.%s = (%s)array%s[i];", self.Name, self.Complex.Name, self.Name)
		}

		return fmt.Sprintf("ins.%s = array%s[i];", self.Name, self.Name)
	}
}

func (self tsField) ReadArrayCode() string {

	var baseType string

	var descHandlerCode string

	var fullTypeName string

	switch self.Type {
	case model.FieldType_Int32:
		baseType = "Int32"
	case model.FieldType_UInt32:
		baseType = "UInt32"
	case model.FieldType_Int64:
		baseType = "Int64"
	case model.FieldType_UInt64:
		baseType = "UInt64"
	case model.FieldType_String:
		baseType = "String"
	case model.FieldType_Float:
		baseType = "Float"
	case model.FieldType_Bool:
		baseType = "Bool"
	case model.FieldType_Enum:

		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("(%s)reader.ReadInt32()", self.Complex.Name)
	case model.FieldType_Struct:
		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("Struct<%s>", self.Complex.Name)
		fullTypeName = fmt.Sprintf(", \"%s%s\"", "Data.", self.Complex.Name)
	}

	if self.Type == model.FieldType_Struct {
		descHandlerCode = fmt.Sprintf("%sDeserializeHandler", self.Complex.Name)
	}

	if self.IsRepeated {
		if self.Type == model.FieldType_Enum {
			return fmt.Sprintf("ins.%s.Add (%s);", self.Name, baseType)
		}
		return fmt.Sprintf("ins.%s.Add( reader.Read%s(%s%s) );", self.Name, baseType, descHandlerCode, fullTypeName)
	} else {
		if self.Type == model.FieldType_Enum {
			log.Infof("enum")
			//return fmt.Sprintf("ins.%s = %s ;", self.Name, baseType)
			return fmt.Sprintf("var array%s = reader.ReadInt32Array(len);", self.Name)
		}

		return fmt.Sprintf("var array%s = reader.Read%sArray(len);", self.Name, baseType)
	}

}

func init() {

	RegisterPrinter("ts", &tsPrinter{})

}