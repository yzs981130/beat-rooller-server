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

type SearchResult struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name string
	Musician string
	ImgLink string
	MusicLink string
	LevelData []LevelData
}

type LevelData struct {
	Difficulty int
	MapLink string
}


func generateJSONResult(searchName string) SearchResult {
	wdpathname, _ := os.Getwd()
	pathname := path.Join(wdpathname, "osz")
	fmt.Println("pathname:" + pathname)
	var r SearchResult
	var files []string
	filepath.Walk(pathname, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.Contains(info.Name(), searchName) && info.Name() != "osz" {
			//fmt.Println(info.Name())
			files = append(files, info.Name())
		}
		return nil
	})

	if len(files) == 0 {
		fmt.Println("no matches for " + searchName)
		return r
	}

	var name, musician, imgLink, musicLink, mapLink string
	var difficulty int

	for _, v := range files {
		fmt.Println(v)
		name = strings.Join(strings.Split(v, " ")[1:], " ")
		fmt.Println("name:" + name)
		nPathname := path.Join(pathname, v)
		fmt.Println("new pathname:" + nPathname)

		var singleResult Result
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
				difficulty = js.Get("OverallDifficulty").MustInt()
				mapLink = path.Join("data", strings.Split(info.Name(), ".")[0] + ".csv")
				fmt.Printf("musician:%s\timgLink:%s\tmusicLink:%s\tdifficulty:%d\tmapLink:%s\n", musician, imgLink, musicLink, difficulty, mapLink)
				singleResult.Name = name
				singleResult.Musician = musician
				singleResult.ImgLink = imgLink
				singleResult.MusicLink = musicLink
				var singleResultLevelData LevelData
				singleResultLevelData.Difficulty = difficulty
				singleResultLevelData.MapLink = mapLink
				singleResult.LevelData = append(singleResult.LevelData, singleResultLevelData)
			}
			return nil
		})
		r.Results = append(r.Results, singleResult)
	}
	return r
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
	fmt.Println("searchname:" + string(t))
	js, err := json.Marshal(generateJSONResult(string(t)))
	if err != nil {
		fmt.Println("error converting:" + err.Error())
	}
	fmt.Println("post response begins")
	os.Stdout.Write(js)
	fmt.Println()
	fmt.Println("post response ends")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func main(){
	fmt.Println("starting")
	http.HandleFunc("/rank", getRank)
	http.HandleFunc("/search", getSearchResult)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil{
		fmt.Println("listen error")
	}
}