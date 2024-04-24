package storage

func artistsMatch(existingArtists []tableArtist, newArtistNames []string) bool {
	if len(existingArtists) != len(newArtistNames) {
		return false
	}

	existingNames := make(map[string]bool)
	for _, artist := range existingArtists {
		existingNames[artist.Name] = true
	}

	for _, name := range newArtistNames {
		if !existingNames[name] {
			return false
		}
	}

	return true
}
