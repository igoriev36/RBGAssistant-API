package main

import (
	"time"
)

type Character struct {
	ID         string `gorm:"primary_key"`
	Name       string
	Realm      string
	Guild      string
	Faction    string
	Race       string
	Class      string
	BattlesLed []Battle `gorm:"has_many:Leader"`
}

type Battle struct {
	ID           uint `gorm:"primary_key"`
	PlayedAt     time.Time
	Map          string
	Winner       string
	LeaderID     string
	RecordedByID string
	Leader       Character `sql:foreign_key("leader_id")`
	RecordedBy   Character `sql:foreign_key("recorded_by_id")`
	IsRated      bool
	Scores       []Score
}

type Score struct {
	ID             uint `gorm:"primary_key"`
	BattleID       uint
	Battle         Battle
	CharacterID    string
	Character      Character
	KillingBlows   int
	HonorableKills int
	Deaths         int
	HonorGained    int
	Damage         int
	Healing        int
	BgRating       int
	BgRatingChange int
	PrematchMmr    int
	MmrChange      int
	TalentSpec     string
}
