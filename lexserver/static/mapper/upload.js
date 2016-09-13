var MAPPER = {};

MAPPER.baseURL = window.location.origin;

// From http://stackoverflow.com/a/8809472
MAPPER.generateUUID = function() {
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


MAPPER.UploadFileModel = function () {
    var self = this; 
    
    
    self.uuid = MAPPER.generateUUID();

    self.serverMessage = ko.observable("_");
    
    self.maxMessages = 10;
    self.serverMessages = ko.observableArray();
    
    self.connectWebSock = function() {
	var zock = new WebSocket(MAPPER.baseURL.replace("http://", "ws://") + "/websockreg" );
	zock.onopen = function() {
	    console.log("MAPPER.connectWebSock: sending uuid over zock: "+ self.uuid);
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
    
    self.message = ko.observable("");
    
    self.selectedFile = ko.observable(null);
    self.hasSelectedFile = ko.observable(false);   
    
    self.setSelectedFile = function(symbolsetFile) {
	self.selectedFile(symbolsetFile);
	console.log("selected file: ", self.selectedFile())
	self.hasSelectedFile(true);
    }

    self.uploadFile = function() {
	console.log("uploading file: ", self.selectedFile())
	var url = MAPPER.baseURL + "/mapper_do_upload"
	var xhr = new XMLHttpRequest();
	var fd = new FormData();
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function() {
            if (xhr.readyState === 4 && xhr.status === 200) {
		// Every thing ok, file uploaded
		console.log("uploadFile return response text", xhr.responseText); // handle response.
		//} //else { // TODO  This doesn't seem to be the right way to handle errors here
		//alert(xhr.responseText);
	    };// else { // TODO this doesn't work
 	//	console.log("uploadLexiconFile return status", xhr.statusText);
	 //   };
	};
	fd.append("client_uuid", self.uuid);
	fd.append("upload_file", self.selectedFile());
	xhr.send(fd);
	self.message("File sent!");
    };
    
};

var upload = new MAPPER.UploadFileModel();
ko.applyBindings(upload);
upload.connectWebSock();

console.log("UUID: "+ upload.uuid);
