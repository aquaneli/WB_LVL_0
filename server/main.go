package main
import (
    "fmt"
    "net/http"
)



func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        http.ServeFile(w, r, "home.html")
    })

    http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request){
     
        name := r.FormValue("username")
        age := r.FormValue("userage")
         
        fmt.Fprintf(w, "Имя: %s Возраст: %s", name, age)
    })

    fmt.Println("Server is listening...")
    http.ListenAndServe("localhost:8181", nil)
}