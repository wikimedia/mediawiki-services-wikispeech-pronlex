namespace Entry {
    export class Transcription {
        readonly id: number;
        entryId: number;
        strn: string;
        language: string;
        sources: string[];
    }

    export class Lemma {
        readonly id: number;
        strn: string;
        reading: string;
        paradigm: string;
    }

    export class Status {
        readonly id: number;
        name: string;
        source: string;
        timestamp: string;
        current: boolean;
    }

    export class EntryValidation {
        readonly id: number;
        level: string;
        ruleName: string;
        message: string;
        timestamp: string;
    }

    export class Entry {
        constructor(strn: string) {
            this.strn = strn;
        }
        readonly id: number;
        readonly lexiconId: number;
        strn: string;
        language: string;
        partOfSpeech: string;
        morphology: string;
        wordParts: string;
        lemma: Lemma;
        transcriptions: Transcription[];
        status: Status;
        entryValidations: EntryValidation[];
    }

    function json2entry(json: JSON): Entry {
        return <Entry>(<any>json); // todo: proper parsing, error handling, etc
    }

    export function json2entries(json: JSON): Entry[] {
        return <Entry[]>(<any>json); // todo: proper parsing, error handling, etc
    }

    function entry2html(e: Entry): HTMLElement {
        let result = document.createElement("tr");

        for (let t of e.transcriptions) {
            result.innerHTML = `<td>${e.strn}</b></td><td>${e.wordParts}</td><td>${t.strn}</td>`;
        }

        return result;
    }

    function entries2html(entries: Entry[]): HTMLElement {
        let table = document.createElement("table");

        table.innerHTML = '<thead><tr><td><b>Strn</b></td><td><b>WordParts</b></td><td><b>Transcription</b></td></tr></thead>';

        let body = document.createElement("tbody");
        for (let entry of entries) {
            body.appendChild(entry2html(entry));
        }
        table.appendChild(body);

        return <HTMLElement>table;
    }
}
