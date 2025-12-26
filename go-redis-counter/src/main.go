// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

// context: 요청이 취소되거나 타임아웃 났을 때 신호 주는 관리자
var ctx = context.Background()

func main() {
	// Docker Compose에서 REDIS_ADDR=redis:6379 라고 넣어줌.
	// 하드코딩 안 하려고 os.Getenv 씀.
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // 로컬 테스트용 안전장치
	}

	// Redis 연결 객체 생성 (이때 실제 연결되는 건 아니고 설정만 함)
	// &redis.Options -> 구조체 포인터 넘기는 거 C랑 똑같음.
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 메모리 변수(count++) 쓰면 서버 재시작할 때 날아가니까 Redis 씀.
		// Incr 함수는 Atomic 해서 동시 접속 꼬임 문제 해결해줌.
		count, err := rdb.Incr(ctx, "visitors").Result()

		// Go는 try-catch 없고 if err != nil 로 잡아야 함.
		if err != nil {
			fmt.Fprintf(w, "Redis 죽음: %v", err)
			return
		}

		fmt.Fprintf(w, "방문자 수: %d", count)
	})

	// 리셋 기능 추가
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		// .Err() 를 꼭 붙여줘야 합니다!
		// 마지막 인자 0은 "만료 시간 없음(영구 저장)"이라는 뜻입니다.
		err := rdb.Set(ctx, "visitors", 0, 0).Err()

		if err != nil {
			fmt.Fprintf(w, "초기화 실패: %v", err)
			return
		}
		fmt.Fprintf(w, "방문자 수가 0으로 초기화되었습니다.")
		fmt.Println("관리자가 방문자 수를 리셋했습니다.") // 서버 로그에도 남기기
	})

	// 서버 시작 (Listen & Serve)
	// C 소켓 프로그래밍의 bind + listen + accept 무한루프를 한 방에 해줌.
	http.ListenAndServe(":8080", nil)
}

// CI/CD 테스트용 주석입니다. (삭제 예정)
