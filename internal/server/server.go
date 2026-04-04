package server

import (
	"encoding/json"
	"net/http"

	"github.com/stockyard-dev/stockyard-almanac/internal/store"
)

type Server struct{ db *store.DB; mux *http.ServeMux; limits Limits }

func New(db *store.DB, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits}
	s.mux.HandleFunc("GET /api/entries", s.list)
	s.mux.HandleFunc("POST /api/entries", s.create)
	s.mux.HandleFunc("GET /api/entries/{id}", s.get)
	s.mux.HandleFunc("PUT /api/entries/{id}", s.update)
	s.mux.HandleFunc("DELETE /api/entries/{id}", s.del)
	s.mux.HandleFunc("GET /api/search", s.search)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"tier": s.limits.Tier, "upgrade_url": "https://stockyard.dev/almanac/"}) })
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
	return s
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func wj(w http.ResponseWriter, c int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(c); json.NewEncoder(w).Encode(v) }
func we(w http.ResponseWriter, c int, m string) { wj(w, c, map[string]string{"error": m}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) { if r.URL.Path != "/" { http.NotFound(w, r); return }; http.Redirect(w, r, "/ui", 302) }
func oe(e []store.Entry) []store.Entry { if e == nil { return []store.Entry{} }; return e }

func (s *Server) list(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"entries": oe(s.db.List(r.URL.Query().Get("month"), 100))}) }
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 && len(s.db.List("", 9999)) >= s.limits.MaxItems { we(w, 402, "Free tier limit reached"); return }
	var e store.Entry; json.NewDecoder(r.Body).Decode(&e); if e.Body == "" { we(w, 400, "body required"); return }
	s.db.Create(&e); wj(w, 201, s.db.Get(e.ID))
}
func (s *Server) get(w http.ResponseWriter, r *http.Request) { e := s.db.Get(r.PathValue("id")); if e == nil { we(w, 404, "not found"); return }; wj(w, 200, e) }
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	existing := s.db.Get(r.PathValue("id")); if existing == nil { we(w, 404, "not found"); return }
	var e store.Entry; json.NewDecoder(r.Body).Decode(&e)
	if e.Body == "" { e.Body = existing.Body }; if e.Date == "" { e.Date = existing.Date }
	s.db.Update(existing.ID, &e); wj(w, 200, s.db.Get(existing.ID))
}
func (s *Server) del(w http.ResponseWriter, r *http.Request) { s.db.Delete(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }
func (s *Server) search(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"entries": oe(s.db.Search(r.URL.Query().Get("q")))}) }
func (s *Server) stats(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) }
func (s *Server) health(w http.ResponseWriter, r *http.Request) { st := s.db.Stats(); wj(w, 200, map[string]any{"service": "almanac", "status": "ok", "entries": st.Entries, "streak": st.Streak}) }
