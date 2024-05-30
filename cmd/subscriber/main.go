package main

import (
	"fmt"
	"html/template"
	"module_0/internal/databases"
	"module_0/internal/models"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/maxchagin/go-memorycache-example"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	cache := memorycache.New(1*time.Hour, 1*time.Hour)
	databases.LoadInCache(cache)
	go OrderServer(cache)
	go databases.NatsSub(cache, &wg)
	wg.Wait()
}

func OrderServer(cache *memorycache.Cache) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("web/home.html")
		id := r.FormValue("Id")
		value, _ := cache.Get(id)
		val, _ := value.(models.Orders)
		tmpl.Execute(w, val)
	})
	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:8181", nil)
}
