package main

import (
	//"fmt"
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
			Leaders []struct {
				Leader string
				Count  uint
			}
		}

		db.Table("battles").Select("leader_id as leader, count(*) as count").Scan(&retData.Leaders)

		out.HTML(200, "index", retData)
	})

	// Upload
	m.Get("/upload", func(r render.Render) {
		var retData struct{}
		r.HTML(200, "upload", retData)
	})

	m.Post("/upload", upload_battle)

	m.Get("/leaders/:leader_id/battles", func(params martini.Params, r render.Render) {
		var retData struct {
			Battles []Battle
			Leader  Character
		}

		retData.Leader = Character{ID: params["leader_id"]}

		db.Find(&retData.Leader)

		db.Where("leader_id = ?", params["leader_id"]).Find(&retData.Battles)

		r.HTML(200, "battles", retData)
	})

	m.Run()
}
