package data

type PlayableRace struct {
	ID                           int                            `json:"id"`
	Name                         string                         `json:"name"`
	Speed                        int                            `json:"speed"`
	Ability_Bonuses              []ability_bonus                `json:"ability_bonuses"`
	Starting_Languages           []string                       `json:"starting_languages"`
	Starting_Proficiencies       []string                       `json:"starting_proficiencies"`
	Starting_Proficiency_Options []starting_proficiency_options `json:"starting_proficiency_options"`
	Traits                       []string                       `json:"traits"`
	Sub_Races                    []string                       `json:"sub_races"`
}

type ability_bonus struct {
	Ability string `json:"ability"`
	Bonus   int    `json:"bonus"`
}

type trait struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type starting_proficiency_options struct {
	Count   int      `json:"count"`
	Options []string `json:"options"`
}
