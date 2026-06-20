package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func atoiDefault(s string, def int) int {
	n, err := strconv.Atoi(s)
	if err != nil || n == 0 {
		return def
	}
	return n
}

func pagination(c *gin.Context) (int, int) {
	page := atoiDefault(c.Query("page"), 1)
	pageSize := atoiDefault(c.Query("pageSize"), 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return page, pageSize
}

func handleListBooks(c *gin.Context) {
	page, pageSize := pagination(c)
	keyword := strings.TrimSpace(c.Query("keyword"))
	category := c.Query("category")

	where := "1=1"
	args := []interface{}{}
	if keyword != "" {
		like := "%" + keyword + "%"
		where += " AND (title LIKE ? OR author LIKE ? OR isbn LIKE ?)"
		args = append(args, like, like, like)
	}
	if category != "" {
		where += " AND category = ?"
		args = append(args, category)
	}

	var total int64
	db.Get(&total, "SELECT COUNT(*) FROM books WHERE "+where, args...)

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)
	var books []Book
	db.Select(&books, "SELECT * FROM books WHERE "+where+" ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?", listArgs...)
	for i := range books {
		books[i].CreatedAt = normTime(books[i].CreatedAt)
	}
	ok(c, PageData{List: books, Total: total, Page: page, PageSize: pageSize})
}

func handleCategories(c *gin.Context) {
	var cats []CategoryCount
	db.Select(&cats, "SELECT category, COUNT(*) as count FROM books GROUP BY category ORDER BY count DESC, category ASC")
	ok(c, cats)
}

func handleGetBook(c *gin.Context) {
	var b Book
	err := db.Get(&b, "SELECT * FROM books WHERE id=?", c.Param("id"))
	if err != nil {
		fail(c, 404, "图书不存在")
		return
	}
	b.CreatedAt = normTime(b.CreatedAt)
	ok(c, b)
}

type bookInput struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	ISBN          string `json:"isbn"`
	Category      string `json:"category"`
	Description   string `json:"description"`
	CoverColor     string `json:"cover_color"`
	Publisher     string `json:"publisher"`
	PublishedYear int    `json:"published_year"`
	TotalCopies   int    `json:"total_copies"`
}

func handleCreateBook(c *gin.Context) {
	var in bookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if in.Title == "" || in.Author == "" {
		fail(c, 400, "书名与作者必填")
		return
	}
	if in.TotalCopies < 1 {
		in.TotalCopies = 1
	}
	if in.CoverColor == "" {
		in.CoverColor = "#1F3D2B"
	}
	res, err := db.Exec(
		"INSERT INTO books (title,author,isbn,category,description,cover_color,publisher,published_year,total_copies,available_copies) VALUES (?,?,?,?,?,?,?,?,?,?)",
		in.Title, in.Author, in.ISBN, in.Category, in.Description, in.CoverColor, in.Publisher, in.PublishedYear, in.TotalCopies, in.TotalCopies,
	)
	if err != nil {
		fail(c, 500, "新增失败")
		return
	}
	id, _ := res.LastInsertId()
	var b Book
	db.Get(&b, "SELECT * FROM books WHERE id=?", id)
	b.CreatedAt = normTime(b.CreatedAt)
	ok(c, b)
}

func handleUpdateBook(c *gin.Context) {
	id := c.Param("id")
	var in bookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if in.TotalCopies < 1 {
		in.TotalCopies = 1
	}
	if in.CoverColor == "" {
		in.CoverColor = "#1F3D2B"
	}
	var old struct {
		Total int `db:"total_copies"`
		Avail int `db:"available_copies"`
	}
	if err := db.Get(&old, "SELECT total_copies, available_copies FROM books WHERE id=?", id); err != nil {
		fail(c, 404, "图书不存在")
		return
	}
	borrowed := old.Total - old.Avail
	if in.TotalCopies < borrowed {
		fail(c, 400, fmt.Sprintf("总馆藏不能小于已借出数量（%d 本）", borrowed))
		return
	}
	newAvail := old.Avail + (in.TotalCopies - old.Total)
	if newAvail < 0 {
		newAvail = 0
	}
	if newAvail > in.TotalCopies {
		newAvail = in.TotalCopies
	}
	db.MustExec(
		"UPDATE books SET title=?,author=?,isbn=?,category=?,description=?,cover_color=?,publisher=?,published_year=?,total_copies=?,available_copies=? WHERE id=?",
		in.Title, in.Author, in.ISBN, in.Category, in.Description, in.CoverColor, in.Publisher, in.PublishedYear, in.TotalCopies, newAvail, id,
	)
	var b Book
	db.Get(&b, "SELECT * FROM books WHERE id=?", id)
	b.CreatedAt = normTime(b.CreatedAt)
	ok(c, b)
}

func handleDeleteBook(c *gin.Context) {
	id := c.Param("id")
	var n int64
	db.Get(&n, "SELECT COUNT(*) FROM borrow_records WHERE book_id=? AND return_date IS NULL", id)
	if n > 0 {
		fail(c, 400, "该图书有未归还的借阅记录，无法删除")
		return
	}
	db.MustExec("DELETE FROM borrow_records WHERE book_id=?", id)
	db.MustExec("DELETE FROM books WHERE id=?", id)
	ok(c, gin.H{"ok": true})
}
