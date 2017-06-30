var LEXIMPORT = {};

LEXIMPORT.baseURL = window.location.origin;

// From http://stackoverflow.com/a/8809472
LEXIMPORT.generateUUID = function() {
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


LEXIMPORT.ImportFileModel = function () {
    var self = this; 
    
    
    self.uuid = LEXIMPORT.generateUUID();

    self.message = ko.observable("_");
    
    self.connectWebSock = function() {
	var zock = new WebSocket(LEXIMPORT.baseURL.replace("http://", "ws://") + "/websockreg" );
	zock.onopen = function() {
	    console.log("LEXIMPORT.connectWebSock: sending uuid over zock: "+ self.uuid);
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
    
    self.hasValidator = ko.observable(false);
    self.validate = ko.observable(true);
    self.lexiconName = ko.observable(null);
    self.lexiconDB = ko.observable(null);
    self.symbolSetName = ko.observable(null);
    self.selectedFile = ko.observable(null);
    self.validForm = ko.computed(function() {
	return (self.lexiconName() != null && self.symbolSetName() != null && self.selectedFile() != null &&
		self.lexiconName().trim() != "" && self.symbolSetName().trim() != "" &&
		self.lexiconDB() != null && self.lexiconDB().trim() != "");
    });
    
    self.setSelectedFile = function(lexiconFile) {
	self.selectedFile(lexiconFile);
	console.log("selected file: ", self.selectedFile())
    }
    
    self.symbolSetNames = ko.observableArray();
    self.lexiconDBs = ko.observableArray();

    self.loadSymbolSetNames = function () {
	$.getJSON(LEXIMPORT.baseURL +"/symbolset/list")
	    .done(function (data) {
		self.symbolSetNames(data.symbol_set_names);
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert("loadSymbolSetNames says: "+ xhr.responseText);
	    });
    };
    
    self.loadLexiconDBs = function () {
	$.getJSON(LEXIMPORT.baseURL +"/admin/list_dbs")
	    .done(function (data) {
		console.log(data)
		self.lexiconDBs(data);
	    })
    	    .fail(function (xhr, textStatus, errorThrown) {
		alert("loadLexiconDBs says: "+ xhr.responseText);
	    });
    };
    
    self.hasValidatorFunc = ko.computed(function() {
	var url = LEXIMPORT.baseURL + "/validation/has_validator/" + self.SymbolSetName;
	var xhr = new XMLHttpRequest();
	var fd = new FormData();
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function() {
            if (xhr.readyState === 4 && xhr.status === 200) {
		var response = xhr.responseText;
		console.log("hasValidator returned response text ", response);
		var hasV = response == 'true';
		self.validate(hasV);
		self.hasValidator(hasV);
		self.message("");
	    } else {
		self.message("Request failed: " + xhr.responseText + " " + xhr.readyState);
	    };
	};
	fd.append("client_uuid", self.uuid);
	fd.append("symbolset", self.symbolSetName());
	xhr.send(fd);	    
    });

    
    self.importLexiconFile = function() {
	console.log("uploading file: ", self.selectedFile())
	var url = LEXIMPORT.baseURL + "/admin/lex_import"
	var xhr = new XMLHttpRequest();
	var fd = new FormData();
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function() {
            if (xhr.readyState === 4 && xhr.status === 200) {
		// Every thing ok, file uploaded
		console.log("importLexiconFile returned response text ", xhr.responseText); // handle response.
		self.message("Import completed without errors: " + xhr.responseText);
	    } else {
		self.message("Import failed: " + xhr.responseText);
	    };
	};
	fd.append("client_uuid", self.uuid);
	fd.append("symbolset_name", self.symbolSetName());
	fd.append("lexicon_name", self.lexiconDB() + ":" + self.lexiconName());
	fd.append("validate", self.validate());
	fd.append("file", self.selectedFile());
	self.message("Importing, please wait ...");
	xhr.send(fd);
    };
    
};

var upload = new LEXIMPORT.ImportFileModel();
upload.loadSymbolSetNames();
upload.loadLexiconDBs();
ko.applyBindings(upload);
upload.connectWebSock();

console.log("UUID: "+ upload.uuid);
