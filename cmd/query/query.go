package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"learn/internal"
)

var bg = context.Background()

func main() {
	// q_1()
	q_2()
}

func q_1() {
	n := 10
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			q1()
			wg.Done()
		}()
	}
	wg.Wait()
}

func q1() {

	url := "http://localhost:2000/ww"

	payload := strings.NewReader("id=1&amount=1")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	fmt.Println(res.StatusCode)

}

func q_2() {
	n := 10
	var wg sync.WaitGroup
	wg.Add(n)
	internal.Rdb.Set(bg, "key2", 11, -1)

	for i := 0; i < n; i++ {
		go func() {
			luaScript := `
				local key = KEYS[1]
				local change = ARGV[1]
				
				local oldvalue = redis.call("GET", key)
				if not oldvalue then
				  oldvalue = 0
				end
				
				local value = oldvalue - change
				if value < 0 then
					value = 0
					return value
				end
				redis.call("SET", key, value)
				return oldvalue
				`

			// 编译脚本
			decrByScript := redis.NewScript(luaScript)
			ticks, err := decrByScript.Run(bg, internal.Rdb, []string{"key2"}, []any{1}...).Int()
			if err != nil {
				log.Printf("err: %v\n", err)
				wg.Done()
				return
			}

			if ticks <= 0 {
				wg.Done()
				return
			}

			log.Printf("call with: %v res: %v\n", ticks, q2())
			wg.Done()
		}()
	}
	wg.Wait()
}

func q2() int {

	url := "http://localhost:2000/ww2"

	payload := strings.NewReader("id=2&amount=1")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	return res.StatusCode

}
