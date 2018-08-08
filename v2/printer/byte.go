package printer

import (
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
)

const combineQFileVersion = 2

type qBinaryPrinter struct {
}

func (self *qBinaryPrinter) Run(g *Globals) *Stream {

	fileStresam := NewStream()
	fileStresam.WriteString("TON")
	fileStresam.WriteInt32(combineQFileVersion)

	for index, tab := range g.Tables {

		if !tab.LocalFD.MatchTag(".bin") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), tab.Name())
			continue
		}

		if !writeTableQBinary(fileStresam, tab, int32(index)) {
			return nil
		}

	}

	return fileStresam
}

func writeTableQBinary(tabStream *Stream, tab *model.Table, index int32) bool {

	var cols int
	var rows int
	rows = len(tab.Recs)
	if rows > 0 {
		cols = len(tab.LocalFD.Descriptors[len(tab.LocalFD.Descriptors)-1].Fields)
	} else {
		cols = 0
	}
	// 写入总行数
	tabStream.WriteInt32(int32(rows))

	for colIndex := 0; colIndex <= cols; colIndex++ {
		rowArrayCnt := make([]*Stream, rows)
		colStream := NewStream()
		arrayStreams := make([]*Stream, 100)

		isArrayStream := false
		structElementCnt := 0
		isRepeated := false

		log.Infof("====>>>>>>>>%d \n\r", colIndex)

		for rowIndex := 0; rowIndex < rows; rowIndex++ {
			if rowArrayCnt[rowIndex] == nil {
				rowArrayCnt[rowIndex] = NewStream()
			}

			node := tab.Recs[rowIndex].FindNodeByOrder(int32(colIndex))

			if node == nil {
				colStream.WriteInt32(0)
				rowArrayCnt[rowIndex].WriteInt32(0)
				log.Infof("检测到空行,写入单行数据长度: %d\n\r", 0)
				continue
			}
			if node.SugguestIgnore {
				colStream.WriteInt32(0)
				rowArrayCnt[rowIndex].WriteInt32(0)
				log.Infof("检测到结构体空行: %d\n\r", 0)
			}
			// 子节点数量
			if node.IsRepeated {
				colStream.WriteInt32(int32(len(node.Child)))
			}
			// 普通值
			if node.Type != model.FieldType_Struct {
				for _, valueNode := range node.Child {
					// 写入字段索引
					if !node.SugguestIgnore {
						colStream.WriteNodeValue(node.Type, valueNode)
					}
				}
			} else {
				isArrayStream = true
				if node.IsRepeated {
					isRepeated = true
					// 写入数组长度
					rowArrayCnt[rowIndex].WriteInt32(int32(len(node.Child)))
					log.Infof("写入单行数据长度: %d\n\r", len(node.Child))
				}
				// 记录数组元素个数
				structElementCnt = len(node.FieldDescriptor.Complex.Fields)
				// 遍历repeated的结构体
				for _, structNode := range node.Child {
					allField := structNode.FieldDescriptor.Complex.Fields
					for i := 0; i < len(allField); i++ {
						if arrayStreams[i] == nil {
							arrayStreams[i] = NewStream()
						}
						// 获取所有属性的名字和实际填写的属性做比较,查看是否有属性缺失
						hasContains := false
						structNodeIndex := 0
						for j := 0; j < len(structNode.Child); j++ {
							if structNode.Child[j].FieldDescriptor.Name == allField[i].Name {
								structNodeIndex = j
								hasContains = true
							}
						}
						if !hasContains {
							// 说明属性缺失
							// 写入默认值
							curFieldType := allField[i].Type

							defaultValue := &model.Node{}
							// 去读取@type表格的默认值
							defaultValue.Value = allField[i].Meta.GetString("Default")
							arrayStreams[i].WriteNodeValue(curFieldType, defaultValue)
							log.Infof("属性缺失 : %s \n\r", defaultValue.Value)
							continue
						}
						if structNode.Child[structNodeIndex].SugguestIgnore {
							// 说明属性缺失
							// 写入默认值
							curFieldType := structNode.Child[structNodeIndex].Type

							defaultValue := &model.Node{}
							// 去读取@type表格的默认值
							defaultValue.Value = structNode.Child[structNodeIndex].Meta.GetString("Default")
							arrayStreams[i].WriteNodeValue(curFieldType, defaultValue)
							log.Infof("属性缺失 : %s \n\r", defaultValue.Value)
							continue
						}
						// 写入字段索引
						// 值节点总是在第一个
						valueNode := structNode.Child[structNodeIndex].Child[0]
						childType := structNode.Child[structNodeIndex].Type
						arrayStreams[i].WriteNodeValue(childType, valueNode)
						log.Infof("获取真实值 : %s \n\r", valueNode.Value)
					}
					log.Infoln("--------------------")
				}
			}
		}
		if isArrayStream {
			if isRepeated {
				for i := 0; i < rows; i++ {
					tabStream.WriteBytes(rowArrayCnt[i].Buffer().Bytes())
				}
			}
			for i := 0; i < structElementCnt; i++ {

				tabStream.WriteBytes(arrayStreams[i].Buffer().Bytes())
			}
		} else {
			tabStream.WriteBytes(colStream.Buffer().Bytes())
		}
	}

	return true

}

func init() {

	RegisterPrinter("bytes", &qBinaryPrinter{})

}
