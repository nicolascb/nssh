package main

var (
	foundEntries []ListEntry
)

func ExecSearch(substr string) []ListEntry {
	entries := CreateList()

	for _, x := range entries {
		if Contains(substr, x.AnotherOptions, x.Hostname, x.Name, x.Port, x.User) {
			foundEntries = append(foundEntries, x)
		}
	}

	return foundEntries
}
