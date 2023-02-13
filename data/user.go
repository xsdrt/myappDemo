package data

import (
	"time"

	up "github.com/upper/db/v4"
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

func (u *User) GetAll(condition up.Cond) ([]*User, error) { //Passing the condition to this func "GetAll Users from the DB using a slice of a pointer to user... "
	collection := upper.Collection(u.Table()) //Using upper/db conventions; things stored in a database are called Collections...
	var all []*User                           //Store the info in a variable "all"...

	res := collection.Find(condition) //get the results for the condition with what ever parameters were passed to the func...
	err := res.All(&all)              //Read the results into the "all" var....
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
