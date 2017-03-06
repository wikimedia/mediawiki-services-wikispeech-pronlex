/// <reference path="../node_modules/@types/knockout/index.d.ts"/>

interface Decomp {
    //word: string;
    parts: string[];
}

const baseURL: string = window.location.origin;

class NIZZE {

    langs: KnockoutObservableArray<string[]>;

    word: KnockoutObservable<string>;

    zmurf: KnockoutComputed<string>;

    decomps: KnockoutObservable<Decomp[]>;

    constructor() {
        this.word = <KnockoutObservable<string>>ko.observable("skrotbil");
        this.decomps = <KnockoutObservable<Decomp[]>>ko.observable([]);
        this.zmurf = ko.computed({
            read: () => {

                this.decomp0(baseURL, 'sv', this.word())
                return "";//this.word().toUpperCase();
            }
        });
        //this.word("bilskrot");
    }

    decomp0(baseURL: string, lang: string, word: string): string { //Decomp[] {
        let itself = this;
        let res: string = "";// [{ parts: [] }];

        console.log("BAseURL: " + baseURL);
        // TODO sanitize input data
        let url = baseURL + "/decomp/decomp?lang=" + lang + "&word=" + word
        console.log("url: " + url);


        let r = new XMLHttpRequest();
        r.open("GET", url, true); // Hur decinficerar man 'word'-strängen?'
        r.onload = function () {

            if (r.status === 200) {

                console.log(r.responseText);
                //let d: string[] = [r.responseText];
                let d: Decomp[] = <Decomp[]>JSON.parse(r.responseText);
                itself.decomps(d);
                res = r.responseText

            } else {
                console.log("readyState: " + r.readyState);
                console.log("statusText: " + r.statusText);
            }


        }

        r.send();
        return res;
    }

    decomp(): void {
        let decompInputElem = <HTMLInputElement>document.getElementById("decomp_word");
        let word = decompInputElem.value

        this.decomp0(baseURL, 'sv', word);

        console.log(word);


    }

}



let n = new NIZZE();

function decomp(): void {
    n.decomp();
}

ko.applyBindings(n);
