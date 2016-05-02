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
    // Sample lexicon object: {"id":0,"name":"nisse2","symbolSetName":"kvack2"}
    self.selectedLexicon = ko.observable();
    // An object/hash with symbol set name as key and a list of symbol objects as value
    self.symbolSets = ko.observable({});
    
    // A sample symbol: {"symbol":"O","category":"Phoneme","subcat":"Syllabic","description":"h(å)ll","ipa":"ɔ"}
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

    self.readSymbolSetFile = function (symbolSetfile) {
	console.log(symbolSetfile.name);
	var reader = new FileReader();

	reader.onloadend = function(evt) {      
	    // Currently expecting hard wired tab separated format: 
	    // DESC/EXAMPLE	NST-XSAMPA	WS-SAMPA	IPA	TYPE
	    // Lines starting with # are descarded 
	    
            lines = evt.target.result.split(/\r?\n/);
            lines.forEach(function (line) {
		
		console.log(line);
            }); 
	};
	
	reader.readAsText(symbolSetfile,"UTF-8");
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
    
    self.selectedIPA = ko.observable({'symbol': ''});
    self.setSelectedIPA = function(symbol) {
	//console.log(">>>>> " + JSON.stringify(symbol));
	self.selectedIPA(symbol);
    };
    
    self.nColumns = ko.observable(6);
    
    self.createIPATableRows = function (nColumns, ipaList ) {
	var res = [];
	var row = [];
	//var tr = document.createElement("tr");
	var j = 0;
	for(var i = 0; i < ipaList.length; i++) {
	    // var td = document.createElement("td")
	    var ipaChar = {'symbol': ipaList[i]};
	    //td.setAttribute("data-bind", "click: $root.setSelectedIPA");
	    //td.setAttribute("text", ipaList[i]);
	    //td.innerHTML = ipaList[i];
	    //ko.applyBindingsToNode(td);
	    //tr.appendChild(td);
	    row.push(ipaChar);
	    j++;
	    if ( j === nColumns) {
		res.push(row);
		row = [];// document.createElement("tr");
		j = 0;
	    };
	}; // <- for
	// "flush":
	if ( j !== nColumns) {
	    res.push(row);
	};
	return res;
    }; 


    // TODO remove hard wired test
    self.dummyIPA = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'x', 'y', 'z', 'å', 'ä', 'ö'];
    
    self.ipaTableRows = ko.computed(function() {
	var n = self.nColumns();
	return self.createIPATableRows(n, self.dummyIPA);
    });// ko.observableArray();
    
    
    // TODO remove hard-wired test
    
    
    //self.ipaTableRows();
    


    // self.ipaTable = ko.computed(function() {
    // 	var tbody = document.createElement("tbody");
    // 	for(var i = 0; i < self.ipaTableRows().length; i++) {
    // 	    tbody.appendChild( self.ipaTableRows()[i] );
    // 	}; 
    // 	return tbody.outerHTML;
    // }, this);
    
};


var adm = new ADMLD.AdminLexDefModel();
ko.applyBindings(adm);



adm.addLexicon("nisse1", "kvack1");
adm.addLexicon("nisse2", "kvack2");

adm.addSymbolToSet("kvack1", {'symbol': 'a:', 'category': 'Phoneme', 'subcat' : 'Syllabic', 'description': 'h(a)t', 'ipa' : 'ɒː'});
adm.addSymbolToSet("kvack1", {'symbol': 'b', 'category': 'Phoneme', 'subcat' : 'NonSyllabic', 'description': '(b)il', 'ipa' : 'b'});

adm.addSymbolToSet("kvack2", {'symbol': 'O', 'category': 'Phoneme', 'subcat' : 'Syllabic', 'description': 'h(å)ll', 'ipa' : 'ɔ'});
adm.addSymbolToSet("kvack2", {'symbol': 'p', 'category': 'Phoneme', 'subcat' : 'NonSyllabic', 'description': '(p)il', 'ipa' : 'p'});

adm.showSymbolSet({"id":0,"name":"nisse2","symbolSetName":"kvack2"});

$(document).on('click', '.selectable', (function(){
    $(this).addClass("selected").siblings().removeClass("selected");
}));
