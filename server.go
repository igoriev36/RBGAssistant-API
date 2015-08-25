package main

import (
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
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
	db.Model(&Battle{}).AddUniqueIndex("idx_battle_playedat_leader", "played_at", "leader_id")

	if err != nil {
		panic(err)
		return
	}

	m := martini.Classic()

	m.Use(martini.Static("assets"))
	m.Use(render.Renderer(render.Options{Layout: "layout"}))

	// Home
	m.Get("/", func(out render.Render) {
		var retData struct {
			Characters []Character
		}

		db.Joins("inner join battles on battles.leader_id = characters.id").Where("battles.is_rated = ?", true).Find(&retData.Characters)

		out.HTML(200, "index", retData)
	})

	// Upload
	m.Get("/upload", func(r render.Render) {
		var retData struct{}
		r.HTML(200, "upload", retData)
	})

	m.Post("/upload", upload_battle)

	m.Run()
}
