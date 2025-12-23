
# Go 컴파일러가 포한된 무거운 이미지 사용
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# CGO_ENABLED=0: C 라이브러리(glibc) 의존성을 제거하고 정적 링크 수행
# Scratch나 Alpine 같은 텅 빈 리눅스에서도 실행되게 하기 위함
# -ldflags "-s -w": 디버그 정보를 제거하여 바이너리 용량 최소화
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main src/main.go


# 1GB가 넘는 빌드 이미지 버리고, 5MB짜리 Alpine 리눅스로 갈아탐 (Multi-stage Build)
FROM alpine:latest
WORKDIR /root/

# 빌드 단계에서 생성한 바이너리 파일만 복사
COPY --from=builder /app/main .

CMD ["./main"]