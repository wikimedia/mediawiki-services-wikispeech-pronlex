var DMCRLX = {};

DMCRLX.baseURL = "http://localhost:8787/"

DMCRLX.CreateLexModel = function() {
    var self = this;
    
    DMCRLX.lexicons = ko.observableArray();
    
    DMCRLX.loadLexiconNames = function () {

	console.log("DGDGDGD");

	$.getJSON(DMCRLX.baseURL +"listlexicons")
	    .done(function (data) {
		DMCRLX.lexicons(data);
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert(xhr.responseText);
	    });
    };
    
    // On pageload:
    DMCRLX.loadLexiconNames();
};

ko.applyBindings(new DMCRLX.CreateLexModel());
