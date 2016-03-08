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
    
    DMCRLX.updateLexicon = function () {
	
	if ( DMCRLX.selectedLexicon().name === "" || DMCRLX.selectedLexicon().symbolSetName === "" ) {
	    alert("Name and Symbol set name field must not be empty")
	    return;
	}
	
	var params = {'id' : DMCRLX.selectedLexicon().id, 'name' : DMCRLX.selectedLexicon().name, 'symbolsetname' : DMCRLX.selectedLexicon().symbolSetName}
	
	
	$.get(DMCRLX.baseURL + "insertorupdatelexicon", params)
	    .done(function(data){
		DMCRLX.loadLexiconNames();
		DMCRLX.selectedLexicon(data);
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
    
    // On pageload:
    DMCRLX.loadLexiconNames();
};

ko.applyBindings(new DMCRLX.CreateLexModel());
