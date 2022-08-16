package queries

const (
	QuerySelectTokenByUserId       string = `SELECT * FROM tokens WHERE userId=?;`
	QuerySelectTokenByRefreshToken string = `SELECT * FROM tokens WHERE token=?;`
	QueryUpdateToken               string = `UPDATE tokens SET token=? WHERE Id=?;`
	QueryInsertToken               string = `INSERT INTO tokens (userId, token) VALUES (?, ?);`
	QueryDeleteToken               string = `DELETE FROM tokens WHERE token=?;`
)
