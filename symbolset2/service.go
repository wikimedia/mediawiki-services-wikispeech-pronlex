package symbolset2

import (
	"fmt"
	"log"
	"strings"
)

// functions for use by the mapper http service

type MapperService struct {
	SymbolSets map[string]SymbolSet
	Mappers    map[string]Mapper
}

func (m MapperService) MapperNames() []string {
	var names = make([]string, 0)
	for name, _ := range m.Mappers {
		names = append(names, name)
	}
	return names
}

func (m MapperService) Delete(ssName string) error {
	_, ok := m.SymbolSets[ssName]
	if !ok {
		return fmt.Errorf("no existing symbol set named %s", ssName)
	}
	delete(m.SymbolSets, ssName)
	log.Printf("Deleted symbol set %v from cache", ssName)
	for mName, _ := range m.Mappers {
		if strings.HasPrefix(mName, ssName+" ") ||
			strings.HasSuffix(mName, " "+ssName) {
			delete(m.Mappers, mName)
			log.Printf("Deleted mapper %v from cache", mName)
		}
	}
	return nil
}

func (m MapperService) Load(fName string) error {
	ss, err := LoadSymbolSet(fName)
	if err != nil {
		return fmt.Errorf("couldn't load symbol set : %v", err)
	}
	m.SymbolSets[ss.Name] = ss
	log.Printf("Loaded symbol set %v into cache", ss.Name)
	return nil
}

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
