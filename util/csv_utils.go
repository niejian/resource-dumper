// util doc

package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"resource-dumper/vo"
	"strconv"
	"time"
)

func CsvWrite(vos []vo.DumpVo )  {
	pwd, _ := os.Getwd()
	timeUnix := time.Now().Unix()

	targetPath := fmt.Sprintf("%s%s%s%d%s", pwd, "/","resource-",timeUnix, ".csv")
	fmt.Println(targetPath)
	f, err := os.Create(targetPath)
	if err != nil{
		fmt.Println(err)
		return
	}

	defer f.Close()
	var data = make([][]string, len(vos))
	//data[0] = []string{"标题", "作者", "时间"}
	//data[1] = []string{"羊皮卷", "鲁迅", "2008"}
	//data[2] = []string{"易筋经", "唐生", "665"}

	for i, vo := range vos {
		fmt.Printf("%v  \n",vo)
		var cpuP float64
		var memP float64
		if "" != vo.LimitCpu && "" != vo.UsageCpu  {
			// 使用率计算
			limitCpu, _ := strconv.ParseFloat(vo.LimitCpu, 64)
			usageCpu, _ := strconv.ParseFloat(vo.UsageCpu, 64)
			cpuP = usageCpu/limitCpu
		}

		if "" != vo.LimitMem && "" != vo.UsageMem {
			limitMem, _ := strconv.ParseFloat(vo.LimitMem, 64)
			usageMem, _ := strconv.ParseFloat(vo.UsageMem, 64)
			memP = usageMem / limitMem
		}
		data[i] = []string{vo.WorkSpace, vo.AppName, vo.PodName, vo.RequestCpu, vo.RequestMem, vo.LimitCpu, vo.LimitMem,
			vo.UsageCpu, vo.UsageMem, fmt.Sprintf("%f%s", cpuP*100, "%"), fmt.Sprintf("%f%s", memP * 100, "%")}
	}

	f.WriteString("\xEF\xBB\xBF")  // 写入一个UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()

}
