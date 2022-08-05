// resource_dumper doc

package main

import (
	"fmt"
	"resource-dumper/api"
	"resource-dumper/util"
	"resource-dumper/vo"
)

func main()  {
	var voList []vo.DumpVo
	fmt.Println("===>")
	//api.GetPodDetail("mh-pdm-goodsbiz-v1-f7db6d7c5-ptdbw", "mh-pdm-parent")
	//if 1 == 1 {
	//	return
	//}
	// 获取所有workspace
	userNs, _ := api.GetUserNs()
	for workSpace, nsList := range userNs {

		for _, ns := range nsList {
			// 获取deploy 信息
			podList, _ := api.GetUserPod(ns)
			// 获取pod资源情况
			for _, pod := range podList.Items {
				dumpVo := api.GetPodDetail(pod.Name, ns)
				dumpVo.WorkSpace = workSpace
				dumpVo.PodName = pod.Name
				dumpVo.AppName = ns
				voList = append(voList, dumpVo)

			}
		}
	}

	util.CsvWrite(voList)


}
