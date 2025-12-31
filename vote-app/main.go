package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/vote":
			menu := r.URL.Query().Get("menu")

			if menu == "" {
				fmt.Fprintf(w, "투표할 메뉴를 입력해주세요! (?memu=xxx)")
				return
			}
			if menu == "jjajang" || menu == "jjamppong" {
				newScore, _ := rdb.Incr(ctx, menu).Result()

				fmt.Fprintf(w, "[투표완료] %s의 현재 점수: %d표 🔥", menu, newScore)
			} else {
				http.Error(w, "그럼 메뉴는 없습니다!", http.StatusBadRequest)
			}
		case "/result":
			score1, _ := rdb.Get(ctx, "jjajang").Result()
			score2, _ := rdb.Get(ctx, "jjamppong").Result()

			if score1 == "" {
				score1 = "0"
			}
			if score2 == "" {
				score2 = "0"
			}

			fmt.Fprintf(w, "=== 🏆 현재 스코어 ===\n")
			fmt.Fprintf(w, "짜장면: %s표\n", score1)
			fmt.Fprintf(w, "짬뽕: %s표\n", score2)
		case "/reset":
			rdb.Del(ctx, "jjajang", "jjamppong")
			fmt.Fprintf(w, "투표함을 비웠습니다! 🗑️")
		default:
			fmt.Fprintf(w, "잘못된 주소입니다. (/vote, /result, /reset)")
		}
	})

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
