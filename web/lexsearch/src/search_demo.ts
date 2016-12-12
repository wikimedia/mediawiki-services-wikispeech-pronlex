class Lexicon {
    name: string
    symbolSetName: string
}

class SearchModel {

    selectedLexicon: KnockoutObservable<Lexicon>
    availableLexicons: KnockoutObservableArray<Lexicon>
    words: KnockoutObservable<string>
    validForm: KnockoutComputed<boolean>
    entries: KnockoutObservableArray<Entry.Entry>

    constructor(lexicon: string, words: string) {
        this.availableLexicons = ko.observableArray([]);
        this.loadLexicons(lexicon);
        this.selectedLexicon = ko.observable(null);
        this.words = ko.observable(words);
        this.entries = ko.observableArray([]);
        this.validForm = ko.computed({
            read: () => {
                return (this.selectedLexicon() != null && this.words() != null && this.words().length > 0);
            }
        });
    }

    loadLexicons(initialChoice: string): void {
        let itself = this;
        let url = window.location.origin + "/lexicon/list";
        let r = new XMLHttpRequest();
        r.open("GET", url);
        r.onload = function () {
            if (r.status === 200) {
                let lexicons = <Lexicon[]>JSON.parse(r.responseText);
                itself.availableLexicons(lexicons);
                for (let lex of itself.availableLexicons()) {
                    if (lex.name == initialChoice) {
                        itself.selectedLexicon(lex);
                    }
                }
                //console.log(JSON.stringify(lexicons));
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            };
        };
        r.send();
    };

    ipa(transcription: string, symbolSetName: string): string {
        console.log(symbolSetName);
        let base_url = window.location.origin;
        let itself = this;
        let ipa = ""
        if (itself.words().length === 0) {
            alert("Cannot search for empty string!");
            return;
        };
        if (itself.selectedLexicon() == null) {
            alert("Lexicon(s) not specified, cannot search!");
            return;
        };
        let r = new XMLHttpRequest();
        let url = base_url + "/mapper/map?from=" + encodeURIComponent(symbolSetName) + "&to=ipa&trans=" + encodeURIComponent(transcription);
        r.open("GET", url, false);
        r.onload = function () {
            if (r.status === 200) {
                ipa = JSON.parse(r.responseText).Result;
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            };
        };
        r.send();
        return ipa;
    }


    search(): void {
        let base_url = window.location.origin;
        let itself = this;
        itself.entries([]); // clear previous search
        if (itself.words().length === 0) {
            alert("Cannot search for empty string!");
            return;
        };
        if (itself.selectedLexicon() == null) {
            alert("Lexicon(s) not specified, cannot search!");
            return;
        };
        let r = new XMLHttpRequest();
        let url = base_url + "/lexicon/lookup?lexicons=" + encodeURIComponent(itself.selectedLexicon().name) + "&words=" + encodeURIComponent(itself.words());
        r.open("GET", url);
        r.onload = function () {
            if (r.status === 200) {
                let json: JSON = JSON.parse(r.responseText);
                if (json != null) {
                    let tmp = Entry.json2entries(json);
                    for (let e of tmp) {
                        for (let t of e.transcriptions) {
                            t.symbolSetName = itself.selectedLexicon().symbolSetName;
                        }
                    }
                    itself.entries(tmp);
                }
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            };
        };
        r.send();
    }

}

let model = new SearchModel("sv-se.nst", "hund, häst")
let x = model.availableLexicons();
ko.applyBindings(model);
