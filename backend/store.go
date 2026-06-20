package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB

func initDB(path string) {
	var err error
	db, err = sqlx.Connect("sqlite", path)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	db.SetMaxOpenConns(1)
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA busy_timeout=5000;")
	db.Exec("PRAGMA foreign_keys=ON;")
	migrate()
	seed()
}

func migrate() {
	schema := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'student',
    email TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    isbn TEXT DEFAULT '',
    category TEXT DEFAULT '',
    description TEXT DEFAULT '',
    cover_color TEXT DEFAULT '#1F3D2B',
    publisher TEXT DEFAULT '',
    published_year INTEGER DEFAULT 0,
    total_copies INTEGER NOT NULL DEFAULT 1,
    available_copies INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS borrow_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    borrow_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date DATETIME NOT NULL,
    return_date DATETIME,
    status TEXT NOT NULL DEFAULT 'borrowed',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);
CREATE INDEX IF NOT EXISTS idx_books_category ON books(category);
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_borrows_user ON borrow_records(user_id);
CREATE INDEX IF NOT EXISTS idx_borrows_status ON borrow_records(status);
`
	db.MustExec(schema)
}

func seed() {
	ensureUser("admin", "admin123", "图书管理员", "admin", "admin@library.edu")
	ensureUser("zhangsan", "student123", "张三", "student", "zhangsan@school.edu")
	ensureUser("lisi", "student123", "李四", "student", "lisi@school.edu")
	ensureUser("wangwu", "student123", "王五", "student", "wangwu@school.edu")

	var n int64
	db.Get(&n, "SELECT COUNT(*) FROM books")
	if n == 0 {
		seedBooks()
	}
}

func ensureUser(username, password, name, role, email string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var id int64
	err := db.Get(&id, "SELECT id FROM users WHERE username=?", username)
	if err != nil {
		db.MustExec(
			"INSERT INTO users (username, password_hash, name, role, email) VALUES (?,?,?,?,?)",
			username, string(hash), name, role, email,
		)
		return
	}
	db.MustExec("UPDATE users SET password_hash=? WHERE id=?", string(hash), id)
}

func seedBooks() {
	books := []Book{
		{Title: "百年孤独", Author: "加西亚·马尔克斯", ISBN: "978-7-5442-5399-4", Category: "文学", Description: "马孔多镇布恩迪亚家族七代人的兴衰史，魔幻现实主义的开山之作。", CoverColor: "#1F3D2B", Publisher: "南海出版公司", PublishedYear: 2017, TotalCopies: 5},
		{Title: "三体", Author: "刘慈欣", ISBN: "978-7-5366-9293-0", Category: "科幻", Description: "地球文明在宇宙中的生存博弈，硬科幻的里程碑式作品。", CoverColor: "#0B1F3A", Publisher: "重庆出版社", PublishedYear: 2008, TotalCopies: 8},
		{Title: "人类简史", Author: "尤瓦尔·赫拉利", ISBN: "978-7-5086-6068-4", Category: "历史", Description: "从认知革命到科学革命，重新审视人类如何成为地球的主宰。", CoverColor: "#7A4A1B", Publisher: "中信出版社", PublishedYear: 2014, TotalCopies: 6},
		{Title: "活着", Author: "余华", ISBN: "978-7-5327-6201-4", Category: "文学", Description: "福贵一生的苦难与坚韧，关于生命与命运的当代经典。", CoverColor: "#3A0E0E", Publisher: "上海译文出版社", PublishedYear: 2014, TotalCopies: 4},
		{Title: "深度工作", Author: "卡尔·纽波特", ISBN: "978-7-210-09320-0", Category: "自我提升", Description: "在分心的世界里培养专注力的方法论。", CoverColor: "#1B3A2F", Publisher: "江西人民出版社", PublishedYear: 2017, TotalCopies: 3},
		{Title: "明朝那些事儿", Author: "当年明月", ISBN: "978-7-80211-679-3", Category: "历史", Description: "用诙谐笔法讲述明朝三百年风云变幻。", CoverColor: "#5B2A0E", Publisher: "中国海关出版社", PublishedYear: 2009, TotalCopies: 5},
		{Title: "算法导论", Author: "Thomas H. Cormen", ISBN: "978-7-111-40701-0", Category: "计算机", Description: "算法领域公认的权威教材，覆盖全面、论证严谨。", CoverColor: "#10243F", Publisher: "机械工业出版社", PublishedYear: 2013, TotalCopies: 3},
		{Title: "设计模式", Author: "Erich Gamma", ISBN: "978-7-111-07575-2", Category: "计算机", Description: "面向对象软件设计模式的奠基之作，俗称『四人帮』。", CoverColor: "#2A2A2A", Publisher: "机械工业出版社", PublishedYear: 2002, TotalCopies: 4},
		{Title: "苏菲的世界", Author: "乔斯坦·贾德", ISBN: "978-7-5442-5397-0", Category: "哲学", Description: "以小说形式串起西方哲学史，写给所有人的哲学入门。", CoverColor: "#3A1F4A", Publisher: "南海出版公司", PublishedYear: 2007, TotalCopies: 4},
		{Title: "艺术的故事", Author: "贡布里希", ISBN: "978-7-5494-0820-2", Category: "艺术", Description: "一部跨越数千年的人类视觉艺术史，图文并茂的经典。", CoverColor: "#6E1F2A", Publisher: "广西美术出版社", PublishedYear: 2014, TotalCopies: 3},
		{Title: "围城", Author: "钱钟书", ISBN: "978-7-020-02940-2", Category: "文学", Description: "方鸿渐的人生围城，讽刺文学的巅峰之作。", CoverColor: "#40301A", Publisher: "人民文学出版社", PublishedYear: 1991, TotalCopies: 4},
		{Title: "时间简史", Author: "史蒂芬·霍金", ISBN: "978-7-5357-7910-9", Category: "科普", Description: "从大爆炸到黑洞，物理学巨匠写给大众的宇宙学读本。", CoverColor: "#0E1A33", Publisher: "湖南科学技术出版社", PublishedYear: 2014, TotalCopies: 5},
	}
	for _, b := range books {
		db.MustExec(
			"INSERT INTO books (title,author,isbn,category,description,cover_color,publisher,published_year,total_copies,available_copies) VALUES (?,?,?,?,?,?,?,?,?,?)",
			b.Title, b.Author, b.ISBN, b.Category, b.Description, b.CoverColor, b.Publisher, b.PublishedYear, b.TotalCopies, b.TotalCopies,
		)
	}
}
