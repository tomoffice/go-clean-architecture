package mcsqlite

const (
	queryInsertMember         = `INSERT INTO members (name, email, password) VALUES (?, ?, ?)`
	querySelectByID           = `SELECT * FROM members WHERE id = ?`
	querySelectByEmail        = `SELECT * FROM members WHERE email = ?`
	querySelectAllBase        = `SELECT * FROM members ORDER BY %s %s LIMIT ? OFFSET ?`
	queryUpdateMemberProfile  = `UPDATE members SET name = ? WHERE id = ?`
	queryUpdateMemberEmail    = `UPDATE members SET email = ? WHERE id = ?`
	queryUpdateMemberPassword = `UPDATE members SET password = ? WHERE id = ?`
	queryDeleteMember         = `DELETE FROM members WHERE id = ?`
	queryCountMembers         = `SELECT COUNT(*) FROM members`
)
