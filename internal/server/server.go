package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stockyard-dev/stockyard-almanac/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	port   int
	limits Limits
}

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.routes()
	return s
}

func (s *Server) routes() {
	// Entries
	s.mux.HandleFunc("POST /api/entries", s.handleCreateEntry)
	s.mux.HandleFunc("GET /api/entries", s.handleListEntries)
	s.mux.HandleFunc("GET /api/entries/{id}", s.handleGetEntry)
	s.mux.HandleFunc("PUT /api/entries/{id}", s.handleUpdateEntry)
	s.mux.HandleFunc("DELETE /api/entries/{id}", s.handleDeleteEntry)

	// Subscribers
	s.mux.HandleFunc("POST /api/subscribers", s.handleAddSubscriber)
	s.mux.HandleFunc("GET /api/subscribers", s.handleListSubscribers)
	s.mux.HandleFunc("GET /unsubscribe", s.handleUnsubscribe)

	// Public pages
	s.mux.HandleFunc("GET /changelog", s.handleChangelogPage)
	s.mux.HandleFunc("GET /api/changelog", s.handleChangelogAPI)
	s.mux.HandleFunc("GET /changelog/rss", s.handleRSS)
	s.mux.HandleFunc("POST /changelog/subscribe", s.handlePublicSubscribe)

	// Status
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"product": "stockyard-almanac", "version": "0.1.0"})
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("[almanac] listening on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

// --- Entry handlers ---

func (s *Server) handleCreateEntry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title     string `json:"title"`
		Version   string `json:"version"`
		Content   string `json:"content"`
		Tags      string `json:"tags"`
		Published bool   `json:"published"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		writeJSON(w, 400, map[string]string{"error": "title is required"})
		return
	}
	if s.limits.MaxEntries > 0 {
		total := s.db.TotalEntries()
		if LimitReached(s.limits.MaxEntries, total) {
			writeJSON(w, 402, map[string]string{"error": fmt.Sprintf("free tier limit: %d entries — upgrade to Pro", s.limits.MaxEntries), "upgrade": "https://stockyard.dev/almanac/"})
			return
		}
	}
	e, err := s.db.CreateEntry(req.Title, req.Version, req.Content, req.Tags, req.Published)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, map[string]any{"entry": e})
}

func (s *Server) handleListEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := s.db.ListEntries(false)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if entries == nil {
		entries = []store.Entry{}
	}
	writeJSON(w, 200, map[string]any{"entries": entries, "count": len(entries)})
}

func (s *Server) handleGetEntry(w http.ResponseWriter, r *http.Request) {
	e, err := s.db.GetEntry(r.PathValue("id"))
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "entry not found"})
		return
	}
	writeJSON(w, 200, map[string]any{"entry": e})
}

func (s *Server) handleUpdateEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := s.db.GetEntry(id); err != nil {
		writeJSON(w, 404, map[string]string{"error": "entry not found"})
		return
	}
	var req struct {
		Title     *string `json:"title"`
		Version   *string `json:"version"`
		Content   *string `json:"content"`
		Tags      *string `json:"tags"`
		Published *bool   `json:"published"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	e, err := s.db.UpdateEntry(id, req.Title, req.Version, req.Content, req.Tags, req.Published)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"entry": e})
}

func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteEntry(r.PathValue("id"))
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// --- Subscribers ---

func (s *Server) handleAddSubscriber(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		writeJSON(w, 400, map[string]string{"error": "email is required"})
		return
	}
	sub, err := s.db.AddSubscriber(req.Email)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, map[string]any{"subscriber": sub})
}

func (s *Server) handleListSubscribers(w http.ResponseWriter, r *http.Request) {
	subs, err := s.db.ListSubscribers()
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if subs == nil {
		subs = []store.Subscriber{}
	}
	writeJSON(w, 200, map[string]any{"subscribers": subs, "count": len(subs)})
}

