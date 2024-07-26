package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 加载 kubeconfig 文件
	kubeconfig := filepath.Join(
		homeDir(), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 定义 Job
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "echo-job",
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "echo-container",
							Image:   "busybox",
							Command: []string{"echo", "Hello, Kubernetes! 你好$小明"},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	// 创建 Job // corev1.NamespaceDefault
	jobClient := clientset.BatchV1().Jobs("pay")
	result, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Created job %q.\n", result.GetObjectMeta().GetName())
}

// homeDir 返回用户的主目录
func homeDir() string {
	h := os.Getenv("HOME")
	return h
	//f := flag.Lookup("home")
	//h := f.Value.String()
	//if h != "" {
	//	return h
	//}
	//return "/root"
}
