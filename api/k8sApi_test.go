// api doc

package api

import (
	"fmt"
	"testing"
)

func TestListNodes(t *testing.T) {
	t.Run("获取节点信息", func(t *testing.T) {
		ListNodes()
	})

}

func TestGetPodDetail(t *testing.T) {
	t.Run("获取pod信息", func(t *testing.T) {
		fmt.Println("===>")
		podDetail := GetPodDetail("xxl-job-admin-6df677cd5f-57jll", "xxl-job")
		fmt.Println(podDetail)
	})
}
