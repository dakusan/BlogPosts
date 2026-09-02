document.querySelectorAll('*').forEach(el => el.style.display='none');
for(let c=document.getElementsByTagName('canvas')[0];c.parentNode;c=c.parentNode) c.style.display='';