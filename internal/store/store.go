package store
import ("database/sql";"fmt";"os";"path/filepath";"strings";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Entry struct{ID string `json:"id"`;Date string `json:"date"`;Title string `json:"title,omitempty"`;Body string `json:"body"`;Mood string `json:"mood,omitempty"`;Tags string `json:"tags,omitempty"`;WordCount int `json:"word_count"`;CreatedAt string `json:"created_at"`;UpdatedAt string `json:"updated_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"almanac.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS entries(id TEXT PRIMARY KEY,date TEXT NOT NULL,title TEXT DEFAULT '',body TEXT DEFAULT '',mood TEXT DEFAULT '',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')),updated_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func today()string{return time.Now().Format("2006-01-02")}
func(d *DB)Create(e *Entry)error{e.ID=genID();e.CreatedAt=now();e.UpdatedAt=e.CreatedAt;if e.Date==""{e.Date=today()};e.WordCount=len(strings.Fields(e.Body))
_,err:=d.db.Exec(`INSERT INTO entries(id,date,title,body,mood,tags,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Date,e.Title,e.Body,e.Mood,e.Tags,e.CreatedAt,e.UpdatedAt);return err}
func(d *DB)Get(id string)*Entry{var e Entry;if d.db.QueryRow(`SELECT id,date,title,body,mood,tags,created_at,updated_at FROM entries WHERE id=?`,id).Scan(&e.ID,&e.Date,&e.Title,&e.Body,&e.Mood,&e.Tags,&e.CreatedAt,&e.UpdatedAt)!=nil{return nil};e.WordCount=len(strings.Fields(e.Body));return &e}
func(d *DB)List(month string,limit int)[]Entry{if limit<=0{limit=50};q:=`SELECT id,date,title,body,mood,tags,created_at,updated_at FROM entries`;args:=[]any{}
if month!=""{q+=` WHERE date LIKE ?`;args=append(args,month+"%")};q+=` ORDER BY date DESC LIMIT ?`;args=append(args,limit)
rows,_:=d.db.Query(q,args...);if rows==nil{return nil};defer rows.Close()
var o []Entry;for rows.Next(){var e Entry;rows.Scan(&e.ID,&e.Date,&e.Title,&e.Body,&e.Mood,&e.Tags,&e.CreatedAt,&e.UpdatedAt);e.WordCount=len(strings.Fields(e.Body));o=append(o,e)};return o}
func(d *DB)Update(id string,e *Entry)error{e.UpdatedAt=now();e.WordCount=len(strings.Fields(e.Body));_,err:=d.db.Exec(`UPDATE entries SET date=?,title=?,body=?,mood=?,tags=?,updated_at=? WHERE id=?`,e.Date,e.Title,e.Body,e.Mood,e.Tags,e.UpdatedAt,id);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM entries WHERE id=?`,id);return err}
func(d *DB)Search(q string)[]Entry{s:="%"+q+"%";rows,_:=d.db.Query(`SELECT id,date,title,body,mood,tags,created_at,updated_at FROM entries WHERE title LIKE ? OR body LIKE ? ORDER BY date DESC`,s,s);if rows==nil{return nil};defer rows.Close()
var o []Entry;for rows.Next(){var e Entry;rows.Scan(&e.ID,&e.Date,&e.Title,&e.Body,&e.Mood,&e.Tags,&e.CreatedAt,&e.UpdatedAt);e.WordCount=len(strings.Fields(e.Body));o=append(o,e)};return o}
func(d *DB)Streak()int{rows,_:=d.db.Query(`SELECT DISTINCT date FROM entries ORDER BY date DESC`);if rows==nil{return 0};defer rows.Close()
streak:=0;expected:=today()
for rows.Next(){var dt string;rows.Scan(&dt);if dt==expected{streak++;t,_:=time.Parse("2006-01-02",dt);expected=t.AddDate(0,0,-1).Format("2006-01-02")}else{break}};return streak}
type Stats struct{Entries int `json:"entries"`;Words int `json:"words"`;Streak int `json:"streak"`;Months int `json:"months"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM entries`).Scan(&s.Entries)
rows,_:=d.db.Query(`SELECT body FROM entries`);if rows!=nil{defer rows.Close();for rows.Next(){var b string;rows.Scan(&b);s.Words+=len(strings.Fields(b))}}
s.Streak=d.Streak();d.db.QueryRow(`SELECT COUNT(DISTINCT substr(date,1,7)) FROM entries`).Scan(&s.Months);return s}
