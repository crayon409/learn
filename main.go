package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"learn/internal"
	"learn/model"
)

type param struct {
	ID     int `json:"id"`
	Amount int `json:"amount"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ww", ww)
	mux.HandleFunc("/ww2", ww2)
	mux.HandleFunc("/rr", rr)

	log.Fatal(http.ListenAndServe(":2000", mux))
}

var bg = context.Background()

func rr(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	// val := r.FormValue("val")
	res, _ := internal.Rdb.Get(bg, key).Result()
	fmt.Printf("%s\n", res)
}

func ww(w http.ResponseWriter, r *http.Request) {
	user_id, _ := strconv.Atoi(r.FormValue("id"))
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	result := addCoin(user_id, amount)
	fmt.Printf("res: %v\n", result)
	if result {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}
func ww2(w http.ResponseWriter, r *http.Request) {
	user_id, _ := strconv.Atoi(r.FormValue("id"))
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	if addCoin2(user_id, amount) {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}

func addCoin(id int, amount int) bool {
	return internal.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			user.ID = uint(id)
			user.Name = fmt.Sprint("user_%d", id)
			if err = tx.Create(&user).Error; err != nil {
				return err
			}
		}
		res := tx.Model(&user).
			Where("id = ? and coin = ?", user.ID, user.Coin).
			Update("coin", user.Coin+amount)
		if res.RowsAffected == 0 {
			return errors.New("no updated")
		}
		return res.Error
	}) == nil
}

func addCoin2(id int, amount int) bool {
	var user model.User
	return internal.DB.Transaction(func(tx *gorm.DB) error {
		tx.Exec("select * from users where id = ? for update", id)
		if err := tx.First(&user, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			user.ID = uint(id)
			user.Name = fmt.Sprintf("user_%d", id)
			if err = tx.Create(&user).Error; err != nil {
				return err
			}
		}
		res := tx.Model(&user).
			Where("id = ?", user.ID).
			Update("coin", user.Coin+amount)
		if res.RowsAffected == 0 {
			return errors.New("no updated")
		}
		return res.Error
	}) == nil
}
