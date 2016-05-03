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
    
    // List of Symbol objects
    self.selectedSymbolSet = ko.observableArray();
    // A sample symbol: {"symbol":"O","category":"Phoneme","description":"h(å)ll","ipa":"ɔ"}
    self.selectedSymbol = ko.observable({});
    
    self.showSymbolSet = function(lexicon) {
	self.selectedSymbol({});
	self.selectedLexicon(lexicon);
     	var symbolSetName = lexicon.symbolSetName;
	if (! self.symbolSets().hasOwnProperty(symbolSetName)) {
	    self.selectedSymbolSet().removeAll();
	} else {
	    self.selectedSymbolSet(self.symbolSets()[symbolSetName]);
	};
    };
    
    self.setSelectedSymbol= function (symbol) {
	self.selectedSymbol(symbol);
    };

    // TODO hard wired list of symbol set file header field names  
    self.headerFields = {'DESCRIPTION' : true, 'SYMBOL': true, 'IPA': true, 'CATEGORY': true};
    self.readSymbolSetFile = function (symbolSetfile) {
	
	// returns hash of header field name -> field index
	function headerIndexMap(header) {
	    var rez = {};
	    
	    var fields = header.trim().split(/\t/); 
	    // TODO hard wired number of fields
	    // TODO proper error handling
	    if (fields.length !== 4) { 
		alert("Wrong number of fields in header: "+ header);
		return;
	    };
	    for(var i = 0; i < fields.length; i++) {
		if (! self.headerFields.hasOwnProperty(fields[i])) {
		    // TODO proper error handling
		    alert("Unknown header field: "+ fields[i]);
		}  
		rez[fields[i]] = i;
	    };
	    return rez;
	};

	var reader = new FileReader();
	reader.onloadend = function(evt) {      
	    // Currently expecting hard wired tab separated format: 
	    // DESC/EXAMPLE	NST-XSAMPA	WS-SAMPA	IPA	TYPE
	    // Lines starting with # are descarded
	    
            lines = evt.target.result.split(/\r?\n/);
            if( lines.length > 0 ) {
		var header = lines.shift();
		var headerIndexes = headerIndexMap(header);
		
	    } else {  // TODO How do you do error handling when asynchronously reading a file?
		alert("Empty input file: "+ symbolSetfile.name)
		return; // ?
	    }
	    lines.forEach(function (line) {
		if (line.trim() === "") return; // "continue"
		if (line.trim().startsWith("#")) return; // "continue"
		
		var fs = line.split(/\t/);
		// TODO hard wired
		if (fs.length !== 4 ) alert("Wrong number of fields in line: "+ line);
		var symbol = {'symbol': fs[headerIndexes['SYMBOL']],
			      'category': fs[headerIndexes['CATEGORY']],
			      'description': fs[headerIndexes['DESCRIPTION']],
			      'ipa': fs[headerIndexes['IPA']]
			     };
		self.selectedSymbolSet.push(symbol);
		console.log(JSON.stringify(symbol));
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

//adm.showSymbolSet({"id":0,"name":"nisse2","symbolSetName":"kvack2"});

$(document).on('click', '.selectable', (function(){
    $(this).addClass("selected").siblings().removeClass("selected");
}));
