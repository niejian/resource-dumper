// api doc

package api

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"path/filepath"
	"resource-dumper/vo"
	"strconv"
	"strings"
)

var clientSet *kubernetes.Clientset
var mc *metrics.Clientset

func init() {
	var kubeConfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
	} else {
		kubeConfig = flag.String("kubeConfig", "", "absolute path to the kubeConfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)

	clientSetInit, err := kubernetes.NewForConfig(config)
	mc, err = metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if nil != clientSetInit {
		clientSet = clientSetInit
	} else {
		panic("k8s connect failed")
	}
	if err != nil {
		log.Printf("k8s connect failed %v \n", err.Error())
	} else {
		log.Printf("connect k8s success \n")

	}
}

func InitK8s() *kubernetes.Clientset {
	return clientSet
}

//DescribePod doc
//@Description: 获取pod的详细信息
//@Author niejian
//@Date 2021-05-08 11:33:27
//@param podName
//@param ns
//@return *v1.Pod
//@return error
func DescribePod(podName, ns string) (*v1.Pod, error) {
	return clientSet.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
}

//func DescribeDeploy(deployName, ns string, labels map[string]string) string {
//	deploy, _ := clientSet.AppsV1().Deployments(ns).Get(context.TODO(), deployName, metav1.GetOptions{})
//	// deploy文件中指定的label信息
//	podLabels := deploy.Spec.Template.Labels
//	for key, val := range labels {
//		// podLabels必须全部包含
//		data, isExist := podLabels[key]
//		if  !isExist {
//			break
//			return ""
//		}
//		if val != data {
//			break
//			return ""
//		}
//	}
//	return de
//}

func ListNodes()  {
	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("获取k8s节点失败%v \n", err.Error())
	}
	for _, node := range nodes.Items {
		addresses := node.Status.Addresses
		//fmt.Printf("node addresses: %v \n", addresses)

		//isWorkerNode := false
		var addressType v1.NodeAddressType
		var addressName string

		for _, address := range addresses {
			addressType = address.Type
			addressName = address.Address

			if v1.NodeInternalIP == addressType {
				fmt.Printf("node ip: %v \n", addressName)
			}
		}
	}
}

func GetUserNs()  (map[string][]string, error) {

	var data  = make(map[string][]string, 2)

	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("获取命名空间失败")
		return nil, err
	}

	items := namespaces.Items
	for _, ns := range items {
		namespace := ns.Name
		ws, exist := ns.Labels["kubesphere.io/workspace"]
		if namespace == "system-workspace" {
			continue
		}

		if ws == "system-workspace" || ws == "cloud-platform"{
			continue
		}

		if !exist {
			continue
		}
		_, ok := data[ws]
		// 不存在
		if !ok {
			var namespaces []string
			//data = map[string][]string{}
			data[ws] = namespaces
		}
		data[ws] = append(data[ws], namespace);

		//return data, err
	}
	return data, err
}

func GetUserDeploys(ns string) ([]string, error) {
	var deploys []string
	deploymentList, err := clientSet.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("获取deployment资源失败%v \n", err.Error())
		return nil, err
	}
	for _, deploy := range deploymentList.Items {
		deploys = append(deploys, deploy.Name)
	}
	return deploys, nil
}

func GetUserPod(ns string) (*v1.PodList, error)  {
	podList, err := clientSet.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	return podList, err
}

func GetPodDetail(podName, ns string) vo.DumpVo {
	var dumpvo vo.DumpVo
	pod, _ := clientSet.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
	for _, con := range pod.Spec.Containers {
		name := con.Name
		if !strings.Contains(podName, name) {
			continue
		}
		memory := con.Resources.Limits.Memory()
		cpuQ := con.Resources.Limits.Cpu()
		if nil == memory || nil == cpuQ {
			continue
		}
		mem := memory.String()
		cpu := cpuQ.String()
		// request
		requestCpuQ := con.Resources.Requests.Cpu()
		requestMemQ := con.Resources.Requests.Memory()
		requestCpu := requestCpuQ.String()
		requestMem := requestMemQ.String()
		mem = strings.ReplaceAll(mem, "Gi", "")
		requestMem = strings.ReplaceAll(requestMem, "Gi", "")
		dumpvo.LimitCpu = cpu
		dumpvo.LimitMem = mem
		dumpvo.RequestMem = requestMem
		dumpvo.RequestCpu = requestCpu
	}

	podMetrics, err := mc.MetricsV1beta1().PodMetricses(ns).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error:", err)
		return dumpvo
	}
	for _, container := range podMetrics.Containers {
		containerName := container.Name
		if !strings.Contains(podName, containerName) {
			continue
		}

		//dumpvo.UsageCpu = container.Usage.Cpu().String()
		//dumpvo.UsageMem = container.Usage.Memory().String()
		usageCpu := container.Usage.Cpu().AsApproximateFloat64()
		cpu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", usageCpu), 64)

		//usageMem := container.Usage.Memory().AsApproximateFloat64()
		//fmt.Printf("===> %v, %v \n",  usageMem, container.Usage.Memory().String())
		usageMemStr := container.Usage.Memory().String()
		var data float64
		if strings.Contains(usageMemStr, "Ki") {
			split := strings.Split(usageMemStr, "Ki")
			data, _ = strconv.ParseFloat(split[0], 32)
			//fmt.Printf("data1==>%v\n", data)
			data = data / (1024 * 1024)
			data, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", data), 64)
			//fmt.Printf("data2==>%v\n", data)


		}else {
			split := strings.Split(usageMemStr, "Gi")
			data, _ = strconv.ParseFloat(split[0], 32)
		}
		//fmt.Println(container.Usage.Memory().String())
		//fmt.Println(container.Usage.Memory().Neg)
		//usageMem = usageMem/1024/1024/1024
		//usageMemStr := strconv.FormatFloat(usageMem,'E',-1,64)
		//fmt.Println(usageMemStr[0:4])

		dumpvo.UsageCpu = fmt.Sprintf("%.2f",cpu)
		dumpvo.UsageMem = fmt.Sprintf("%.2f",data)

		//fmt.Println(container.Usage.Cpu().AsApproximateFloat64())
	}

	return dumpvo
}
