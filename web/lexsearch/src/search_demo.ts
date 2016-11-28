class SearchModel {

    lexicons: KnockoutObservable<string>
    words: KnockoutObservable<string>
    result: KnockoutObservableArray<Entry.Entry>
    
    constructor(lexicons: string, words: string) {
        this.lexicons = ko.observable(lexicons);
        this.words = ko.observable(words);
	this.result = ko.observableArray([]);
    }

    search(): void { 
	let base_url  = window.location.origin;
	let itself = this;
	if (itself.words().length === 0) {
            alert("Cannot search for empty string!");
            return;
	};
	if (itself.lexicons().length === 0) {
            alert("Lexicon(s) not specified, cannot search!");
            return;
	};
	let r = new XMLHttpRequest();
	let url = base_url + "/lexicon/lookup?lexicons=" + itself.lexicons() + "&words=" + itself.words(); // TODO: urlencode
	//console.log("search_demo URL=" + url);
	r.open("GET", url);
	r.onload = function () {
            if (r.status === 200) {
		let json: JSON = JSON.parse(r.responseText);
		console.log(JSON.stringify(json));
		itself.result(Entry.json2entries(json));
            }
            else {
		alert("ERROR\n" + r.status + "\n" + r.responseText);
            };
	};
	r.send();
    }

}

ko.applyBindings(new SearchModel("sv-se.nst", "hund, häst"));
