/// <reference path="../node_modules/@types/knockout/index.d.ts"/>

interface Decomp {
    //word: string;
    parts: string[];
}







const baseURL: string = window.location.origin;


class NIZZE {

    decomps: KnockoutObservable<string>

    constructor() {
        this.decomps = <KnockoutObservable<string>>ko.observable();
    }

    decomp0(baseURL: string, word: string): void { //Decomp[] {
        let itself = this;
        let res = [{ parts: [] }];

        console.log("BAseURL: " + baseURL);
        let url = baseURL + "/decomp/decomp?word=" + word
        console.log("url: " + url);


        let r = new XMLHttpRequest();
        r.open("GET", url, true); // Hur decinficerar man 'word'-strängen?'
        r.onload = function () {

            if (r.status === 200) {

                console.log(r.responseText);
                itself.decomps(r.responseText);

            } else {
                console.log("readyState: " + r.readyState);
                console.log("statusText: " + r.statusText);
            }


        }

        r.send();

    }

    decomp(): void {
        let decompInputElem = <HTMLInputElement>document.getElementById("decomp_word");
        let word = decompInputElem.value

        this.decomp0(baseURL, word);

        console.log(word);


    }

}



let n = new NIZZE();

function decomp(): void {
    n.decomp();
}

ko.applyBindings(n);
