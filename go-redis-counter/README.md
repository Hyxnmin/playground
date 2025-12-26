# Go Redis visitor Counter

Go와 Redis를 Docker compose로 엮어서 만든 방문자 수 카운터
C언어/커널 베이스에서 도커 이미지 최적화(Multi-stage build) 실습용

## Key Features
- **Multi-stage Build**: Docker 이미지를 1.2GB -> 15MB로 경량화
- **Container Orchestration**: Docker Compose를 이용한 App-DB 연결
- **Optimized Binary**: CGO 비활성화 및 정적 링크 빌드

## How to Run
```bash
# 실행
docker compose up --build

# 접속
http://localhost:8080