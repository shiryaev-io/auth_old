package queries

const (
	QuerySelectTokenByUserId       string = `SELECT * FROM tokens WHERE user_id=$1;`
	QuerySelectTokenByRefreshToken string = `SELECT * FROM tokens WHERE token=$1;`
	QueryUpdateToken               string = `UPDATE tokens SET token=$1 WHERE id=$2;`
	QueryInsertToken               string = `INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING id;`
	QueryDeleteToken               string = `DELETE FROM tokens WHERE token=$1;`
)
