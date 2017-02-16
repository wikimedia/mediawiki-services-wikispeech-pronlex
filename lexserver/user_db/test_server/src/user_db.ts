// TODO list possible roles
// TODO list available DBs

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

    
    // These are knockout.js variables tied to HTML input fields,
    // together making up a new user
    newUserName:      KnockoutObservable<string> = ko.observable("");
    newUserPassword:  KnockoutObservable<string> = ko.observable("");
    newUserRoles:     KnockoutObservable<string> = ko.observable("");
    newUserDBs:       KnockoutObservable<string> = ko.observable("");

    //message: string = "";

    
    // List of the users in a knockout.js variable, tied to a HTML
    // table.  The contents of the list is obtained from the DB, by
    // calling this.getUsers() below.

    users: KnockoutObservableArray<User> = ko.observableArray<User>([]);

    // The side effect of this call is to fill in this.users.
    getUsers(): void {

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

    };
    

    // Deletes a user from the database. Forever gone.
    deleteUser = (user: User) => { 
	// Using the arrow syntax above, "this" appears to work in a
	// slightly more sane way, and referes to the embedding class
	// (UserDB).
	let zelf = this;
	

	
	
        let baseURL = window.location.origin;
	
	// TODO Error check user iput
	// TODO Sanitize user input
        let url = baseURL + "/admin/user_db/delete_user?name=" + user.name;
        console.log(user.name);

        let r = new XMLHttpRequest();
        r.open("GET", url);
        r.onload = function () {
            if (r.status === 200) {
                //userDB.getUsers();
		zelf.getUsers();
            }
            else {
		// TODO Better error handling. For example by writing
		// to a message text area.
                alert("ERROR\n" + r.status + "\n" + r.responseText);
            };
        };
        r.send();
    };

    addUser = (): void => {
	
	let zelf = this;
	
	// TODO Error check user iput
	let newUser = {"name": this.newUserName(), "password": this.newUserPassword(), "roles": this.newUserRoles(), "dbs": this.newUserDBs()};
	
        console.error("Adding new user " + JSON.stringify(newUser));
	
	let baseURL = window.location.origin;
	
	
	// TODO Sanitize user input
        let url = baseURL + "/admin/user_db/add_user?name=" + newUser.name +
	    "&password="+ newUser.password +
	    "&roles="+ newUser.roles + "&dbs="+ newUser.dbs;
        
	let r = new XMLHttpRequest();
	r.open("GET", url);
	r.onload = function() {
	    if (r.status === 200) {
		console.error("Added user '"+ newUser.name + "'");
		zelf.getUsers();
	    } else {
		alert("ERROR\n"+ r.status + "\n"+ r.responseText);
	    };
	};
	r.send();
    };
}

let udb = new UserDB();
udb.getUsers();
ko.applyBindings(udb);

