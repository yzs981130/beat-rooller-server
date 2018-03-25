package main

import (
	"fmt"
	"net/http"
)

func getRank(w http.ResponseWriter, req *http.Request){
	fmt.Println("get rank begin")
	fmt.Fprintf(w, "Ranking:")
}

func main(){
	fmt.Println("starting")
	http.HandleFunc("/rank", getRank)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil{
		fmt.Println("listen error")
	}
}