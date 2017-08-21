var Lexicon = (function () {
    function Lexicon() {
    }
    return Lexicon;
}());
var SearchModel = (function () {
    function SearchModel(lexicon, words) {
        var _this = this;
        this.availableLexicons = ko.observableArray([]);
        this.loadLexicons(lexicon);
        this.selectedLexicon = ko.observable(null);
        this.words = ko.observable(words);
        this.entries = ko.observableArray([]);
        this.validForm = ko.computed({
            read: function () {
                return (_this.selectedLexicon() != null && _this.words() != null && _this.words().length > 0);
            }
        });
    }
    SearchModel.prototype.loadLexicons = function (initialChoice) {
        var itself = this;
        var url = window.location.origin + "/lexicon/list";
        var r = new XMLHttpRequest();
        r.open("GET", url);
        r.onload = function () {
            if (r.status === 200) {
                var lexicons = JSON.parse(r.responseText);
                itself.availableLexicons(lexicons);
                for (var _i = 0, _a = itself.availableLexicons(); _i < _a.length; _i++) {
                    var lex = _a[_i];
                    if (lex.name == initialChoice) {
                        itself.selectedLexicon(lex);
                    }
                }
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            }
            ;
        };
        r.send();
    };
    ;
    SearchModel.prototype.ipa = function (transcription, symbolSetName) {
        console.log(symbolSetName);
        var base_url = window.location.origin;
        var itself = this;
        var ipa = "";
        if (itself.words().length === 0) {
            alert("Cannot search for empty string!");
            return;
        }
        ;
        if (itself.selectedLexicon() == null) {
            alert("Lexicon(s) not specified, cannot search!");
            return;
        }
        ;
        var r = new XMLHttpRequest();
        //let url = base_url + "/mapper/map?from=" + encodeURIComponent(symbolSetName) + "&to=ipa&trans=" + encodeURIComponent(transcription);
        var url = base_url + "/mapper/map/" + encodeURIComponent(symbolSetName) + "/ipa/" + encodeURIComponent(transcription);
        r.open("GET", url, false);
        r.onload = function () {
            if (r.status === 200) {
                ipa = JSON.parse(r.responseText).Result;
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            }
            ;
        };
        r.send();
        return ipa;
    };
    SearchModel.prototype.search = function () {
        var base_url = window.location.origin;
        var itself = this;
        itself.entries([]); // clear previous search
        if (itself.words().length === 0) {
            alert("Cannot search for empty string!");
            return;
        }
        ;
        if (itself.selectedLexicon() == null) {
            alert("Lexicon(s) not specified, cannot search!");
            return;
        }
        ;
        var r = new XMLHttpRequest();
        var url = base_url + "/lexicon/lookup?lexicons=" + encodeURIComponent(itself.selectedLexicon().name) + "&words=" + encodeURIComponent(itself.words());
        r.open("GET", url);
        r.onload = function () {
            if (r.status === 200) {
                var json = JSON.parse(r.responseText);
                if (json != null) {
                    var tmp = Entry.json2entries(json);
                    for (var _i = 0, tmp_1 = tmp; _i < tmp_1.length; _i++) {
                        var e = tmp_1[_i];
                        for (var _a = 0, _b = e.transcriptions; _a < _b.length; _a++) {
                            var t = _b[_a];
                            t.symbolSetName = itself.selectedLexicon().symbolSetName;
                        }
                    }
                    itself.entries(tmp);
                }
            }
            else {
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            }
            ;
        };
        r.send();
    };
    return SearchModel;
}());
var model = new SearchModel("sv-se.nst", "hund, häst");
var x = model.availableLexicons();
ko.applyBindings(model);
var Entry;
(function (Entry_1) {
    var Transcription = (function () {
        function Transcription() {
        }
        return Transcription;
    }());
    Entry_1.Transcription = Transcription;
    var Lemma = (function () {
        function Lemma() {
        }
        return Lemma;
    }());
    Entry_1.Lemma = Lemma;
    var Status = (function () {
        function Status() {
        }
        return Status;
    }());
    Entry_1.Status = Status;
    var EntryValidation = (function () {
        function EntryValidation() {
        }
        return EntryValidation;
    }());
    Entry_1.EntryValidation = EntryValidation;
    var Entry = (function () {
        function Entry(strn) {
            this.strn = strn;
        }
        return Entry;
    }());
    Entry_1.Entry = Entry;
    function json2entry(json) {
        return json; // todo: proper parsing, error handling, etc
    }
    function json2entries(json) {
        return json; // todo: proper parsing, error handling, etc
    }
    Entry_1.json2entries = json2entries;
    function entry2html(e) {
        var result = document.createElement("tr");
        for (var _i = 0, _a = e.transcriptions; _i < _a.length; _i++) {
            var t = _a[_i];
            result.innerHTML = "<td>" + e.strn + "</b></td><td>" + e.wordParts + "</td><td>" + t.strn + "</td>";
        }
        return result;
    }
    function entries2html(entries) {
        var table = document.createElement("table");
        table.innerHTML = '<thead><tr><td><b>Strn</b></td><td><b>WordParts</b></td><td><b>Transcription</b></td></tr></thead>';
        var body = document.createElement("tbody");
        for (var _i = 0, entries_1 = entries; _i < entries_1.length; _i++) {
            var entry = entries_1[_i];
            body.appendChild(entry2html(entry));
        }
        table.appendChild(body);
        return table;
    }
})(Entry || (Entry = {}));
//# sourceMappingURL=search_demo.js.map