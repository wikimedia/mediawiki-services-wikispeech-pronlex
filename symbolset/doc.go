/*

Package symbolset is used to define symbol sets and mapping between different sets, such as IPA to SAMPA, and so on.

The basis of this package is a symbol set file on this tab-separated format:
	DESCRIPTION	SYMBOL	IPA	IPA UNICODE	CATEGORY

Sample lines (Swedish Wikispeech SAMPA):
	DESCRIPTION	SYMBOL	IPA	IPA UNICODE	CATEGORY
	sil	i:	iː	U+0069U+02D0	Syllabic
	bok	b	b	U+0062	NonSyllabic
	forna	rn	ɳ	U+0273	NonSyllabic

Note that the header is required on the first line. For real world examples (used for unit tests), see the test_data folder in this package.

**/
package symbolset
