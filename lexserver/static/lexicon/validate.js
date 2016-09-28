var LEXVDATE = {};

LEXVDATE.baseURL = window.location.origin;

// From http://stackoverflow.com/a/8809472
LEXVDATE.generateUUID = function() {
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


LEXVDATE.VdateModel = function () {
    var self = this; 
    
    
    self.uuid = LEXVDATE.generateUUID();

    self.message = ko.observable("_");
    
    self.connectWebSock = function() {
	var zock = new WebSocket(LEXVDATE.baseURL.replace("http://", "ws://") + "/websockreg" );
	zock.onopen = function() {
	    console.log("LEXVDATE.connectWebSock: sending uuid over zock: "+ self.uuid);
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
		// self.message(msg);
	    }
	    else {
		//console.log("Websocket got: "+ e.data)
		self.message(e.data);
	    };
	};
	zock.onerror = function(e){
	    console.log("websocket error: " + e.data);
	};
	zock.onclose = function (e) {
	    console.log("websocket got close event: "+ e.code)
	};
    };
    
    self.selectedLexicon = ko.observable(null);
    self.validForm = ko.computed(function() {
	return (self.selectedLexicon() != null);
    });
    
    self.availableLexicons = ko.observableArray();

    self.loadLexicons = function () {
	$.getJSON(LEXVDATE.baseURL +"/lexicon/listlexicons")
	    .done(function (data) {
		self.availableLexicons(data);
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert("loadLexicons says: "+ xhr.responseText);
	    });
    };
    
    
    self.runValidation = function() {
	console.log("validating lexicon: ", self.selectedLexicon())
	var url = LEXVDATE.baseURL + "/lex_do_validate"
	var xhr = new XMLHttpRequest();
	var fd = new FormData();
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function() {
            if (xhr.readyState === 4 && xhr.status === 200) {
		// Every thing ok
		console.log("runValidation returned response text ", xhr.responseText);
		self.message("Validation completed without errors: " + xhr.responseText);
	    } else {
		self.message("Validation failed: " + xhr.responseText);
	    };
	};
	fd.append("client_uuid", self.uuid);
	fd.append("lexicon_name", self.selectedLexicon().name);
	self.message("Validating, please wait ...");
	xhr.send(fd);
    };
    
};

var vdate = new LEXVDATE.VdateModel();
vdate.loadLexicons();
ko.applyBindings(vdate);
vdate.connectWebSock();

console.log("UUID: "+ vdate.uuid);
