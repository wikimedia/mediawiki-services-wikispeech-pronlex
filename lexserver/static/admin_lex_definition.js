var ADMLD = {};

ADMLD.baseURL = window.location.origin;

// From http://stackoverflow.com/a/8809472
ADMLD.generateUUID = function() {
    var d = new Date().getTime();
    if(window.performance && typeof window.performance.now === "function"){
        d += performance.now(); //use high-precision timer if available
    }
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (d + Math.random()*16)%16 | 0;
        d = Math.floor(d/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
};


ADMLD.AdminLexDefModel = function () {
    var self = this; 
    
    
    self.uuid = ADMLD.generateUUID();

    self.serverMessage = ko.observable("_");
    
    self.maxMessages = 10;
    self.serverMessages = ko.observableArray();
    
    self.connectWebSock = function() {
	var zock = new WebSocket(ADMLD.baseURL.replace("http://", "ws://") + "/websockreg" );
	zock.onopen = function() {
	    console.log("ADMLD.connectWebSock: sending uuid over zock: "+ self.uuid);
	    zock.send("CLIENT_ID: "+ self.uuid);
	};
	zock.onmessage = function(e) {
	    // Just drop the keepalive message
	    if(e.data === "WS_KEEPALIVE") {
		// var d = new Date();
		// var h = d.getHours();
		// var m = d.getMinutes();
		// var s = d.getSeconds();
		// var msg = "Websocket keepalive recieved "+ h +":"+ m +":"+ s;
		// self.serverMessage(msg);
	    }
	    else {
		//console.log("Websocket got: "+ e.data)
		self.serverMessage(e.data);
	    };
	};
	zock.onerror = function(e){
	    console.log("websocket error: " + e.data);
	};
	zock.onclose = function (e) {
	    console.log("websocket got close event: "+ e.code)
	};
    };
    
    
    
    // self.connectServerMessageWS = function() {
    // 	var zock = new WebSocket( ADMLD.baseURL.replace("http://", "ws://") + "/adminmessages" );
    // 	zock.onmessage = function(e) {
    // 	    if (e.data === "WS_KEEPALIVE") {
    // 		// Dumdeedum
    // 	    } else {
    // 		self.serverMessage(e.data);
    // 		while (self.serverMessages().length >= self.maxMessages) {
    // 		    self.serverMessages().pop();
    // 		};
    // 		self.serverMessages().push(e.data);
    // 	    };
    // 	};
    // 	zock.onerror = function(e){
    // 	    console.log("websocket error: " + e.data);
    // 	};
    // 	zock.onclose = function (e) {
    // 	    console.log("websocket got close event: "+ e.code)
    //     };
    // };
    

    // TODO hard wired names. Fetch from somewhere?
    self.symbolCategories = 
	["syllabic", "non syllabic", "stress", "phoneme delimiter", "syllable delimiter", "morpheme delimiter", "word delimiter"];
    
    
    
    self.uploadLexiconFile = function(lexiconFile) {
	var url = ADMLD.baseURL + "/admin/lexiconfileupload" ;//'server/index.php';
	var xhr = new XMLHttpRequest();
	var fd = new FormData();
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function() {
            if (xhr.readyState === 4 && xhr.status === 200) {
		// Every thing ok, file uploaded
		console.log("uploadLexiconFile return response text", xhr.responseText); // handle response.
		//} //else { // TODO  This doesn't seem to be the right way to handle errors here
		//alert(xhr.responseText);
	    };// else { // TODO this doesn't work
 	//	console.log("uploadLexiconFile return status", xhr.statusText);
	 //   };
	};
	fd.append("client_uuid", self.uuid);
	fd.append("lexicon_id", self.selectedLexicon().id);
	fd.append("lexicon_name", self.selectedLexicon().name);
	fd.append("symbolset_name", self.selectedLexicon().symbolSetName);
	fd.append("upload_file", lexiconFile);
	xhr.send(fd);
    };
    
    
    self.lexicons = ko.observableArray();
    
    // selectedLexicon is a trigger for different things
    // Sample lexicon object: {"id":0,"name":"nisse2","symbolSetName":"kvack2"}
    self.selectedLexicon = ko.observable({'id': 0, 'name': '', 'symbolSetName': ''});
    //self.selectedLexicon = ko.observable();

    self.addLexiconName = ko.observable("");
    self.addSymbolSetName = ko.observable("");

    self.loadLexiconNames = function () {
	
	$.getJSON(ADMLD.baseURL +"/listlexicons")
	    .done(function (data) {
		self.lexicons(data);
		self.loadSymbolSetsIfEmpty(data); 
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert("loadLexiconNames says: "+ xhr.responseText);
	    });
    };
    
    // TODO update also name used in symbol set table, etc
    //      probably also name of selected lexicon?
    // Maybe a bit chaotic. 
    // Seems like an update of lexicon/symbol set name triggers adding already existing symbols to symbol set?! 
    self.updateLexicon = function () {
	
    	if ( self.selectedLexicon().name.trim() === "" || self.selectedLexicon().symbolSetName.trim() === "" ) {
    	    alert("Name and Symbol set name field must not be empty")
    	    return;
    	}
	
    	var params = {'id' : self.selectedLexicon().id, 'name' : self.selectedLexicon().name, 'symbolsetname' : self.selectedLexicon().symbolSetName}
	
    	// var xhr = new XMLHTTPRequest();
	// xhr.
	
	$.get(ADMLD.baseURL + "/admin/insertorupdatelexicon", params)
    	    .done(function(data, status, xhr){
		//console.log("status: "+ status);
		//console.log("xhr: "+ JSON.stringify(xhr));
    		self.loadLexiconNames();
    		self.selectedLexicon(data); // ?
    	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		console.log("updateLexicon fail xhr: "+ JSON.stringify(xhr));
		console.log("updateLexicon xhr.responseText: "+ xhr.responseText);
    		console.log("updateLexicon fail textStatus: "+ textStatus);
		console.log("updateLexicon fail errorThrown: "+ errorThrown);
		alert("updateLexicon says: "+ xhr.responseText);
    	    });	
    };
    
    self.deleteLexicon = function (lexicon) {
	
	console.log("deleteLexicon lexicon was called: "+ JSON.stringify(lexicon));
	
	var params = {'id' : lexicon.id}
    	$.get(ADMLD.baseURL + "/admin/deletelexicon", params)
    	    .done(function(data){
    		console.log("deleteLexicon got response: "+ JSON.stringify(data));
    		self.loadLexiconNames();
    	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
    		console.log("deleteLexicon fail");
		alert(xhr.responseText);	    
    	    });
    };
    
    // TODO Might not be safe to keep. Mostly for develpoment.
    // Whipes lexicon and matching symbolset totally from the DB 
    self.superDeleteLexicon = function (lexicon) {
	
	console.log("superDeleteLexicon called : "+ JSON.stringify(lexicon));
	var params = {'id' : lexicon.id, 'client_uuid': self.uuid};
    	$.get(ADMLD.baseURL + "/admin/superdeletelexicon", params)
    	    .done(function(data){
    		// dumdelidum
    		self.loadLexiconNames();
    	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
    		alert(xhr.responseText);	    
    	    });
    };
    

    self.exportSelectedLexicon = function() {
	var lex = self.selectedLexicon();
	lex.client_uuid = self.uuid;
	$.get(ADMLD.baseURL + "/admin/exportlexicon", lex)
	    .done(function(exportFName) {
		// build download URL
		var url = ADMLD.baseURL +"/download?file="+ exportFName;
		//console.log("THIS: "+ url);
		// add download URL to element
		$("#export_download_link").append("<a href="+ url +">"+ exportFName +"<a>"); // TODO This is not the way to build up elements safely?

	    })
	    .fail(function (xhr, textStatus, errorThrown) {
    		alert(xhr.responseText);	    
    	    });
    };
 


    
    // An object/hash with symbol set name as key and a list of symbol objects as value
    self.symbolSets = ko.observable({});
    
    self.deleteSymbol = function(zymbl) {
	var currSyms = self.symbolSets()[self.selectedLexicon().symbolSetName];
	
	// This appears to be JavaScript's way of removing an entry from an array: 
	var i = currSyms.indexOf(zymbl);
	if(i != -1) {
	    currSyms.splice(i, 1);
	}
	
	// update to trigger event
	// TODO why is this needed?
	self.selectedLexicon(self.selectedLexicon());
	
    };

    self.showSymbolSet = function(lexicon) {
	
	// update to trigger event
	// TODO why is this needed?
	self.selectedLexicon(lexicon);
    };
    

    // List of Symbol objects
    self.selectedSymbolSet = ko.computed(function() {
	
	if (self.symbolSets().hasOwnProperty(self.selectedLexicon().symbolSetName)) {
	    return self.symbolSets()[self.selectedLexicon().symbolSetName];
	} else {
	    return [];
	};
    }, this);
    
    self.saveSymbolSetToDB = function () {
	
	var ssName = self.selectedLexicon().symbolSetName;
	if("" === ssName) {
	    console.log("saveSymbolSetToDB: no symbol set name");
	    return;
	};
	var ss = self.symbolSets()[ssName];
	if (typeof ss === 'undefined' ) {
	    console.log("saveSymbolSetToDB: no symbol set to save");
	    return;
	};
	
	// for(var i = 0; i < ss.length; i++) {
	//     console.log(JSON.stringify(ss[i]));
	//  };	
	
	// TODO signal to user that something is happening/has happened

	$.post(ADMLD.baseURL + "/admin/savesymbolset", JSON.stringify(ss))
	    .fail(function (xhr, textStatus, errorThrown) {
		console.log("saveSymbolSetToDB fail xhr: "+ JSON.stringify(xhr));
		console.log("saveSymbolSetToDB xhr.responseText: "+ xhr.responseText);
    		console.log("saveSymbolSetToDB fail textStatus: "+ textStatus);
		console.log("saveSymbolSetToDB fail errorThrown: "+ errorThrown);
		alert("saveSymbolSetToDB says: "+ xhr.responseText);
    	    });	
	
    };


    // A sample symbol: {"symbol":"O","category":"Phoneme","description":"h(å)ll","ipa":"ɔ"}
    self.selectedSymbol = ko.observable({});
    
    self.symbolSetEmpty = function(lex) {
	if( ! self.symbolSets().hasOwnProperty(lex.symbolSetName) ) {
	    return true;
	};
	if (self.symbolSets()[lex.symbolSetName].length === 0) {
	    return true;
	};
	return false;
    };

    self.loadSymbolSetsIfEmpty = function(lexicons) {
	lexicons.forEach(function(lex) {
	    if(self.symbolSetEmpty(lex)) {
		self.loadSymbolSet(lex);
	    };
	}); 
    };

    self.loadSymbolSets = function(lexicons) {
	lexicons.forEach(function(lex) {
	    self.loadSymbolSet(lex);
	}); 
    };
    self.loadSymbolSet = function (lexicon) {
	console.log("loadSymbolSet: "+ JSON.stringify(lexicon))
	$.getJSON(ADMLD.baseURL +"/admin/listsymbolset", {lexiconid: lexicon.id}, function (data) {
	    if(data) {
		data.forEach(function(s) {
		    var sym = {'lexiconId': s.lexiconId, 'symbol': s.symbol, 'category': s.category, 'description': s.description, 'ipa': s.ipa};
		    self.addSymbolToSet(lexicon.symbolSetName, sym);
		});
		
	    } else {
		// Something fishy: this one gets called on pageload, _after_ first succeeding above
		console.log("self.loadSymbolSet: call to /listsymbolset/ returned null for lexiconid: "+ lexicon.id);
	    };
	    //self.selectedLexicon(self.selectedLexicon());
	})
	    .fail(function (xhr, textStatus, errorThrown) {
		console.log("loadSymbolSet fail xhr: "+ JSON.stringify(xhr));
		console.log("loadSymbolSet xhr.responseText: "+ xhr.responseText);
    		console.log("loadSymbolSet fail textStatus: "+ textStatus);
		console.log("loadSymbolSet fail errorThrown: "+ errorThrown);
		    
		alert(xhr.responseText);
	    });	
    };
    
    
    self.saveSymbolSet = function () {
	$.post(ADMLD.baseURL +"/admin/savesymbolset", JSON.stringify(self.selectedSymbolSet))
	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);
	    });
    };
    

    self.setSelectedSymbol= function (symbol) {
	self.selectedSymbol(symbol);
    };
    
    // TODO hard wired list of symbol set file header field names
    // Fetch from somewhere?
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
	    //var newSyms = [];
	    lines.forEach(function (line) {
		if (line.trim() === "") return; // "continue"
		if (line.trim().startsWith("#")) return; // "continue"
		
		var fs = line.split(/\t/);
		// TODO hard wired
		if (fs.length !== 4 ) alert("Wrong number of fields in line: "+ line);
		var symbol = {'symbol': fs[headerIndexes['SYMBOL']],
			      'category': fs[headerIndexes['CATEGORY']],
			      'description': fs[headerIndexes['DESCRIPTION']],
			      'ipa': fs[headerIndexes['IPA']],
			      'lexiconId': self.selectedLexicon().id
			     };
		
		if(! self.symbolSets().hasOwnProperty(self.selectedLexicon().symbolSetName)) {
		    self.symbolSets()[self.selectedLexicon().symbolSetName] = [];
		};
		
		if (symbol.symbol === "" ) {
		    console.log("No symbol value, skipping line: "+ line);
		} else { 
		    self.addSymbolToSet(self.selectedLexicon().symbolSetName, symbol);
		};
            });
	    
	    // update to trigger event
	    // TODO why is this needed?
	    self.selectedLexicon(self.selectedLexicon());
	};
	
	reader.readAsText(symbolSetfile,"UTF-8");
    };
    

    
    self.addLexicon = function() {
		
	var newLex = {'id': 0, 'name' :  self.addLexiconName(), 'symbolSetName': self.addSymbolSetName()};
	self.selectedLexicon(newLex);
	
	self.updateLexicon();
	
	self.addLexiconName("");
	self.addSymbolSetName("");
    };
    
    // These fields maka up the definition of a symbol
    self.symbolToAdd = ko.observable();
    self.categoryToAdd = ko.observable();
    self.descriptionToAdd = ko.observable();
    self.ipaToAdd = ko.observable();


    //TODO add input validation
    self.addSymbol = function() {
	
	var newSymbol = {'symbol': self.symbolToAdd(), 
			 'category': self.categoryToAdd(), 
			 'description': self.descriptionToAdd(), 
			 'ipa': self.ipaToAdd(),
			 'lexiconId': self.selectedLexicon().id};
	
	if( self.symbolToAdd() !== "") { 
	    self.addSymbolToSet(self.selectedLexicon().symbolSetName, newSymbol);
	};
	// empty input after adding new symbol
	self.symbolToAdd("");
	//self.categoryToAdd("");
	self.descriptionToAdd("");
	self.ipaToAdd("");

	// update to trigger event
	// TODO why is this needed?
	self.selectedLexicon(self.selectedLexicon());
    };
    
    self.addSymbolToSet = function(symbolSetName, symbol) {	
	if ( ! self.symbolSets().hasOwnProperty(symbolSetName) ) {
	    var ss = self.symbolSets();		
	    ss[symbolSetName] = [];
	};
	if(symbol.symbol === "" || symbol.symbol === undefined) {
	    var msg = "addSymbolToSet: Symbol field cannot be empty: "+ JSON.stringify(symbol); 
	    console.log(msg);
	    alert(msg); // TODO
	    return;
	};
	if(symbol.description === "" || symbol.description === undefined) {
	    var msg = "addSymbolToSet: Description field cannot be empty"; 
	    console.log(msg);
	    alert(msg); // TODO 
	    return;
	};
	

	// TODO validate that syllabic/non... etc has an IPA symbol
	
	
	// There is an uniqueness constraint in the database 
	var dupes = self.symbolSets()[symbolSetName].filter(function(obj) {
	    return obj.symbol === symbol.symbol;
	});
	
	if(dupes.length > 0) {
	    dupes.push(symbol);
	    var msg = "addSymbolToSet: duplicate symbols are not allowed: "+ JSON.stringify(dupes);
	    console.log(msg);
	    alert(msg); // TODO
	} else {    
	    self.symbolSets()[symbolSetName].push(symbol);
	}
    };
    
    self.setSelectedIPA = function(symbol) {
	self.ipaToAdd(symbol.symbol);
	self.descriptionToAdd(symbol.description);
    };
    
    
    self.nColumns = ko.observable(15);
    
    self.createIPATableRows = function (nColumns, ipaList ) {
	var res = [];
	var row = [];
	var j = 0;
	for(var i = 0; i < self.ipaTable.length; i++) {
	    
	    var ipaChar = {'symbol': self.ipaTable[i].symbol, 'description': self.ipaTable[i].description};
	    row.push(ipaChar);
	    j++;
	    if ( j === nColumns) {
		res.push(row);
		row = [];
		j = 0;
	    };
	}; // <- for
	// "flush":
	if ( j !== nColumns) {
	    res.push(row);
	};
	return res;
    }; 
    

    
    // TODO remove hard wired IPA table
    // This should be downloaded from lexserver: ipa_table.txt
    self.ipaTable = [{"symbol": "ɐ", "description":  "Near-open central vowel"},
		     {"symbol": "ɑ", "description":  "Open back unrounded vowel"},
		     {"symbol": "ɒ", "description":  "Open back rounded vowel"},
		     {"symbol": "ɓ", "description":  "Voiced bilabial implosive"},
		     {"symbol": "ɔ", "description":  "Open-mid back rounded vowel"},
		     {"symbol": "ɕ", "description":  "Voiceless alveolo-palatal fricative"},
		     {"symbol": "ɖ", "description":  "Voiced retroflex plosive"},
		     {"symbol": "ɗ", "description":  "Voiced alveolar implosive"},
		     {"symbol": "ɘ", "description":  "Close-mid central unrounded vowel"},
		     {"symbol": "ə", "description":  "Mid central vowel"},
		     {"symbol": "ɚ", "description":  "Rhotacized Mid central vowel"},
		     {"symbol": "ɛ", "description":  "Open-mid front unrounded vowel"},
		     {"symbol": "ɜ", "description":  "Open-mid central unrounded vowel"},
		     {"symbol": "ɝ", "description":  "Rhotacized Open-mid central unrounded vowel"},
		     {"symbol": "ɞ", "description":  "Open-mid central rounded vowel"},
		     {"symbol": "ɟ", "description":  "Voiced palatal plosive"},
		     {"symbol": "ɠ", "description":  "Voiced velar implosive"},
		     {"symbol": "ɡ", "description":  "Voiced velar plosive"},
		     {"symbol": "ɢ", "description":  "Voiced uvular plosive"},
		     {"symbol": "ɣ", "description":  "Voiced velar fricative"},
		     {"symbol": "ɤ", "description":  "Close-mid back unrounded vowel"},
		     {"symbol": "ɥ", "description":  "Labial-palatal approximant"},
		     {"symbol": "ɦ", "description":  "Voiced glottal fricative"},
		     {"symbol": "ɧ", "description":  "Swedish sj-sound. Similar to: Voiceless postalveolar fricative, Voiceless velar fricative"},
		     {"symbol": "ɨ", "description":  "Close central unrounded vowel"},
		     {"symbol": "ɩ", "description":  "pre-1989 form of 'ɪ' (obsolete)"},
		     {"symbol": "ɪ", "description":  "Near-close near-front unrounded vowel"},
		     {"symbol": "ɫ", "description":  "Velar/pharyngeal Alveolar lateral approximant"},
		     {"symbol": "ɬ", "description":  "Voiceless alveolar lateral fricative"},
		     {"symbol": "ɭ", "description":  "Retroflex lateral approximant"},
		     {"symbol": "ɮ", "description":  "Voiced alveolar lateral fricative"},
		     {"symbol": "ɯ", "description":  "Close back unrounded vowel"},
		     {"symbol": "ɰ", "description":  "Velar approximant"},
		     {"symbol": "ɱ", "description":  "Labiodental nasal"},
		     {"symbol": "ɲ", "description":  "Palatal nasal"},
		     {"symbol": "ɳ", "description":  "Retroflex nasal"},
		     {"symbol": "ɴ", "description":  "Uvular nasal"},
		     {"symbol": "ɵ", "description":  "Close-mid central rounded vowel"},
		     {"symbol": "ɶ", "description":  "Open front rounded vowel"},
		     {"symbol": "ɷ", "description":  "pre-1989 form of 'ʊ' (obsolete)"},
		     {"symbol": "ɸ", "description":  "Voiceless bilabial fricative"},
		     {"symbol": "ɹ", "description":  "Alveolar approximant"},
		     {"symbol": "ɺ", "description":  "Alveolar lateral flap"},
		     {"symbol": "ɻ", "description":  "Retroflex approximant"},
		     {"symbol": "ɼ", "description":  "Alveolar trill"},
		     {"symbol": "ɽ", "description":  "Retroflex flap"},
		     {"symbol": "ɾ", "description":  "Alveolar tap"},
		     {"symbol": "ɿ", "description":  "Syllabic voiced alveolar fricative (Sinologist usage)"},
		     {"symbol": "ʀ", "description":  "Uvular trill"},
		     {"symbol": "ʁ", "description":  "Voiced uvular fricative"},
		     {"symbol": "ʂ", "description":  "Voiceless retroflex fricative"},
		     {"symbol": "ʃ", "description":  "Voiceless postalveolar fricative"},
		     {"symbol": "ʄ", "description":  "Voiced palatal implosive"},
		     {"symbol": "ʅ", "description":  "Syllabic voiced retroflex fricative (Sinologist usage)"},
		     {"symbol": "ʆ", "description":  "Voiceless alveolo-palatal fricative (obsolete)"},
		     {"symbol": "ʇ", "description":  "Dental click (obsolete)"},
		     {"symbol": "ʈ", "description":  "Voiceless retroflex plosive"},
		     {"symbol": "ʉ", "description":  "Close central rounded vowel"},
		     {"symbol": "ʊ", "description":  "Near-close near-back rounded vowel"},
		     {"symbol": "ʋ", "description":  "Labiodental approximant"},
		     {"symbol": "ʌ", "description":  "Open-mid back unrounded vowel"},
		     {"symbol": "ʍ", "description":  "Voiceless labiovelar approximant"},
		     {"symbol": "ʎ", "description":  "Palatal lateral approximant"},
		     {"symbol": "ʏ", "description":  "Near-close near-front rounded vowel"},
		     {"symbol": "ʐ", "description":  "Voiced retroflex fricative"},
		     {"symbol": "ʑ", "description":  "Voiced alveolo-palatal fricative"},
		     {"symbol": "ʒ", "description":  "Voiced postalveolar fricative"},
		     {"symbol": "ʓ", "description":  "Voiced alveolo-palatal fricative (obsolete)"},
		     {"symbol": "ʔ", "description":  "Glottal stop"},
		     {"symbol": "ʕ", "description":  "Voiced pharyngeal fricative"},
		     {"symbol": "ʖ", "description":  "Alveolar lateral click (obsolete)"},
		     {"symbol": "ʗ", "description":  "Postalveolar click (obsolete)"},
		     {"symbol": "ʘ", "description":  "Bilabial click"},
		     {"symbol": "ʙ", "description":  "Bilabial trill"},
		     {"symbol": "ʚ", "description":  "A mistake for [œ]"},
		     {"symbol": "ʛ", "description":  "Voiced uvular implosive"},
		     {"symbol": "ʜ", "description":  "Voiceless epiglottal fricative"},
		     {"symbol": "ʝ", "description":  "Voiced palatal fricative"},
		     {"symbol": "ʞ", "description":  "Velar click (obsolete)"},
		     {"symbol": "ʟ", "description":  "Velar lateral approximant"},
		     {"symbol": "ʠ", "description":  "'Voiceless' uvular implosive (obsolete)"},
		     {"symbol": "ʡ", "description":  "Epiglottal plosive"},
		     {"symbol": "ʢ", "description":  "Voiced epiglottal fricative"},
		     {"symbol": "ʣ", "description":  "Voiced alveolar affricate (obsolete)"},
		     {"symbol": "ʤ", "description":  "Voiced postalveolar affricate (obsolete)"},
		     {"symbol": "ʥ", "description":  "Voiced alveolo-palatal affricate (obsolete)"},
		     {"symbol": "ʦ", "description":  "Voiceless alveolar affricate (obsolete)"},
		     {"symbol": "ʧ", "description":  "Voiceless postalveolar affricate (obsolete)"},
		     {"symbol": "ʨ", "description":  "Voiceless alveolo-palatal affricate (obsolete)"},
		     {"symbol": "ʩ", "description":  "velopharyngeal fricative"},
		     {"symbol": "ʪ", "description":  "voiceless lateral alveolar fricative"},
		     {"symbol": "ʫ", "description":  "voiced lateral alveolar fricative"},
		     {"symbol": "ʬ", "description":  "Bilabial percussive"},
		     {"symbol": "ʭ", "description":  "Bidental percussive"},
		     {"symbol": "ʮ", "description":  "Syllabic labialized voiced alveolar fricative (Sinologist usage)"},
		     {"symbol": "ʯ", "description":  "Syllabic labialized voiced retroflex fricative (Sinologist usage)"},
		     {"symbol": "a", "description":  "Open front unrounded vowel"},
		     {"symbol": "b", "description":  "bilabial plosive"},
		     {"symbol": "c", "description":  "palatal plosive"},
		     {"symbol": "d", "description":  "alveolar plosive"},
		     {"symbol": "e", "description":  "close-mid front unrounded vowel"},
		     {"symbol": "f", "description":  "labiodental fricative"},
		     {"symbol": "g", "description":  "velar plosive Ascii g"},
		     {"symbol": "h", "description":  "glottal fricative"},
		     {"symbol": "i", "description":  "close front unrounded vowel"},
		     {"symbol": "j", "description":  "palatal approximant"},
		     {"symbol": "k", "description":  "velar plosive"},
		     {"symbol": "l", "description":  "lateral alveolar approximant"},
		     {"symbol": "l̩̩̩", "description":  "syllabic l"},
		     {"symbol": "m", "description":  "bilabial nasal"},
		     {"symbol": "m̩̩", "description":  "syllabic m"},
		     {"symbol": "n", "description":  "alveolar nasal"},
		     {"symbol": "n̩̩̩", "description":  "syllabic n"},
		     {"symbol": "o", "description":  "close-mid back rounded vowel"},
		     {"symbol": "p", "description":  "bilabial plosive"},
		     {"symbol": "q", "description":  "uvular plosive"},
		     {"symbol": "r", "description":  "alveolar trill"},
		     {"symbol": "s", "description":  "alveolar fricative"},
		     {"symbol": "t", "description":  "alveolar plosive"},
		     {"symbol": "u", "description":  "close back rounded vowel"},
		     {"symbol": "v", "description":  "labiodental fricative"},
		     {"symbol": "w", "description":  "labial-velar approximant"},
		     {"symbol": "x", "description":  "velar fricative"},
		     {"symbol": "y", "description":  "close front rounded vowel"},
		     {"symbol": "z", "description":  "alveolar fricative"},
		     {"symbol": "æ", "description":  "raised-open front unrounded vowel"},
		     {"symbol": "ç", "description":  "palatal fricative"},
		     {"symbol": "ð", "description":  "dental fricative"},
		     {"symbol": "ø", "description":  "close-mid front rounded vowel"},
		     {"symbol": "ħ", "description":  "pharyngeal fricative"},
		     {"symbol": "ŋ", "description":  "velar nasal"},
		     {"symbol": "œ", "description":  "Open-mid front rounded vowel"},
		     {"symbol": "β", "description":  "bilabial fricative"},
		     {"symbol": "θ", "description":  "dental fricative"},
		     {"symbol": "χ", "description":  "uvular fricative"},
		     {"symbol": "aː", "description":  "Open front unrounded vowel (long)"},
		     {"symbol": "eː", "description":  "close-mid front unrounded vowel (long)"},
		     {"symbol": "iː", "description":  "close front unrounded vowel (long)"},
		     {"symbol": "oː", "description":  "close-mid back rounded vowel (long)"},
		     {"symbol": "uː", "description":  "close back rounded vowel (long)"},
		     {"symbol": "yː", "description":  "close front rounded vowel (long)"},
		     {"symbol": "æː", "description":  "raised-open front unrounded vowel (long)"},
		     {"symbol": "øː", "description":  "close-mid front rounded vowel (long)"},
		     {"symbol": "œː", "description":  "Open-mid front rounded vowel (long)"},
		     {"symbol": "ɐː", "description":  "Near-open central vowel (long)"},
		     {"symbol": "ɑː", "description":  "Open back unrounded vowel (long)"},
		     {"symbol": "ɒː", "description":  "Open back rounded vowel (long)"},
		     {"symbol": "ɔː", "description":  "Open-mid back rounded vowel (long)"},
		     {"symbol": "ɘː", "description":  "Close-mid central unrounded vowel (long)"},
		     {"symbol": "əː", "description":  "Mid central vowel (long)"},
		     {"symbol": "ɚː", "description":  "Rhotacized Mid central vowel (long)"},
		     {"symbol": "ɛː", "description":  "Open-mid front unrounded vowel (long)"},
		     {"symbol": "ɜː", "description":  "Open-mid central unrounded vowel (long)"},
		     {"symbol": "ɝː", "description":  "Rhotacized Open-mid central unrounded vowel (long)"},
		     {"symbol": "ɞː", "description":  "Open-mid central rounded vowel (long)"},
		     {"symbol": "ɤː", "description":  "Close-mid back unrounded vowel (long)"},
		     {"symbol": "ɨː", "description":  "Close central unrounded vowel (long)"},
		     {"symbol": "ɪː", "description":  "Near-close near-front unrounded vowel (long)"},
		     {"symbol": "ɯː", "description":  "Close back unrounded vowel (long)"},
		     {"symbol": "ɵː", "description":  "Close-mid central rounded vowel (long)"},
		     {"symbol": "ɶː", "description":  "Open front rounded vowel (long)"},
		     {"symbol": "ʉː", "description":  "Close central rounded vowel (long)"},
		     {"symbol": "ʊː", "description":  "Near-close near-back rounded vowel (long)"},
		     {"symbol": "ʌː", "description":  "Open-mid back unrounded vowel (long)"},
		     {"symbol": "ʏː", "description":  "Near-close near-front rounded vowel (long)"}];

    self.ipaSymbols = ko.observableArray();
    self.loadIPASymbols = function () {
	self.ipaSymbols.push("");
	for(var i = 0; i < self.ipaTable.length; i++) {
	    self.ipaSymbols.push(self.ipaTable[i].symbol);
	}
    };
    
    
    self.ipaTableRows = ko.computed(function() {
	var n = self.nColumns();
	return self.createIPATableRows(n, self.dummyIPA);
    }); 

};


// For marking the selected row in a table
$(document).on('click', '.selectable', (function(){
    $(this).addClass("selected").siblings().removeClass("selected");
}));


var adm = new ADMLD.AdminLexDefModel();
adm.loadIPASymbols();
ko.applyBindings(adm);
adm.loadLexiconNames();
//adm.connectWebSock();
adm.connectWebSock();

console.log("UUID: "+ adm.uuid);
