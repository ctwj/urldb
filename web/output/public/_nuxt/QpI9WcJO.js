import{n as c,m as p,s as $,v as y,ae as k,ab as _,ak as C,ah as B,ap as S,j as v,k as s,b0 as I,f as l,b1 as N,g as T,h as u,aG as x,aq as F}from"./DmHPR5lg.js";function q(e){const n=c(e),i=p(n.value);return $(n,o=>{i.value=o}),typeof e=="function"?i:{__v_isRef:!0,get value(){return i.value},set value(o){e.set(o)}}}function R(){const e=p(!1);return y(()=>{e.value=!0}),k(e)}const f=typeof document<"u"&&typeof window<"u",g=B("n-form-item");function A(e,{defaultSize:n="medium",mergedSize:i,mergedDisabled:o}={}){const t=C(g,null);S(g,null);const a=c(i?()=>i(t):()=>{const{size:r}=e;if(r)return r;if(t){const{mergedSize:m}=t;if(m.value!==void 0)return m.value}return n}),b=c(o?()=>o(t):()=>{const{disabled:r}=e;return r!==void 0?r:t?t.disabled.value:!1}),w=c(()=>{const{status:r}=e;return r||(t==null?void 0:t.mergedValidationStatus.value)});return _(()=>{t&&t.restoreValidation()}),{mergedSizeRef:a,mergedDisabledRef:b,mergedStatusRef:w,nTriggerFormBlur(){t&&t.handleContentBlur()},nTriggerFormChange(){t&&t.handleContentChange()},nTriggerFormFocus(){t&&t.handleContentFocus()},nTriggerFormInput(){t&&t.handleContentInput()}}}const z=v({name:"BaseIconSwitchTransition",setup(e,{slots:n}){const i=R();return()=>s(I,{name:"icon-switch-transition",appear:i.value},n)}}),{cubicBezierEaseInOut:j}=N;function h({originalTransform:e="",left:n=0,top:i=0,transition:o=`all .3s ${j} !important`}={}){return[l("&.icon-switch-transition-enter-from, &.icon-switch-transition-leave-to",{transform:`${e} scale(0.75)`,left:n,top:i,opacity:0}),l("&.icon-switch-transition-enter-to, &.icon-switch-transition-leave-from",{transform:`scale(1) ${e}`,left:n,top:i,opacity:1}),l("&.icon-switch-transition-enter-active, &.icon-switch-transition-leave-active",{transformOrigin:"center",position:"absolute",left:n,top:i,transition:o})]}const M=l([l("@keyframes rotator",`
 0% {
 -webkit-transform: rotate(0deg);
 transform: rotate(0deg);
 }
 100% {
 -webkit-transform: rotate(360deg);
 transform: rotate(360deg);
 }`),T("base-loading",`
 position: relative;
 line-height: 0;
 width: 1em;
 height: 1em;
 `,[u("transition-wrapper",`
 position: absolute;
 width: 100%;
 height: 100%;
 `,[h()]),u("placeholder",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[h({left:"50%",top:"50%",originalTransform:"translateX(-50%) translateY(-50%)"})]),u("container",`
 animation: rotator 3s linear infinite both;
 `,[u("icon",`
 height: 1em;
 width: 1em;
 `)])])]),d="1.6s",P={strokeWidth:{type:Number,default:28},stroke:{type:String,default:void 0}},E=v({name:"BaseLoading",props:Object.assign({clsPrefix:{type:String,required:!0},show:{type:Boolean,default:!0},scale:{type:Number,default:1},radius:{type:Number,default:100}},P),setup(e){x("-base-loading",M,F(e,"clsPrefix"))},render(){const{clsPrefix:e,radius:n,strokeWidth:i,stroke:o,scale:t}=this,a=n/t;return s("div",{class:`${e}-base-loading`,role:"img","aria-label":"loading"},s(z,null,{default:()=>this.show?s("div",{key:"icon",class:`${e}-base-loading__transition-wrapper`},s("div",{class:`${e}-base-loading__container`},s("svg",{class:`${e}-base-loading__icon`,viewBox:`0 0 ${2*a} ${2*a}`,xmlns:"http://www.w3.org/2000/svg",style:{color:o}},s("g",null,s("animateTransform",{attributeName:"transform",type:"rotate",values:`0 ${a} ${a};270 ${a} ${a}`,begin:"0s",dur:d,fill:"freeze",repeatCount:"indefinite"}),s("circle",{class:`${e}-base-loading__icon`,fill:"none",stroke:"currentColor","stroke-width":i,"stroke-linecap":"round",cx:a,cy:a,r:n-i/2,"stroke-dasharray":5.67*n,"stroke-dashoffset":18.48*n},s("animateTransform",{attributeName:"transform",type:"rotate",values:`0 ${a} ${a};135 ${a} ${a};450 ${a} ${a}`,begin:"0s",dur:d,fill:"freeze",repeatCount:"indefinite"}),s("animate",{attributeName:"stroke-dashoffset",values:`${5.67*n};${1.42*n};${5.67*n}`,begin:"0s",dur:d,fill:"freeze",repeatCount:"indefinite"})))))):s("div",{key:"placeholder",class:`${e}-base-loading__placeholder`},this.$slots)}))}}),V=f&&"chrome"in window;f&&navigator.userAgent.includes("Firefox");const K=f&&navigator.userAgent.includes("Safari")&&!V;export{z as N,q as a,f as b,E as c,K as d,R as e,g as f,h as i,A as u};
