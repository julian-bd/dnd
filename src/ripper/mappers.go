package ripper

import "github.com/julian-bd/dnd/data"

func toDomainRace(raw raceResponse) data.PlayableRace {
	return data.PlayableRace{
		Name:                       raw.Name,
		Speed:                      raw.Speed,
		AbilityBonuses:             mapAbilityBonuses(raw.AbilityBonuses),
		StartingLanguages:          mapNames(raw.Languages),
		StartingProficiencies:      mapNames(raw.StartingProficiencies),
		StartingProficiencyOptions: mapProficiencyOptions(raw.StartingProficiencyOptions),
		Traits:                     mapNames(raw.Traits),
		SubRaces:                   mapNames(raw.SubRaces),
	}
}

func mapProficiencyOptions(raw startingProficiencyOptions) []data.StartingProficiencyOption {
	var os []string
	for _, o := range raw.From.Options {
		os = append(os, o.Item.Name)
	}
	// TODO: This does not need to be an array
	return []data.StartingProficiencyOption{
		{
			Count:   raw.Choose,
			Options: os,
		},
	}
}

func mapAbilityBonuses(raw []abilityBonus) []data.AbilityBonus {
	var abs []data.AbilityBonus
	for _, r := range raw {
		ab := data.AbilityBonus{
			Ability: r.AbilityScore.Name,
			Bonus:   r.Bonus,
		}
		abs = append(abs, ab)
	}
	return abs
}

func mapNames(raw []indexNameUrl) []string {
	var ls []string
	for _, r := range raw {
		ls = append(ls, r.Name)
	}
	return ls
}
