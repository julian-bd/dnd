package ripper

var baseUrl = "https://www.dnd5eapi.co"

func Seed() {
	err := seedTraits()
	if err != nil {
		return
	}
	err = seedProficiencies()
	if err != nil {
		return
	}
	err = seedLanguages()
	if err != nil {
		return
	}
	err = seedRaces("/api/subraces")
	if err != nil {
		return
	}
	err = seedRaces("/api/races")
	if err != nil {
		return
	}
}
