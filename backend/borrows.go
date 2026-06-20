package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func keysInt64(m map[int64]bool) []int64 {
	ks := make([]int64, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func deriveStatus(dueDate, returnDate string) string {
	if returnDate != "" {
		return "returned"
	}
	t, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return "borrowed"
	}
	if t.Before(time.Now()) {
		return "overdue"
	}
	return "borrowed"
}

func handleListBorrows(c *gin.Context) {
	page, pageSize := pagination(c)
	status := c.Query("status")
	uid := getUID(c)
	role := getRole(c)

	where := "1=1"
	args := []interface{}{}
	if role != "admin" {
		where += " AND user_id=?"
		args = append(args, uid)
	} else if uidStr := c.Query("userId"); uidStr != "" {
		where += " AND user_id=?"
		args = append(args, uidStr)
	}
	switch status {
	case "borrowed":
		where += " AND return_date IS NULL AND due_date >= datetime('now')"
	case "overdue":
		where += " AND return_date IS NULL AND due_date < datetime('now')"
	case "returned":
		where += " AND return_date IS NOT NULL"
	}

	var total int64
	db.Get(&total, "SELECT COUNT(*) FROM borrow_records WHERE "+where, args...)

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)
	var recs []Borrow
	db.Select(&recs, "SELECT * FROM borrow_records WHERE "+where+" ORDER BY borrow_date DESC, id DESC LIMIT ? OFFSET ?", listArgs...)

	bookIDs := map[int64]bool{}
	userIDs := map[int64]bool{}
	for _, r := range recs {
		bookIDs[r.BookID] = true
		userIDs[r.UserID] = true
	}
	books := map[int64]*Book{}
	users := map[int64]*User{}
	if len(bookIDs) > 0 {
		var bs []Book
		query, qargs, _ := sqlx.In("SELECT * FROM books WHERE id IN (?)", keysInt64(bookIDs))
		db.Select(&bs, query, qargs...)
		for i := range bs {
			bs[i].CreatedAt = normTime(bs[i].CreatedAt)
			books[bs[i].ID] = &bs[i]
		}
	}
	if len(userIDs) > 0 {
		var us []User
		query, qargs, _ := sqlx.In("SELECT id, username, name, role, email, created_at FROM users WHERE id IN (?)", keysInt64(userIDs))
		db.Select(&us, query, qargs...)
		for i := range us {
			us[i].CreatedAt = normTime(us[i].CreatedAt)
			users[us[i].ID] = &us[i]
		}
	}
	for i := range recs {
		recs[i].BorrowDate = normTime(recs[i].BorrowDate)
		recs[i].DueDate = normTime(recs[i].DueDate)
		recs[i].ReturnDate = normTime(recs[i].ReturnDate)
		recs[i].Status = deriveStatus(recs[i].DueDate, recs[i].ReturnDate)
		if b, ok := books[recs[i].BookID]; ok {
			recs[i].Book = b
		}
		if u, ok := users[recs[i].UserID]; ok {
			recs[i].User = u
		}
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
