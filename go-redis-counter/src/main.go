// main.go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count, err := rdb.Incr(ctx, "count").Result()
		if err != nil {
			fmt.Fprintf(w, "DB 연결 실패: %v", err)
			return
		}

		fmt.Fprintf(w, "당신은 %d번째 방문자입니다!", count)
		fmt.Printf("방문자 증가: %d\n", count)
	})

	fmt.Println("서버 시작 (8080)")
	http.ListenAndServe(":8080", nil)
}
