package main

import "github.com/julian-bd/dnd/data"

type headerContent struct {
	PlayableRaceNames []string
}

type indexContent struct {
	HeaderContent headerContent
}

type playableRaceContent struct {
	HeaderContent  headerContent
	AbilityBonuses []data.AbilityBonus
	Languages      []string
	Name           string
	Proficiencies  []string
	SubRaces       []string
	Traits         []string
	Speed          int
}

type createPlayableRaceContent struct {
	HeaderContent headerContent
	Abilities     []string
	Languages     []string
	Proficiencies []string
	SubRaces      []string
	Traits        []string
}
