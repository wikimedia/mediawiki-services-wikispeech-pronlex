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

    //zelf; // = this; // = null; // = this;
    //constructor() {
    // 	this.zelf = this;
	//zelf.getUsers = getUsers();
    //};
 
    
    //deleteUser;

    public itself: any = this;

    message: string = "";

    users: KnockoutObservableArray<User> = ko.observableArray<User>([]);

     constructor() {
     	this.itself = this;
     }
    
    public  getUsers = () => {
        //this.users.push( {name: "nils", roles: "thingy", dbs: "lexdb"})
        //this.users.push( {name: "nuls", roles: "thungy", dbs: "loxdb"} )
	
        console.log("GETTING USER LIST");
	
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
    
    deleteUser(user: User): void {
	//let f = itself.getUsers;
        console.log(user.name);
        let itself = this;
        let url = baseURL + "/admin/user_db/delete_user?name="+ user.name;
        let r = new XMLHttpRequest();
	r.open("GET", url);


	r.onload = () => {
	    if(r.status == 200) {
		//let x = new UserDB();
		//x.getUsers();
		itself.getUsers()
	    }
	}
	r.send();
    }
    
    addUser(): void {
        console.log("YEAH, new user");
    }
}

let udb = new UserDB();
udb.getUsers();
ko.applyBindings(udb);
