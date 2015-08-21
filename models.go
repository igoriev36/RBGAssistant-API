package main

import (
	"time"
)

type Character struct {
	ID      string `gorm:"primary_key"`
	Name    string
	Realm   string
	Guild   string
	Faction string
	Race    string
	Class   string
}

type Battle struct {
	ID         uint `gorm:"primary_key"`
	PlayedAt   time.Time
	Map        string
	Winner     string
	Leader     Character
	RecordedBy Character
	IsRated    bool
	Scores     []Score
}

type Score struct {
	ID             uint `gorm:"primary_key"`
	Battle         Battle
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
