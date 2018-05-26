package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/natsukagami/go-osu-parser"
	"encoding/json"
	"github.com/bitly/go-simplejson"
)


func generateJSONResult(searchName string) {
	fmt.Println("entering")
	wdpathname, _ := os.Getwd()
	pathname := path.Join(wdpathname, "osz")
	fmt.Println("pathname:" + pathname)
	var files []string
	filepath.Walk(pathname, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.Contains(info.Name(), searchName) && info.Name() != "osz" {
			//fmt.Println(info.Name())
			files = append(files, info.Name())
		}
		return nil
	})

	var name, musician, imgLink, musicLink string

	for _, v := range files {
		fmt.Println(v)
		name = strings.Join(strings.Split(v, " ")[1:], " ")
		fmt.Println("name:" + name)
		nPathname := path.Join(pathname, v)
		fmt.Println("new pathname:" + nPathname)

		filepath.Walk(nPathname, func(p string, info os.FileInfo, err error) error {
			if !info.IsDir() && filepath.Ext(info.Name()) == ".osu" {
				fmt.Println("osu file:" + info.Name())
				b, err := parser.ParseFile(path.Join(nPathname, info.Name()))
				tbytes, err := json.MarshalIndent(b, "", "\t")
				if err != nil {
					fmt.Println("converting to json error:" + err.Error())
				}
				js, err := simplejson.NewJson(tbytes)
				musician = js.Get("Artist").MustString()
				imgLink = path.Join("osz", v, js.Get("bgFilename").MustString())
				musicLink = path.Join("osz", v, js.Get("AudioFilename").MustString())
				fmt.Printf("musician:%s\timgLink:%s\tmusicLink:%s\n", musician, imgLink, musicLink)
			}
			return nil
		})
	}
}


func getRank(w http.ResponseWriter, req *http.Request){
	fmt.Println("get rank begin")
	r := "Ranking"
	fmt.Fprintf(w, r)
	fmt.Println("sending: " + r)
}

func getSearchResult(w http.ResponseWriter, req *http.Request){
	t, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	result := string(t)
	fmt.Println("result:" + result)

}

func main(){
	/*
	fmt.Println("starting")
	http.HandleFunc("/rank", getRank)
	http.HandleFunc("/search", getSearchResult)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil{
		fmt.Println("listen error")
	}
	*/
	generateJSONResult("o")
}