import{e as Y,ag as R,ah as re,g as u,ai as ee,aj as te,f as S,h as $,j as W,k as d,V as de,l as J,m as A,ak as Z,n as F,u as D,al as X,ad as Q,v as ne,am as ce,ab as ue,s as fe,i as ge,an as I,r as pe,ao as ie,ap as me,aq as be,ar as ve,F as le,H as he,c as U,b as g,a as i,w as n,o as K,t as M,z as w,a4 as H,d as P,A as xe,N as ye,X as _e,_ as Ce}from"./DmHPR5lg.js";import{c as ze,B as we}from"./DO8alW5h.js";import{t as Se,a as $e,_ as ke}from"./DMTo9cht.js";import{i as Pe,o as Re}from"./DEoDX60x.js";import{r as Oe,a as Te}from"./CoaUF789.js";import{f as je}from"./Hg3Q8c_0.js";import{g as Ee}from"./Bk_rJcZu.js";import{u as Le}from"./D__gE1zc.js";import{_ as Be}from"./HCHZ96f3.js";import{_ as Me}from"./DQDQW_eH.js";import"./QpI9WcJO.js";import"./CbvP3fHg.js";import"./g0YHQayI.js";function oe(t,o="default",a=[]){const{children:f}=t;if(f!==null&&typeof f=="object"&&!Array.isArray(f)){const p=f[o];if(typeof p=="function")return p()}return a}function Ie(t){const{borderRadius:o,avatarColor:a,cardColor:f,fontSize:p,heightTiny:b,heightSmall:v,heightMedium:x,heightLarge:h,heightHuge:l,modalColor:e,popoverColor:c}=t;return{borderRadius:o,fontSize:p,border:`2px solid ${f}`,heightTiny:b,heightSmall:v,heightMedium:x,heightLarge:h,heightHuge:l,color:R(f,a),colorModal:R(e,a),colorPopover:R(c,a)}}const Fe={common:Y,self:Ie},He=re("n-avatar-group"),De=u("avatar",`
 width: var(--n-merged-size);
 height: var(--n-merged-size);
 color: #FFF;
 font-size: var(--n-font-size);
 display: inline-flex;
 position: relative;
 overflow: hidden;
 text-align: center;
 border: var(--n-border);
 border-radius: var(--n-border-radius);
 --n-merged-color: var(--n-color);
 background-color: var(--n-merged-color);
 transition:
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
`,[ee(S("&","--n-merged-color: var(--n-color-modal);")),te(S("&","--n-merged-color: var(--n-color-popover);")),S("img",`
 width: 100%;
 height: 100%;
 `),$("text",`
 white-space: nowrap;
 display: inline-block;
 position: absolute;
 left: 50%;
 top: 50%;
 `),u("icon",`
 vertical-align: bottom;
 font-size: calc(var(--n-merged-size) - 6px);
 `),$("text","line-height: 1.25")]),Ve=Object.assign(Object.assign({},D.props),{size:[String,Number],src:String,circle:{type:Boolean,default:void 0},objectFit:String,round:{type:Boolean,default:void 0},bordered:{type:Boolean,default:void 0},onError:Function,fallbackSrc:String,intersectionObserverOptions:Object,lazy:Boolean,onLoad:Function,renderPlaceholder:Function,renderFallback:Function,imgProps:Object,color:String}),Ae=W({name:"Avatar",props:Ve,slots:Object,setup(t){const{mergedClsPrefixRef:o,inlineThemeDisabled:a}=J(t),f=A(!1);let p=null;const b=A(null),v=A(null),x=()=>{const{value:r}=b;if(r&&(p===null||p!==r.innerHTML)){p=r.innerHTML;const{value:m}=v;if(m){const{offsetWidth:_,offsetHeight:z}=m,{offsetWidth:s,offsetHeight:E}=r,L=.9,B=Math.min(_/s*L,z/E*L,1);r.style.transform=`translateX(-50%) translateY(-50%) scale(${B})`}}},h=Z(He,null),l=F(()=>{const{size:r}=t;if(r)return r;const{size:m}=h||{};return m||"medium"}),e=D("Avatar","-avatar",De,Fe,t,o),c=Z(Se,null),y=F(()=>{if(h)return!0;const{round:r,circle:m}=t;return r!==void 0||m!==void 0?r||m:c?c.roundRef.value:!1}),C=F(()=>h?!0:t.bordered||!1),O=F(()=>{const r=l.value,m=y.value,_=C.value,{color:z}=t,{self:{borderRadius:s,fontSize:E,color:L,border:B,colorModal:V,colorPopover:q},common:{cubicBezierEaseInOut:G}}=e.value;let N;return typeof r=="number"?N=`${r}px`:N=e.value.self[X("height",r)],{"--n-font-size":E,"--n-border":_?B:"none","--n-border-radius":m?"50%":s,"--n-color":z||L,"--n-color-modal":z||V,"--n-color-popover":z||q,"--n-bezier":G,"--n-merged-size":`var(--n-avatar-size-override, ${N})`}}),k=a?Q("avatar",F(()=>{const r=l.value,m=y.value,_=C.value,{color:z}=t;let s="";return r&&(typeof r=="number"?s+=`a${r}`:s+=r[0]),m&&(s+="b"),_&&(s+="c"),z&&(s+=ze(z)),s}),O,t):void 0,T=A(!t.lazy);ne(()=>{if(t.lazy&&t.intersectionObserverOptions){let r;const m=ce(()=>{r==null||r(),r=void 0,t.lazy&&(r=Re(v.value,t.intersectionObserverOptions,T))});ue(()=>{m(),r==null||r()})}}),fe(()=>{var r;return t.src||((r=t.imgProps)===null||r===void 0?void 0:r.src)},()=>{f.value=!1});const j=A(!t.lazy);return{textRef:b,selfRef:v,mergedRoundRef:y,mergedClsPrefix:o,fitTextTransform:x,cssVars:a?void 0:O,themeClass:k==null?void 0:k.themeClass,onRender:k==null?void 0:k.onRender,hasLoadError:f,shouldStartLoading:T,loaded:j,mergedOnError:r=>{if(!T.value)return;f.value=!0;const{onError:m,imgProps:{onError:_}={}}=t;m==null||m(r),_==null||_(r)},mergedOnLoad:r=>{const{onLoad:m,imgProps:{onLoad:_}={}}=t;m==null||m(r),_==null||_(r),j.value=!0}}},render(){var t,o;const{$slots:a,src:f,mergedClsPrefix:p,lazy:b,onRender:v,loaded:x,hasLoadError:h,imgProps:l={}}=this;v==null||v();let e;const c=!x&&!h&&(this.renderPlaceholder?this.renderPlaceholder():(o=(t=this.$slots).placeholder)===null||o===void 0?void 0:o.call(t));return this.hasLoadError?e=this.renderFallback?this.renderFallback():Oe(a.fallback,()=>[d("img",{src:this.fallbackSrc,style:{objectFit:this.objectFit}})]):e=Te(a.default,y=>{if(y)return d(de,{onResize:this.fitTextTransform},{default:()=>d("span",{ref:"textRef",class:`${p}-avatar__text`},y)});if(f||l.src){const C=this.src||l.src;return d("img",Object.assign(Object.assign({},l),{loading:Pe&&!this.intersectionObserverOptions&&b?"lazy":"eager",src:b&&this.intersectionObserverOptions?this.shouldStartLoading?C:void 0:C,"data-image-src":C,onLoad:this.mergedOnLoad,onError:this.mergedOnError,style:[l.style||"",{objectFit:this.objectFit},c?{height:"0",width:"0",visibility:"hidden",position:"absolute"}:""]}))}}),d("span",{ref:"selfRef",class:[`${p}-avatar`,this.themeClass],style:this.cssVars},e,b&&c)}}),We={thPaddingBorderedSmall:"8px 12px",thPaddingBorderedMedium:"12px 16px",thPaddingBorderedLarge:"16px 24px",thPaddingSmall:"0",thPaddingMedium:"0",thPaddingLarge:"0",tdPaddingBorderedSmall:"8px 12px",tdPaddingBorderedMedium:"12px 16px",tdPaddingBorderedLarge:"16px 24px",tdPaddingSmall:"0 0 8px 0",tdPaddingMedium:"0 0 12px 0",tdPaddingLarge:"0 0 16px 0"};function Ne(t){const{tableHeaderColor:o,textColor2:a,textColor1:f,cardColor:p,modalColor:b,popoverColor:v,dividerColor:x,borderRadius:h,fontWeightStrong:l,lineHeight:e,fontSizeSmall:c,fontSizeMedium:y,fontSizeLarge:C}=t;return Object.assign(Object.assign({},We),{lineHeight:e,fontSizeSmall:c,fontSizeMedium:y,fontSizeLarge:C,titleTextColor:f,thColor:R(p,o),thColorModal:R(b,o),thColorPopover:R(v,o),thTextColor:f,thFontWeight:l,tdTextColor:a,tdColor:p,tdColorModal:b,tdColorPopover:v,borderColor:R(p,x),borderColorModal:R(b,x),borderColorPopover:R(v,x),borderRadius:h})}const qe={common:Y,self:Ne},Ge=S([u("descriptions",{fontSize:"var(--n-font-size)"},[u("descriptions-separator",`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),u("descriptions-table-wrapper",[u("descriptions-table",[u("descriptions-table-row",[u("descriptions-table-header",{padding:"var(--n-th-padding)"}),u("descriptions-table-content",{padding:"var(--n-td-padding)"})])])]),ge("bordered",[u("descriptions-table-wrapper",[u("descriptions-table",[u("descriptions-table-row",[S("&:last-child",[u("descriptions-table-content",{paddingBottom:0})])])])])]),I("left-label-placement",[u("descriptions-table-content",[S("> *",{verticalAlign:"top"})])]),I("left-label-align",[S("th",{textAlign:"left"})]),I("center-label-align",[S("th",{textAlign:"center"})]),I("right-label-align",[S("th",{textAlign:"right"})]),I("bordered",[u("descriptions-table-wrapper",`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[u("descriptions-table",[u("descriptions-table-row",[S("&:not(:last-child)",[u("descriptions-table-content",{borderBottom:"1px solid var(--n-merged-border-color)"}),u("descriptions-table-header",{borderBottom:"1px solid var(--n-merged-border-color)"})]),u("descriptions-table-header",`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[S("&:not(:last-child)",{borderRight:"1px solid var(--n-merged-border-color)"})]),u("descriptions-table-content",[S("&:not(:last-child)",{borderRight:"1px solid var(--n-merged-border-color)"})])])])])]),u("descriptions-header",`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),u("descriptions-table-wrapper",`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[u("descriptions-table",`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[u("descriptions-table-row",`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[u("descriptions-table-header",`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),u("descriptions-table-content",`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[$("content",`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),$("label",`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),u("descriptions-table-wrapper",`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),ee(u("descriptions-table-wrapper",`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),te(u("descriptions-table-wrapper",`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),se="DESCRIPTION_ITEM_FLAG";function Ke(t){return typeof t=="object"&&t&&!Array.isArray(t)?t.type&&t.type[se]:!1}const Ue=Object.assign(Object.assign({},D.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:"top"},labelAlign:{type:String,default:"left"},separator:{type:String,default:":"},size:{type:String,default:"medium"},bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),Xe=W({name:"Descriptions",props:Ue,slots:Object,setup(t){const{mergedClsPrefixRef:o,inlineThemeDisabled:a}=J(t),f=D("Descriptions","-descriptions",Ge,qe,t,o),p=F(()=>{const{size:v,bordered:x}=t,{common:{cubicBezierEaseInOut:h},self:{titleTextColor:l,thColor:e,thColorModal:c,thColorPopover:y,thTextColor:C,thFontWeight:O,tdTextColor:k,tdColor:T,tdColorModal:j,tdColorPopover:r,borderColor:m,borderColorModal:_,borderColorPopover:z,borderRadius:s,lineHeight:E,[X("fontSize",v)]:L,[X(x?"thPaddingBordered":"thPadding",v)]:B,[X(x?"tdPaddingBordered":"tdPadding",v)]:V}}=f.value;return{"--n-title-text-color":l,"--n-th-padding":B,"--n-td-padding":V,"--n-font-size":L,"--n-bezier":h,"--n-th-font-weight":O,"--n-line-height":E,"--n-th-text-color":C,"--n-td-text-color":k,"--n-th-color":e,"--n-th-color-modal":c,"--n-th-color-popover":y,"--n-td-color":T,"--n-td-color-modal":j,"--n-td-color-popover":r,"--n-border-radius":s,"--n-border-color":m,"--n-border-color-modal":_,"--n-border-color-popover":z}}),b=a?Q("descriptions",F(()=>{let v="";const{size:x,bordered:h}=t;return h&&(v+="a"),v+=x[0],v}),p,t):void 0;return{mergedClsPrefix:o,cssVars:a?void 0:p,themeClass:b==null?void 0:b.themeClass,onRender:b==null?void 0:b.onRender,compitableColumn:Le(t,["columns","column"]),inlineThemeDisabled:a}},render(){const t=this.$slots.default,o=t?je(t()):[];o.length;const{contentClass:a,labelClass:f,compitableColumn:p,labelPlacement:b,labelAlign:v,size:x,bordered:h,title:l,cssVars:e,mergedClsPrefix:c,separator:y,onRender:C}=this;C==null||C();const O=o.filter(r=>Ke(r)),k={span:0,row:[],secondRow:[],rows:[]},j=O.reduce((r,m,_)=>{const z=m.props||{},s=O.length-1===_,E=["label"in z?z.label:oe(m,"label")],L=[oe(m)],B=z.span||1,V=r.span;r.span+=B;const q=z.labelStyle||z["label-style"]||this.labelStyle,G=z.contentStyle||z["content-style"]||this.contentStyle;if(b==="left")h?r.row.push(d("th",{class:[`${c}-descriptions-table-header`,f],colspan:1,style:q},E),d("td",{class:[`${c}-descriptions-table-content`,a],colspan:s?(p-V)*2+1:B*2-1,style:G},L)):r.row.push(d("td",{class:`${c}-descriptions-table-content`,colspan:s?(p-V)*2:B*2},d("span",{class:[`${c}-descriptions-table-content__label`,f],style:q},[...E,y&&d("span",{class:`${c}-descriptions-separator`},y)]),d("span",{class:[`${c}-descriptions-table-content__content`,a],style:G},L)));else{const N=s?(p-V)*2:B*2;r.row.push(d("th",{class:[`${c}-descriptions-table-header`,f],colspan:N,style:q},E)),r.secondRow.push(d("td",{class:[`${c}-descriptions-table-content`,a],colspan:N,style:G},L))}return(r.span>=p||s)&&(r.span=0,r.row.length&&(r.rows.push(r.row),r.row=[]),b!=="left"&&r.secondRow.length&&(r.rows.push(r.secondRow),r.secondRow=[])),r},k).rows.map(r=>d("tr",{class:`${c}-descriptions-table-row`},r));return d("div",{style:e,class:[`${c}-descriptions`,this.themeClass,`${c}-descriptions--${b}-label-placement`,`${c}-descriptions--${v}-label-align`,`${c}-descriptions--${x}-size`,h&&`${c}-descriptions--bordered`]},l||this.$slots.header?d("div",{class:`${c}-descriptions-header`},l||Ee(this,"header")):null,d("div",{class:`${c}-descriptions-table-wrapper`},d("table",{class:`${c}-descriptions-table`},d("tbody",null,b==="top"&&d("tr",{class:`${c}-descriptions-table-row`,style:{visibility:"collapse"}},pe(p*2,d("td",null))),j))))}}),Ye={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},Je=W({name:"DescriptionsItem",[se]:!0,props:Ye,slots:Object,render(){return null}});function Qe(t){const{textColor2:o,cardColor:a,modalColor:f,popoverColor:p,dividerColor:b,borderRadius:v,fontSize:x,hoverColor:h}=t;return{textColor:o,color:a,colorHover:h,colorModal:f,colorHoverModal:R(f,h),colorPopover:p,colorHoverPopover:R(p,h),borderColor:b,borderColorModal:R(f,b),borderColorPopover:R(p,b),borderRadius:v,fontSize:x}}const Ze={common:Y,self:Qe};function et(t){const{textColor1:o,textColor2:a,fontWeightStrong:f,fontSize:p}=t;return{fontSize:p,titleTextColor:o,textColor:a,titleFontWeight:f}}const tt={common:Y,self:et},ot=S([u("list",`
 --n-merged-border-color: var(--n-border-color);
 --n-merged-color: var(--n-color);
 --n-merged-color-hover: var(--n-color-hover);
 margin: 0;
 font-size: var(--n-font-size);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 padding: 0;
 list-style-type: none;
 color: var(--n-text-color);
 background-color: var(--n-merged-color);
 `,[I("show-divider",[u("list-item",[S("&:not(:last-child)",[$("divider",`
 background-color: var(--n-merged-border-color);
 `)])])]),I("clickable",[u("list-item",`
 cursor: pointer;
 `)]),I("bordered",`
 border: 1px solid var(--n-merged-border-color);
 border-radius: var(--n-border-radius);
 `),I("hoverable",[u("list-item",`
 border-radius: var(--n-border-radius);
 `,[S("&:hover",`
 background-color: var(--n-merged-color-hover);
 `,[$("divider",`
 background-color: transparent;
 `)])])]),I("bordered, hoverable",[u("list-item",`
 padding: 12px 20px;
 `),$("header, footer",`
 padding: 12px 20px;
 `)]),$("header, footer",`
 padding: 12px 0;
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[S("&:not(:last-child)",`
 border-bottom: 1px solid var(--n-merged-border-color);
 `)]),u("list-item",`
 position: relative;
 padding: 12px 0; 
 box-sizing: border-box;
 display: flex;
 flex-wrap: nowrap;
 align-items: center;
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[$("prefix",`
 margin-right: 20px;
 flex: 0;
 `),$("suffix",`
 margin-left: 20px;
 flex: 0;
 `),$("main",`
 flex: 1;
 `),$("divider",`
 height: 1px;
 position: absolute;
 bottom: 0;
 left: 0;
 right: 0;
 background-color: transparent;
 transition: background-color .3s var(--n-bezier);
 pointer-events: none;
 `)])]),ee(u("list",`
 --n-merged-color-hover: var(--n-color-hover-modal);
 --n-merged-color: var(--n-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),te(u("list",`
 --n-merged-color-hover: var(--n-color-hover-popover);
 --n-merged-color: var(--n-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),rt=Object.assign(Object.assign({},D.props),{size:{type:String,default:"medium"},bordered:Boolean,clickable:Boolean,hoverable:Boolean,showDivider:{type:Boolean,default:!0}}),ae=re("n-list"),nt=W({name:"List",props:rt,slots:Object,setup(t){const{mergedClsPrefixRef:o,inlineThemeDisabled:a,mergedRtlRef:f}=J(t),p=ie("List",f,o),b=D("List","-list",ot,Ze,t,o);me(ae,{showDividerRef:be(t,"showDivider"),mergedClsPrefixRef:o});const v=F(()=>{const{common:{cubicBezierEaseInOut:h},self:{fontSize:l,textColor:e,color:c,colorModal:y,colorPopover:C,borderColor:O,borderColorModal:k,borderColorPopover:T,borderRadius:j,colorHover:r,colorHoverModal:m,colorHoverPopover:_}}=b.value;return{"--n-font-size":l,"--n-bezier":h,"--n-text-color":e,"--n-color":c,"--n-border-radius":j,"--n-border-color":O,"--n-border-color-modal":k,"--n-border-color-popover":T,"--n-color-modal":y,"--n-color-popover":C,"--n-color-hover":r,"--n-color-hover-modal":m,"--n-color-hover-popover":_}}),x=a?Q("list",void 0,v,t):void 0;return{mergedClsPrefix:o,rtlEnabled:p,cssVars:a?void 0:v,themeClass:x==null?void 0:x.themeClass,onRender:x==null?void 0:x.onRender}},render(){var t;const{$slots:o,mergedClsPrefix:a,onRender:f}=this;return f==null||f(),d("ul",{class:[`${a}-list`,this.rtlEnabled&&`${a}-list--rtl`,this.bordered&&`${a}-list--bordered`,this.showDivider&&`${a}-list--show-divider`,this.hoverable&&`${a}-list--hoverable`,this.clickable&&`${a}-list--clickable`,this.themeClass],style:this.cssVars},o.header?d("div",{class:`${a}-list__header`},o.header()):null,(t=o.default)===null||t===void 0?void 0:t.call(o),o.footer?d("div",{class:`${a}-list__footer`},o.footer()):null)}}),it=W({name:"ListItem",slots:Object,setup(){const t=Z(ae,null);return t||ve("list-item","`n-list-item` must be placed in `n-list`."),{showDivider:t.showDividerRef,mergedClsPrefix:t.mergedClsPrefixRef}},render(){const{$slots:t,mergedClsPrefix:o}=this;return d("li",{class:`${o}-list-item`},t.prefix?d("div",{class:`${o}-list-item__prefix`},t.prefix()):null,t.default?d("div",{class:`${o}-list-item__main`},t):null,t.suffix?d("div",{class:`${o}-list-item__suffix`},t.suffix()):null,this.showDivider&&d("div",{class:`${o}-list-item__divider`}))}}),lt=u("thing",`
 display: flex;
 transition: color .3s var(--n-bezier);
 font-size: var(--n-font-size);
 color: var(--n-text-color);
`,[u("thing-avatar",`
 margin-right: 12px;
 margin-top: 2px;
 `),u("thing-avatar-header-wrapper",`
 display: flex;
 flex-wrap: nowrap;
 `,[u("thing-header-wrapper",`
 flex: 1;
 `)]),u("thing-main",`
 flex-grow: 1;
 `,[u("thing-header",`
 display: flex;
 margin-bottom: 4px;
 justify-content: space-between;
 align-items: center;
 `,[$("title",`
 font-size: 16px;
 font-weight: var(--n-title-font-weight);
 transition: color .3s var(--n-bezier);
 color: var(--n-title-text-color);
 `)]),$("description",[S("&:not(:last-child)",`
 margin-bottom: 4px;
 `)]),$("content",[S("&:not(:first-child)",`
 margin-top: 12px;
 `)]),$("footer",[S("&:not(:first-child)",`
 margin-top: 12px;
 `)]),$("action",[S("&:not(:first-child)",`
 margin-top: 12px;
 `)])])]),st=Object.assign(Object.assign({},D.props),{title:String,titleExtra:String,description:String,descriptionClass:String,descriptionStyle:[String,Object],content:String,contentClass:String,contentStyle:[String,Object],contentIndented:Boolean}),at=W({name:"Thing",props:st,slots:Object,setup(t,{slots:o}){const{mergedClsPrefixRef:a,inlineThemeDisabled:f,mergedRtlRef:p}=J(t),b=D("Thing","-thing",lt,tt,t,a),v=ie("Thing",p,a),x=F(()=>{const{self:{titleTextColor:l,textColor:e,titleFontWeight:c,fontSize:y},common:{cubicBezierEaseInOut:C}}=b.value;return{"--n-bezier":C,"--n-font-size":y,"--n-text-color":e,"--n-title-font-weight":c,"--n-title-text-color":l}}),h=f?Q("thing",void 0,x,t):void 0;return()=>{var l;const{value:e}=a,c=v?v.value:!1;return(l=h==null?void 0:h.onRender)===null||l===void 0||l.call(h),d("div",{class:[`${e}-thing`,h==null?void 0:h.themeClass,c&&`${e}-thing--rtl`],style:f?void 0:x.value},o.avatar&&t.contentIndented?d("div",{class:`${e}-thing-avatar`},o.avatar()):null,d("div",{class:`${e}-thing-main`},!t.contentIndented&&(o.header||t.title||o["header-extra"]||t.titleExtra||o.avatar)?d("div",{class:`${e}-thing-avatar-header-wrapper`},o.avatar?d("div",{class:`${e}-thing-avatar`},o.avatar()):null,o.header||t.title||o["header-extra"]||t.titleExtra?d("div",{class:`${e}-thing-header-wrapper`},d("div",{class:`${e}-thing-header`},o.header||t.title?d("div",{class:`${e}-thing-header__title`},o.header?o.header():t.title):null,o["header-extra"]||t.titleExtra?d("div",{class:`${e}-thing-header__extra`},o["header-extra"]?o["header-extra"]():t.titleExtra):null),o.description||t.description?d("div",{class:[`${e}-thing-main__description`,t.descriptionClass],style:t.descriptionStyle},o.description?o.description():t.description):null):null):d(le,null,o.header||t.title||o["header-extra"]||t.titleExtra?d("div",{class:`${e}-thing-header`},o.header||t.title?d("div",{class:`${e}-thing-header__title`},o.header?o.header():t.title):null,o["header-extra"]||t.titleExtra?d("div",{class:`${e}-thing-header__extra`},o["header-extra"]?o["header-extra"]():t.titleExtra):null):null,o.description||t.description?d("div",{class:[`${e}-thing-main__description`,t.descriptionClass],style:t.descriptionStyle},o.description?o.description():t.description):null),o.default||t.content?d("div",{class:[`${e}-thing-main__content`,t.contentClass],style:t.contentStyle},o.default?o.default():t.content):null,o.footer?d("div",{class:`${e}-thing-main__footer`},o.footer()):null,o.action?d("div",{class:`${e}-thing-main__action`},o.action()):null))}}}),dt={class:"space-y-8"},ct={class:"flex items-center justify-between"},ut={class:"text-3xl font-bold mb-2 text-white"},ft={class:"grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"},gt={class:"flex items-center"},pt={class:"ml-4"},mt={class:"text-2xl font-bold text-gray-900 dark:text-white"},bt={class:"flex items-center"},vt={class:"ml-4"},ht={class:"text-2xl font-bold text-gray-900 dark:text-white"},xt={class:"flex items-center"},yt={class:"ml-4"},_t={class:"text-2xl font-bold text-gray-900 dark:text-white"},Ct={class:"flex items-center"},zt={class:"ml-4"},wt={class:"text-2xl font-bold text-gray-900 dark:text-white"},St={class:"grid grid-cols-1 lg:grid-cols-2 gap-8"},$t={key:0,class:"text-center py-8"},kt={key:1,class:"space-y-4"},Pt={class:"grid grid-cols-2 gap-4"},Rt=W({__name:"index",setup(t){const o=he(),a=A({resources:0,favorites:0,history:0,recent:0}),f=A([]),p=l=>l?new Date(l).toLocaleDateString("zh-CN"):"未知",b=l=>{l.url&&window.open(l.url,"_blank")},v=()=>{_e().info({content:"导出功能开发中...",duration:3e3})},x=async()=>{try{a.value={resources:12,favorites:8,history:25,recent:3}}catch(l){console.error("获取用户统计数据失败:",l)}},h=async()=>{try{f.value=[{id:1,title:"示例资源1",description:"这是一个示例资源描述",url:"https://example.com"},{id:2,title:"示例资源2",description:"这是另一个示例资源描述",url:"https://example.com"}]}catch(l){console.error("获取最近资源失败:",l)}};return ne(async()=>{await x(),await h()}),(l,e)=>{const c=Me,y=we,C=ke,O=$e,k=Ae,T=at,j=it,r=nt,m=Be,_=Je,z=Xe;return K(),U("div",dt,[g(c,{class:"bg-gradient-to-r from-blue-500 to-purple-600 text-white border-0"},{default:n(()=>{var s;return[i("div",ct,[i("div",null,[i("h1",ut," 欢迎回来，"+M(((s=w(o).user)==null?void 0:s.username)||"用户")+"！ ",1),e[8]||(e[8]=i("p",{class:"text-blue-100 text-lg"}," 这里是您的个人中心，您可以管理您的资源、收藏和历史记录。 ",-1))]),e[9]||(e[9]=i("div",{class:"hidden md:block"},[i("div",{class:"w-16 h-16 bg-white/20 rounded-full flex items-center justify-center"},[i("i",{class:"fas fa-user text-2xl text-white"})])],-1))])]}),_:1}),i("div",ft,[g(c,null,{footer:n(()=>[g(y,{text:"",type:"primary",onClick:e[0]||(e[0]=s=>("navigateTo"in l?l.navigateTo:w(H))("/user/resources"))},{icon:n(()=>e[12]||(e[12]=[i("i",{class:"fas fa-arrow-right"},null,-1)])),default:n(()=>[e[13]||(e[13]=P(" 查看详情 ",-1))]),_:1,__:[13]})]),default:n(()=>[i("div",gt,[e[11]||(e[11]=i("div",{class:"p-3 bg-blue-100 dark:bg-blue-900 rounded-lg"},[i("i",{class:"fas fa-cloud text-blue-600 dark:text-blue-400 text-xl"})],-1)),i("div",pt,[e[10]||(e[10]=i("p",{class:"text-sm font-medium text-gray-600 dark:text-gray-400"},"我的资源",-1)),i("p",mt,M(w(a).resources||0),1)])])]),_:1}),g(c,null,{footer:n(()=>[g(y,{text:"",type:"error",onClick:e[1]||(e[1]=s=>("navigateTo"in l?l.navigateTo:w(H))("/user/favorites"))},{icon:n(()=>e[16]||(e[16]=[i("i",{class:"fas fa-arrow-right"},null,-1)])),default:n(()=>[e[17]||(e[17]=P(" 查看详情 ",-1))]),_:1,__:[17]})]),default:n(()=>[i("div",bt,[e[15]||(e[15]=i("div",{class:"p-3 bg-red-100 dark:bg-red-900 rounded-lg"},[i("i",{class:"fas fa-heart text-red-600 dark:text-red-400 text-xl"})],-1)),i("div",vt,[e[14]||(e[14]=i("p",{class:"text-sm font-medium text-gray-600 dark:text-gray-400"},"收藏夹",-1)),i("p",ht,M(w(a).favorites||0),1)])])]),_:1}),g(c,null,{footer:n(()=>[g(y,{text:"",type:"success",onClick:e[2]||(e[2]=s=>("navigateTo"in l?l.navigateTo:w(H))("/user/history"))},{icon:n(()=>e[20]||(e[20]=[i("i",{class:"fas fa-arrow-right"},null,-1)])),default:n(()=>[e[21]||(e[21]=P(" 查看详情 ",-1))]),_:1,__:[21]})]),default:n(()=>[i("div",xt,[e[19]||(e[19]=i("div",{class:"p-3 bg-green-100 dark:bg-green-900 rounded-lg"},[i("i",{class:"fas fa-history text-green-600 dark:text-green-400 text-xl"})],-1)),i("div",yt,[e[18]||(e[18]=i("p",{class:"text-sm font-medium text-gray-600 dark:text-gray-400"},"浏览历史",-1)),i("p",_t,M(w(a).history||0),1)])])]),_:1}),g(c,null,{footer:n(()=>[g(y,{text:"",type:"info",onClick:e[3]||(e[3]=s=>("navigateTo"in l?l.navigateTo:w(H))("/user/activity"))},{icon:n(()=>e[24]||(e[24]=[i("i",{class:"fas fa-arrow-right"},null,-1)])),default:n(()=>[e[25]||(e[25]=P(" 查看详情 ",-1))]),_:1,__:[25]})]),default:n(()=>[i("div",Ct,[e[23]||(e[23]=i("div",{class:"p-3 bg-purple-100 dark:bg-purple-900 rounded-lg"},[i("i",{class:"fas fa-clock text-purple-600 dark:text-purple-400 text-xl"})],-1)),i("div",zt,[e[22]||(e[22]=i("p",{class:"text-sm font-medium text-gray-600 dark:text-gray-400"},"最近活动",-1)),i("p",wt,M(w(a).recent||0),1)])])]),_:1})]),i("div",St,[g(c,{title:"最近资源",bordered:!1},{"header-extra":n(()=>[g(C,{type:"info",size:"small"},{default:n(()=>e[26]||(e[26]=[P("最新",-1)])),_:1,__:[26]})]),default:n(()=>[w(f).length===0?(K(),U("div",$t,[g(O,{description:"暂无最近资源"},{icon:n(()=>e[27]||(e[27]=[i("i",{class:"fas fa-cloud text-gray-400 text-3xl"},null,-1)])),extra:n(()=>[g(y,{type:"primary",onClick:e[4]||(e[4]=s=>("navigateTo"in l?l.navigateTo:w(H))("/"))},{default:n(()=>e[28]||(e[28]=[P(" 去发现资源 ",-1)])),_:1,__:[28]})]),_:1})])):(K(),U("div",kt,[g(r,null,{default:n(()=>[(K(!0),U(le,null,xe(w(f),s=>(K(),ye(j,{key:s.id},{prefix:n(()=>[g(k,{round:"",size:"small"},{default:n(()=>e[29]||(e[29]=[i("i",{class:"fas fa-file-alt"},null,-1)])),_:1,__:[29]})]),suffix:n(()=>[g(y,{text:"",type:"primary",onClick:E=>b(s)},{icon:n(()=>e[31]||(e[31]=[i("i",{class:"fas fa-external-link-alt"},null,-1)])),_:2},1032,["onClick"])]),default:n(()=>[g(T,{title:s.title,description:s.description},{avatar:n(()=>[g(k,{round:"",size:"small"},{default:n(()=>e[30]||(e[30]=[i("i",{class:"fas fa-file-alt"},null,-1)])),_:1,__:[30]})]),_:2},1032,["title","description"])]),_:2},1024))),128))]),_:1})]))]),_:1}),g(c,{title:"快速操作",bordered:!1},{"header-extra":n(()=>[g(C,{type:"success",size:"small"},{default:n(()=>e[32]||(e[32]=[P("快捷",-1)])),_:1,__:[32]})]),default:n(()=>[i("div",Pt,[g(y,{quaternary:"",type:"primary",size:"large",class:"h-24",onClick:e[5]||(e[5]=s=>("navigateTo"in l?l.navigateTo:w(H))("/user/profile"))},{icon:n(()=>e[33]||(e[33]=[i("i",{class:"fas fa-user-edit text-2xl"},null,-1)])),default:n(()=>[e[34]||(e[34]=i("div",{class:"flex flex-col items-center"},[i("span",{class:"text-sm font-medium"},"个人资料")],-1))]),_:1,__:[34]}),g(y,{quaternary:"",type:"success",size:"large",class:"h-24",onClick:e[6]||(e[6]=s=>("navigateTo"in l?l.navigateTo:w(H))("/user/settings"))},{icon:n(()=>e[35]||(e[35]=[i("i",{class:"fas fa-cog text-2xl"},null,-1)])),default:n(()=>[e[36]||(e[36]=i("div",{class:"flex flex-col items-center"},[i("span",{class:"text-sm font-medium"},"设置")],-1))]),_:1,__:[36]}),g(y,{quaternary:"",type:"info",size:"large",class:"h-24",onClick:e[7]||(e[7]=s=>("navigateTo"in l?l.navigateTo:w(H))("/"))},{icon:n(()=>e[37]||(e[37]=[i("i",{class:"fas fa-search text-2xl"},null,-1)])),default:n(()=>[e[38]||(e[38]=i("div",{class:"flex flex-col items-center"},[i("span",{class:"text-sm font-medium"},"搜索资源")],-1))]),_:1,__:[38]}),g(y,{quaternary:"",type:"warning",size:"large",class:"h-24",onClick:v},{icon:n(()=>e[39]||(e[39]=[i("i",{class:"fas fa-download text-2xl"},null,-1)])),default:n(()=>[e[40]||(e[40]=i("div",{class:"flex flex-col items-center"},[i("span",{class:"text-sm font-medium"},"导出数据")],-1))]),_:1,__:[40]})])]),_:1})]),g(c,{title:"账户信息",bordered:!1},{"header-extra":n(()=>[g(C,{type:"primary",size:"small"},{default:n(()=>e[41]||(e[41]=[P("账户",-1)])),_:1,__:[41]})]),default:n(()=>[g(z,{column:2,bordered:""},{default:n(()=>[g(_,{label:"用户名"},{default:n(()=>[g(m,null,{default:n(()=>{var s;return[P(M((s=w(o).user)==null?void 0:s.username),1)]}),_:1})]),_:1}),g(_,{label:"邮箱"},{default:n(()=>[g(m,null,{default:n(()=>{var s;return[P(M(((s=w(o).user)==null?void 0:s.email)||"未设置"),1)]}),_:1})]),_:1}),g(_,{label:"注册时间"},{default:n(()=>[g(m,null,{default:n(()=>{var s;return[P(M(p(((s=w(o).user)==null?void 0:s.created_at)||"")),1)]}),_:1})]),_:1}),g(_,{label:"最后登录"},{default:n(()=>[g(m,null,{default:n(()=>{var s;return[P(M(p(((s=w(o).user)==null?void 0:s.last_login_at)||"")),1)]}),_:1})]),_:1}),g(_,{label:"账户状态"},{default:n(()=>[g(C,{type:"success",size:"small"},{icon:n(()=>e[42]||(e[42]=[i("i",{class:"fas fa-check-circle"},null,-1)])),default:n(()=>[e[43]||(e[43]=P(" 正常 ",-1))]),_:1,__:[43]})]),_:1}),g(_,{label:"用户角色"},{default:n(()=>[g(C,{type:"primary",size:"small"},{default:n(()=>{var s;return[P(M(((s=w(o).user)==null?void 0:s.role)==="admin"?"管理员":"普通用户"),1)]}),_:1})]),_:1}),g(_,{label:"账户类型"},{default:n(()=>[g(m,null,{default:n(()=>e[44]||(e[44]=[P("免费账户",-1)])),_:1,__:[44]})]),_:1})]),_:1})]),_:1})])}}}),Wt=Ce(Rt,[["__scopeId","data-v-0f4fff64"]]);export{Wt as default};
