package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "almanac.db"))
	if err != nil {
		return nil, err
	}
	conn.Exec("PRAGMA journal_mode=WAL")
	conn.Exec("PRAGMA busy_timeout=5000")
	conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS entries (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    version TEXT DEFAULT '',
    content TEXT DEFAULT '',
    tags TEXT DEFAULT '',
    published INTEGER DEFAULT 0,
    published_at TEXT DEFAULT '',
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_entries_pub ON entries(published, published_at);

CREATE TABLE IF NOT EXISTS subscribers (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    token TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_subs_token ON subscribers(token);
`)
	return err
}

// --- Entries ---

type Entry struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Version     string `json:"version"`
	Content     string `json:"content"`
	Tags        string `json:"tags"`
	Published   bool   `json:"published"`
	PublishedAt string `json:"published_at,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (db *DB) CreateEntry(title, version, content, tags string, published bool) (*Entry, error) {
	id := "ent_" + genID(8)
	now := time.Now().UTC().Format(time.RFC3339)
	pub := 0
	pubAt := ""
	if published {
		pub = 1
		pubAt = now
	}
	_, err := db.conn.Exec("INSERT INTO entries (id,title,version,content,tags,published,published_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?)",
		id, title, version, content, tags, pub, pubAt, now, now)
	if err != nil {
		return nil, err
	}
	return &Entry{ID: id, Title: title, Version: version, Content: content, Tags: tags,
		Published: published, PublishedAt: pubAt, CreatedAt: now, UpdatedAt: now}, nil
}

func (db *DB) ListEntries(publishedOnly bool) ([]Entry, error) {
	var query string
	if publishedOnly {
		query = "SELECT id,title,version,content,tags,published,published_at,created_at,updated_at FROM entries WHERE published=1 ORDER BY published_at DESC"
	} else {
		query = "SELECT id,title,version,content,tags,published,published_at,created_at,updated_at FROM entries ORDER BY created_at DESC"
	}
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Entry
	for rows.Next() {
		var e Entry
		var pub int
		rows.Scan(&e.ID, &e.Title, &e.Version, &e.Content, &e.Tags, &pub, &e.PublishedAt, &e.CreatedAt, &e.UpdatedAt)
		e.Published = pub == 1
		out = append(out, e)
	}
	return out, rows.Err()
}

func (db *DB) GetEntry(id string) (*Entry, error) {
	var e Entry
	var pub int
	err := db.conn.QueryRow("SELECT id,title,version,content,tags,published,published_at,created_at,updated_at FROM entries WHERE id=?", id).
		Scan(&e.ID, &e.Title, &e.Version, &e.Content, &e.Tags, &pub, &e.PublishedAt, &e.CreatedAt, &e.UpdatedAt)
	e.Published = pub == 1
	return &e, err
}

func (db *DB) UpdateEntry(id string, title, version, content, tags *string, published *bool) (*Entry, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	if title != nil {
		db.conn.Exec("UPDATE entries SET title=?, updated_at=? WHERE id=?", *title, now, id)
	}
	if version != nil {
		db.conn.Exec("UPDATE entries SET version=?, updated_at=? WHERE id=?", *version, now, id)
	}
	if content != nil {
		db.conn.Exec("UPDATE entries SET content=?, updated_at=? WHERE id=?", *content, now, id)
	}
	if tags != nil {
		db.conn.Exec("UPDATE entries SET tags=?, updated_at=? WHERE id=?", *tags, now, id)
	}
	if published != nil {
		pub := 0
		if *published {
			pub = 1
			// Set published_at if first publish
			var existing string
			db.conn.QueryRow("SELECT published_at FROM entries WHERE id=?", id).Scan(&existing)
			if existing == "" {
				db.conn.Exec("UPDATE entries SET published_at=? WHERE id=?", now, id)
			}
		}
		db.conn.Exec("UPDATE entries SET published=?, updated_at=? WHERE id=?", pub, now, id)
	}
	return db.GetEntry(id)
}

func (db *DB) DeleteEntry(id string) error {
	_, err := db.conn.Exec("DELETE FROM entries WHERE id=?", id)
	return err
}

func (db *DB) TotalEntries() int {
	var count int
	db.conn.QueryRow("SELECT COUNT(*) FROM entries").Scan(&count)
	return count
}

// --- Subscribers ---

type Subscriber struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (db *DB) AddSubscriber(email string) (*Subscriber, error) {
	id := "sub_" + genID(6)
	token := genID(16)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO subscribers (id,email,token,created_at) VALUES (?,?,?,?)", id, email, token, now)
	if err != nil {
		return nil, err
	}
	return &Subscriber{ID: id, Email: email, Status: "active", CreatedAt: now}, nil
}

func (db *DB) ListSubscribers() ([]Subscriber, error) {
	rows, err := db.conn.Query("SELECT id,email,status,created_at FROM subscribers WHERE status='active' ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subscriber
	for rows.Next() {
		var s Subscriber
		rows.Scan(&s.ID, &s.Email, &s.Status, &s.CreatedAt)
		out = append(out, s)
	}
	return out, rows.Err()
}

func (db *DB) Unsubscribe(token string) error {
	_, err := db.conn.Exec("UPDATE subscribers SET status='unsubscribed' WHERE token=?", token)
	return err
}

// --- RSS ---

func (db *DB) RSSItems() string {
	entries, _ := db.ListEntries(true)
	var items strings.Builder
	for _, e := range entries {
		pubDate := e.PublishedAt
		if t, err := time.Parse(time.RFC3339, pubDate); err == nil {
			pubDate = t.Format(time.RFC1123Z)
		}
		title := e.Title
		if e.Version != "" {
			title += " (" + e.Version + ")"
		}
		items.WriteString("<item>")
		items.WriteString("<title>" + xmlEscape(title) + "</title>")
		items.WriteString("<description>" + xmlEscape(e.Content) + "</description>")
		items.WriteString("<pubDate>" + pubDate + "</pubDate>")
		items.WriteString("<guid>" + e.ID + "</guid>")
		items.WriteString("</item>\n")
	}
	return items.String()
}

// --- Stats ---

func (db *DB) Stats() map[string]any {
	var entries, published, drafts, subscribers int
	db.conn.QueryRow("SELECT COUNT(*) FROM entries").Scan(&entries)
	db.conn.QueryRow("SELECT COUNT(*) FROM entries WHERE published=1").Scan(&published)
	db.conn.QueryRow("SELECT COUNT(*) FROM entries WHERE published=0").Scan(&drafts)
	db.conn.QueryRow("SELECT COUNT(*) FROM subscribers WHERE status='active'").Scan(&subscribers)
	return map[string]any{"entries": entries, "published": published, "drafts": drafts, "subscribers": subscribers}
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func genID(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
