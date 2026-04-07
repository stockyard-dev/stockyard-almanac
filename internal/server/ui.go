package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Almanac</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.7}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{max-width:720px;margin:0 auto;padding:1.5rem}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1.2rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center;font-family:var(--mono)}
.st-v{font-size:1.3rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.entry{border-bottom:1px solid var(--bg3);padding:1.2rem 0;cursor:pointer;transition:background .15s}
.entry:hover{background:var(--bg2);margin:0 -1rem;padding:1.2rem 1rem}
.entry-date{font-family:var(--mono);font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;display:flex;align-items:center;gap:.5rem;flex-wrap:wrap}
.entry-mood{font-size:.75rem;color:var(--cream)}
.entry-title{font-size:1.15rem;margin:.3rem 0;font-weight:700}
.entry-body{font-size:.88rem;color:var(--cd);display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;overflow:hidden;line-height:1.6}
.entry-meta{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.entry-extra{font-family:var(--mono);font-size:.58rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.entry-extra-row{display:flex;gap:.4rem}
.entry-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.entry-extra-val{color:var(--cream)}
.tag{font-size:.55rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cd);font-family:var(--mono)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .7rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.editor{display:none;position:fixed;inset:0;background:var(--bg);z-index:50;overflow-y:auto}
.editor.open{display:block}
.editor-inner{max-width:720px;margin:0 auto;padding:1.5rem}
.editor-bar{display:flex;gap:.5rem;margin-bottom:.8rem;flex-wrap:wrap;align-items:center}
.editor-bar label{font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px}
.editor-bar input,.editor-bar select{padding:.35rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.editor textarea{width:100%;min-height:420px;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--serif);font-size:1rem;padding:1.2rem;line-height:1.8;resize:vertical}
.editor textarea:focus{outline:none;border-color:var(--leather)}
.editor-extras{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.editor-extras-label{font-family:var(--mono);font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.editor-foot{display:flex;justify-content:space-between;margin-top:.8rem;align-items:center}
.wc{font-family:var(--mono);font-size:.6rem;color:var(--cm)}
.editor-acts{display:flex;gap:.4rem}
.editor-del{font-family:var(--mono);font-size:.6rem;color:var(--red);background:transparent;border:1px solid #3a1a1a;padding:.3rem .6rem;cursor:pointer;margin-right:auto}
.editor-del:hover{border-color:var(--red)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
.count-label{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-bottom:.5rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}.toolbar{flex-direction:column;align-items:stretch}.search{min-width:100%}.editor textarea{min-height:300px}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> ALMANAC</h1>
<button class="btn btn-p" onclick="newEntry()">+ New Entry</button>
</div>

<div class="main" id="main-view">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search title, body, tags..." oninput="debouncedSearch()">
<select class="filter-sel" id="month-filter" onchange="changeMonth()">
<option value="">All Time</option>
</select>
</div>
<div class="count-label" id="count"></div>
<div id="list"></div>
</div>

<div class="editor" id="editor">
<div class="editor-inner">
<div class="editor-bar">
<label>Date</label><input id="e-date" type="date">
<label>Title</label><input id="e-title" placeholder="optional" style="width:200px">
<label>Mood</label><select id="e-mood"></select>
<label>Tags</label><input id="e-tags" placeholder="comma separated" style="width:160px">
</div>
<textarea id="e-body" placeholder="Write freely..."></textarea>
<div id="e-extras-section"></div>
<div class="editor-foot">
<button class="editor-del" id="e-del" onclick="delEntry()">Delete</button>
<div class="wc" id="e-wc">0 words</div>
<div class="editor-acts">
<button class="btn" onclick="cancelEdit()">Cancel</button>
<button class="btn btn-p" onclick="saveEntry()">Save</button>
</div>
</div>
</div>
</div>

<script>
var A='/api';
var RESOURCE='entries';

// Default mood options. Overridable from /api/config.moods.
var defaultMoods=[
{value:'',label:'—'},
{value:'😊',label:'😊 great'},
{value:'🙂',label:'🙂 good'},
{value:'😐',label:'😐 ok'},
{value:'😔',label:'😔 low'},
{value:'😤',label:'😤 frustrated'}
];
var moods=defaultMoods.slice();

// Custom fields injected by personalization (extras keys)
var customFields=[];

var entries=[],editId=null,curMonth='',curSearch='',taskExtras={},searchTimer=null;

// ─── Helpers ──────────────────────────────────────────────────────

function fmtDate(d){
if(!d)return'';
try{
var dt=new Date(d+'T12:00:00');
if(isNaN(dt.getTime()))return d;
return dt.toLocaleDateString('en-US',{weekday:'long',year:'numeric',month:'long',day:'numeric'});
}catch(e){return d}
}

function fmtMonthLabel(ym){
try{
var dt=new Date(ym+'-01T12:00:00');
return dt.toLocaleDateString('en-US',{month:'long',year:'numeric'});
}catch(e){return ym}
}

function wordCount(s){
if(!s)return 0;
return s.trim().split(/\s+/).filter(function(w){return w}).length;
}

function debouncedSearch(){
clearTimeout(searchTimer);
searchTimer=setTimeout(function(){
curSearch=document.getElementById('search').value;
load();
},300);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var url;
if(curSearch){
url=A+'/search?q='+encodeURIComponent(curSearch);
}else{
url=A+'/entries';
if(curMonth)url+='?month='+encodeURIComponent(curMonth);
}
var resps=await Promise.all([
fetch(url).then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
entries=resps[0].entries||[];
var stats=resps[1]||{};

// Merge extras
try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
taskExtras=ex||{};
entries.forEach(function(e){
var x=taskExtras[e.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(e[k]===undefined)e[k]=x[k]});
});
}catch(e){taskExtras={}}

renderStats(stats);
populateMonthFilter();
}catch(e){
console.error('load failed',e);
entries=[];
}
render();
}

function populateMonthFilter(){
var sel=document.getElementById('month-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var months=[];
entries.forEach(function(e){
if(!e.date)return;
var m=String(e.date).substring(0,7);
if(!seen[m]){seen[m]=true;months.push(m)}
});
months.sort().reverse();
sel.innerHTML='<option value="">All Time</option>'+months.map(function(m){return'<option value="'+esc(m)+'"'+(m===current?' selected':'')+'>'+esc(fmtMonthLabel(m))+'</option>'}).join('');
}

function changeMonth(){
curMonth=document.getElementById('month-filter').value;
curSearch='';
document.getElementById('search').value='';
load();
}

function renderStats(s){
var entries=s.entries||0;
var words=s.words||0;
var streak=s.streak||0;
var months=s.months||0;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+entries+'</div><div class="st-l">Entries</div></div>'+
'<div class="st"><div class="st-v">'+words.toLocaleString()+'</div><div class="st-l">Words</div></div>'+
'<div class="st"><div class="st-v">'+streak+'</div><div class="st-l">Day Streak</div></div>'+
'<div class="st"><div class="st-v">'+months+'</div><div class="st-l">Months</div></div>';
}

function render(){
document.getElementById('count').textContent=entries.length+' entr'+(entries.length!==1?'ies':'y');
if(!entries.length){
var msg=window._emptyMsg||'Your journal is empty. Write your first entry.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}
var h='';
entries.forEach(function(e){h+=entryHTML(e)});
document.getElementById('list').innerHTML=h;
}

function entryHTML(e){
var h='<div class="entry" onclick="editEntry(\''+e.id+'\')">';
h+='<div class="entry-date">'+esc(fmtDate(e.date));
if(e.mood)h+='<span class="entry-mood">'+esc(e.mood)+'</span>';
h+='</div>';
if(e.title)h+='<div class="entry-title">'+esc(e.title)+'</div>';
h+='<div class="entry-body">'+esc(e.body)+'</div>';
h+='<div class="entry-meta">';
h+='<span>'+wordCount(e.body)+' words</span>';
if(e.tags){
String(e.tags).split(',').forEach(function(t){
t=t.trim();
if(t)h+='<span class="tag">#'+esc(t)+'</span>';
});
}
h+='</div>';

// Custom field display
var customRows='';
customFields.forEach(function(f){
var v=e[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="entry-extra-row">';
customRows+='<span class="entry-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="entry-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="entry-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Editor ───────────────────────────────────────────────────────

function renderMoodSelect(){
var sel=document.getElementById('e-mood');
sel.innerHTML='';
moods.forEach(function(m){
var opt=document.createElement('option');
opt.value=m.value;
opt.textContent=m.label;
sel.appendChild(opt);
});
}

function renderExtrasSection(values){
var section=document.getElementById('e-extras-section');
if(!customFields.length){
section.innerHTML='';
return;
}
var label=window._customSectionLabel||'Additional Details';
var h='<div class="editor-extras"><div class="editor-extras-label">'+esc(label)+'</div>';
customFields.forEach(function(f){
var v=values&&values[f.name]!==undefined?values[f.name]:'';
h+='<div class="fr"><label>'+esc(f.label)+'</label>';
if(f.type==='select'){
h+='<select id="ex-'+f.name+'">';
h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='textarea'){
h+='<textarea id="ex-'+f.name+'" rows="2">'+esc(String(v))+'</textarea>';
}else if(f.type==='number'||f.type==='integer'){
h+='<input type="number" id="ex-'+f.name+'" value="'+esc(String(v))+'">';
}else{
h+='<input type="text" id="ex-'+f.name+'" value="'+esc(String(v))+'">';
}
h+='</div>';
});
h+='</div>';
section.innerHTML=h;
}

function newEntry(){
editId=null;
document.getElementById('e-date').value=new Date().toISOString().split('T')[0];
document.getElementById('e-title').value='';
document.getElementById('e-body').value='';
document.getElementById('e-mood').value='';
document.getElementById('e-tags').value='';
document.getElementById('e-wc').textContent='0 words';
document.getElementById('e-del').style.display='none';
renderExtrasSection({});
document.getElementById('editor').classList.add('open');
document.getElementById('e-body').focus();
}

function editEntry(id){
var e=null;
for(var i=0;i<entries.length;i++){if(entries[i].id===id){e=entries[i];break}}
if(!e)return;
editId=id;
document.getElementById('e-date').value=e.date||'';
document.getElementById('e-title').value=e.title||'';
document.getElementById('e-body').value=e.body||'';
document.getElementById('e-mood').value=e.mood||'';
document.getElementById('e-tags').value=e.tags||'';
document.getElementById('e-wc').textContent=wordCount(e.body)+' words';
document.getElementById('e-del').style.display='inline-block';
renderExtrasSection(e);
document.getElementById('editor').classList.add('open');
document.getElementById('e-body').focus();
}

function cancelEdit(){
document.getElementById('editor').classList.remove('open');
editId=null;
}

async function saveEntry(){
var bodyText=document.getElementById('e-body').value;
if(!bodyText.trim()){alert('Body is required');return}

var data={
date:document.getElementById('e-date').value,
title:document.getElementById('e-title').value,
body:bodyText,
mood:document.getElementById('e-mood').value,
tags:document.getElementById('e-tags').value
};

var extras={};
customFields.forEach(function(f){
var el=document.getElementById('ex-'+f.name);
if(!el)return;
var val;
if(f.type==='number'||f.type==='integer')val=parseFloat(el.value)||0;
else val=el.value.trim();
extras[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/entries/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/entries',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Save failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
cancelEdit();
load();
}

async function delEntry(){
if(!editId)return;
if(!confirm('Delete this entry?'))return;
await fetch(A+'/entries/'+editId,{method:'DELETE'});
cancelEdit();
load();
}

document.getElementById('e-body').addEventListener('input',function(){
document.getElementById('e-wc').textContent=wordCount(this.value)+' words';
});

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape'&&document.getElementById('editor').classList.contains('open'))cancelEdit()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

// Custom moods replace the default emoji set. Each mood can be a simple
// string ("focused") or an object {value, label}.
if(Array.isArray(cfg.moods)&&cfg.moods.length){
moods=[{value:'',label:'—'}];
cfg.moods.forEach(function(m){
if(typeof m==='string'){
moods.push({value:m,label:m});
}else if(m&&m.value){
moods.push({value:m.value,label:m.label||m.value});
}
});
}
renderMoodSelect();

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
customFields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[]
});
});
}
}).catch(function(){
}).finally(function(){
renderMoodSelect();
load();
});
})();
</script>
</body>
</html>`
