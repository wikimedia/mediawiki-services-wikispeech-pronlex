var DMCRLX = {};

DMCRLX.baseURL = "http://localhost:8787/"

DMCRLX.CreateLexModel = function() {
    var self = this;
    
    DMCRLX.lexicons = ko.observableArray();
    DMCRLX.selectedLexicon = ko.observable();
    
    DMCRLX.loadLexiconNames = function () {

	$.getJSON(DMCRLX.baseURL +"listlexicons")
	    .done(function (data) {
		DMCRLX.lexicons(data);
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);
	    });
    };
    
    DMCRLX.updateLexiconName = function () {
	console.log("EN SNÄLL APA ÄTER SPENAT: "+ JSON.stringify(DMCRLX.selectedLexicon()));
    }

    
    DMCRLX.addLexicon = function () {
	var newLex = {'id': 0, 'name': " ", 'symbolSetName': " "};
	//var newLex2 = {'id': 0, 'name': " ", 'symbolSetName': " "};
	console.log(JSON.stringify(newLex));
	DMCRLX.lexicons.push(newLex);
	DMCRLX.selectedLexicon(newLex);
    }
    


    // On pageload:
    DMCRLX.loadLexiconNames();
};

ko.applyBindings(new DMCRLX.CreateLexModel());
