const baseURL: string = window.location.origin;	

interface User {

    name: string;
    roles: string;
    dbs: string;

}


// class User {

//     name: string;
//     roles: string;
//     dbs: string;

//     constructor(name: string, roles: string, dbs: string) {
//         this.name = name;
//         this.roles = roles;
//         this.dbs = dbs;
//     }
// }

class UserDB {
    //nisse = "NIZZE"

    users: KnockoutObservableArray<User> = ko.observableArray<User>([]);

    getUsers(): void {
        //this.users.push( {name: "nils", roles: "thingy", dbs: "lexdb"})
        //this.users.push( {name: "nuls", roles: "thungy", dbs: "loxdb"} )

	let itself = this;
	let url = baseURL + "/admin/user_db/list_users"
	let r = new XMLHttpRequest();
	console.log(url);
	r.open("GET", url);
	r.onload = function () {
	    if (r.status == 200) {
		// TODO How do you handle errors?
		let u: User[] = <User[]>JSON.parse(r.responseText); 
		itself.users(u);
		
	    } else {
                console.log("readyState: " + r.readyState);
                console.log("statusText: " + r.statusText);
            }
	}
	r.send();
    }

}

let udb = new UserDB();
udb.getUsers();
ko.applyBindings(udb);
