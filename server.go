package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"net/http"
	"strings"
)

var (
	db            gorm.DB
	sqlConnection string
)

func main() {
	var err error
	sqlConnection = "root@tcp(127.0.0.1:3306)/rbg_assistant?parseTime=True"
	db, err = gorm.Open("mysql", sqlConnection)

	db.AutoMigrate(&Character{}, &Score{}, &Battle{})

	if err != nil {
		panic(err)
		return
	}

	m := martini.Classic()

	m.Use(martini.Static("assets"))
	m.Use(render.Renderer(render.Options{Layout: "layout"}))

	// Home
	m.Get("/", func(r render.Render) {
		var retData struct {
			Characters []Character
		}

		db.Find(&retData.Characters)

		r.HTML(200, "index", retData)
	})

	// Upload
	m.Get("/upload", func(r render.Render) {
		var retData struct{}
		r.HTML(200, "upload", retData)
	})

	m.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("txtUpload")

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()

			//fmt.Println(str)

			beg := strings.Index(str, "\"")
			end := strings.LastIndex(str, "\"")
			if beg != -1 && end != -1 {
				jsonString := strings.Replace(str[beg+1:end-1], "\\\"", "\"", -1)
				fmt.Println(jsonString)
			}
		}

		if err = scanner.Err(); err != nil {
			fmt.Fprintln(w, err)
			return
		}

	})

	m.Run()
}
