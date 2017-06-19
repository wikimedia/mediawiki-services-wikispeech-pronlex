package symbolset

import (
	"fmt"
	"log"
	"strings"
)

// functions for use by the mapper http service

// MapperService is a container for maintaining 'cached' mappers and their symbol sets. Please note that currently, MapperService need to be used as mutex, see lexserver/mapper.go
type MapperService struct {
	SymbolSets map[string]SymbolSet
	Mappers    map[string]Mapper
}

// MapperNames lists the names for all loaded mappers
func (m MapperService) MapperNames() []string {
	var names = make([]string, 0)
	for name := range m.Mappers {
		names = append(names, name)
	}
	return names
}

// Delete is used to delete a named symbol set from the cache. Deletes the named symbol set, and all mappers using this symbol set.
func (m MapperService) Delete(ssName string) error {
	_, ok := m.SymbolSets[ssName]
	if !ok {
		return fmt.Errorf("no existing symbol set named %s", ssName)
	}
	delete(m.SymbolSets, ssName)
	log.Printf("Deleted symbol set %v from cache", ssName)
	for mName := range m.Mappers {
		if strings.HasPrefix(mName, ssName+" ") ||
			strings.HasSuffix(mName, " "+ssName) {
			delete(m.Mappers, mName)
			log.Printf("Deleted mapper %v from cache", mName)
		}
	}
	return nil
}

// DeleteMapper is used to delete a mapper the cache.
func (m MapperService) DeleteMapper(fromName string, toName string) error {
	name := fromName + " to " + toName
	for mName := range m.Mappers {
		if mName == name {
			delete(m.Mappers, mName)
			log.Printf("Deleted mapper %v from cache", mName)
		}
	}
	return nil
}

// Load is used to load a symbol set from file
func (m MapperService) Load(symbolSetFile string) error {
	ss, err := LoadSymbolSet(symbolSetFile)
	if err != nil {
		return fmt.Errorf("couldn't load symbol set : %v", err)
	}
	m.SymbolSets[ss.Name] = ss
	log.Printf("Loaded symbol set %v into cache", ss.Name)
	return nil
}

// Clear is used to clear the cache (all loaded symbol sets and mappers)
func (m MapperService) Clear() {
	// TODO: MapperService need to be used as mutex, see lexserver/mapper.go
	m.SymbolSets = make(map[string]SymbolSet)
	m.Mappers = make(map[string]Mapper)
}

func (m MapperService) getOrCreateMapper(fromName string, toName string) (Mapper, error) {
	name := fromName + " to " + toName
	mapper, ok := m.Mappers[name]
	if ok {
		return mapper, nil
	}

	var nilRes Mapper
	var from, to SymbolSet
	from, okFrom := m.SymbolSets[fromName]
	if !okFrom {
		return nilRes, fmt.Errorf("couldn't find left hand symbol set named '%s'", fromName)
	}
	to, okTo := m.SymbolSets[toName]
	if !okTo {
		return nilRes, fmt.Errorf("couldn't find right hand symbol set named '%s'", toName)
	}
	mapper, err := LoadMapper(from, to)
	if err == nil {
		m.Mappers[name] = mapper
	}
	return mapper, err
}

// Map is used by the server to map a transcription from one symbol set to another
func (m MapperService) Map(fromName string, toName string, trans string) (string, error) {
	if toName == "ipa" {
		ss, ok := m.SymbolSets[fromName]
		if !ok {
			return "", fmt.Errorf("couldn't create mapper from %s to %s", fromName, toName)
		}
		return ss.ConvertToIPA(trans)
	} else if fromName == "ipa" {
		ss, ok := m.SymbolSets[toName]
		if !ok {
			return "", fmt.Errorf("couldn't create mapper from %s to %s", fromName, toName)
		}
		return ss.ConvertFromIPA(trans)
	} else {
		mapper, err := m.getOrCreateMapper(fromName, toName)
		if err != nil {
			return "", fmt.Errorf("couldn't create mapper from %s to %s : %v", fromName, toName, err)
		}
		return mapper.MapTranscription(trans)
	}
}

// GetMapTable is used by the server to show/get a mapping table between two symbol sets
func (m MapperService) GetMapTable(fromName string, toName string) (Mapper, error) {
	mapper, err := m.getOrCreateMapper(fromName, toName)
	if err != nil {
		return Mapper{}, fmt.Errorf("couldn't create mapper from %s to %s : %v", fromName, toName, err)
	}
	return mapper, nil
}
