package cross_compiling

import (
    "net/http"
    "github.com/gorilla/mux"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
}

func YourHandler2(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Shit, this binary is too big!!\n"))
}

func main() {
    r := mux.NewRouter()
    // Routes consist of a path and a handler function.
    r.HandleFunc("/", YourHandler)
    r.HandleFunc("/size", YourHandler2)

    // Bind to a port and pass our router in
    http.ListenAndServe(":8000", r)
}
