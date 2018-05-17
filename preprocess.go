package main

import (
	"fmt"
	"os"
	"path/filepath"
	"path"
	"archive/zip"
	"io"
	"github.com/natsukagami/go-osu-parser"
	"encoding/json"
	"bytes"
)

func main(){
	fmt.Println("starting")
	wdpathname, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Mkdir(path.Join(wdpathname, "data/"), 0777)
	pathname := path.Join(wdpathname, "osz/")
	fmt.Println("pathname:" + pathname)
	var files []string
	filepath.Walk(pathname, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".osz" {
			files = append(files, info.Name())
			fmt.Println("info name:" + info.Name())

			r, err := zip.OpenReader(path.Join(pathname, info.Name()))
			if err != nil {
				fmt.Println(err)
			}
			defer r.Close()

			t := path.Join(pathname, info.Name()[0:len(info.Name()) - len(filepath.Ext(info.Name()))])

			err = os.Mkdir(t, 0777)
			if err != nil {
				fmt.Println(err)
			}
			for _, f := range r.File {
				fmt.Println(info.Name() + " contains:" + f.Name)
				rc, err := f.Open()
				if err != nil {
					fmt.Println(err)
				}
				tpath := filepath.Join(pathname, f.Name)
				fmt.Println("tpath:" + tpath)


				tf, err := os.OpenFile(path.Join(t, f.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					fmt.Println(err)
				}
				_, err = io.Copy(tf, rc)
				if err != nil {
					fmt.Println(err)
				}

				if filepath.Ext(f.Name) == ".osu" {
					fmt.Println("converting:" + f.Name)
					b, err := parser.ParseFile(path.Join(t, f.Name))
					tbytes, err := json.MarshalIndent(b, "", "\t")
					if err != nil {
						fmt.Println("converting to json error:" + err.Error())
					}
					fdata, err := os.OpenFile(filepath.Join(wdpathname, "data", f.Name[0:len(f.Name) - len(filepath.Ext(f.Name))] + ".json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
					if err != nil {
						fmt.Println(err)
					}
					io.Copy(fdata, bytes.NewReader(tbytes))
				}

			}

		}
		return nil
	})
}
