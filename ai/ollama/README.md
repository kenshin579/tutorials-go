# Ollama on Kubernetes with Terraform

이 프로젝트는 Terraform을 사용하여 Kind 클러스터에 Ollama를 배포하는 예제입니다.

## 빠른 시작

```bash
# 1. Terraform 초기화
make init

# 2. 클러스터 및 Ollama 배포
make apply

# 3. 배포 상태 확인
make status

# 4. API 테스트
make test

# 5. 모델 다운로드
make model-pull

# 6. 채팅 테스트
make chat
```

## 파일 구조

- `k8s.tf` - Kind 클러스터 구성
- `infra.tf` - Ollama Helm Chart 배포
- `Makefile` - 편의 명령어 모음

## 주요 명령어

```bash
make help       # 사용 가능한 명령어 보기
make logs       # Ollama 로그 확인
make destroy    # 모든 리소스 삭제
```

## 설정 변경

### 포트 변경
`k8s.tf`와 `infra.tf`에서 30025 포트를 원하는 포트로 변경

### 리소스 조정
`infra.tf`의 resources 섹션에서 CPU/메모리 조정

### 스토리지 크기
`infra.tf`의 persistentVolume.size 조정

## 문제 해결

### Pod가 시작되지 않을 때
```bash
kubectl describe pod -n ollama
kubectl get events -n ollama
```

### 메모리 부족 시
`infra.tf`의 메모리 요청/제한 값을 줄이기

## 참고 링크

- [상세 가이드](../content.md)
- [Ollama 공식 문서](https://ollama.ai/)
- [Terraform 문서](https://www.terraform.io/) 