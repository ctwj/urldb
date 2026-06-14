import{e as T,j as z,k as r,ax as N,bj as j,bh as q,bg as O,bi as G,n as x,f as P,g as s,an as $,l as _,u as L,ad as A,al as W}from"./DmHPR5lg.js";import{f as w}from"./B-p6aW7q.js";function M(e){const{infoColor:d,successColor:u,warningColor:f,errorColor:n,textColor2:o,progressRailColor:t,fontSize:l,fontWeight:a}=e;return{fontSize:l,fontSizeCircle:"28px",fontWeightCircle:a,railColor:t,railHeight:"8px",iconSizeCircle:"36px",iconSizeLine:"18px",iconColor:d,iconColorInfo:d,iconColorSuccess:u,iconColorWarning:f,iconColorError:n,textColorCircle:o,textColorLineInner:"rgb(255, 255, 255)",textColorLineOuter:o,fillColor:d,fillColorInfo:d,fillColorSuccess:u,fillColorWarning:f,fillColorError:n,lineBgProcessing:"linear-gradient(90deg, rgba(255, 255, 255, .3) 0%, rgba(255, 255, 255, .5) 100%)"}}const X={name:"Progress",common:T,self:M},Y={success:r(G,null),error:r(O,null),warning:r(q,null),info:r(j,null)},H=z({name:"ProgressCircle",props:{clsPrefix:{type:String,required:!0},status:{type:String,required:!0},strokeWidth:{type:Number,required:!0},fillColor:[String,Object],railColor:String,railStyle:[String,Object],percentage:{type:Number,default:0},offsetDegree:{type:Number,default:0},showIndicator:{type:Boolean,required:!0},indicatorTextColor:String,unit:String,viewBoxWidth:{type:Number,required:!0},gapDegree:{type:Number,required:!0},gapOffsetDegree:{type:Number,default:0}},setup(e,{slots:d}){function u(n,o,t,l){const{gapDegree:a,viewBoxWidth:h,strokeWidth:m}=e,g=50,p=0,y=g,i=0,c=2*g,C=50+m/2,S=`M ${C},${C} m ${p},${y}
      a ${g},${g} 0 1 1 ${i},${-c}
      a ${g},${g} 0 1 1 ${-i},${c}`,b=Math.PI*2*g,v={stroke:l==="rail"?t:typeof e.fillColor=="object"?"url(#gradient)":t,strokeDasharray:`${n/100*(b-a)}px ${h*8}px`,strokeDashoffset:`-${a/2}px`,transformOrigin:o?"center":void 0,transform:o?`rotate(${o}deg)`:void 0};return{pathString:S,pathStyle:v}}const f=()=>{const n=typeof e.fillColor=="object",o=n?e.fillColor.stops[0]:"",t=n?e.fillColor.stops[1]:"";return n&&r("defs",null,r("linearGradient",{id:"gradient",x1:"0%",y1:"100%",x2:"100%",y2:"0%"},r("stop",{offset:"0%","stop-color":o}),r("stop",{offset:"100%","stop-color":t})))};return()=>{const{fillColor:n,railColor:o,strokeWidth:t,offsetDegree:l,status:a,percentage:h,showIndicator:m,indicatorTextColor:g,unit:p,gapOffsetDegree:y,clsPrefix:i}=e,{pathString:c,pathStyle:C}=u(100,0,o,"rail"),{pathString:S,pathStyle:b}=u(h,l,n,"fill"),v=100+t;return r("div",{class:`${i}-progress-content`,role:"none"},r("div",{class:`${i}-progress-graph`,"aria-hidden":!0},r("div",{class:`${i}-progress-graph-circle`,style:{transform:y?`rotate(${y}deg)`:void 0}},r("svg",{viewBox:`0 0 ${v} ${v}`},f(),r("g",null,r("path",{class:`${i}-progress-graph-circle-rail`,d:c,"stroke-width":t,"stroke-linecap":"round",fill:"none",style:C})),r("g",null,r("path",{class:[`${i}-progress-graph-circle-fill`,h===0&&`${i}-progress-graph-circle-fill--empty`],d:S,"stroke-width":t,"stroke-linecap":"round",fill:"none",style:b}))))),m?r("div",null,d.default?r("div",{class:`${i}-progress-custom-content`,role:"none"},d.default()):a!=="default"?r("div",{class:`${i}-progress-icon`,"aria-hidden":!0},r(N,{clsPrefix:i},{default:()=>Y[a]})):r("div",{class:`${i}-progress-text`,style:{color:g},role:"none"},r("span",{class:`${i}-progress-text__percentage`},h),r("span",{class:`${i}-progress-text__unit`},p))):null)}}}),E={success:r(G,null),error:r(O,null),warning:r(q,null),info:r(j,null)},V=z({name:"ProgressLine",props:{clsPrefix:{type:String,required:!0},percentage:{type:Number,default:0},railColor:String,railStyle:[String,Object],fillColor:[String,Object],status:{type:String,required:!0},indicatorPlacement:{type:String,required:!0},indicatorTextColor:String,unit:{type:String,default:"%"},processing:{type:Boolean,required:!0},showIndicator:{type:Boolean,required:!0},height:[String,Number],railBorderRadius:[String,Number],fillBorderRadius:[String,Number]},setup(e,{slots:d}){const u=x(()=>w(e.height)),f=x(()=>{var t,l;return typeof e.fillColor=="object"?`linear-gradient(to right, ${(t=e.fillColor)===null||t===void 0?void 0:t.stops[0]} , ${(l=e.fillColor)===null||l===void 0?void 0:l.stops[1]})`:e.fillColor}),n=x(()=>e.railBorderRadius!==void 0?w(e.railBorderRadius):e.height!==void 0?w(e.height,{c:.5}):""),o=x(()=>e.fillBorderRadius!==void 0?w(e.fillBorderRadius):e.railBorderRadius!==void 0?w(e.railBorderRadius):e.height!==void 0?w(e.height,{c:.5}):"");return()=>{const{indicatorPlacement:t,railColor:l,railStyle:a,percentage:h,unit:m,indicatorTextColor:g,status:p,showIndicator:y,processing:i,clsPrefix:c}=e;return r("div",{class:`${c}-progress-content`,role:"none"},r("div",{class:`${c}-progress-graph`,"aria-hidden":!0},r("div",{class:[`${c}-progress-graph-line`,{[`${c}-progress-graph-line--indicator-${t}`]:!0}]},r("div",{class:`${c}-progress-graph-line-rail`,style:[{backgroundColor:l,height:u.value,borderRadius:n.value},a]},r("div",{class:[`${c}-progress-graph-line-fill`,i&&`${c}-progress-graph-line-fill--processing`],style:{maxWidth:`${e.percentage}%`,background:f.value,height:u.value,lineHeight:u.value,borderRadius:o.value}},t==="inside"?r("div",{class:`${c}-progress-graph-line-indicator`,style:{color:g}},d.default?d.default():`${h}${m}`):null)))),y&&t==="outside"?r("div",null,d.default?r("div",{class:`${c}-progress-custom-content`,style:{color:g},role:"none"},d.default()):p==="default"?r("div",{role:"none",class:`${c}-progress-icon ${c}-progress-icon--as-text`,style:{color:g}},h,m):r("div",{class:`${c}-progress-icon`,"aria-hidden":!0},r(N,{clsPrefix:c},{default:()=>E[p]}))):null)}}});function D(e,d,u=100){return`m ${u/2} ${u/2-e} a ${e} ${e} 0 1 1 0 ${2*e} a ${e} ${e} 0 1 1 0 -${2*e}`}const F=z({name:"ProgressMultipleCircle",props:{clsPrefix:{type:String,required:!0},viewBoxWidth:{type:Number,required:!0},percentage:{type:Array,default:[0]},strokeWidth:{type:Number,required:!0},circleGap:{type:Number,required:!0},showIndicator:{type:Boolean,required:!0},fillColor:{type:Array,default:()=>[]},railColor:{type:Array,default:()=>[]},railStyle:{type:Array,default:()=>[]}},setup(e,{slots:d}){const u=x(()=>e.percentage.map((o,t)=>`${Math.PI*o/100*(e.viewBoxWidth/2-e.strokeWidth/2*(1+2*t)-e.circleGap*t)*2}, ${e.viewBoxWidth*8}`)),f=(n,o)=>{const t=e.fillColor[o],l=typeof t=="object"?t.stops[0]:"",a=typeof t=="object"?t.stops[1]:"";return typeof e.fillColor[o]=="object"&&r("linearGradient",{id:`gradient-${o}`,x1:"100%",y1:"0%",x2:"0%",y2:"100%"},r("stop",{offset:"0%","stop-color":l}),r("stop",{offset:"100%","stop-color":a}))};return()=>{const{viewBoxWidth:n,strokeWidth:o,circleGap:t,showIndicator:l,fillColor:a,railColor:h,railStyle:m,percentage:g,clsPrefix:p}=e;return r("div",{class:`${p}-progress-content`,role:"none"},r("div",{class:`${p}-progress-graph`,"aria-hidden":!0},r("div",{class:`${p}-progress-graph-circle`},r("svg",{viewBox:`0 0 ${n} ${n}`},r("defs",null,g.map((y,i)=>f(y,i))),g.map((y,i)=>r("g",{key:i},r("path",{class:`${p}-progress-graph-circle-rail`,d:D(n/2-o/2*(1+2*i)-t*i,o,n),"stroke-width":o,"stroke-linecap":"round",fill:"none",style:[{strokeDashoffset:0,stroke:h[i]},m[i]]}),r("path",{class:[`${p}-progress-graph-circle-fill`,y===0&&`${p}-progress-graph-circle-fill--empty`],d:D(n/2-o/2*(1+2*i)-t*i,o,n),"stroke-width":o,"stroke-linecap":"round",fill:"none",style:{strokeDasharray:u.value[i],strokeDashoffset:0,stroke:typeof a[i]=="object"?`url(#gradient-${i})`:a[i]}})))))),l&&d.default?r("div",null,r("div",{class:`${p}-progress-text`},d.default())):null)}}}),K=P([s("progress",{display:"inline-block"},[s("progress-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 `),$("line",`
 width: 100%;
 display: block;
 `,[s("progress-content",`
 display: flex;
 align-items: center;
 `,[s("progress-graph",{flex:1})]),s("progress-custom-content",{marginLeft:"14px"}),s("progress-icon",`
 width: 30px;
 padding-left: 14px;
 height: var(--n-icon-size-line);
 line-height: var(--n-icon-size-line);
 font-size: var(--n-icon-size-line);
 `,[$("as-text",`
 color: var(--n-text-color-line-outer);
 text-align: center;
 width: 40px;
 font-size: var(--n-font-size);
 padding-left: 4px;
 transition: color .3s var(--n-bezier);
 `)])]),$("circle, dashboard",{width:"120px"},[s("progress-custom-content",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 justify-content: center;
 `),s("progress-text",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 color: inherit;
 font-size: var(--n-font-size-circle);
 color: var(--n-text-color-circle);
 font-weight: var(--n-font-weight-circle);
 transition: color .3s var(--n-bezier);
 white-space: nowrap;
 `),s("progress-icon",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 color: var(--n-icon-color);
 font-size: var(--n-icon-size-circle);
 `)]),$("multiple-circle",`
 width: 200px;
 color: inherit;
 `,[s("progress-text",`
 font-weight: var(--n-font-weight-circle);
 color: var(--n-text-color-circle);
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 display: flex;
 align-items: center;
 justify-content: center;
 transition: color .3s var(--n-bezier);
 `)]),s("progress-content",{position:"relative"}),s("progress-graph",{position:"relative"},[s("progress-graph-circle",[P("svg",{verticalAlign:"bottom"}),s("progress-graph-circle-fill",`
 stroke: var(--n-fill-color);
 transition:
 opacity .3s var(--n-bezier),
 stroke .3s var(--n-bezier),
 stroke-dasharray .3s var(--n-bezier);
 `,[$("empty",{opacity:0})]),s("progress-graph-circle-rail",`
 transition: stroke .3s var(--n-bezier);
 overflow: hidden;
 stroke: var(--n-rail-color);
 `)]),s("progress-graph-line",[$("indicator-inside",[s("progress-graph-line-rail",`
 height: 16px;
 line-height: 16px;
 border-radius: 10px;
 `,[s("progress-graph-line-fill",`
 height: inherit;
 border-radius: 10px;
 `),s("progress-graph-line-indicator",`
 background: #0000;
 white-space: nowrap;
 text-align: right;
 margin-left: 14px;
 margin-right: 14px;
 height: inherit;
 font-size: 12px;
 color: var(--n-text-color-line-inner);
 transition: color .3s var(--n-bezier);
 `)])]),$("indicator-inside-label",`
 height: 16px;
 display: flex;
 align-items: center;
 `,[s("progress-graph-line-rail",`
 flex: 1;
 transition: background-color .3s var(--n-bezier);
 `),s("progress-graph-line-indicator",`
 background: var(--n-fill-color);
 font-size: 12px;
 transform: translateZ(0);
 display: flex;
 vertical-align: middle;
 height: 16px;
 line-height: 16px;
 padding: 0 10px;
 border-radius: 10px;
 position: absolute;
 white-space: nowrap;
 color: var(--n-text-color-line-inner);
 transition:
 right .2s var(--n-bezier),
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `)]),s("progress-graph-line-rail",`
 position: relative;
 overflow: hidden;
 height: var(--n-rail-height);
 border-radius: 5px;
 background-color: var(--n-rail-color);
 transition: background-color .3s var(--n-bezier);
 `,[s("progress-graph-line-fill",`
 background: var(--n-fill-color);
 position: relative;
 border-radius: 5px;
 height: inherit;
 width: 100%;
 max-width: 0%;
 transition:
 background-color .3s var(--n-bezier),
 max-width .2s var(--n-bezier);
 `,[$("processing",[P("&::after",`
 content: "";
 background-image: var(--n-line-bg-processing);
 animation: progress-processing-animation 2s var(--n-bezier) infinite;
 `)])])])])])]),P("@keyframes progress-processing-animation",`
 0% {
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 right: 100%;
 opacity: 1;
 }
 66% {
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 right: 0;
 opacity: 0;
 }
 100% {
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 right: 0;
 opacity: 0;
 }
 `)]),Z=Object.assign(Object.assign({},L.props),{processing:Boolean,type:{type:String,default:"line"},gapDegree:Number,gapOffsetDegree:Number,status:{type:String,default:"default"},railColor:[String,Array],railStyle:[String,Array],color:[String,Array,Object],viewBoxWidth:{type:Number,default:100},strokeWidth:{type:Number,default:7},percentage:[Number,Array],unit:{type:String,default:"%"},showIndicator:{type:Boolean,default:!0},indicatorPosition:{type:String,default:"outside"},indicatorPlacement:{type:String,default:"outside"},indicatorTextColor:String,circleGap:{type:Number,default:1},height:Number,borderRadius:[String,Number],fillBorderRadius:[String,Number],offsetDegree:Number}),U=z({name:"Progress",props:Z,setup(e){const d=x(()=>e.indicatorPlacement||e.indicatorPosition),u=x(()=>{if(e.gapDegree||e.gapDegree===0)return e.gapDegree;if(e.type==="dashboard")return 75}),{mergedClsPrefixRef:f,inlineThemeDisabled:n}=_(e),o=L("Progress","-progress",K,X,e,f),t=x(()=>{const{status:a}=e,{common:{cubicBezierEaseInOut:h},self:{fontSize:m,fontSizeCircle:g,railColor:p,railHeight:y,iconSizeCircle:i,iconSizeLine:c,textColorCircle:C,textColorLineInner:S,textColorLineOuter:b,lineBgProcessing:v,fontWeightCircle:B,[W("iconColor",a)]:R,[W("fillColor",a)]:k}}=o.value;return{"--n-bezier":h,"--n-fill-color":k,"--n-font-size":m,"--n-font-size-circle":g,"--n-font-weight-circle":B,"--n-icon-color":R,"--n-icon-size-circle":i,"--n-icon-size-line":c,"--n-line-bg-processing":v,"--n-rail-color":p,"--n-rail-height":y,"--n-text-color-circle":C,"--n-text-color-line-inner":S,"--n-text-color-line-outer":b}}),l=n?A("progress",x(()=>e.status[0]),t,e):void 0;return{mergedClsPrefix:f,mergedIndicatorPlacement:d,gapDeg:u,cssVars:n?void 0:t,themeClass:l==null?void 0:l.themeClass,onRender:l==null?void 0:l.onRender}},render(){const{type:e,cssVars:d,indicatorTextColor:u,showIndicator:f,status:n,railColor:o,railStyle:t,color:l,percentage:a,viewBoxWidth:h,strokeWidth:m,mergedIndicatorPlacement:g,unit:p,borderRadius:y,fillBorderRadius:i,height:c,processing:C,circleGap:S,mergedClsPrefix:b,gapDeg:v,gapOffsetDegree:B,themeClass:R,$slots:k,onRender:I}=this;return I==null||I(),r("div",{class:[R,`${b}-progress`,`${b}-progress--${e}`,`${b}-progress--${n}`],style:d,"aria-valuemax":100,"aria-valuemin":0,"aria-valuenow":a,role:e==="circle"||e==="line"||e==="dashboard"?"progressbar":"none"},e==="circle"||e==="dashboard"?r(H,{clsPrefix:b,status:n,showIndicator:f,indicatorTextColor:u,railColor:o,fillColor:l,railStyle:t,offsetDegree:this.offsetDegree,percentage:a,viewBoxWidth:h,strokeWidth:m,gapDegree:v===void 0?e==="dashboard"?75:0:v,gapOffsetDegree:B,unit:p},k):e==="line"?r(V,{clsPrefix:b,status:n,showIndicator:f,indicatorTextColor:u,railColor:o,fillColor:l,railStyle:t,percentage:a,processing:C,indicatorPlacement:g,unit:p,fillBorderRadius:i,railBorderRadius:y,height:c},k):e==="multiple-circle"?r(F,{clsPrefix:b,strokeWidth:m,railColor:o,fillColor:l,railStyle:t,viewBoxWidth:h,percentage:a,showIndicator:f,circleGap:S},k):null)}});export{U as _,X as p};
