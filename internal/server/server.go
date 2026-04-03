package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-almanac/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/entries",s.list);s.mux.HandleFunc("POST /api/entries",s.create);s.mux.HandleFunc("GET /api/entries/{id}",s.get);s.mux.HandleFunc("PUT /api/entries/{id}",s.update);s.mux.HandleFunc("DELETE /api/entries/{id}",s.del)
s.mux.HandleFunc("GET /api/search",s.search);s.mux.HandleFunc("GET /api/streak",s.streak)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/almanac/"})})
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"entries":oe(s.db.List(r.URL.Query().Get("month"),100))})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){if s.limits.MaxItems>0{items:=s.db.List();if len(items)>=s.limits.MaxItems{we(w,402,"Free tier limit reached. Upgrade at https://stockyard.dev/almanac/");return}};var e store.Entry;json.NewDecoder(r.Body).Decode(&e);if e.Body==""{we(w,400,"body required");return};s.db.Create(&e);wj(w,201,s.db.Get(e.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){e:=s.db.Get(r.PathValue("id"));if e==nil{we(w,404,"not found");return};wj(w,200,e)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var e store.Entry;json.NewDecoder(r.Body).Decode(&e);if e.Date==""{e.Date=ex.Date};s.db.Update(id,&e);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)search(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"entries":oe(s.db.Search(r.URL.Query().Get("q")))})}
func(s *Server)streak(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]int{"streak":s.db.Streak()})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"almanac","entries":st.Entries,"streak":st.Streak})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
