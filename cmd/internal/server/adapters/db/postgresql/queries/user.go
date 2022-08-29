package queries

const (
	QuerySelectUserByEmail string = `SELECT * FROM users WHERE email=$1;`
	QuerySelectUserById    string = `SELECT * FROM users WHERE id=$1;`
)
