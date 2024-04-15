package main

type createPlayableRace struct {
	Name           string `form:"name"`
	Speed          int    `form:"speed"`
	AbilityBonuses string `form:"abilities"`
	Languages      string `form:"languages"`
	Proficiencies  string `form:"proficiencies"`
	Traits         string `form:"traits"`
	SubRaces       string `form:"subraces"`
}
