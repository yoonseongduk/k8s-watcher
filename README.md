# k8s-watcher Deployment on Kubernetes

이 문서는 `main.go` 애플리케이션을 Kubernetes 클러스터에 배포하는 절차를 설명합니다.

## 요구 사항

- Kubernetes 클러스터
- `kubectl` CLI 도구
- Docker 이미지가 레지스트리에 푸시되어 있어야 함

## 배포 절차

1. **Docker 이미지 빌드 및 푸시**

   먼저, 애플리케이션의 Docker 이미지를 빌드하고 Docker 레지스트리에 푸시합니다.

   ```bash
   docker build -t yoonseongduk/k8s-watcher:latest .
   docker push yoonseongduk/k8s-watcher:latest
   ```

2. **YAML 파일 생성**

   다음 YAML 파일들을 생성합니다:

   - `deploy.yaml`: Deployment 설정
   - `service.yaml`: Service 설정
   - `configmap.yaml`: ConfigMap 설정 (필요시)
   - `secret.yaml`: Secret 설정 (필요시)

   예시 파일은 다음과 같습니다:

   - **deploy.yaml**
     ```yaml
     apiVersion: apps/v1
     kind: Deployment
     metadata:
       name: k8s-watcher
       labels:
         app: k8s-watcher
     spec:
       replicas: 1
       selector:
         matchLabels:
           app: k8s-watcher
       template:
         metadata:
           labels:
             app: k8s-watcher
         spec:
           containers:
           - name: k8s-watcher-container
             image: yoonseongduk/k8s-watcher:latest  # Docker 이미지 이름 수정
             ports:
             - containerPort: 8080  # 애플리케이션이 사용하는 포트
             env:
             - name: ENV_VAR_NAME  # 환경 변수 설정 (필요시)
               value: "value"
     ```

   - **service.yaml**
     ```yaml
     apiVersion: v1
     kind: Service
     metadata:
       name: k8s-watcher-service
     spec:
       type: ClusterIP  # 또는 LoadBalancer, NodePort 등
       ports:
       - port: 80  # 서비스가 노출할 포트
         targetPort: 8080  # 컨테이너의 포트
       selector:
         app: k8s-watcher
     ```

   - **configmap.yaml** (필요시)
     ```yaml
     apiVersion: v1
     kind: ConfigMap
     metadata:
       name: k8s-watcher-config
     data:
       CONFIG_KEY: "config_value"  # 애플리케이션 설정
     ```

   - **secret.yaml** (필요시)
     ```yaml
     apiVersion: v1
     kind: Secret
     metadata:
       name: k8s-watcher-secret
     type: Opaque
     data:
       SECRET_KEY: c2VjcmV0X3ZhbHVl  # base64 인코딩된 비밀 값
     ```

3. **Kubernetes 클러스터에 배포**

   생성한 YAML 파일들을 사용하여 Kubernetes 클러스터에 배포합니다.

   ```bash
   kubectl apply -f deploy.yaml
   kubectl apply -f service.yaml
   kubectl apply -f configmap.yaml  # 필요시
   kubectl apply -f secret.yaml      # 필요시
   ```

4. **배포 확인**

   배포가 완료되면, 다음 명령어로 Pod와 Service의 상태를 확인합니다.

   ```bash
   kubectl get pods
   kubectl get services
   ```

5. **애플리케이션 접근**

   Service의 타입에 따라 애플리케이션에 접근하는 방법이 다릅니다. `ClusterIP` 타입의 경우 클러스터 내에서만 접근 가능하며, `LoadBalancer` 또는 `NodePort` 타입의 경우 외부에서 접근할 수 있습니다.

## 추가 정보

- 애플리케이션의 로그를 확인하려면 다음 명령어를 사용합니다:

  ```bash
  kubectl logs <pod-name>
  ```

- Pod의 상세 정보를 보려면 다음 명령어를 사용합니다:

  ```bash
  kubectl describe pod <pod-name>
  ```

이 문서를 통해 `main.go` 애플리케이션을 Kubernetes 클러스터에 성공적으로 배포할 수 있습니다.
