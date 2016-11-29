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
                return (this.selectedLexicon() != null && this.words().length > 0);
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

    ipa(transcription: string): string {
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
        let url = base_url + "/mapper/map?from=" + encodeURIComponent(itself.selectedLexicon().symbolSetName) + "&to=ipa&trans=" + encodeURIComponent(transcription);
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
        if (itself.words().length === 0) {
            alert("Cannot search for empty string!");
            return;
        };
        if (itself.selectedLexicon() == null) {
            alert("Lexicon(s) not specified, cannot search!");
            return;
        };
        let r = new XMLHttpRequest();
        let url = base_url + "/lexicon/lookup?lexicons=" + itself.selectedLexicon().name + "&words=" + encodeURIComponent(itself.words()); // TODO: urlencode
        r.open("GET", url);
        r.onload = function () {
            if (r.status === 200) {
                let json: JSON = JSON.parse(r.responseText);
                //console.log(JSON.stringify(json));
                itself.entries(Entry.json2entries(json));
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            };
        };
        r.send();
    }

}

let model = new SearchModel("sv-se.nst", "hund, häst")
ko.applyBindings(model);
