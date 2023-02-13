package data

import (
	"time"

	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct { //This will be exported add the fields in/from the database hiSpeed (testing)
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
}

func (u *User) Table() string { // func available every time there is a type User; this is a means to override the table name...
	return "users"
}

func (u *User) GetAll() ([]*User, error) { //Passing the condition to this func "GetAll Users" from the DB using a slice of a pointer to user... "
	collection := upper.Collection(u.Table()) //Using upper/db conventions; things stored in a database are called Collections...
	var all []*User                           //Store the info in a variable "all"...

	res := collection.Find().OrderBy("last_name") //get the user results for the condition with what ever parameters were passed to the func and order by their last name...
	err := res.All(&all)                          //Read the results into the "all" var....
	if err != nil {
		return nil, err // return the err if one, otherwise
	}

	return all, nil

}

func (u *User) GetByEmail(email string) (*User, error) { //Look the users up by email...
	var theUser User //Store the data in var "theUser"
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email =": email}) // Ok find the specified condition "email" so all records in the db where the email is == to email supplied to the func...
	err := res.One(&theUser)                          // Only return 1 user per email (usually; would have to change if allows multiple emails per user...) and read it into theUser var...
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id = ": theUser.ID, "expiry <": time.Now()}).OrderBy("created_at desc") // get the most recent token if any...
	err = res.One(&token)                                                                                       // Read into token var...
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows { //Make sure if valid user but no token still show the user without a token value...
			return nil, err
		}
	}

	theUser.Token = token // Read into the User even with an empty token, so we can still return a valid user with an empty token...

	return &theUser, nil //so the pointer to the user and no error
}

func (u *User) Get(id int) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id =": id}) //Find a user by their ID...
	err := res.One(&theUser)                    // Should only be one user per ID so read it into the var "theUser"...
	if err != nil {
		return nil, err
	}

	var token Token //Get the tokens, if any, and check for expired tokens...
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id = ": theUser.ID, "expiry <": time.Now()}).OrderBy("created_at desc") // get the most recent token if any...
	err = res.One(&token)                                                                                       // Read into token var...
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows { //Make sure if valid user but no token still show the user without a token value...
			return nil, err
		}
	}

	theUser.Token = token // Read into the User even with an empty token, so we can still return a valid user with an empty token...

	return &theUser, nil

}

func (u *User) Update(theUser User) error { //Update the user record
	theUser.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table()) //Again, get the collection....
	res := collection.Find(theUser.ID)        //Since updating an existing user, ID should be there...
	err := res.Update(&theUser)               // Update using a reference to "theUser" which received thru the call to this func...
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert(theUser User) (int, error) { //should return an int of the new record or an error...
	newHash, err := bcrypt.GenerateFromPassword([]byte(theUser.Password), 12)
	if err != nil {
		return 0, err
	}

	theUser.CreatedAt = time.Now()
	theUser.UpdatedAt = time.Now()
	theUser.Password = string(newHash)

	collection := upper.Collection(u.Table())
	res, err := collection.Insert(theUser)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
