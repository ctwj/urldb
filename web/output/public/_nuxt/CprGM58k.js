import{_ as ie}from"./D9zP7iJN.js";import{ak as ce,n as q,ba as ue,e as de,f as B,g as me,an as G,h as fe,j as ee,k as X,l as pe,m as g,v as te,s as V,aq as Y,u as se,ad as ge,X as he,c as $,b as d,w as c,Q as A,z as l,F as Z,c7 as _e,o as k,a as s,A as ve,y as F,t as p,d as z,P as xe,_ as ye}from"./DmHPR5lg.js";import{B as be}from"./DO8alW5h.js";import{u as we}from"./sJus4ovQ.js";import{_ as je}from"./DYSVe14K.js";import{_ as $e}from"./BnsUzXib.js";import{_ as ke}from"./CXBNmlmj.js";import{_ as Ce}from"./DMTo9cht.js";import{_ as ze}from"./BK0zW0Pa.js";import{_ as Se}from"./BF7UzXbo.js";import"./QpI9WcJO.js";import"./CoaUF789.js";import"./BzzwCNSt.js";import"./DQDQW_eH.js";import"./C-4Oazno.js";import"./Hg3Q8c_0.js";import"./CGr2hHf7.js";import"./bmQVJIR_.js";import"./BD_kU2HR.js";import"./C3WqSwOh.js";import"./CUndgRjY.js";import"./B-p6aW7q.js";import"./Bjh6Msyv.js";import"./D__gE1zc.js";import"./CM8LO42l.js";import"./CbvP3fHg.js";import"./0vdbVihP.js";import"./Dfqi6zXz.js";import"./BQQJOfr4.js";function Pe(r,t){const f=ce(ue,null);return q(()=>r.hljs||(f==null?void 0:f.mergedHljsRef.value))}function Ne(r){const{textColor2:t,fontSize:f,fontWeightStrong:x,textColor3:h}=r;return{textColor:t,fontSize:f,fontWeightStrong:x,"mono-3":"#a0a1a7","hue-1":"#0184bb","hue-2":"#4078f2","hue-3":"#a626a4","hue-4":"#50a14f","hue-5":"#e45649","hue-5-2":"#c91243","hue-6":"#986801","hue-6-2":"#c18401",lineNumberTextColor:h}}const Te={common:de,self:Ne},Le=B([me("code",`
 font-size: var(--n-font-size);
 font-family: var(--n-font-family);
 `,[G("show-line-numbers",`
 display: flex;
 `),fe("line-numbers",`
 user-select: none;
 padding-right: 12px;
 text-align: right;
 transition: color .3s var(--n-bezier);
 color: var(--n-line-number-text-color);
 `),G("word-wrap",[B("pre",`
 white-space: pre-wrap;
 word-break: break-all;
 `)]),B("pre",`
 margin: 0;
 line-height: inherit;
 font-size: inherit;
 font-family: inherit;
 `),B("[class^=hljs]",`
 color: var(--n-text-color);
 transition: 
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `)]),({props:r})=>{const t=`${r.bPrefix}code`;return[`${t} .hljs-comment,
 ${t} .hljs-quote {
 color: var(--n-mono-3);
 font-style: italic;
 }`,`${t} .hljs-doctag,
 ${t} .hljs-keyword,
 ${t} .hljs-formula {
 color: var(--n-hue-3);
 }`,`${t} .hljs-section,
 ${t} .hljs-name,
 ${t} .hljs-selector-tag,
 ${t} .hljs-deletion,
 ${t} .hljs-subst {
 color: var(--n-hue-5);
 }`,`${t} .hljs-literal {
 color: var(--n-hue-1);
 }`,`${t} .hljs-string,
 ${t} .hljs-regexp,
 ${t} .hljs-addition,
 ${t} .hljs-attribute,
 ${t} .hljs-meta-string {
 color: var(--n-hue-4);
 }`,`${t} .hljs-built_in,
 ${t} .hljs-class .hljs-title {
 color: var(--n-hue-6-2);
 }`,`${t} .hljs-attr,
 ${t} .hljs-variable,
 ${t} .hljs-template-variable,
 ${t} .hljs-type,
 ${t} .hljs-selector-class,
 ${t} .hljs-selector-attr,
 ${t} .hljs-selector-pseudo,
 ${t} .hljs-number {
 color: var(--n-hue-6);
 }`,`${t} .hljs-symbol,
 ${t} .hljs-bullet,
 ${t} .hljs-link,
 ${t} .hljs-meta,
 ${t} .hljs-selector-id,
 ${t} .hljs-title {
 color: var(--n-hue-2);
 }`,`${t} .hljs-emphasis {
 font-style: italic;
 }`,`${t} .hljs-strong {
 font-weight: var(--n-font-weight-strong);
 }`,`${t} .hljs-link {
 text-decoration: underline;
 }`]}]),Ae=Object.assign(Object.assign({},se.props),{language:String,code:{type:String,default:""},trim:{type:Boolean,default:!0},hljs:Object,uri:Boolean,inline:Boolean,wordWrap:Boolean,showLineNumbers:Boolean,internalFontSize:Number,internalNoHighlight:Boolean}),qe=ee({name:"Code",props:Ae,setup(r,{slots:t}){const{internalNoHighlight:f}=r,{mergedClsPrefixRef:x,inlineThemeDisabled:h}=pe(),y=g(null),S=f?{value:void 0}:Pe(r),P=(a,i,u)=>{const{value:m}=S;return!m||!(a&&m.getLanguage(a))?null:m.highlight(u?i.trim():i,{language:a}).value},C=q(()=>r.inline||r.wordWrap?!1:r.showLineNumbers),b=()=>{if(t.default)return;const{value:a}=y;if(!a)return;const{language:i}=r,u=r.uri?window.decodeURIComponent(r.code):r.code;if(i){const v=P(i,u,r.trim);if(v!==null){if(r.inline)a.innerHTML=v;else{const T=a.querySelector(".__code__");T&&a.removeChild(T);const L=document.createElement("pre");L.className="__code__",L.innerHTML=v,a.appendChild(L)}return}}if(r.inline){a.textContent=u;return}const m=a.querySelector(".__code__");if(m)m.textContent=u;else{const v=document.createElement("pre");v.className="__code__",v.textContent=u,a.innerHTML="",a.appendChild(v)}};te(b),V(Y(r,"language"),b),V(Y(r,"code"),b),f||V(S,b);const N=se("Code","-code",Le,Te,r,x),w=q(()=>{const{common:{cubicBezierEaseInOut:a,fontFamilyMono:i},self:{textColor:u,fontSize:m,fontWeightStrong:v,lineNumberTextColor:T,"mono-3":L,"hue-1":E,"hue-2":R,"hue-3":D,"hue-4":I,"hue-5":H,"hue-5-2":O,"hue-6":U,"hue-6-2":M}}=N.value,{internalFontSize:n}=r;return{"--n-font-size":n?`${n}px`:m,"--n-font-family":i,"--n-font-weight-strong":v,"--n-bezier":a,"--n-text-color":u,"--n-mono-3":L,"--n-hue-1":E,"--n-hue-2":R,"--n-hue-3":D,"--n-hue-4":I,"--n-hue-5":H,"--n-hue-5-2":O,"--n-hue-6":U,"--n-hue-6-2":M,"--n-line-number-text-color":T}}),_=h?ge("code",q(()=>`${r.internalFontSize||"a"}`),w,r):void 0;return{mergedClsPrefix:x,codeRef:y,mergedShowLineNumbers:C,lineNumbers:q(()=>{let a=1;const i=[];let u=!1;for(const m of r.code)m===`
`?(u=!0,i.push(a++)):u=!1;return u||i.push(a++),i.join(`
`)}),cssVars:h?void 0:w,themeClass:_==null?void 0:_.themeClass,onRender:_==null?void 0:_.onRender}},render(){var r,t;const{mergedClsPrefix:f,wordWrap:x,mergedShowLineNumbers:h,onRender:y}=this;return y==null||y(),X("code",{class:[`${f}-code`,this.themeClass,x&&`${f}-code--word-wrap`,h&&`${f}-code--show-line-numbers`],style:this.cssVars,ref:"codeRef"},h?X("pre",{class:`${f}-code__line-numbers`},this.lineNumbers):null,(t=(r=this.$slots).default)===null||t===void 0?void 0:t.call(r))}}),Re={class:"flex space-x-3"},Be={class:"bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4"},Ee={class:"grid grid-cols-1 md:grid-cols-4 gap-4"},De={class:"flex items-center justify-between"},Ie={class:"flex items-center space-x-4"},He={class:"text-sm text-gray-500 dark:text-gray-400"},Oe={class:"flex space-x-6"},Ue={class:"text-center flex items-base"},Me={class:"text-2xl font-bold text-blue-600"},Ve={class:"text-center flex items-base"},Fe={class:"text-2xl font-bold text-green-600"},We={class:"text-center flex items-base"},Ke={class:"text-2xl font-bold text-purple-600"},Je={class:"text-center flex items-base"},Qe={class:"text-2xl font-bold text-red-600"},Ge={key:0,class:"flex items-center justify-center py-12"},Xe={key:1,class:"flex flex-col items-center justify-center py-12"},Ye={key:2,class:"space-y-2 h-full overflow-y-auto"},Ze={class:"flex items-start justify-between"},et={class:"flex-1 min-w-0"},tt={class:"flex items-center space-x-3 mb-2 w-full"},st={class:"text-xs text-gray-500 dark:text-gray-400"},nt={class:"text-sm bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded"},ot={class:"flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 float-right"},rt={class:"text-sm bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded"},at={class:"flex items-center"},lt={key:0,class:"flex items-center"},it={key:0,class:"mt-2 text-xs text-gray-600 dark:text-gray-400 flex items-center justify-between"},ct={class:"flex items-center flex-1 min-w-0"},ut={class:"truncate"},dt={class:"flex items-center space-x-1 ml-2"},mt={key:1,class:"mt-2 text-xs text-red-600 dark:text-red-400"},ft={class:"p-4"},pt={class:"flex justify-center"},gt=ee({__name:"api-access-logs",setup(r){const t=he(),f=we(),x=_e(),h=g(!1),y=g(!1),S=g([]),P=g({total_requests:0,today_requests:0,week_requests:0,month_requests:0,error_requests:0,unique_ips:0}),C=g(""),b=g(null),N=g(null),w=g(1),_=g(20),a=g(0),i=g(!1),u=g(""),m=async()=>{h.value=!0;try{const n={page:w.value,page_size:_.value};if(b.value){const j=new Date(b.value);n.start_date=j.toISOString().split("T")[0]}if(N.value){const j=new Date(N.value);n.end_date=j.toISOString().split("T")[0]}C.value&&(n.endpoint=C.value,n.ip=C.value);const e=await x.getApiAccessLogs(n);S.value=e.data||[],a.value=e.total||0}catch(n){console.error("获取API访问日志失败:",n),t.error({content:"获取API访问日志失败",duration:3e3}),S.value=[],a.value=0}finally{h.value=!1}},v=async()=>{try{const n=await x.getApiAccessLogSummary();P.value=n}catch(n){console.error("获取统计汇总失败:",n)}},T=()=>{w.value=1,m()},L=n=>{w.value=n,m()},E=n=>{_.value=n,w.value=1,m()},R=()=>{m(),v()},D=n=>({GET:"info",POST:"success",PUT:"warning",DELETE:"error",PATCH:"warning"})[n]||"default",I=n=>n>=200&&n<300?"success":n>=400&&n<500?"warning":n>=500?"error":"default",H=n=>n?new Date(n).toLocaleString("zh-CN",{year:"numeric",month:"2-digit",day:"2-digit",hour:"2-digit",minute:"2-digit",second:"2-digit"}):"-",O=n=>{navigator.clipboard.writeText(n).then(()=>{t.success({content:"已复制到剪贴板",duration:2e3})}).catch(()=>{t.error({content:"复制失败",duration:2e3})})},U=n=>{try{u.value=JSON.stringify(JSON.parse(n),null,2)}catch{u.value=n}i.value=!0},M=async()=>{f.warning({title:"清理旧日志",content:"确定要清理30天前的旧日志吗？此操作不可恢复。",positiveText:"确定",negativeText:"取消",onPositiveClick:async()=>{try{y.value=!0,await x.clearApiAccessLogs(30),t.success({content:"旧日志清理成功",duration:3e3}),R()}catch(n){console.error("清理旧日志失败:",n),t.error({content:"清理旧日志失败",duration:3e3})}finally{y.value=!1}}})};return te(async()=>{await Promise.all([m(),v()])}),(n,e)=>{const j=be,ne=ze,W=Se,oe=ke,K=Ce,re=$e,ae=ie,J=qe,Q=je;return k(),$(Z,null,[d(ae,null,{"page-header":c(()=>[e[11]||(e[11]=s("div",null,[s("h1",{class:"text-2xl font-bold text-gray-900 dark:text-white"},"公开API访问日志"),s("p",{class:"text-gray-600 dark:text-gray-400"},"查看公开API的访问记录和统计信息")],-1)),s("div",Re,[d(j,{type:"primary",onClick:R,loading:l(h)},{icon:c(()=>e[7]||(e[7]=[s("i",{class:"fas fa-refresh"},null,-1)])),default:c(()=>[e[8]||(e[8]=z(" 刷新 ",-1))]),_:1,__:[8]},8,["loading"]),d(j,{type:"warning",onClick:M,loading:l(y)},{icon:c(()=>e[9]||(e[9]=[s("i",{class:"fas fa-trash-alt"},null,-1)])),default:c(()=>[e[10]||(e[10]=z(" 清理旧日志 ",-1))]),_:1,__:[10]},8,["loading"])])]),"filter-bar":c(()=>[s("div",Be,[s("div",Ee,[d(ne,{value:l(C),"onUpdate:value":e[0]||(e[0]=o=>A(C)?C.value=o:null),placeholder:"搜索接口路径或IP...",onKeyup:xe(T,["enter"]),clearable:""},{prefix:c(()=>e[12]||(e[12]=[s("i",{class:"fas fa-search"},null,-1)])),_:1},8,["value"]),d(W,{value:l(b),"onUpdate:value":e[1]||(e[1]=o=>A(b)?b.value=o:null),type:"date",placeholder:"开始日期",clearable:""},null,8,["value"]),d(W,{value:l(N),"onUpdate:value":e[2]||(e[2]=o=>A(N)?N.value=o:null),type:"date",placeholder:"结束日期",clearable:""},null,8,["value"]),d(j,{type:"primary",onClick:T,class:"w-20"},{icon:c(()=>e[13]||(e[13]=[s("i",{class:"fas fa-search"},null,-1)])),default:c(()=>[e[14]||(e[14]=z(" 搜索 ",-1))]),_:1,__:[14]})])])]),"content-header":c(()=>[s("div",De,[s("div",Ie,[e[15]||(e[15]=s("span",{class:"text-lg font-semibold"},"访问日志列表",-1)),s("div",He," 共 "+p(l(a))+" 条日志 ",1)]),s("div",Oe,[s("div",Ue,[s("div",Me,p(l(P).total_requests),1),e[16]||(e[16]=s("div",{class:"text-xs text-gray-500"},"总请求",-1))]),s("div",Ve,[s("div",Fe,p(l(P).today_requests),1),e[17]||(e[17]=s("div",{class:"text-xs text-gray-500"},"今日请求",-1))]),s("div",We,[s("div",Ke,p(l(P).week_requests),1),e[18]||(e[18]=s("div",{class:"text-xs text-gray-500"},"本周请求",-1))]),s("div",Je,[s("div",Qe,p(l(P).error_requests),1),e[19]||(e[19]=s("div",{class:"text-xs text-gray-500"},"错误请求",-1))])])])]),content:c(()=>[l(h)?(k(),$("div",Ge,[d(oe,{size:"large"})])):l(S).length===0?(k(),$("div",Xe,e[20]||(e[20]=[s("i",{class:"fas fa-file-alt text-4xl text-gray-400 mb-4"},null,-1),s("p",{class:"text-gray-500 dark:text-gray-400"},"暂无访问日志",-1)]))):(k(),$("div",Ye,[(k(!0),$(Z,null,ve(l(S),o=>(k(),$("div",{key:o.id,class:"border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"},[s("div",Ze,[s("div",et,[s("div",tt,[s("span",st,p(o.id),1),d(K,{type:D(o.method),size:"small"},{default:c(()=>[z(p(o.method),1)]),_:2},1032,["type"]),d(K,{type:I(o.response_status),size:"small"},{default:c(()=>[z(p(o.response_status),1)]),_:2},1032,["type"]),s("code",nt,p(o.endpoint),1),s("div",ot,[s("code",rt,p(o.ip),1),s("span",at,[e[21]||(e[21]=s("i",{class:"fas fa-clock mr-1"},null,-1)),z(" "+p(H(o.created_at)),1)]),o.processing_time>0?(k(),$("span",lt,[e[22]||(e[22]=s("i",{class:"fas fa-tachometer-alt mr-1"},null,-1)),z(" "+p(o.processing_time)+"ms ",1)])):F("",!0)])]),o.request_params?(k(),$("div",it,[s("div",ct,[e[23]||(e[23]=s("strong",{class:"mr-2 flex-0 whitespace-nowrap"},"请求参数:",-1)),s("span",ut,p(o.request_params),1)]),s("div",dt,[d(j,{size:"tiny",onClick:le=>O(o.request_params)},{icon:c(()=>e[24]||(e[24]=[s("i",{class:"fas fa-copy"},null,-1)])),_:2},1032,["onClick"]),d(j,{size:"tiny",onClick:le=>U(o.request_params)},{icon:c(()=>e[25]||(e[25]=[s("i",{class:"fas fa-eye"},null,-1)])),_:2},1032,["onClick"])])])):F("",!0),o.error_message?(k(),$("div",mt,[e[26]||(e[26]=s("strong",null,"错误信息:",-1)),z(" "+p(o.error_message),1)])):F("",!0)])])]))),128))]))]),"content-footer":c(()=>[s("div",ft,[s("div",pt,[d(re,{page:l(w),"onUpdate:page":[e[3]||(e[3]=o=>A(w)?w.value=o:null),L],"page-size":l(_),"onUpdate:pageSize":[e[4]||(e[4]=o=>A(_)?_.value=o:null),E],"item-count":l(a),"page-sizes":[10,20,50,100],"show-size-picker":""},null,8,["page","page-size","item-count"])])])]),_:1}),d(Q,{show:l(i),"onUpdate:show":e[5]||(e[5]=o=>A(i)?i.value=o:null),preset:"card",title:"请求参数详情",style:{"min-width":"600px"}},{default:c(()=>[d(J,{code:l(u),language:"json",folding:!0,"show-line-numbers":!0,class:"bg-gray-100 dark:bg-gray-700 p-4 rounded max-h-96 overflow-auto"},null,8,["code"])]),_:1},8,["show"]),d(Q,{show:l(i),"onUpdate:show":e[6]||(e[6]=o=>A(i)?i.value=o:null),preset:"card",title:"请求参数详情",style:{"min-width":"600px"}},{default:c(()=>[d(J,{code:l(u),language:"json",folding:!0,"show-line-numbers":!0,class:"bg-gray-100 dark:bg-gray-700 p-4 rounded max-h-96 overflow-auto"},null,8,["code"])]),_:1},8,["show"])],64)}}}),Ft=ye(gt,[["__scopeId","data-v-75e55a9c"]]);export{Ft as default};
