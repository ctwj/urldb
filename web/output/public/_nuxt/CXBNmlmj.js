import{e as S,f,g as p,b6 as C,an as h,j as x,k as r,b0 as k,l as T,u as v,n as u,ad as $,m as w,am as R,bp as _,al as j}from"./DmHPR5lg.js";import{c as B}from"./QpI9WcJO.js";import{u as O}from"./D__gE1zc.js";function N(e){const{opacityDisabled:a,heightTiny:t,heightSmall:i,heightMedium:l,heightLarge:s,heightHuge:c,primaryColor:o,fontSize:n}=e;return{fontSize:n,textColor:o,sizeTiny:t,sizeSmall:i,sizeMedium:l,sizeLarge:s,sizeHuge:c,color:o,opacitySpinning:a}}const P={common:S,self:N},V=f([f("@keyframes spin-rotate",`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),p("spin-container",`
 position: relative;
 `,[p("spin-body",`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[C()])]),p("spin-body",`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),p("spin",`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[h("rotate",`
 animation: spin-rotate 2s linear infinite;
 `)]),p("spin-description",`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),p("spin-content",`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[h("spinning",`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),W={small:20,medium:18,large:16},H=Object.assign(Object.assign({},v.props),{contentClass:String,contentStyle:[Object,String],description:String,stroke:String,size:{type:[String,Number],default:"medium"},show:{type:Boolean,default:!0},strokeWidth:Number,rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),E=x({name:"Spin",props:H,slots:Object,setup(e){const{mergedClsPrefixRef:a,inlineThemeDisabled:t}=T(e),i=v("Spin","-spin",V,P,e,a),l=u(()=>{const{size:n}=e,{common:{cubicBezierEaseInOut:d},self:m}=i.value,{opacitySpinning:g,color:y,textColor:b}=m,z=typeof n=="number"?_(n):m[j("size",n)];return{"--n-bezier":d,"--n-opacity-spinning":g,"--n-size":z,"--n-color":y,"--n-text-color":b}}),s=t?$("spin",u(()=>{const{size:n}=e;return typeof n=="number"?String(n):n[0]}),l,e):void 0,c=O(e,["spinning","show"]),o=w(!1);return R(n=>{let d;if(c.value){const{delay:m}=e;if(m){d=window.setTimeout(()=>{o.value=!0},m),n(()=>{clearTimeout(d)});return}}o.value=c.value}),{mergedClsPrefix:a,active:o,mergedStrokeWidth:u(()=>{const{strokeWidth:n}=e;if(n!==void 0)return n;const{size:d}=e;return W[typeof d=="number"?"medium":d]}),cssVars:t?void 0:l,themeClass:s==null?void 0:s.themeClass,onRender:s==null?void 0:s.onRender}},render(){var e,a;const{$slots:t,mergedClsPrefix:i,description:l}=this,s=t.icon&&this.rotate,c=(l||t.description)&&r("div",{class:`${i}-spin-description`},l||((e=t.description)===null||e===void 0?void 0:e.call(t))),o=t.icon?r("div",{class:[`${i}-spin-body`,this.themeClass]},r("div",{class:[`${i}-spin`,s&&`${i}-spin--rotate`],style:t.default?"":this.cssVars},t.icon()),c):r("div",{class:[`${i}-spin-body`,this.themeClass]},r(B,{clsPrefix:i,style:t.default?"":this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,class:`${i}-spin`}),c);return(a=this.onRender)===null||a===void 0||a.call(this),t.default?r("div",{class:[`${i}-spin-container`,this.themeClass],style:this.cssVars},r("div",{class:[`${i}-spin-content`,this.active&&`${i}-spin-content--spinning`,this.contentClass],style:this.contentStyle},t),r(k,{name:"fade-in-transition"},{default:()=>this.active?o:null})):o}});export{E as _};
