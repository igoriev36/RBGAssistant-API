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
	"time"
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
				jsonString := strings.Replace(str[beg+1:end], "\\\"", "\"", -1)
				saveBattle(jsonString)
			}
		}

		if err = scanner.Err(); err != nil {
			fmt.Fprintln(w, err)
			return
		}

	})

	m.Run()
}

func saveBattle(s string) {
	fmt.Println(s)

	type BG struct {
		Time    string `json:"time"`
		Map     string `json:"map"`
		Leader  string `json:"leader"`
		Winner  string `json:"winnger"`
		Player  string `json:"player"`
		IsRated bool   `json:"is_rated"`
		Scores  []struct {
			Name           string `json:"name"`
			Kb             int    `json:"kb"`
			Hk             int    `json:"hk"`
			Deaths         int    `json:"deaths"`
			Honor          int    `json:"honor"`
			Faction        string `json:"faction"`
			Race           string `json:"race"`
			Class          string `json:"class"`
			Damage         int    `json:"damage"`
			Healing        int    `json:"healing"`
			BgRating       int    `json:"bg_rating"`
			BgRatingChange int    `json:"bg_rating_change"`
			PreMmr         int    `json:"pre_mmr"`
			MmrChange      int    `json:"mmr_change"`
			TalentSpec     string `json:"talent_spec"`
		}
	}

	bg := &BG{}
	err := json.Unmarshal([]byte(s), bg)
	if err != nil {
		fmt.Println(err)
	}

	scores := []Score{}

	for _, score := range bg.Scores {
		var name, realm string

		if strings.Index(score.Name, "-") > 0 {
			s := strings.Split(score.Name, "-")
			name = s[0]
			realm = s[1]
		}

		//Create/Update Character
		c := Character{}
		db.FirstOrCreate(&c, Character{ID: score.Name})

		c.Name = name
		c.Realm = realm
		c.Faction = score.Faction
		c.Race = score.Race
		c.Class = score.Class

		db.Save(&c)

		score := Score{
			Character:      c,
			KillingBlows:   score.Kb,
			HonorableKills: score.Hk,
			Deaths:         score.Deaths,
			HonorGained:    score.Honor,
			Damage:         score.Damage,
			Healing:        score.Healing,
			BgRating:       score.BgRating,
			BgRatingChange: score.BgRatingChange,
			PrematchMmr:    score.PreMmr,
			MmrChange:      score.MmrChange,
			TalentSpec:     score.TalentSpec,
		}
	}

	leader := Character{}
	db.FirstOrCreate(&leader, Character{ID: bg.Leader})

	player := Character{}
	db.FirstOrCreate(&player, Character{ID: bg.Leader})

	pt, _ := time.Parse("2006-01-02 15:04", bg.Time)

	b := Battle{
		PlayedAt:   pt,
		Map:        bg.Map,
		Winner:     bg.Winner,
		Leader:     leader,
		RecordedBy: player,
		IsRated:    bg.IsRated,
	}

	fmt.Println(bg)

}
