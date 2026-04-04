package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Almanac</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.7}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.streak{font-family:var(--mono);font-size:.7rem;color:var(--gold)}
.stats-row{display:grid;grid-template-columns:repeat(4,1fr);gap:.6rem;padding:1rem 1.5rem;border-bottom:1px solid var(--bg3)}
.stat{text-align:center}.stat-val{font-family:var(--mono);font-size:1.3rem;color:var(--cream)}.stat-label{font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px}
.wrap{display:grid;grid-template-columns:280px 1fr;max-width:1000px;margin:0 auto;min-height:calc(100vh - 150px)}
@media(max-width:700px){.wrap{grid-template-columns:1fr}.sidebar{border-bottom:1px solid var(--bg3);border-right:none;max-height:300px;overflow-y:auto}}
.sidebar{border-right:1px solid var(--bg3);padding:1rem;overflow-y:auto}
.search{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem;margin-bottom:.8rem}
.entry-item{padding:.6rem;border-bottom:1px solid var(--bg3);cursor:pointer;transition:background .1s}
.entry-item:hover{background:var(--bg2)}
.entry-item.active{background:var(--bg2);border-left:2px solid var(--rust)}
.entry-date{font-family:var(--mono);font-size:.6rem;color:var(--cm)}
.entry-preview{font-size:.78rem;color:var(--cd);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.entry-mood{font-size:.7rem}
.editor{padding:1.5rem}
.editor-date{font-family:var(--mono);font-size:.7rem;color:var(--cm);margin-bottom:.3rem}
.editor-title{width:100%;background:transparent;border:none;color:var(--cream);font-family:var(--serif);font-size:1.3rem;padding:0;margin-bottom:.5rem;outline:none}
.editor-body{width:100%;min-height:300px;background:transparent;border:none;color:var(--cd);font-family:var(--serif);font-size:.95rem;line-height:1.8;resize:none;outline:none}
.editor-footer{display:flex;gap:.8rem;align-items:center;margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.mood-btn{font-size:1.1rem;cursor:pointer;opacity:.4;transition:opacity .15s;background:none;border:none}
.mood-btn:hover,.mood-btn.active{opacity:1}
.tag-input{font-family:var(--mono);font-size:.7rem;padding:.2rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);flex:1}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.btn-new{position:fixed;bottom:1.5rem;right:1.5rem;width:44px;height:44px;border-radius:50%;background:var(--rust);color:var(--bg);border:none;font-size:1.2rem;cursor:pointer;box-shadow:0 4px 12px rgba(0,0,0,.4)}
.wc{font-family:var(--mono);font-size:.6rem;color:var(--cm)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic}
</style></head><body>
<div class="hdr"><h1>ALMANAC</h1><div class="streak" id="streak"></div></div>
<div class="stats-row" id="statsRow"></div>
<div class="wrap">
<div class="sidebar"><input class="search" id="search" placeholder="Search entries..." oninput="doSearch(this.value)"><div id="list"></div></div>
<div class="editor" id="editor"><div class="empty">Select an entry or create a new one</div></div>
</div>
<button class="btn-new" onclick="newEntry()">+</button>
<script>
const A='/api';let entries=[],current=null;const moods=['😊','😐','😔','😤','🤔','💪','😴'];
async function load(){const[e,s]=await Promise.all([fetch(A+'/entries?limit=100').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
entries=e.entries||[];
document.getElementById('streak').innerHTML='🔥 '+s.streak+' day streak';
document.getElementById('statsRow').innerHTML='<div class="stat"><div class="stat-val">'+s.entries+'</div><div class="stat-label">Entries</div></div><div class="stat"><div class="stat-val">'+fmt(s.words)+'</div><div class="stat-label">Words</div></div><div class="stat"><div class="stat-val">'+s.streak+'</div><div class="stat-label">Streak</div></div><div class="stat"><div class="stat-val">'+s.months+'</div><div class="stat-label">Months</div></div>';
renderList();if(!current&&entries.length)selectEntry(entries[0].id);}
function renderList(){const l=document.getElementById('list');
if(!entries.length){l.innerHTML='<div class="empty">No entries yet</div>';return;}
l.innerHTML=entries.map(e=>'<div class="entry-item'+(current&&current.id===e.id?' active':'')+'" onclick="selectEntry(\''+e.id+'\')"><div class="entry-date">'+(e.mood?'<span class="entry-mood">'+e.mood+' </span>':'')+e.date+'</div><div class="entry-preview">'+(e.title||esc(e.body.substring(0,60))||'(empty)')+'</div></div>').join('');}
async function selectEntry(id){const e=await fetch(A+'/entries/'+id).then(r=>r.json());current=e;renderList();renderEditor();}
function renderEditor(){if(!current){document.getElementById('editor').innerHTML='<div class="empty">Select an entry</div>';return;}
const ed=document.getElementById('editor');
ed.innerHTML='<div class="editor-date">'+current.date+' <span class="wc">'+current.word_count+' words</span></div><input class="editor-title" id="ed-title" value="'+esc(current.title||'')+'" placeholder="Title (optional)"><textarea class="editor-body" id="ed-body" placeholder="Write...">'+esc(current.body)+'</textarea><div class="editor-footer"><div>'+moods.map(m=>'<button class="mood-btn'+(current.mood===m?' active':'')+'" onclick="setMood(\''+m+'\')">'+m+'</button>').join('')+'</div><input class="tag-input" id="ed-tags" placeholder="tags (comma sep)" value="'+(current.tags||'')+'"><button class="btn btn-primary" onclick="save()">Save</button><button class="btn" onclick="del()" style="color:#c94444">Delete</button></div>';}
function setMood(m){current.mood=m;renderEditor();}
async function save(){await fetch(A+'/entries/'+current.id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({date:current.date,title:document.getElementById('ed-title').value,body:document.getElementById('ed-body').value,mood:current.mood||'',tags:document.getElementById('ed-tags').value})});load();}
async function del(){if(confirm('Delete this entry?')){await fetch(A+'/entries/'+current.id,{method:'DELETE'});current=null;load();}}
async function newEntry(){await fetch(A+'/entries',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({date:new Date().toISOString().slice(0,10),title:'',body:'',mood:''})});load();}
async function doSearch(q){if(!q){load();return;}const r=await fetch(A+'/entries/search?q='+encodeURIComponent(q)).then(r=>r.json());entries=r.entries||[];renderList();}
function fmt(n){return n>=1000?(n/1000).toFixed(1)+'k':n.toString();}
function esc(s){if(!s)return'';return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');}
load();
</script></body></html>`
