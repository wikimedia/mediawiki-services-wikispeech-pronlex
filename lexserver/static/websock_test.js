WebSockTest = function () {
    var self = this;
    
    self.message = ko.observable("_");
    
    self.connectWebSock = function() {
	var zock = new WebSocket("ws://localhost:8787/websocktick") 
	zock.onmessage = function(e) {
	    if(e.data === "WS_KEEPALIVE") {
		var d = new Date();
		var h = d.getHours();
		var m = d.getMinutes();
		var s = d.getSeconds();
		var msg = "Websocket keepalive recieved "+ h +":"+ m +":"+ s;
		self.message(msg);
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

    self.timeToWsChan = function() {
	var url = window.location.origin + "/timetochan"
	//console.log("URL: "+ url);
	$.get(url).done(function(data){
	    console.log("timeToChan respons: "+ data);
	})
	    .fail(function (xhr, textStatus, errorThrown) {
		console.log("timeToWsChan xhr.responseText: "+ xhr.responseText);
		console.log("timeToWsChan textStatus: "+ textStatus);
		console.log("timeToWsChan errorThrown: "+ errorThrown);
		//alert("timeToChan says: "+ xhr.responseText);
	    });
    };
};

var wsTst = new WebSockTest();
ko.applyBindings(wsTst);
wsTst.connectWebSock();
