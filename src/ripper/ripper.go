package ripper

import "fmt"

var baseUrl = "https://www.dnd5eapi.co"

func Seed() {
	err := seedTraits()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = seedProficiencies()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = seedLanguages()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = seedRaces("/api/subraces")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = seedRaces("/api/races")
	if err != nil {
		fmt.Println(err)
		return
	}
}
