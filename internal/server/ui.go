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
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.main{max-width:700px;margin:0 auto;padding:1.5rem}
.entry{border-bottom:1px solid var(--bg3);padding:1.2rem 0;cursor:pointer}
.entry:hover{background:var(--bg2);margin:0 -1rem;padding:1.2rem 1rem}
.entry-date{font-family:var(--mono);font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px}
.entry-title{font-size:1.05rem;margin:.2rem 0}
.entry-body{font-size:.88rem;color:var(--cd);display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;overflow:hidden}
.entry-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.8rem}
.mood{font-size:.75rem}
.tag{font-size:.55rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cm);font-family:var(--mono)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.editor{display:none;max-width:700px;margin:0 auto;padding:1.5rem}
.editor.open{display:block}
.editor textarea{width:100%;min-height:300px;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--serif);font-size:.95rem;padding:1rem;line-height:1.8;resize:vertical}
.editor input,.editor select{padding:.4rem .6rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.75rem}
.editor-bar{display:flex;gap:.5rem;margin-bottom:.8rem;flex-wrap:wrap;align-items:center}
.editor-bar label{font-family:var(--mono);font-size:.6rem;color:var(--cm)}
.wc{font-family:var(--mono);font-size:.6rem;color:var(--cm);text-align:right;margin-top:.3rem}
.stats{display:flex;gap:1.5rem;font-family:var(--mono);font-size:.65rem;color:var(--cm);margin-bottom:1.5rem;padding-bottom:.8rem;border-bottom:1px solid var(--bg3)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic}
</style></head><body>
<div class="hdr"><h1>ALMANAC</h1><button class="btn btn-p" onclick="newEntry()">+ New Entry</button></div>
<div class="editor" id="editor">
  <div class="editor-bar">
    <label>Date</label><input id="e-date" type="date">
    <label>Title</label><input id="e-title" placeholder="optional" style="width:150px">
    <label>Mood</label><select id="e-mood"><option value="">—</option><option value="😊">😊 great</option><option value="🙂">🙂 good</option><option value="😐">😐 ok</option><option value="😔">😔 low</option><option value="😤">😤 frustrated</option></select>
    <label>Tags</label><input id="e-tags" placeholder="comma separated" style="width:120px">
  </div>
  <textarea id="e-body" placeholder="Write freely..."></textarea>
  <div style="display:flex;justify-content:space-between;margin-top:.5rem">
    <div class="wc" id="e-wc">0 words</div>
    <div style="display:flex;gap:.4rem"><button class="btn" onclick="cancelEdit()">Cancel</button><button class="btn btn-p" onclick="saveEntry()">Save</button></div>
  </div>
</div>
<div class="main" id="main"></div>
<script>
const A='/api';let entries=[],editId=null;
async function load(){
  const r=await fetch(A+'/entries').then(r=>r.json());
  entries=r.entries||[];render();
}
function render(){
  const m=document.getElementById('main');
  if(!entries.length){m.innerHTML='<div class="empty">Your journal is empty. Write your first entry.</div>';return;}
  let h='<div class="stats"><span>'+entries.length+' entries</span><span>'+entries.reduce((s,e)=>s+(e.body?e.body.split(/\s+/).length:0),0)+' total words</span></div>';
  entries.forEach(e=>{
    const moods={'😊':'great','🙂':'good','😐':'ok','😔':'low','😤':'frustrated'};
    h+='<div class="entry" onclick="editEntry(\''+e.id+'\')">';
    h+='<div class="entry-date">'+formatDate(e.date)+(e.mood?' <span class="mood">'+e.mood+'</span>':'')+'</div>';
    if(e.title)h+='<div class="entry-title">'+esc(e.title)+'</div>';
    h+='<div class="entry-body">'+esc(e.body)+'</div>';
    h+='<div class="entry-meta"><span>'+wordCount(e.body)+' words</span>';
    if(e.tags){e.tags.split(',').forEach(t=>{if(t.trim())h+='<span class="tag">'+esc(t.trim())+'</span>';});}
    h+='<span onclick="event.stopPropagation();delEntry(\''+e.id+'\')" style="color:var(--rust);cursor:pointer">delete</span>';
    h+='</div></div>';
  });
  m.innerHTML=h;
}
function newEntry(){
  editId=null;
  document.getElementById('e-date').value=new Date().toISOString().split('T')[0];
  document.getElementById('e-title').value='';
  document.getElementById('e-body').value='';
  document.getElementById('e-mood').value='';
  document.getElementById('e-tags').value='';
  document.getElementById('e-wc').textContent='0 words';
  document.getElementById('editor').classList.add('open');
  document.getElementById('e-body').focus();
}
function editEntry(id){
  const e=entries.find(x=>x.id===id);if(!e)return;
  editId=id;
  document.getElementById('e-date').value=e.date;
  document.getElementById('e-title').value=e.title||'';
  document.getElementById('e-body').value=e.body||'';
  document.getElementById('e-mood').value=e.mood||'';
  document.getElementById('e-tags').value=e.tags||'';
  document.getElementById('e-wc').textContent=wordCount(e.body)+' words';
  document.getElementById('editor').classList.add('open');
}
function cancelEdit(){document.getElementById('editor').classList.remove('open');editId=null;}
async function saveEntry(){
  const data={date:document.getElementById('e-date').value,title:document.getElementById('e-title').value,body:document.getElementById('e-body').value,mood:document.getElementById('e-mood').value,tags:document.getElementById('e-tags').value};
  if(editId){await fetch(A+'/entries/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)});}
  else{await fetch(A+'/entries',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)});}
  cancelEdit();load();
}
async function delEntry(id){if(confirm('Delete this entry?')){await fetch(A+'/entries/'+id,{method:'DELETE'});load();}}
document.getElementById('e-body').addEventListener('input',function(){document.getElementById('e-wc').textContent=wordCount(this.value)+' words';});
function wordCount(s){return s?s.trim().split(/\s+/).filter(w=>w).length:0;}
function formatDate(d){const dt=new Date(d+'T12:00:00');return dt.toLocaleDateString('en-US',{weekday:'long',year:'numeric',month:'long',day:'numeric'});}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
