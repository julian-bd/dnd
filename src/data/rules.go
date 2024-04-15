package data

func AbilityNames() ([]string, error) {
	var abilities []string
	rows, err := db.Query("SELECT name FROM ability")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		abilities = append(abilities, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return abilities, nil
}

func TraitNames() ([]string, error) {
	var traits []string
	rows, err := db.Query("SELECT name FROM trait")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		traits = append(traits, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return traits, nil
}

func LanguageNames() ([]string, error) {
	var languages []string
	rows, err := db.Query("SELECT name FROM language")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		languages = append(languages, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return languages, nil
}

func ProficiencyNames() ([]string, error) {
	var proficiencies []string
	rows, err := db.Query("SELECT name FROM proficiency")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		proficiencies = append(proficiencies, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return proficiencies, nil
}
