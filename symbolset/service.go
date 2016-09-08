package symbolset

import "fmt"

// functions for use by the mapper http service

type MapperService struct {
	SymbolSets map[string]SymbolSet
	mappers    map[string]Mapper
}

func (m MapperService) getOrCreateMapper(fromName string, toName string) (Mapper, error) {
	if m.mappers == nil {
		m.mappers = make(map[string]Mapper)
	}
	name := fromName + " to " + toName
	mapper, ok := m.mappers[name]
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
	return LoadMapper(from, to)
}

// Map is used by the server to map a transcription from one symbol set to another
func (m MapperService) Map(fromName string, trans string, toName string) (string, error) {
	mapper, err := m.getOrCreateMapper(fromName, toName)
	if err != nil {
		return "", fmt.Errorf("couldn't create mapper from %s to %s", fromName, toName)
	}
	return mapper.MapTranscription(trans)
}

// GetMapTable is used by the server to show/get a mapping table between two symbol sets
func (m MapperService) GetMapTable(fromName string, toName string) (Mapper, error) {
	mapper, err := m.getOrCreateMapper(fromName, toName)
	if err != nil {
		return Mapper{}, fmt.Errorf("couldn't create mapper from %s to %s", fromName, toName)
	}
	return mapper, nil
}
