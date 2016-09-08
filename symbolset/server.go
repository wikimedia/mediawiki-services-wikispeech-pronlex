package symbolset

import "fmt"

// functions for use by the mapper http service

func createMapper(sets []SymbolSet, fromName string, toName string) (Mapper, error) {
	var from, to SymbolSet
	for _, set := range sets {
		if set.Name == fromName {
			from = set
		} else if set.Name == toName {
			to = set
		}
	}
	return LoadMapper(from, to)
}

// CanMap is used by the server
func CanMap(sets []SymbolSet, fromName string, toName string) bool {
	_, err := createMapper(sets, fromName, toName)
	return (err == nil)
}

// CanMap is used by the server
func Map(sets []SymbolSet, fromName string, trans string, toName string) (string, error) {
	m, err := createMapper(sets, fromName, toName)
	if err != nil {
		return "", fmt.Errorf("couldn't create mapper from %s to %s", fromName, toName)
	}
	return m.MapTranscription(trans)
}