func (s *Server) handleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "" {
		s.db.Unsubscribe(token)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html><html><head><title>Unsubscribed</title><style>body{font-family:system-ui;display:flex;justify-content:center;align-items:center;min-height:100vh;background:#1a1410;color:#f0e6d3}h1{margin-bottom:.5rem}</style></head><body><div><h1>Unsubscribed</h1><p>You've been removed.</p></div></body></html>`))
}

func (s *Server) handlePublicSubscribe(w http.ResponseWriter, r *http.Request) {
	var email string
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		var req struct{ Email string `json:"email"` }
		json.NewDecoder(r.Body).Decode(&req)
		email = req.Email
	} else {
		r.ParseForm()
		email = r.FormValue("email")
	}
	if email == "" {
		http.Error(w, "Email required", 400)
		return
	}
	s.db.AddSubscriber(email)
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		writeJSON(w, 200, map[string]string{"status": "subscribed"})
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html><html><head><title>Subscribed</title><style>body{font-family:system-ui;display:flex;justify-content:center;align-items:center;min-height:100vh;background:#1a1410;color:#f0e6d3}h1{margin-bottom:.5rem}</style></head><body><div><h1>Subscribed!</h1><p>You'll be notified of new releases.</p></div></body></html>`))
}

// --- Public changelog ---

func (s *Server) handleChangelogAPI(w http.ResponseWriter, r *http.Request) {
	entries, err := s.db.ListEntries(true)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if entries == nil {
		entries = []store.Entry{}
	}
	writeJSON(w, 200, map[string]any{"entries": entries, "count": len(entries)})
}

func (s *Server) handleChangelogPage(w http.ResponseWriter, r *http.Request) {
	entries, _ := s.db.ListEntries(true)
	var content strings.Builder
	for _, e := range entries {
		content.WriteString(`<div style="margin-bottom:2rem;padding-bottom:1.5rem;border-bottom:1px solid #2e261e">`)
		content.WriteString(`<div style="display:flex;justify-content:space-between;align-items:baseline;margin-bottom:.5rem">`)
		content.WriteString(`<h2 style="font-family:'JetBrains Mono',monospace;font-size:1rem;color:#f0e6d3;margin:0">`)
		content.WriteString(htmlEscape(e.Title))
		content.WriteString(`</h2>`)
		if e.Version != "" {
			content.WriteString(`<span style="font-family:'JetBrains Mono',monospace;font-size:.7rem;color:#c45d2c;border:1px solid #c45d2c;padding:.1rem .5rem">` + htmlEscape(e.Version) + `</span>`)
		}
		content.WriteString(`</div>`)
		if e.PublishedAt != "" {
			content.WriteString(`<div style="font-family:'JetBrains Mono',monospace;font-size:.65rem;color:#7a7060;margin-bottom:.5rem">` + e.PublishedAt[:10] + `</div>`)
		}
		if e.Tags != "" {
			for _, tag := range strings.Split(e.Tags, ",") {
				tag = strings.TrimSpace(tag)
				content.WriteString(`<span style="font-family:'JetBrains Mono',monospace;font-size:.55rem;background:#241e18;color:#a0845c;padding:.1rem .4rem;margin-right:.3rem">` + htmlEscape(tag) + `</span>`)
			}
		}
		content.WriteString(`<div style="margin-top:.8rem;color:#bfb5a3;line-height:1.6;white-space:pre-wrap">` + htmlEscape(e.Content) + `</div>`)
		content.WriteString(`</div>`)
	}
	if len(entries) == 0 {
		content.WriteString(`<p style="color:#7a7060;text-align:center;padding:3rem;font-style:italic">No entries yet.</p>`)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, changelogPageTemplate, content.String())
}

func (s *Server) handleRSS(w http.ResponseWriter, r *http.Request) {
	items := s.db.RSSItems()
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
<title>Changelog</title>
<link>http://localhost:%d/changelog</link>
<description>Product changelog and release notes</description>
%s
</channel>
</rss>`, s.port, items)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Stats())
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

const changelogPageTemplate = `<!DOCTYPE html><html lang="en"><head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Changelog</title>
<link rel="alternate" type="application/rss+xml" title="Changelog RSS" href="/changelog/rss">
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>
body{background:#1a1410;color:#f0e6d3;font-family:'Libre Baskerville',Georgia,serif;margin:0;padding:0;min-height:100vh}
.container{max-width:700px;margin:0 auto;padding:2rem 1.5rem}
.header{text-align:center;margin-bottom:2.5rem}
.header h1{font-family:'JetBrains Mono',monospace;font-size:1.2rem;letter-spacing:2px;color:#f0e6d3}
.header p{font-size:.8rem;color:#7a7060;margin-top:.3rem}
.subscribe{text-align:center;margin-bottom:2rem}
.subscribe form{display:inline-flex;gap:.5rem}
.subscribe input{font-family:'JetBrains Mono',monospace;font-size:.75rem;background:#2e261e;border:1px solid #2e261e;color:#f0e6d3;padding:.4rem .8rem;outline:none;width:200px}
.subscribe button{font-family:'JetBrains Mono',monospace;font-size:.72rem;background:transparent;border:1px solid #c45d2c;color:#e8753a;padding:.4rem .8rem;cursor:pointer}
.footer{text-align:center;margin-top:2rem;font-size:.6rem;color:#7a7060}
.footer a{color:#e8753a;text-decoration:none}
</style></head><body>
<div class="container">
<div class="header"><h1>Changelog</h1><p>What's new and what changed</p></div>
<div class="subscribe">
<form method="POST" action="/changelog/subscribe"><input name="email" type="email" placeholder="your@email.com" required><button type="submit">Subscribe</button></form>
</div>
%s
<div class="footer">Powered by <a href="https://stockyard.dev/almanac/">Stockyard Almanac</a> · <a href="/changelog/rss">RSS</a></div>
</div></body></html>`
