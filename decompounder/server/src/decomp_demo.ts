/// <reference path="../node_modules/@types/knockout/index.d.ts"/>

interface Decomp {
    //word: string;
    parts: string[];
}

const baseURL: string = window.location.origin;

class NIZZE {
    
    langs: KnockoutObservableArray<string>;
    selectedLang: KnockoutObservable<string>;
    
    word: KnockoutObservable<string>;

    addPrefix: KnockoutObservable<string>;
    addSuffix: KnockoutObservable<string>;
    
    zmurf: KnockoutComputed<string>;

    addPrefixMsg: KnockoutComputed<string>;
    addSuffixMsg: KnockoutComputed<string>;
    
    decomps: KnockoutObservable<Decomp[]>;
    
    constructor() {
	
        this.langs = <KnockoutObservableArray<string>>ko.observableArray([]);
	
	// TODO remove hard wired value
	this.selectedLang = <KnockoutObservable<string>>ko.observable("sv");

	// TODO remove hard wired value
	this.word = <KnockoutObservable<string>>ko.observable("");
	
	this.addPrefix = <KnockoutObservable<string>>ko.observable("");
	this.addSuffix = <KnockoutObservable<string>>ko.observable("");
	
	this.decomps = <KnockoutObservable<Decomp[]>>ko.observable([]);
	
	// Only here for side effect of calling decomp server
	this.zmurf = ko.computed({
            read: () => { // <-- What _is_ this...?!

		if (this.selectedLang() != "" && this.word() != "") { 
                    this.decomp0(baseURL, this.selectedLang(), this.word());
		}
                return "";
            }
        });

	// When addPrefix is updated, send this value to server
	this.addPrefixMsg = ko.computed({
	    read: () => { // <-- What _is_ this...?!
		let itself = this;
		let res = "";
		console.error("addPrefix: "+ itself.addPrefix());
		
		if (itself.addPrefix() != "" && itself.selectedLang() != "") {
		    
		    let r = new XMLHttpRequest();
		    let url = baseURL + "/decomp/add_prefix?lang="+ itself.selectedLang() + "&prefix="+ itself.addPrefix();
		    r.open("GET", url);
		    r.onload = function () {
			if (r.status === 200) {

			    console.log("added prefix: "+ r.responseText);
			    res = r.responseText;
			    //return r.responseText;
			    
			    
			} else {
			    //TODO error handling
			    console.log("failed to add prefix: "+ r.responseText);
			    //return "";
			    res = r.responseText;
			};
		    };
		    r.send();
		};
		
		return res;
	    }
	});


	// When addSuffix is updated, send this value to server
	this.addSuffixMsg = ko.computed({
	    read: () => { // <-- What _is_ this...?!
		let itself = this;
		let res = "";
		console.error("addSuffix: "+ itself.addSuffix());
		
		if (itself.addSuffix() != "" && itself.selectedLang() != "") {
		    
		    let r = new XMLHttpRequest();
		    let url = baseURL + "/decomp/add_suffix?lang="+ itself.selectedLang() + "&suffix="+ itself.addSuffix();
		    r.open("GET", url);
		    r.onload = function () {
			if (r.status === 200) {

			    console.log("added suffix: "+ r.responseText);
			    res = r.responseText;
			    //return r.responseText;
			    
			    
			} else {
			    //TODO error handling
			    console.log("failed to add suffix: "+ r.responseText);
			    //return "";
			    res = r.responseText;
			};
		    };
		    r.send();
		};
		
		return res;
	    }
	});

	
    }

    
    loadLangs = () => {
	
		
	let itself = this;
        let r = new XMLHttpRequest();
        let url = baseURL + "/decomp/list_languages";
        r.open("GET", url);
	r.onload = function () {
	    if (r.status === 200) {
		// TODO how to handle error?		
		let ls: string[] = <string[]>JSON.parse(r.responseText);
		itself.langs.removeAll();
		for(let i =0; i<ls.length; i++) {
		    itself.langs.push(ls[i]);
		};

		
		
	    } else {
		//TODO error handling
		console.log("failed to load languages: "+ r.statusText);
	    };
	};
        r.send();
	//return;
    };

    decomp0(baseURL: string, lang: string, word: string): string { //Decomp[] {
        let itself = this;
        let res: string = "";// [{ parts: [] }];

	if ("" === lang) {
	    console.error("decomp0 called without lang value");
	    return res;
	};
	if ("" === word) {
	    console.error("decomp0 called without word value");
	    return res;
	};

	
        // TODO sanitize input data
        let url = baseURL + "/decomp/decomp?lang=" + lang + "&word=" + word

        let r = new XMLHttpRequest();
        r.open("GET", url, true); // Hur decinficerar man 'word'-strängen?'
        r.onload = function () {
	    
            if (r.status === 200) {
		
                let d: Decomp[] = <Decomp[]>JSON.parse(r.responseText);
                itself.decomps(d);
                res = r.responseText
		
            } else { //TODO error handling
                console.log("readyState: " + r.readyState);
                console.log("statusText: " + r.statusText);
            }
        }
	
        r.send();
        return res;
    }
    
    // decomp(): void {
    //     let decompInputElem = <HTMLInputElement>document.getElementById("decomp_word");
    //     let word = decompInputElem.value
	
    //     this.decomp0(baseURL, this.selectedLang(), word);
			
    // }

}



let n = new NIZZE();

// function decomp(): void {
//     n.decomp();
// }

ko.applyBindings(n);
n.loadLangs();
//TODO remove:
n.word("skrotbil");
