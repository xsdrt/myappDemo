package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"net/http"
	"strings"
	"time"

	up "github.com/upper/db/v4"
)

type Token struct {
	ID        int       `db:"id" json:"id"`
	UserId    int       `db:"user_id" json:"user_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	Email     string    `db:"email" json:"email"`
	PlainText string    `db:"token" json:"token"`
	Hash      []byte    `db:"token_hash" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Expires   time.Time `db:"expiry" json:"expiry"`
}

func (t *Token) Table() string { // func to get the token in the user.go
	return "tokens"
}

// Get a user by a token...will need  when authenticating thru an API...
func (t *Token) GetUserForToken(token string) (*User, error) {
	var u User
	var theToken Token

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": token})
	err := res.One(&theToken)
	if err != nil {
		return nil, err
	}

	collection = upper.Collection("users")
	res = collection.Find(up.Cond{"id": theToken.UserId})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}

	u.Token = theToken

	return &u, nil

}

// Get all of a given users tokens...
func (t *Token) GetTokensForUser(id int) ([]*Token, error) { // Return a slice of pointers to Token...
	var tokens []*Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id": id})
	err := res.All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Get a Token by Id...
func (t *Token) Get(id int) (*Token, error) {
	var token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id": id})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Get by a Token itself...
func (t *Token) GetByToken(plainText string) (*Token, error) {
	var token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": plainText})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Delete tokens by id Ie.. user logs out....
func (t *Token) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

// Delete tokens by plain text version of token Ie.. user logs out....
func (t *Token) DeleteByToken(plainText string) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": plainText})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

// Insert tokens...
func (t *Token) Insert(token Token, u User) error {
	collection := upper.Collection(t.Table())

	// delete existing tokens (if we are inserting a token, why have unneeded sessions/tokens hanging around...)
	res := collection.Find(up.Cond{"user_id": u.ID})
	err := res.Delete()
	if err != nil {
		return err
	}

	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.FirstName = u.FirstName
	token.Email = u.Email

	_, err = collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}

// Generate the tokens to be inserted...
func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserId:  userID,
		Expires: time.Now().Add(ttl), //Called expiry in the DB, we are adding the time to live to time.Now...
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes) //So , should give us the token with exact number of characters every time...
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

// Authenticate the token...
func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ") // Expect to find an authorizationHeader consists with the word Bearer, fo/lowed
	// by a blank space, followed by Plaintext version of the token the user is going to authenticate with...
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]

	if len(token) != 26 { //Should always be 26 characters...
		return nil, errors.New("token is the wrong size")
	}

	// get the token from the database...
	tkn, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("no matching token found")
	}

	if tkn.Expires.Before(time.Now()) { //Check to see if the token is expired , because if so, throw an error.
		return nil, errors.New("expired token")
	}
	user, err := t.GetUserForToken(token)
	if err != nil { //make sure there is a user...
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

// Validate the token
func (t *Token) ValidToken(token string) (bool, error) {
	user, err := t.GetUserForToken(token)
	if err != nil { //make sure there is a user...
		return false, errors.New("no matching user found")
	}

	if user.Token.PlainText == "" { // make sur not an empty token...
		return false, errors.New("no matching token found")
	}

	if user.Token.Expires.Before(time.Now()) { //Check to see if the token is expired , because if so, throw an error.
		return false, errors.New("expired token")
	}

	return true, nil
}
