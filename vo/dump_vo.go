// vo doc

package vo

type DumpVo struct {
	WorkSpace string `json:"WorkSpace"`
	AppName string `json:"appName"`
	PodName string `json:"podName"`
	LimitCpu string `json:"LimitCpu"`
	LimitMem string `json:"LimitMem"`
	RequestCpu string `json:"RequestCpu"`
	RequestMem string `json:"RequestMem"`
	UsageCpu string `json:"UsageCpu"`
	UsageMem string `json:"UsageMem"`
}
