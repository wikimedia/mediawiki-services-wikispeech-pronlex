var ADMLD = {};


ADMLD.baseURL = window.location.origin;


ADMLD.AdminLexDefModel = function () {
    var self = this; 
    
    self.symbolCategories = {
	'Phoneme': ["Syllabic", "NonSyllabic", "Stress"]
	, 'Delimiter': ["PhonemeDelimiter", "ExplicitPhonemeDelimiter", "SyllableDelimiter", "MorphemeDelimiter", "WordDelimiter"] 
    };
    
    self.nRead = ko.observable(0);
    
    ADMLD.readLexiconFile = function(fileObject) {
	var i = 0;
	new LineReader(fileObject).readLines(function(line){
	    i = i + 1;
	    if (i % 1000 === 0 ) {
		//console.log(i);
		self.nRead(i);
	    };
	});
    };
    
    
    
    self.lexicons = ko.observableArray();
    self.selectedLexicon = ko.observable();
    self.symbolSets = ko.observable({});
    
    self.selectedSymbolSet = ko.observable();
    self.selectedSymbol = ko.observable({});

    self.showSymbolSet = function(lexicon) {
	self.selectedSymbol({});
	self.selectedLexicon(lexicon);
     	var symbolSetName = lexicon.symbolSetName;
	if (! self.symbolSets().hasOwnProperty(symbolSetName)) {
	    self.selectedSymbolSet({});
	} else {
	    self.selectedSymbolSet(self.symbolSets()[symbolSetName]);
	};
	
    };
    
    self.setSelectedSymbol= function (symbol) {
	self.selectedSymbol(symbol);
    };
 
    // self.importLexiconFile = function() {
    // 	console.log("Opening file");
    // }

    
    self.addLexicon = function(lexiconName, symbolSetName) {
	//self = this;
	//self.lexiconName = lexiconName;
	//self.symbolSetName = symbolSetName;
	
	var newLex = {'id': 0, 'name' :  lexiconName, 'symbolSetName': symbolSetName};
	self.lexicons.push(newLex);
    };
    
    self.addSymbolToSet = function(symbolSetName, symbol) {
	if ( ! self.symbolSets().hasOwnProperty(symbolSetName) ) {
	    var ss = self.symbolSets();		
	    ss[symbolSetName] = [];
	};
	self.symbolSets()[symbolSetName].push(symbol);
   
	//console.log(self.symbolSets());
    };


};


var adm = new ADMLD.AdminLexDefModel();
ko.applyBindings(adm);



adm.addLexicon("nisse1", "kvack1");
adm.addLexicon("nisse2", "kvack2");

adm.addSymbolToSet("kvack1", {'symbol': 'a:', 'category': 'Phoneme', 'subcat' : 'Syllabic', 'description': 'h(a)t', 'ipa' : 'ɒː'});
adm.addSymbolToSet("kvack1", {'symbol': 'b', 'category': 'Phoneme', 'subcat' : 'NonSyllabic', 'description': '(b)il', 'ipa' : 'b'});

adm.addSymbolToSet("kvack2", {'symbol': 'O', 'category': 'Phoneme', 'subcat' : 'Syllabic', 'description': 'h(å)ll', 'ipa' : 'ɔ'});
adm.addSymbolToSet("kvack2", {'symbol': 'p', 'category': 'Phoneme', 'subcat' : 'NonSyllabic', 'description': '(p)il', 'ipa' : 'p'});



$(document).on('click', '.selectable', (function(){
    console.log(this);
    $(this).addClass("selected").siblings().removeClass("selected");
}));
