package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

func handleListUsers(c *gin.Context) {
	page, pageSize := pagination(c)
	keyword := strings.TrimSpace(c.Query("keyword"))

	where := "1=1"
	args := []interface{}{}
	if keyword != "" {
		like := "%" + keyword + "%"
		where += " AND (username LIKE ? OR name LIKE ? OR email LIKE ?)"
		args = append(args, like, like, like)
	}

	var total int64
	db.Get(&total, "SELECT COUNT(*) FROM users WHERE "+where, args...)

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)
	var users []User
	db.Select(&users, "SELECT id, username, name, role, email, created_at FROM users WHERE "+where+" ORDER BY id ASC LIMIT ? OFFSET ?", listArgs...)
	for i := range users {
		users[i].CreatedAt = normTime(users[i].CreatedAt)
	}
	ok(c, PageData{List: users, Total: total, Page: page, PageSize: pageSize})
}

func handleCreateUser(c *gin.Context) {
	var in userInput
	if err := c.ShouldBindJSON(&in); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if in.Username == "" || in.Name == "" {
		fail(c, 400, "用户名与姓名必填")
		return
	}
	if in.Role != "admin" && in.Role != "student" {
		in.Role = "student"
	}
	if in.Password == "" {
		in.Password = "123456"
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	res, err := db.Exec(
		"INSERT INTO users (username, password_hash, name, role, email) VALUES (?,?,?,?,?)",
		in.Username, string(hash), in.Name, in.Role, in.Email,
	)
	if err != nil {
		fail(c, 400, "用户名已存在")
		return
	}
	id, _ := res.LastInsertId()
	var u User
	db.Get(&u, "SELECT id, username, name, role, email, created_at FROM users WHERE id=?", id)
	u.CreatedAt = normTime(u.CreatedAt)
	ok(c, u)
}

func handleUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var in userInput
	if err := c.ShouldBindJSON(&in); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if in.Role != "admin" && in.Role != "student" {
		in.Role = "student"
	}
	if in.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		db.MustExec("UPDATE users SET password_hash=? WHERE id=?", string(hash), id)
	}
	db.MustExec("UPDATE users SET name=?, role=?, email=? WHERE id=?", in.Name, in.Role, in.Email, id)
	var u User
	db.Get(&u, "SELECT id, username, name, role, email, created_at FROM users WHERE id=?", id)
	u.CreatedAt = normTime(u.CreatedAt)
	ok(c, u)
}

func handleDeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == strconv.FormatInt(getUID(c), 10) {
		fail(c, 400, "不能删除当前登录用户")
		return
	}
	var n int64
	db.Get(&n, "SELECT COUNT(*) FROM borrow_records WHERE user_id=? AND return_date IS NULL", id)
	if n > 0 {
		fail(c, 400, "该用户有未归还的借阅记录，无法删除")
		return
	}
	db.MustExec("DELETE FROM borrow_records WHERE user_id=?", id)
	db.MustExec("DELETE FROM users WHERE id=?", id)
	ok(c, gin.H{"ok": true})
}
