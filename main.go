package main

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// In-cluster config 로드
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Error loading in-cluster config:", err)
		panic(err)
	}
	fmt.Println("Successfully loaded in-cluster config")

	// 클라이언트 생성
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		panic(err)
	}
	fmt.Println("Successfully created Kubernetes client")

	// Pod 모니터링을 위한 watcher 설정
	watcher, err := clientset.CoreV1().Pods("default").Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error setting up watcher:", err)
		panic(err)
	}
	fmt.Println("Watcher successfully set up for Pods")

	// 환경 변수에서 SECRET_KEY 읽기
	secretKey := os.Getenv("SECRET_KEY")
	fmt.Println("Secret Key:", secretKey)

	timeZone := os.Getenv("TZ") // 환경 변수에서 SECRET_KEY 읽기
	fmt.Println("timeZone:", timeZone)
	// 이벤트 수신 및 처리
	for event := range watcher.ResultChan() {
		fmt.Printf("Event Type: %v\n", event.Type)
		// 이벤트 처리 로직 구현
	}
}
