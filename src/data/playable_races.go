package data

type PlayableRace struct {
	ID                         int                         `json:"id"`
	Name                       string                      `json:"name"`
	Speed                      int                         `json:"speed"`
	AbilityBonuses             []abilityBonus              `json:"ability_bonuses"`
	StartingLanguages          []string                    `json:"starting_languages"`
	StartingProficiencies      []string                    `json:"starting_proficiencies"`
	StartingProficiencyOptions []startingProficiencyOption `json:"starting_proficiency_options"`
	Traits                     []string                    `json:"traits"`
	SubRaces                   []string                    `json:"sub_races"`
}

type abilityBonus struct {
	Ability string `json:"ability"`
	Bonus   int    `json:"bonus"`
}

type trait struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type startingProficiencyOption struct {
	Count   int      `json:"count"`
	Options []string `json:"options"`
}
