var HEJ = {};

HEJ.ListKnownModel = function() {
    var self = this;
    self.freeText = ko.observable("");
    
    self.wordList = ko.computed( function () {
	var wds = self.freeText().trim().split(/[ ,;."'!?]+/);
	var freqs = _.countBy(_.map(wds, function(s) {return s.toLowerCase();}), _.identity());
	//var found = {};
	
	return _.pairs(freqs);
    }, this				 
			       ); 
}

ko.applyBindings(new HEJ.ListKnownModel());
console.log("VAD I HELA HELVETET?!");
