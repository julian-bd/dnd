package ripper

type indexNameUrl struct {
	Index string
	Name  string
	Url   string
}
type resultsResponse struct {
	Count   int
	Results []indexNameUrl
}
type abilityBonus struct {
	AbilityScore indexNameUrl `json:"ability_score"`
	Bonus        int
}
type option struct {
	Type string `json:"option_type"`
	Item indexNameUrl
}
type optionSet struct {
	Type    string `json:"option_set_type"`
	Options []option
}
type startingProficiencyOptions struct {
	Choose int
	From   optionSet
}
type raceResponse struct {
	Name                       string
	Speed                      int
	Size                       string
	AbilityBonuses             []abilityBonus             `json:"ability_bonuses"`
	StartingProficiencies      []indexNameUrl             `json:"starting_proficiencies"`
	StartingProficiencyOptions startingProficiencyOptions `json:"starting_proficiency_options"`
	Languages                  []indexNameUrl
	Traits                     []indexNameUrl
	SubRaces                   []indexNameUrl `json:"subraces"`
}
