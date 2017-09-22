## File formats

### Wikispeech lexicon file format

The Wikispeech lexicon file format is defined in `ws.go`. Lexicon files are tab separated text files (UTF-8 encoded), and should contain fields listed below. Empty fields are allowed in most positions.

Any lexicon files you want to import into the lexicon database must be in this file format.

1. Orth: The word's orthography
2. Pos: The part of speech tag
3. Morph: Morphological features (gender, number, etc)
4. WordParts: Compound parts, if any, separated by a plus sign (`+`)
5. Lemma: The word's lemma form
6. Paradigm: The name of the paradigm used for inflections
7. Lang: The word's language label
8. Preferred: Takes values 1/0, and is used to defined which reading for a specific orthography should be the standard one (in case of homographs)
9. Trans1: The first transcription (default for TTS)
10. Translang1: The language of the Trans1
11. Trans2: Alternative transcription
12. Translang2: The language of the Trans2
13. Trans3:  Alternative transcription
14. Translang3: The language of the Trans3
15. Trans4:  Alternative transcription
16. Translang4: The language of the Trans4
17. StatusName: Status of the lexicon entry
18. StatusSource: Source of the status

Sample line:

   ``finalspelet	NN	SIN|DEF|NOM|NEU	final+spelet	finalspel	s7n-övriga ex träd	sv-se	false	f I . "" n A: l . % s p e: . l e t	sv-se							imported	nst``



### NST lexicon file format

This format is used for converting NST lexicon files to the Wikispeech lexicon file format
