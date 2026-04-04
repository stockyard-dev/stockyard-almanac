package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Almanac</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.7}
.header{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.header h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.content{display:grid;grid-template-columns:280px 1fr;height:calc(100vh - 56px)}
@media(max-width:700px){.content{grid-template-columns:1fr}.sidebar{max-height:200px;overflow-y:auto;border-bottom:1px solid var(--bg3);border-right:none}}
.sidebar{border-right:1px solid var(--bg3);overflow-y:auto}
.sidebar-header{padding:.8rem 1rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.sidebar-header input{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem;padding:.3rem .5rem;width:120px}
.stats-bar{display:flex;gap:1rem;padding:.6rem 1rem;border-bottom:1px solid var(--bg3);font-family:var(--mono);font-size:.6rem;color:var(--cm)}
.stat-item span{color:var(--cream)}
.entry-item{padding:.6rem 1rem;border-bottom:1px solid var(--bg3);cursor:pointer;transition:background .1s}
.entry-item:hover,.entry-item.active{background:var(--bg2)}
.entry-date{font-family:var(--mono);font-size:.65rem;color:var(--cm)}
.entry-preview{font-size:.78rem;color:var(--cd);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.entry-mood{font-size:.7rem;margin-left:.3rem}
.editor{padding:1.5rem;overflow-y:auto}
.editor-header{display:flex;justify-content:space-between;align-items:flex-start;margin-bottom:1rem}
.editor-date{font-family:var(--mono);font-size:.75rem;color:var(--leather)}
.editor-meta{display:flex;gap:.5rem;margin-bottom:.8rem;flex-wrap:wrap}
.editor-meta input,.editor-meta select{background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem;padding:.3rem .5rem}
.editor textarea{width:100%;min-height:300px;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--serif);font-size:.95rem;padding:1rem;line-height:1.8;resize:vertical}
.editor textarea:focus{outline:none;border-color:var(--leather)}
.word-count{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.3rem;text-align:right}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-primary{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-primary:hover{opacity:.85}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.empty{padding:3rem 2rem;text-align:center;color:var(--cm);font-style:italic}
.tag{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cm);margin-right:.2rem}
</style></head><body>
<div class="header"><h1>ALMANAC</h1><button class="btn btn-primary" onclick="newEntry()">+ New Entry</button></div>
<div class="content">
<div class="sidebar">
  <div class="stats-bar" id="statsBar"></div>
  <div class="sidebar-header"><input id="search" placeholder="Search..." oninput="searchEntries(this.value)"></div>
  <div id="entryList"></div>
</div>
<div class="editor" id="editor"><div class="empty">Select an entry or create a new one.</div></div>
</div>

<script>
const API='/api';let entries=[],current=null;

async function load(){
  const[e,s]=await Promise.all([fetch(API+'/entries').then(r=>r.json()),fetch(API+'/stats').then(r=>r.json())]);
  entries=e.entries||[];
  document.getElementById('statsBar').innerHTML='<div class="stat-item"><span>'+s.entries+'</span> entries</div><div class="stat-item"><span>'+s.words+'</span> words</div><div class="stat-item"><span>'+s.streak+'</span> day streak</div>';
  renderList();
}

function renderList(){
  const el=document.getElementById('entryList');
  if(!entries||!entries.length){el.innerHTML='<div class="empty">No entries yet.</div>';return;}
  let h='';
  (entries||[]).forEach(e=>{
    const moods={'great':'😊','good':'🙂','okay':'😐','bad':'😔','terrible':'😢'};
    const preview=e.title||e.body.slice(0,60)||'(empty)';
    h+='<div class="entry-item'+(current&&current.id===e.id?' active':'')+'" onclick="selectEntry(\''+e.id+'\')"><div class="entry-date">'+e.date+(e.mood?' '+( moods[e.mood]||e.mood):'')+'</div><div class="entry-preview">'+esc(preview)+'</div>';
    if(e.tags){e.tags.split(',').forEach(t=>{if(t.trim())h+='<span class="tag">'+esc(t.trim())+'</span>';});}
    h+='</div>';
  });
  el.innerHTML=h;
}

function selectEntry(id){
  current=(entries||[]).find(e=>e.id===id);
  if(!current)return;
  renderEditor();renderList();
}

function renderEditor(){
  if(!current){document.getElementById('editor').innerHTML='<div class="empty">Select an entry or create a new one.</div>';return;}
  const e=current;
  const ed=document.getElementById('editor');
  ed.innerHTML='<div class="editor-header"><div class="editor-date">'+e.date+'</div><div style="display:flex;gap:.3rem"><button class="btn btn-sm btn-primary" onclick="saveEntry()">Save</button><button class="btn btn-sm" onclick="delEntry(\''+e.id+'\')" style="color:#c94444">Delete</button></div></div><div class="editor-meta"><input id="ed-title" placeholder="Title (optional)" value="'+esc(e.title||'')+'"><select id="ed-mood"><option value="">No mood</option><option value="great"'+(e.mood==='great'?' selected':'')+'>😊 Great</option><option value="good"'+(e.mood==='good'?' selected':'')+'>🙂 Good</option><option value="okay"'+(e.mood==='okay'?' selected':'')+'>😐 Okay</option><option value="bad"'+(e.mood==='bad'?' selected':'')+'>😔 Bad</option><option value="terrible"'+(e.mood==='terrible'?' selected':'')+'>😢 Terrible</option></select><input id="ed-tags" placeholder="Tags (comma sep)" value="'+esc(e.tags||'')+'"></div><textarea id="ed-body" oninput="updateWC()">'+esc(e.body||'')+'</textarea><div class="word-count" id="wc">'+wordCount(e.body||'')+' words</div>';
}

function newEntry(){
  const today=new Date().toISOString().slice(0,10);
  current={id:'new',date:today,title:'',body:'',mood:'',tags:''};
  renderEditor();
  // Replace save to create
  document.querySelector('.btn-primary').onclick=createEntry;
}

async function createEntry(){
  const body={date:current.date,title:document.getElementById('ed-title').value,body:document.getElementById('ed-body').value,mood:document.getElementById('ed-mood').value,tags:document.getElementById('ed-tags').value};
  await fetch(API+'/entries',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  current=null;load();
}

async function saveEntry(){
  if(!current||current.id==='new')return createEntry();
  const body={date:current.date,title:document.getElementById('ed-title').value,body:document.getElementById('ed-body').value,mood:document.getElementById('ed-mood').value,tags:document.getElementById('ed-tags').value};
  await fetch(API+'/entries/'+current.id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  load();
}

async function delEntry(id){if(confirm('Delete this entry?')){await fetch(API+'/entries/'+id,{method:'DELETE'});current=null;load();document.getElementById('editor').innerHTML='<div class="empty">Entry deleted.</div>';}}

async function searchEntries(q){
  if(!q){load();return;}
  const r=await fetch(API+'/entries/search?q='+encodeURIComponent(q)).then(r=>r.json());
  entries=r.entries||[];renderList();
}

function updateWC(){const b=document.getElementById('ed-body');document.getElementById('wc').textContent=wordCount(b.value)+' words';}
function wordCount(s){return s.trim()?s.trim().split(/\s+/).length:0;}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
