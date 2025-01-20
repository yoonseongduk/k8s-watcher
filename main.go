package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	secretKey := os.Getenv("SECRET_KEY")
	fmt.Println("Secret Key:", secretKey)

	timeZone := os.Getenv("TZ") // 환경 변수에서 SECRET_KEY 읽기
	fmt.Println("timeZone:", timeZone)

	// In-cluster config 로드
	config, err := rest.InClusterConfig()
	if err != nil {
		logError("Error loading in-cluster config", err)
		panic(err)
	}
	logInfo("Successfully loaded in-cluster config")

	// 클라이언트 생성
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logError("Error creating Kubernetes client", err)
		panic(err)
	}
	logInfo("Successfully created Kubernetes client")

	// Pod 모니터링을 위한 watcher 설정
	watcher, err := clientset.CoreV1().Pods("default").Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		logError("Error setting up watcher", err)
		panic(err)
	}
	logInfo("Watcher successfully set up for Pods")

	// 이전 상태를 저장할 맵
	podStates := make(map[string]*v1.Pod)

	// 이벤트 수신 및 처리
	for event := range watcher.ResultChan() {
		// 현재 시간 출력
		currentTime := time.Now().Format(time.RFC3339)
		fmt.Printf("Current Time: %s\n", currentTime)

		// 이벤트 타입 출력
		fmt.Printf("Event Type: %v\n", event.Type)

		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			logError("Error asserting event object type", nil)
			continue
		}

		// ADDED 이벤트 처리
		if event.Type == "ADDED" {
			podStates[pod.Name] = pod // 새로운 Pod 상태 저장
			eventDetails, err := json.MarshalIndent(pod, "", "  ")
			if err != nil {
				logError("Error marshaling event details", err)
			} else {
				fmt.Printf("Event Details: %s\n", eventDetails)
			}
			logInfo(fmt.Sprintf("Pod added: %s", pod.Name))
		}

		// MODIFIED 이벤트 처리
		if event.Type == "MODIFIED" {
			previousPod, exists := podStates[pod.Name]
			if exists {
				// 변경된 부분만 출력
				changedFields := comparePods(previousPod, pod)
				if len(changedFields) > 0 {
					fmt.Printf("Modified Pod: %s\n", pod.Name)
					for field, values := range changedFields {
						fmt.Printf("  %s: Previous Value: %v, New Value: %v\n", field, values[0], values[1])
					}
					// 전체 Pod 정보 출력
					podDetails, err := json.MarshalIndent(pod, "", "  ")
					if err != nil {
						logError("Error marshaling pod details", err)
					} else {
						fmt.Printf("  Full Pod Details: %s\n", podDetails)
					}
					logInfo(fmt.Sprintf("Pod modified: %s", pod.Name))
				}
			}
			podStates[pod.Name] = pod // 현재 Pod 상태 저장
		}
	}
}

// comparePods는 두 Pod 객체를 비교하여 변경된 필드를 반환합니다.
func comparePods(oldPod, newPod *v1.Pod) map[string][2]interface{} {
	changedFields := make(map[string][2]interface{})

	// 필드 비교
	if oldPod.Status.Phase != newPod.Status.Phase {
		changedFields["Phase"] = [2]interface{}{oldPod.Status.Phase, newPod.Status.Phase}
	}

	if oldPod.Spec.Containers[0].Image != newPod.Spec.Containers[0].Image {
		changedFields["Image"] = [2]interface{}{oldPod.Spec.Containers[0].Image, newPod.Spec.Containers[0].Image}
	}

	if oldPod.Spec.RestartPolicy != newPod.Spec.RestartPolicy {
		changedFields["RestartPolicy"] = [2]interface{}{oldPod.Spec.RestartPolicy, newPod.Spec.RestartPolicy}
	}

	if oldPod.Status.HostIP != newPod.Status.HostIP {
		changedFields["HostIP"] = [2]interface{}{oldPod.Status.HostIP, newPod.Status.HostIP}
	}

	if oldPod.Status.PodIP != newPod.Status.PodIP {
		changedFields["PodIP"] = [2]interface{}{oldPod.Status.PodIP, newPod.Status.PodIP}
	}

	// 추가적인 필드 비교
	if oldPod.Spec.TerminationGracePeriodSeconds != newPod.Spec.TerminationGracePeriodSeconds {
		changedFields["TerminationGracePeriodSeconds"] = [2]interface{}{oldPod.Spec.TerminationGracePeriodSeconds, newPod.Spec.TerminationGracePeriodSeconds}
	}

	if oldPod.Spec.ActiveDeadlineSeconds != newPod.Spec.ActiveDeadlineSeconds {
		changedFields["ActiveDeadlineSeconds"] = [2]interface{}{oldPod.Spec.ActiveDeadlineSeconds, newPod.Spec.ActiveDeadlineSeconds}
	}

	if oldPod.Spec.NodeName != newPod.Spec.NodeName {
		changedFields["NodeName"] = [2]interface{}{oldPod.Spec.NodeName, newPod.Spec.NodeName}
	}

	if oldPod.Status.Reason != newPod.Status.Reason {
		changedFields["Reason"] = [2]interface{}{oldPod.Status.Reason, newPod.Status.Reason}
	}

	if oldPod.Status.Message != newPod.Status.Message {
		changedFields["Message"] = [2]interface{}{oldPod.Status.Message, newPod.Status.Message}
	}

	// 추가적인 필드 비교를 여기에 추가할 수 있습니다.

	return changedFields
}

// logInfo는 정보 메시지를 출력합니다.
func logInfo(message string) {
	fmt.Printf("[INFO] %s: %s\n", time.Now().Format(time.RFC3339), message)
}

// logError는 오류 메시지를 출력합니다.
func logError(message string, err error) {
	if err != nil {
		fmt.Printf("[ERROR] %s: %s - %v\n", time.Now().Format(time.RFC3339), message, err)
	} else {
		fmt.Printf("[ERROR] %s: %s\n", time.Now().Format(time.RFC3339), message)
	}
}
