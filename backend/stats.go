package main

import "github.com/gin-gonic/gin"

func handleStats(c *gin.Context) {
	var s Stats
	db.Get(&s, `SELECT
		(SELECT COUNT(*) FROM books) AS total_books,
		(SELECT COALESCE(SUM(total_copies),0) FROM books) AS total_copies,
		(SELECT COALESCE(SUM(available_copies),0) FROM books) AS available_copies,
		(SELECT COUNT(*) FROM users) AS total_users,
		(SELECT COUNT(*) FROM borrow_records WHERE return_date IS NULL) AS active_borrows,
		(SELECT COUNT(*) FROM borrow_records WHERE return_date IS NULL AND due_date < datetime('now')) AS overdue_borrows`)
	var cats []CategoryCount
	db.Select(&cats, "SELECT category, COUNT(*) as count FROM books GROUP BY category ORDER BY count DESC, category ASC")
	s.Categories = cats
	ok(c, s)
}
