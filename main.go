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
	"strconv"
	"sort"
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

type ReturnRank struct {
	BestScores []SingleReturnRank
}

type SingleReturnRank struct {
	Player string
	Rank string
	Score int
}

func getRank(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Println("getrank parse error:" + err.Error())
	}
	fmt.Println("getrank post begins")
	musicname := req.Form["musicname"][0]
	difficulty, _ := strconv.Atoi(req.Form["difficulty"][0])
	fmt.Println("getrank post begins")
	fmt.Printf("musicname:%s\ndifficulty:%d\n", musicname, difficulty)
	fmt.Println("getrank post ends")

	var p []Record
	for _, v := range Rank {
		if v.Musicname == musicname && v.Difficulty == difficulty {
			p = append(p, v)
		}
	}
	cnt := len(p)
	if cnt > 5 {
		cnt = 5
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].Score > p[j].Score
	})
	var returnRank ReturnRank
	for i := 0; i < cnt; i++ {
		var tResult SingleReturnRank
		tResult.Player = p[i].Username
		tResult.Rank = p[i].Rank
		tResult.Score= p[i].Score
		returnRank.BestScores = append(returnRank.BestScores, tResult)
	}
	js, err := json.Marshal(returnRank)
	if err != nil {
		fmt.Println("error converting in get rank:" + err.Error())
	}
	fmt.Println("get rank response begins")
	os.Stdout.Write(js)
	fmt.Println()
	fmt.Println("get rank response ends")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
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

var User map[string] string

type Record struct {
	Username string
	Musicname string
	Difficulty int
	Score int
	Rank string
}

var Rank []Record

func upload(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Println("upload parse error:" + err.Error())
	}
	fmt.Println("upload post begins")
	username := req.Form["username"][0]
	password := req.Form["password"][0]
	musicname := req.Form["musicname"][0]
	difficulty, _ := strconv.Atoi(req.Form["difficulty"][0])
	rank := req.Form["rank"][0]
	score, _ := strconv.Atoi(req.Form["score"][0])
	fmt.Printf("username:%s\npassword:%s\nmusicname:%s\ndifficulty:%d\nrank:%s\nscore:%d\n", username, password, musicname, difficulty, rank, score)
	fmt.Println("upload post ends")

	if User[username] == password {
		sRecord := Record{username, musicname, difficulty, score, rank}
		Rank = append(Rank, sRecord)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	} else {
		fmt.Println("username:" + username + " and password:" + password + "mismatch")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("username and password mismatch"))
	}

}


func signup(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Println("signup parse error:" + err.Error())
	}
	fmt.Println("signup post begins")
	username := req.Form["username"][0]
	password := req.Form["password"][0]
	fmt.Println("signup post begins")
	fmt.Printf("username:%s\npassword:%s\n", username, password)
	fmt.Println("signup post ends")

	if _, ok := User[username]; ok {
		fmt.Println("key " + username + " already exist")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("username already exists"))
		return
	}

	User[username] = password
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func main(){
	fmt.Println("starting")
	User = make(map[string]string)
	http.HandleFunc("/rank", getRank)
	http.HandleFunc("/search", getSearchResult)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/signup", signup)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil{
		fmt.Println("listen error")
	}
}