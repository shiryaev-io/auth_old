package queries

const (
	QuerySelectUserByEmail string = `SELECT * FROM users WHERE email=?;`
	QuerySelectUserById    string = `SELECT * FROM users WHERE id=?;`
)
