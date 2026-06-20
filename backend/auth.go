package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecret = "qingtan-library-secret-2026"
	tokenTTL  = 72 * time.Hour
)

type Claims struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func generateToken(u User) (string, error) {
	claims := Claims{
		UID:      u.ID,
		Username: u.Username,
		Name:     u.Name,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(jwtSecret))
}

func parseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !t.Valid {
		return nil, err
	}
	return claims, nil
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return c.Query("token")
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseToken(extractToken(c))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, Response{Code: 401, Msg: "未登录", Data: nil})
			return
		}
		c.Set("uid", claims.UID)
		c.Set("username", claims.Username)
		c.Set("name", claims.Name)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if getRole(c) != "admin" {
			c.AbortWithStatusJSON(http.StatusOK, Response{Code: 403, Msg: "无权限", Data: nil})
			return
		}
		c.Next()
	}
}

func getUID(c *gin.Context) int64 {
	v, _ := c.Get("uid")
	if id, ok := v.(int64); ok {
		return id
	}
	return 0
}

func getRole(c *gin.Context) string {
	return c.GetString("role")
}

func handleLogin(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Username == "" || body.Password == "" {
		fail(c, 400, "请输入账号和密码")
		return
	}
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE username=?", body.Username)
	if err != nil {
		fail(c, 401, "账号或密码错误")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)) != nil {
		fail(c, 401, "账号或密码错误")
		return
	}
	token, err := generateToken(u)
	if err != nil {
		fail(c, 500, "生成令牌失败")
		return
	}
	u.PasswordHash = ""
	u.CreatedAt = normTime(u.CreatedAt)
	ok(c, loginResult{Token: token, User: u})
}

func handleMe(c *gin.Context) {
	var u User
	err := db.Get(&u, "SELECT id, username, name, role, email, created_at FROM users WHERE id=?", getUID(c))
	if err != nil {
		fail(c, 404, "用户不存在")
		return
	}
	u.CreatedAt = normTime(u.CreatedAt)
	ok(c, u)
}

func handleChangePassword(c *gin.Context) {
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		fail(c, 400, "参数错误")
		return
	}
	if len(body.NewPassword) < 6 {
		fail(c, 400, "新密码至少 6 位")
		return
	}
	uid := getUID(c)
	var hash string
	if err := db.Get(&hash, "SELECT password_hash FROM users WHERE id=?", uid); err != nil {
		fail(c, 404, "用户不存在")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.OldPassword)) != nil {
		fail(c, 400, "原密码错误")
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	db.MustExec("UPDATE users SET password_hash=? WHERE id=?", string(newHash), uid)
	ok(c, gin.H{"ok": true})
}
