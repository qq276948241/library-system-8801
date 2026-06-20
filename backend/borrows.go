package main

import (
	"github.com/gin-gonic/gin"
)

type borrowJoinRow struct {
	ID         int64   `db:"id"`
	UserID     int64   `db:"user_id"`
	BookID     int64   `db:"book_id"`
	BorrowDate string  `db:"borrow_date"`
	DueDate    string  `db:"due_date"`
	ReturnDate *string `db:"return_date"`
	Status     string  `db:"status"`

	BookID2          *int64  `db:"b_id"`
	BookTitle        *string `db:"b_title"`
	BookAuthor       *string `db:"b_author"`
	BookISBN         *string `db:"b_isbn"`
	BookCategory     *string `db:"b_category"`
	BookDescription  *string `db:"b_description"`
	BookCoverColor   *string `db:"b_cover_color"`
	BookPublisher    *string `db:"b_publisher"`
	BookPublishedYear *int   `db:"b_published_year"`
	BookTotalCopies  *int    `db:"b_total_copies"`
	BookAvailCopies  *int    `db:"b_available_copies"`
	BookCreatedAt    *string `db:"b_created_at"`

	UserID2       *int64  `db:"u_id"`
	UserUsername  *string `db:"u_username"`
	UserName      *string `db:"u_name"`
	UserRole      *string `db:"u_role"`
	UserEmail     *string `db:"u_email"`
	UserCreatedAt *string `db:"u_created_at"`
}

func (row *borrowJoinRow) toBorrow() Borrow {
	b := Borrow{
		ID:         row.ID,
		UserID:     row.UserID,
		BookID:     row.BookID,
		BorrowDate: row.BorrowDate,
		DueDate:    row.DueDate,
		ReturnDate: row.ReturnDate,
		Status:     row.Status,
	}
	if row.BookID2 != nil {
		b.Book = &Book{
			ID:              *row.BookID2,
			Title:           safeStr(row.BookTitle),
			Author:          safeStr(row.BookAuthor),
			ISBN:            safeStr(row.BookISBN),
			Category:        safeStr(row.BookCategory),
			Description:     safeStr(row.BookDescription),
			CoverColor:      safeStr(row.BookCoverColor),
			Publisher:       safeStr(row.BookPublisher),
			PublishedYear:   safeInt(row.BookPublishedYear),
			TotalCopies:     safeInt(row.BookTotalCopies),
			AvailableCopies: safeInt(row.BookAvailCopies),
			CreatedAt:       safeStr(row.BookCreatedAt),
		}
	}
	if row.UserID2 != nil {
		b.User = &User{
			ID:        *row.UserID2,
			Username:  safeStr(row.UserUsername),
			Name:      safeStr(row.UserName),
			Role:      safeStr(row.UserRole),
			Email:     safeStr(row.UserEmail),
			CreatedAt: safeStr(row.UserCreatedAt),
		}
	}
	enrichBorrow(&b)
	return b
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

const borrowJoinSQL = `
SELECT
  br.id, br.user_id, br.book_id, br.borrow_date, br.due_date, br.return_date, br.status,
  b.id AS b_id, b.title AS b_title, b.author AS b_author, b.isbn AS b_isbn,
  b.category AS b_category, b.description AS b_description, b.cover_color AS b_cover_color,
  b.publisher AS b_publisher, b.published_year AS b_published_year,
  b.total_copies AS b_total_copies, b.available_copies AS b_available_copies,
  b.created_at AS b_created_at,
  u.id AS u_id, u.username AS u_username, u.name AS u_name, u.role AS u_role,
  u.email AS u_email, u.created_at AS u_created_at
FROM borrow_records br
LEFT JOIN books b ON br.book_id = b.id
LEFT JOIN users u ON br.user_id = u.id
`

func handleListBorrows(c *gin.Context) {
	page, pageSize := pagination(c)
	status := c.Query("status")
	uid := getUID(c)
	role := getRole(c)

	where := "WHERE 1=1"
	args := []interface{}{}
	if role != "admin" {
		where += " AND br.user_id=?"
		args = append(args, uid)
	} else if uidStr := c.Query("userId"); uidStr != "" {
		where += " AND br.user_id=?"
		args = append(args, uidStr)
	}
	switch status {
	case "borrowed":
		where += " AND br.return_date IS NULL AND br.due_date >= datetime('now')"
	case "overdue":
		where += " AND br.return_date IS NULL AND br.due_date < datetime('now')"
	case "returned":
		where += " AND br.return_date IS NOT NULL"
	}

	var total int64
	db.Get(&total, "SELECT COUNT(*) FROM borrow_records br "+where, args...)

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)

	var rows []borrowJoinRow
	db.Select(&rows, borrowJoinSQL+where+" ORDER BY br.borrow_date DESC, br.id DESC LIMIT ? OFFSET ?", listArgs...)

	recs := make([]Borrow, len(rows))
	for i := range rows {
		recs[i] = rows[i].toBorrow()
	}

	ok(c, PageData{List: recs, Total: total, Page: page, PageSize: pageSize})
}

func handleBorrow(c *gin.Context) {
	var body struct {
		BookID int64 `json:"book_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.BookID == 0 {
		fail(c, 400, "参数错误")
		return
	}
	uid := getUID(c)
	tx := db.MustBegin()
	var avail int
	err := tx.Get(&avail, "SELECT available_copies FROM books WHERE id=?", body.BookID)
	if err != nil || avail <= 0 {
		tx.Rollback()
		fail(c, 400, "该书暂无可借副本")
		return
	}
	tx.MustExec("INSERT INTO borrow_records (user_id, book_id, due_date) VALUES (?, ?, datetime('now','+30 days'))", uid, body.BookID)
	tx.MustExec("UPDATE books SET available_copies = available_copies - 1 WHERE id=?", body.BookID)
	if err := tx.Commit(); err != nil {
		fail(c, 500, "借阅失败")
		return
	}
	ok(c, gin.H{"ok": true})
}

func handleReturn(c *gin.Context) {
	id := c.Param("id")
	uid := getUID(c)
	role := getRole(c)

	tx := db.MustBegin()
	var rec struct {
		ID         int64   `db:"id"`
		UserID     int64   `db:"user_id"`
		BookID     int64   `db:"book_id"`
		ReturnDate *string `db:"return_date"`
	}
	err := tx.Get(&rec, "SELECT id, user_id, book_id, return_date FROM borrow_records WHERE id=?", id)
	if err != nil {
		tx.Rollback()
		fail(c, 404, "借阅记录不存在")
		return
	}
	if role != "admin" && rec.UserID != uid {
		tx.Rollback()
		fail(c, 403, "无权操作")
		return
	}
	if rec.ReturnDate != nil {
		tx.Rollback()
		fail(c, 400, "该书已归还")
		return
	}
	tx.MustExec("UPDATE borrow_records SET return_date=datetime('now'), status='returned' WHERE id=?", id)
	tx.MustExec("UPDATE books SET available_copies = MIN(available_copies+1, total_copies) WHERE id=?", rec.BookID)
	if err := tx.Commit(); err != nil {
		fail(c, 500, "归还失败")
		return
	}
	ok(c, gin.H{"ok": true})
}
