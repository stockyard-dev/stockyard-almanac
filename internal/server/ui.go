package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Almanac</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.7}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.stats-bar{display:flex;gap:1.5rem;font-family:var(--mono);font-size:.7rem;color:var(--cm)}
.stats-bar strong{color:var(--gold)}
.ct{display:grid;grid-template-columns:280px 1fr;min-height:calc(100vh - 55px)}
@media(max-width:700px){.ct{grid-template-columns:1fr}}
.sidebar{border-right:1px solid var(--bg3);padding:1rem}
.editor{padding:1.5rem}
.entry-item{padding:.6rem .8rem;cursor:pointer;border-bottom:1px solid var(--bg3);transition:background .1s}
.entry-item:hover{background:var(--bg2)}
.entry-item.active{background:var(--bg2);border-left:2px solid var(--rust)}
.entry-date{font-family:var(--mono);font-size:.7rem;color:var(--leather)}
.entry-preview{font-size:.78rem;color:var(--cm);white-space:nowrap;overflow:hidden;text-overflow:ellipsis;margin-top:.1rem}
.entry-mood{font-size:.7rem;margin-top:.1rem}
.ed-date{font-family:var(--mono);font-size:.7rem;color:var(--leather);margin-bottom:.3rem}
.ed-title{width:100%;background:transparent;border:none;color:var(--cream);font-family:var(--serif);font-size:1.2rem;padding:.3rem 0;border-bottom:1px solid var(--bg3);outline:none;margin-bottom:.5rem}
.ed-body{width:100%;min-height:300px;background:transparent;border:none;color:var(--cd);font-family:var(--serif);font-size:.95rem;line-height:1.8;outline:none;resize:vertical}
.ed-meta{display:flex;gap:1rem;margin-top:1rem;font-family:var(--mono);font-size:.7rem;align-items:center;flex-wrap:wrap}
.ed-meta label{color:var(--cm)}
.ed-meta input,.ed-meta select{background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem;padding:.2rem .4rem}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}.btn-p:hover{opacity:.85}
.search{width:100%;padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem;margin-bottom:.5rem;outline:none}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic}
.wc{font-family:var(--mono);font-size:.6rem;color:var(--cm)}
</style></head><body>
<div class="hdr"><h1>ALMANAC</h1><div class="stats-bar" id="st"></div></div>
<div class="ct">
<div class="sidebar">
  <div style="display:flex;gap:.4rem;margin-bottom:.8rem"><button class="btn btn-p" onclick="newEntry()" style="flex:1">+ New Entry</button></div>
  <input class="search" id="sq" placeholder="Search entries..." oninput="doSearch(this.value)">
  <div id="list"></div>
</div>
<div class="editor" id="editor"><div class="empty">Select an entry or create a new one</div></div>
</div>
<script>
const A='/api';let entries=[],current=null;
async function ld(){const[e,s]=await Promise.all([fetch(A+'/entries').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);entries=e.entries||[];
document.getElementById('st').innerHTML='<span><strong>'+s.streak+'</strong> day streak</span><span><strong>'+s.entries+'</strong> entries</span><span><strong>'+s.words+'</strong> words</span>';
rnList();}
function rnList(){
  const m=document.getElementById('list');
  if(!entries.length){m.innerHTML='<div class="empty" style="padding:1rem">No entries yet</div>';return;}
  let h='';entries.forEach(e=>{
    const moods={'great':'😊','good':'🙂','okay':'😐','bad':'😞','awful':'😢'};
    h+='<div class="entry-item'+(current&&current.id===e.id?' active':'')+'" onclick="sel(\''+e.id+'\')"><div class="entry-date">'+e.date+(e.mood?' '+(moods[e.mood]||e.mood):'')+'</div>';
    h+='<div class="entry-preview">'+(e.title||esc(e.body.substring(0,60)))+'</div></div>';
  });m.innerHTML=h;
}
function sel(id){current=entries.find(e=>e.id===id);rnList();rnEditor();}
function rnEditor(){
  if(!current){document.getElementById('editor').innerHTML='<div class="empty">Select an entry</div>';return;}
  const e=current;
  document.getElementById('editor').innerHTML='<div class="ed-date">'+e.date+'</div><input class="ed-title" id="et" value="'+esc(e.title||'')+'" placeholder="Title (optional)"><textarea class="ed-body" id="eb">'+esc(e.body)+'</textarea><div class="ed-meta"><label>Mood</label><select id="em"><option value="">—</option><option value="great"'+(e.mood==='great'?' selected':'')+'>😊 Great</option><option value="good"'+(e.mood==='good'?' selected':'')+'>🙂 Good</option><option value="okay"'+(e.mood==='okay'?' selected':'')+'>😐 Okay</option><option value="bad"'+(e.mood==='bad'?' selected':'')+'>😞 Bad</option><option value="awful"'+(e.mood==='awful'?' selected':'')+'>😢 Awful</option></select><label>Tags</label><input id="eg" value="'+esc(e.tags||'')+'" placeholder="comma separated" style="width:150px"><span class="wc" id="wc">'+wc(e.body)+' words</span><button class="btn btn-p" onclick="save()">Save</button><button class="btn" onclick="del()" style="color:#c94444">Delete</button></div>';
  document.getElementById('eb').addEventListener('input',function(){document.getElementById('wc').textContent=wc(this.value)+' words';});
}
async function save(){
  const body={title:document.getElementById('et').value,body:document.getElementById('eb').value,mood:document.getElementById('em').value,tags:document.getElementById('eg').value};
  await fetch(A+'/entries/'+current.id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});ld();
}
async function del(){if(confirm('Delete this entry?')){await fetch(A+'/entries/'+current.id,{method:'DELETE'});current=null;ld();rnEditor();}}
async function newEntry(){
  const today=new Date().toISOString().split('T')[0];
  const r=await fetch(A+'/entries',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({date:today,body:''})}).then(r=>r.json());
  await ld();sel(r.id);
}
let searchTimeout;
async function doSearch(q){clearTimeout(searchTimeout);searchTimeout=setTimeout(async()=>{
  if(!q){ld();return;}
  const r=await fetch(A+'/search?q='+encodeURIComponent(q)).then(r=>r.json());entries=r.entries||[];rnList();
},300);}
function wc(s){return s?s.trim().split(/\s+/).filter(w=>w).length:0;}
function esc(s){if(!s)return'';return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');}
ld();
</script></body></html>`
