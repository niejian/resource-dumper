// dumper doc

package dumper

import (
	"fmt"
	"log"
	"resource-dumper/api"
)

func GetPodDetail(ns string)  {
	podList, err := api.GetUserPod(ns)
	if err != nil{
		log.Printf("获取pod失败")
	}
	items := podList.Items

	for _, pod := range items {
		fmt.Println(pod)

	}

}
