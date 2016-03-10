var DMCRLX = {};

DMCRLX.baseURL = window.location.origin

DMCRLX.CreateLexModel = function() {
    var self = this;
    
    DMCRLX.lexicons = ko.observableArray();
    DMCRLX.selectedLexicon = ko.observable();
    DMCRLX.symbolSet = ko.observableArray();
    //DMCRLX.symbolSetter = ko.computed();
    //DMCRLX.symbols = ko.observableArray();

    DMCRLX.symbolCategories = {
	'Phoneme': ["Syllabic", "NonSyllabic", "Stress"]
	, 'Delimiter': ["PhonemeDelimiter", "ExplicitPhonemeDelimiter", "SyllableDelimiter", "MorphemeDelimiter", "WordDelimiter"] 
    };

    DMCRLX.Symbol = function(lexiconId, symbol, category, subcat, description, ipa) {
	var self = this;
	self.lexiconId = lexiconId;
	self.symbol = symbol;
	self.category = ko.observable(category);
	self.subcat = subcat;
	self.description = description;
	self.ipa = ko.observable(ipa);
	self.symSubCats = ko.computed(function() {
	    return DMCRLX.symbolCategories[self.category()];
	}, this);
    }

    DMCRLX.ipaCharacters = ko.observableArray();
    DMCRLX.ipaDescriptions = ko.observable({});
    DMCRLX.getSymbolDescription = function (s) {
	return DMCRLX.ipaDescriptions()[s];
    };
    
    DMCRLX.loadIPATable = function() {
	$.get(DMCRLX.baseURL +"/ipa_table.txt", function(data){
	    var lines = data.trim().split(/\n+/g);
	    _.each(lines, function(l) {
		var fs = l.split(/\t/g);
		//SYMSETED.ipaCharacters.push({'ipachar': fs[1], 'desc':fs[4]});
		//console.log(fs[1]);
		//console.log(fs[4]);
		DMCRLX.ipaCharacters.push(fs[1]);
		DMCRLX.ipaDescriptions()[fs[1]] = fs[4];
	    });
	});
    }

    



    DMCRLX.loadLexiconNames = function () {
	
	$.getJSON(DMCRLX.baseURL +"/listlexicons")
	    .done(function (data) {
		DMCRLX.lexicons(data);
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);
	    });
    };
    
    DMCRLX.updateLexicon = function () {
	
	if ( DMCRLX.selectedLexicon().name === "" || DMCRLX.selectedLexicon().symbolSetName === "" ) {
	    alert("Name and Symbol set name field must not be empty")
	    return;
	}
	
	var params = {'id' : DMCRLX.selectedLexicon().id, 'name' : DMCRLX.selectedLexicon().name, 'symbolsetname' : DMCRLX.selectedLexicon().symbolSetName}
	
	
	$.get(DMCRLX.baseURL + "/admin/insertorupdatelexicon", params)
	    .done(function(data){
		DMCRLX.loadLexiconNames();
		DMCRLX.selectedLexicon(data);
	    })
	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);
	    });
	
    }

    DMCRLX.deleteLexicon = function () {
	var params = {'id' : DMCRLX.selectedLexicon().id} // , 'name' : DMCRLX.selectedLexicon().name, 'symbolsetname' : DMCRLX.selectedLexicon().symbolSetName}
	$.get(DMCRLX.baseURL + "/admin/deletelexicon", params)
	    .done(function(data){
		// dumdelidum
		DMCRLX.loadLexiconNames();
	    })
	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);	    
	    });
    }
    


    
    DMCRLX.addLexicon = function () {
	
	function hasNewLex(arr) {
	    for(var i = 0; i < arr.length; i++) {
		var x = arr[i];
		if( x.id === 0 && x.name === "" && x.symbolSetName === "") {
		    return true;
		}
	    }
	    return false;
	}
	
	var newLex = {'id': 0, 'name': "", 'symbolSetName': ""};
	console.log(JSON.stringify(newLex));
	if ( ! hasNewLex(DMCRLX.lexicons()) ) { 
	    DMCRLX.lexicons.push(newLex);
	}
	DMCRLX.selectedLexicon(newLex);
    }
    
    DMCRLX.addSymbol = function () {
	DMCRLX.symbolSet.push(new DMCRLX.Symbol(DMCRLX.selectedLexicon().id, "", "", "", "", ""));
    }

    DMCRLX.loadSymbolSet = ko.computed(
	function () {
	    if(DMCRLX.selectedLexicon() !== undefined) {
		//console.log("HÄR KOMMER GREJOR: "+ DMCRLX.selectedLexicon())
		$.getJSON(DMCRLX.baseURL +"/admin/listphonemesymbols", {lexiconId: DMCRLX.selectedLexicon().id}, function (data) {
		    var syms = _.map(data, function (s) {
			return new DMCRLX.Symbol(s.lexiconId, s.symbol, s.category, s.subcat, s.description, s.ipa);
		    }); 
		    DMCRLX.symbols(syms);
		});
	    }
	}
    , this);
    DMCRLX.saveSymbolSet = function () {
	var post = _.map(DMCRLX.symbols(), function(s) {
	    return {'lexiconId': s.lexiconId, 'symbol': s.symbol, 'category': s.category(), 'subcat': s.subcat, 'description': s.description, 'ipa': s.ipa()};
	});
	//console.log("SENDING: "+ post);
	$.post(DMCRLX.baseURL +"/admin/savesymbolset", JSON.stringify(post))
	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);
	    });
    };




    // On pageload:
    DMCRLX.loadLexiconNames();
    DMCRLX.loadIPATable();
    //console.log("!!!!!!!!"+ )
};

ko.applyBindings(new DMCRLX.CreateLexModel());
