class User {
    name: string;
    roles: string;
    dbs: string;

    constructor(name: string, roles: string, dbs: string) {
        this.name = name;
        this.roles = roles;
        this.dbs = dbs;
    }
}

class UserDB {
    //nisse = "NIZZE"

    users: KnockoutObservableArray<User> = ko.observableArray<User>([]);

    getUsers(): void {
        this.users.push(new User("nils", "thingy", "lexdb"))
        this.users.push(new User("nuls", "thungy", "luxdb"))
    }

}

let udb = new UserDB();
udb.getUsers();
ko.applyBindings(udb);
