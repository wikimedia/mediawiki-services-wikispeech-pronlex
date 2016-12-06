package symbolset2

// Mapper is a struct for package private usage. To create a new instance of Mapper, use LoadMapper.
type Mapper struct {
	Name       string
	SymbolSet1 SymbolSet
	SymbolSet2 SymbolSet
}
