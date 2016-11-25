interface Decomp {
    //word: string;
    parts: string[];
}

const baseURL: string = window.location.origin;

function decomp(): void {
    let decompInputElem = <HTMLInputElement>document.getElementById("decomp_word");
    let word = decompInputElem.value

    let decomps = decomp0(baseURL, word);

    console.log(word);


    //return 

    //return decomp0(word);
}




function decomp0(baseURL: string, word: string): void { //Decomp[] {
    let res = [{ parts: [] }];

    console.log("BAseURL: " + baseURL);
    let url = baseURL + "/decomp/decomp?word=" + word
    console.log("url: " + url);


    let r = new XMLHttpRequest();
    r.open("GET", url, true); // Hur decinficerar man 'word'-strängen?'
    r.onload = function () {

        if (r.status === 200) {

            console.log(r.responseText);
            //insertDecompResult(r.responseText);
            //alert(r.responseText);
        } else {
            console.log("readyState: " + r.readyState);
            console.log("statusText: " + r.statusText);
        }

        //console.log("readyState: " + r.readyState);
        //console.log("status: " + r.status);
    }

    r.send();
    //r.setRequestHeader();


    //return res;
}
