package sqlx

const (
	queryInsertMember  = `INSERT INTO members (name, email, password) VALUES (?, ?, ?)`
	querySelectByID    = `SELECT * FROM members WHERE id = ?`
	querySelectByEmail = `SELECT * FROM members WHERE email = ?`
	querySelectAllBase = `SELECT * FROM members ORDER BY %s %s LIMIT ? OFFSET ?`
	queryUpdateMember  = `UPDATE members SET name = ?, email = ?, password = ? WHERE id = ?`
	queryDeleteMember  = `DELETE FROM members WHERE id = ?`
)
