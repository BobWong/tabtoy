package main

import (
	"fmt"
	path2 "path"
	"strconv"
	"strings"

	"github.com/davyxu/tabtoy/util"
)

func main() {
	files := util.GetFileList(util.GetCurrentDirectory() + "/../xlsx")
	if util.Exists(util.GetCurrentDirectory() + "/../gen/GameConfig") {
		fmt.Println("--------------------- Start clear data ---------------------")
		util.DelDir(util.GetCurrentDirectory() + "/../gen/GameConfig")
	}
	globalClientCmd := `
		cd %s/../xlsx
		./../tabtoy --mode=v2 --g_csharp_out=../gen/GameConfig/CSharp/Globals.cs --combinename=GlobalsConfig --lan=zh_cn --log_enable=false Globals.xlsx
	`
	singleClientCmd := `
		cd %s/../xlsx
		./../tabtoy --mode=v2 --byte_out=../gen/GameConfig/Bytes/%sConfig.bytes --csharp_out=../gen/GameConfig/CSharp/%sConfig.cs --combinename=%sConfig --lan=zh_cn --log_enable=false Globals.xlsx %s
		`
	globalServerCmd := `
		cd %s/../xlsx
		./../tabtoy --mode=v2 --go_out=../gen/GameConfig/GoLang/static_data.go --json_out=../gen/GameConfig/Json/static_data.json --lan=zh_cn --log_enable=false Globals.xlsx `
	fmt.Println("--------------------- Start gen client data ---------------------")
	util.ExecuteCmd(fmt.Sprintf(globalClientCmd, util.GetCurrentDirectory()))
	for _, v := range files {
		fileSingleNameWithSuffix := path2.Base(v)
		fileSingleName := ""
		if strings.HasSuffix(fileSingleNameWithSuffix, "xlsm") {
			fileSingleName = strings.TrimSuffix(fileSingleNameWithSuffix, ".xlsm")
		}
		if strings.HasSuffix(fileSingleNameWithSuffix, "xlsx") {
			fileSingleName = strings.TrimSuffix(fileSingleNameWithSuffix, ".xlsx")
		}
		if strings.HasPrefix(fileSingleNameWithSuffix, "s_") {
			fileSingleName = strings.TrimPrefix(fileSingleName, "s_")
		}
		if fileSingleName == "Globals" {
			continue
		}

		util.ExecuteCmd(fmt.Sprintf(singleClientCmd, util.GetCurrentDirectory(), fileSingleName, fileSingleName, fileSingleName, fileSingleNameWithSuffix))

		if result, ok := util.GetFileSize(fmt.Sprintf("%s/../gen/GameConfig/Bytes/%sConfig.bytes", util.GetCurrentDirectory(), fileSingleName)); ok {
			tmp := float32(result) / 1024 / 1024
			fileSize, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", tmp), 32)
			if fileSize == 0.000 {
				tmp = float32(result) / 1024
				fileSize, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", tmp), 32)
				if fileSize == 0.000 {
					tmp = float32(result)
					fmt.Printf("Generate %sStatic.bytes====>>>%.3fB\n", fileSingleName, tmp)
				} else {
					fmt.Printf("Generate %sStatic.bytes====>>>%.3fK\n", fileSingleName, tmp)
				}
			} else {
				fmt.Printf("Generate %sStatic.bytes====>>>%.3fM\n", fileSingleName, tmp)
			}

		}
		globalServerCmd += fileSingleNameWithSuffix + " "
	}
	fmt.Println("--------------------- Client data gen end ---------------------")
	fmt.Println("--------------------- Start gen server data ---------------------")
	util.ExecuteCmd(fmt.Sprintf(globalServerCmd, util.GetCurrentDirectory()))
	fmt.Println("--------------------- Server data gen end ---------------------")
}
