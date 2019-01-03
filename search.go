package main

func searchAlias(substr string) ([]ListEntry, error) {
	foundEntries := []ListEntry{}
	entries, err := getList()
	if err != nil {
		return foundEntries, err
	}
	for _, x := range entries {
		if Contains(substr, x.AnotherOptions, x.Hostname, x.Name, x.Port, x.User) {
			foundEntries = append(foundEntries, x)
		}
	}

	return foundEntries, err
}
