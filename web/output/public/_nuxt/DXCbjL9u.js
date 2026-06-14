import{av as me,bk as ge,e as ve,j as A,a1 as _,k as s,b0 as q,bm as pe,aH as V,b7 as L,m as H,n as $,l as K,ao as we,ak as G,am as ye,s as Ce,ab as Se,ap as P,f as i,b1 as I,g as v,an as z,h as E,b6 as ze,u as J,aq as N,ad as $e,bf as xe,ar as Be}from"./DmHPR5lg.js";import{F as ke,c as Re,l as U,d as Ee,p as Te,a as He,L as Fe,z as Oe,e as Me}from"./C-4Oazno.js";import{u as Pe,a as Ie}from"./DYSVe14K.js";import{e as De}from"./QpI9WcJO.js";import{u as X}from"./Bjh6Msyv.js";import{c as T}from"./CoaUF789.js";import{f as Y}from"./B-p6aW7q.js";function We(e){const{modalColor:t,textColor1:r,textColor2:m,boxShadow3:u,lineHeight:w,fontWeightStrong:d,dividerColor:h,closeColorHover:g,closeColorPressed:C,closeIconColor:x,closeIconColorHover:B,closeIconColorPressed:k,borderRadius:y,primaryColorHover:f}=e;return{bodyPadding:"16px 24px",borderRadius:y,headerPadding:"16px 24px",footerPadding:"16px 24px",color:t,textColor:m,titleTextColor:r,titleFontSize:"18px",titleFontWeight:d,boxShadow:u,lineHeight:w,headerBorderBottom:`1px solid ${h}`,footerBorderTop:`1px solid ${h}`,closeIconColor:x,closeIconColorHover:B,closeIconColorPressed:k,closeSize:"22px",closeIconSize:"18px",closeColorHover:g,closeColorPressed:C,closeBorderRadius:y,resizableTriggerColorHover:f}}const je=me({name:"Drawer",common:ve,peers:{Scrollbar:ge},self:We}),_e=A({name:"NDrawerContent",inheritAttrs:!1,props:{blockScroll:Boolean,show:{type:Boolean,default:void 0},displayDirective:{type:String,required:!0},placement:{type:String,required:!0},contentClass:String,contentStyle:[Object,String],nativeScrollbar:{type:Boolean,required:!0},scrollbarProps:Object,trapFocus:{type:Boolean,default:!0},autoFocus:{type:Boolean,default:!0},showMask:{type:[Boolean,String],required:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,onClickoutside:Function,onAfterLeave:Function,onAfterEnter:Function,onEsc:Function},setup(e){const t=H(!!e.show),r=H(null),m=G(U);let u=0,w="",d=null;const h=H(!1),g=H(!1),C=$(()=>e.placement==="top"||e.placement==="bottom"),{mergedClsPrefixRef:x,mergedRtlRef:B}=K(e),k=we("Drawer",B,x),y=o,f=n=>{g.value=!0,u=C.value?n.clientY:n.clientX,w=document.body.style.cursor,document.body.style.cursor=C.value?"ns-resize":"ew-resize",document.body.addEventListener("mousemove",S),document.body.addEventListener("mouseleave",y),document.body.addEventListener("mouseup",o)},O=()=>{d!==null&&(window.clearTimeout(d),d=null),g.value?h.value=!0:d=window.setTimeout(()=>{h.value=!0},300)},D=()=>{d!==null&&(window.clearTimeout(d),d=null),h.value=!1},{doUpdateHeight:W,doUpdateWidth:j}=m,F=n=>{const{maxWidth:l}=e;if(l&&n>l)return l;const{minWidth:b}=e;return b&&n<b?b:n},M=n=>{const{maxHeight:l}=e;if(l&&n>l)return l;const{minHeight:b}=e;return b&&n<b?b:n};function S(n){var l,b;if(g.value)if(C.value){let p=((l=r.value)===null||l===void 0?void 0:l.offsetHeight)||0;const R=u-n.clientY;p+=e.placement==="bottom"?R:-R,p=M(p),W(p),u=n.clientY}else{let p=((b=r.value)===null||b===void 0?void 0:b.offsetWidth)||0;const R=u-n.clientX;p+=e.placement==="right"?R:-R,p=F(p),j(p),u=n.clientX}}function o(){g.value&&(u=0,g.value=!1,document.body.style.cursor=w,document.body.removeEventListener("mousemove",S),document.body.removeEventListener("mouseup",o),document.body.removeEventListener("mouseleave",y))}ye(()=>{e.show&&(t.value=!0)}),Ce(()=>e.show,n=>{n||o()}),Se(()=>{o()});const a=$(()=>{const{show:n}=e,l=[[L,n]];return e.showMask||l.push([Re,e.onClickoutside,void 0,{capture:!0}]),l});function c(){var n;t.value=!1,(n=e.onAfterLeave)===null||n===void 0||n.call(e)}return Pe($(()=>e.blockScroll&&t.value)),P(Ee,r),P(Te,null),P(He,null),{bodyRef:r,rtlEnabled:k,mergedClsPrefix:m.mergedClsPrefixRef,isMounted:m.isMountedRef,mergedTheme:m.mergedThemeRef,displayed:t,transitionName:$(()=>({right:"slide-in-from-right-transition",left:"slide-in-from-left-transition",top:"slide-in-from-top-transition",bottom:"slide-in-from-bottom-transition"})[e.placement]),handleAfterLeave:c,bodyDirectives:a,handleMousedownResizeTrigger:f,handleMouseenterResizeTrigger:O,handleMouseleaveResizeTrigger:D,isDragging:g,isHoverOnResizeTrigger:h}},render(){const{$slots:e,mergedClsPrefix:t}=this;return this.displayDirective==="show"||this.displayed||this.show?_(s("div",{role:"none"},s(ke,{disabled:!this.showMask||!this.trapFocus,active:this.show,autoFocus:this.autoFocus,onEsc:this.onEsc},{default:()=>s(q,{name:this.transitionName,appear:this.isMounted,onAfterEnter:this.onAfterEnter,onAfterLeave:this.handleAfterLeave},{default:()=>_(s("div",pe(this.$attrs,{role:"dialog",ref:"bodyRef","aria-modal":"true",class:[`${t}-drawer`,this.rtlEnabled&&`${t}-drawer--rtl`,`${t}-drawer--${this.placement}-placement`,this.isDragging&&`${t}-drawer--unselectable`,this.nativeScrollbar&&`${t}-drawer--native-scrollbar`]}),[this.resizable?s("div",{class:[`${t}-drawer__resize-trigger`,(this.isDragging||this.isHoverOnResizeTrigger)&&`${t}-drawer__resize-trigger--hover`],onMouseenter:this.handleMouseenterResizeTrigger,onMouseleave:this.handleMouseleaveResizeTrigger,onMousedown:this.handleMousedownResizeTrigger}):null,this.nativeScrollbar?s("div",{class:[`${t}-drawer-content-wrapper`,this.contentClass],style:this.contentStyle,role:"none"},e):s(V,Object.assign({},this.scrollbarProps,{contentStyle:this.contentStyle,contentClass:[`${t}-drawer-content-wrapper`,this.contentClass],theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar}),e)]),this.bodyDirectives)})})),[[L,this.displayDirective==="if"||this.displayed||this.show]]):null}}),{cubicBezierEaseIn:Ae,cubicBezierEaseOut:Ue}=I;function Le({duration:e="0.3s",leaveDuration:t="0.2s",name:r="slide-in-from-bottom"}={}){return[i(`&.${r}-transition-leave-active`,{transition:`transform ${t} ${Ae}`}),i(`&.${r}-transition-enter-active`,{transition:`transform ${e} ${Ue}`}),i(`&.${r}-transition-enter-to`,{transform:"translateY(0)"}),i(`&.${r}-transition-enter-from`,{transform:"translateY(100%)"}),i(`&.${r}-transition-leave-from`,{transform:"translateY(0)"}),i(`&.${r}-transition-leave-to`,{transform:"translateY(100%)"})]}const{cubicBezierEaseIn:Ne,cubicBezierEaseOut:Xe}=I;function Ye({duration:e="0.3s",leaveDuration:t="0.2s",name:r="slide-in-from-left"}={}){return[i(`&.${r}-transition-leave-active`,{transition:`transform ${t} ${Ne}`}),i(`&.${r}-transition-enter-active`,{transition:`transform ${e} ${Xe}`}),i(`&.${r}-transition-enter-to`,{transform:"translateX(0)"}),i(`&.${r}-transition-enter-from`,{transform:"translateX(-100%)"}),i(`&.${r}-transition-leave-from`,{transform:"translateX(0)"}),i(`&.${r}-transition-leave-to`,{transform:"translateX(-100%)"})]}const{cubicBezierEaseIn:qe,cubicBezierEaseOut:Ve}=I;function Ke({duration:e="0.3s",leaveDuration:t="0.2s",name:r="slide-in-from-right"}={}){return[i(`&.${r}-transition-leave-active`,{transition:`transform ${t} ${qe}`}),i(`&.${r}-transition-enter-active`,{transition:`transform ${e} ${Ve}`}),i(`&.${r}-transition-enter-to`,{transform:"translateX(0)"}),i(`&.${r}-transition-enter-from`,{transform:"translateX(100%)"}),i(`&.${r}-transition-leave-from`,{transform:"translateX(0)"}),i(`&.${r}-transition-leave-to`,{transform:"translateX(100%)"})]}const{cubicBezierEaseIn:Ge,cubicBezierEaseOut:Je}=I;function Qe({duration:e="0.3s",leaveDuration:t="0.2s",name:r="slide-in-from-top"}={}){return[i(`&.${r}-transition-leave-active`,{transition:`transform ${t} ${Ge}`}),i(`&.${r}-transition-enter-active`,{transition:`transform ${e} ${Je}`}),i(`&.${r}-transition-enter-to`,{transform:"translateY(0)"}),i(`&.${r}-transition-enter-from`,{transform:"translateY(-100%)"}),i(`&.${r}-transition-leave-from`,{transform:"translateY(0)"}),i(`&.${r}-transition-leave-to`,{transform:"translateY(-100%)"})]}const Ze=i([v("drawer",`
 word-break: break-word;
 line-height: var(--n-line-height);
 position: absolute;
 pointer-events: all;
 box-shadow: var(--n-box-shadow);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 background-color: var(--n-color);
 color: var(--n-text-color);
 box-sizing: border-box;
 `,[Ke(),Ye(),Qe(),Le(),z("unselectable",`
 user-select: none; 
 -webkit-user-select: none;
 `),z("native-scrollbar",[v("drawer-content-wrapper",`
 overflow: auto;
 height: 100%;
 `)]),E("resize-trigger",`
 position: absolute;
 background-color: #0000;
 transition: background-color .3s var(--n-bezier);
 `,[z("hover",`
 background-color: var(--n-resize-trigger-color-hover);
 `)]),v("drawer-content-wrapper",`
 box-sizing: border-box;
 `),v("drawer-content",`
 height: 100%;
 display: flex;
 flex-direction: column;
 `,[z("native-scrollbar",[v("drawer-body-content-wrapper",`
 height: 100%;
 overflow: auto;
 `)]),v("drawer-body",`
 flex: 1 0 0;
 overflow: hidden;
 `),v("drawer-body-content-wrapper",`
 box-sizing: border-box;
 padding: var(--n-body-padding);
 `),v("drawer-header",`
 font-weight: var(--n-title-font-weight);
 line-height: 1;
 font-size: var(--n-title-font-size);
 color: var(--n-title-text-color);
 padding: var(--n-header-padding);
 transition: border .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-divider-color);
 border-bottom: var(--n-header-border-bottom);
 display: flex;
 justify-content: space-between;
 align-items: center;
 `,[E("main",`
 flex: 1;
 `),E("close",`
 margin-left: 6px;
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)]),v("drawer-footer",`
 display: flex;
 justify-content: flex-end;
 border-top: var(--n-footer-border-top);
 transition: border .3s var(--n-bezier);
 padding: var(--n-footer-padding);
 `)]),z("right-placement",`
 top: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[E("resize-trigger",`
 width: 3px;
 height: 100%;
 top: 0;
 left: 0;
 transform: translateX(-1.5px);
 cursor: ew-resize;
 `)]),z("left-placement",`
 top: 0;
 bottom: 0;
 left: 0;
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[E("resize-trigger",`
 width: 3px;
 height: 100%;
 top: 0;
 right: 0;
 transform: translateX(1.5px);
 cursor: ew-resize;
 `)]),z("top-placement",`
 top: 0;
 left: 0;
 right: 0;
 border-bottom-left-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[E("resize-trigger",`
 width: 100%;
 height: 3px;
 bottom: 0;
 left: 0;
 transform: translateY(1.5px);
 cursor: ns-resize;
 `)]),z("bottom-placement",`
 left: 0;
 bottom: 0;
 right: 0;
 border-top-left-radius: var(--n-border-radius);
 border-top-right-radius: var(--n-border-radius);
 `,[E("resize-trigger",`
 width: 100%;
 height: 3px;
 top: 0;
 left: 0;
 transform: translateY(-1.5px);
 cursor: ns-resize;
 `)])]),i("body",[i(">",[v("drawer-container",`
 position: fixed;
 `)])]),v("drawer-container",`
 position: relative;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 `,[i("> *",`
 pointer-events: all;
 `)]),v("drawer-mask",`
 background-color: rgba(0, 0, 0, .3);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `,[z("invisible",`
 background-color: rgba(0, 0, 0, 0)
 `),ze({enterDuration:"0.2s",leaveDuration:"0.2s",enterCubicBezier:"var(--n-bezier-in)",leaveCubicBezier:"var(--n-bezier-out)"})])]),et=Object.assign(Object.assign({},J.props),{show:Boolean,width:[Number,String],height:[Number,String],placement:{type:String,default:"right"},maskClosable:{type:Boolean,default:!0},showMask:{type:[Boolean,String],default:!0},to:[String,Object],displayDirective:{type:String,default:"if"},nativeScrollbar:{type:Boolean,default:!0},zIndex:Number,onMaskClick:Function,scrollbarProps:Object,contentClass:String,contentStyle:[Object,String],trapFocus:{type:Boolean,default:!0},onEsc:Function,autoFocus:{type:Boolean,default:!0},closeOnEsc:{type:Boolean,default:!0},blockScroll:{type:Boolean,default:!0},maxWidth:Number,maxHeight:Number,minWidth:Number,minHeight:Number,resizable:Boolean,defaultWidth:{type:[Number,String],default:251},defaultHeight:{type:[Number,String],default:251},onUpdateWidth:[Function,Array],onUpdateHeight:[Function,Array],"onUpdate:width":[Function,Array],"onUpdate:height":[Function,Array],"onUpdate:show":[Function,Array],onUpdateShow:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,drawerStyle:[String,Object],drawerClass:String,target:null,onShow:Function,onHide:Function}),dt=A({name:"Drawer",inheritAttrs:!1,props:et,setup(e){const{mergedClsPrefixRef:t,namespaceRef:r,inlineThemeDisabled:m}=K(e),u=De(),w=J("Drawer","-drawer",Ze,je,e,t),d=H(e.defaultWidth),h=H(e.defaultHeight),g=X(N(e,"width"),d),C=X(N(e,"height"),h),x=$(()=>{const{placement:o}=e;return o==="top"||o==="bottom"?"":Y(g.value)}),B=$(()=>{const{placement:o}=e;return o==="left"||o==="right"?"":Y(C.value)}),k=o=>{const{onUpdateWidth:a,"onUpdate:width":c}=e;a&&T(a,o),c&&T(c,o),d.value=o},y=o=>{const{onUpdateHeight:a,"onUpdate:width":c}=e;a&&T(a,o),c&&T(c,o),h.value=o},f=$(()=>[{width:x.value,height:B.value},e.drawerStyle||""]);function O(o){const{onMaskClick:a,maskClosable:c}=e;c&&F(!1),a&&a(o)}function D(o){O(o)}const W=Ie();function j(o){var a;(a=e.onEsc)===null||a===void 0||a.call(e),e.show&&e.closeOnEsc&&Me(o)&&(W.value||F(!1))}function F(o){const{onHide:a,onUpdateShow:c,"onUpdate:show":n}=e;c&&T(c,o),n&&T(n,o),a&&!o&&T(a,o)}P(U,{isMountedRef:u,mergedThemeRef:w,mergedClsPrefixRef:t,doUpdateShow:F,doUpdateHeight:y,doUpdateWidth:k});const M=$(()=>{const{common:{cubicBezierEaseInOut:o,cubicBezierEaseIn:a,cubicBezierEaseOut:c},self:{color:n,textColor:l,boxShadow:b,lineHeight:p,headerPadding:R,footerPadding:Q,borderRadius:Z,bodyPadding:ee,titleFontSize:te,titleTextColor:re,titleFontWeight:oe,headerBorderBottom:ne,footerBorderTop:ie,closeIconColor:se,closeIconColorHover:ae,closeIconColorPressed:le,closeColorHover:de,closeColorPressed:ce,closeIconSize:ue,closeSize:he,closeBorderRadius:fe,resizableTriggerColorHover:be}}=w.value;return{"--n-line-height":p,"--n-color":n,"--n-border-radius":Z,"--n-text-color":l,"--n-box-shadow":b,"--n-bezier":o,"--n-bezier-out":c,"--n-bezier-in":a,"--n-header-padding":R,"--n-body-padding":ee,"--n-footer-padding":Q,"--n-title-text-color":re,"--n-title-font-size":te,"--n-title-font-weight":oe,"--n-header-border-bottom":ne,"--n-footer-border-top":ie,"--n-close-icon-color":se,"--n-close-icon-color-hover":ae,"--n-close-icon-color-pressed":le,"--n-close-size":he,"--n-close-color-hover":de,"--n-close-color-pressed":ce,"--n-close-icon-size":ue,"--n-close-border-radius":fe,"--n-resize-trigger-color-hover":be}}),S=m?$e("drawer",void 0,M,e):void 0;return{mergedClsPrefix:t,namespace:r,mergedBodyStyle:f,handleOutsideClick:D,handleMaskClick:O,handleEsc:j,mergedTheme:w,cssVars:m?void 0:M,themeClass:S==null?void 0:S.themeClass,onRender:S==null?void 0:S.onRender,isMounted:u}},render(){const{mergedClsPrefix:e}=this;return s(Fe,{to:this.to,show:this.show},{default:()=>{var t;return(t=this.onRender)===null||t===void 0||t.call(this),_(s("div",{class:[`${e}-drawer-container`,this.namespace,this.themeClass],style:this.cssVars,role:"none"},this.showMask?s(q,{name:"fade-in-transition",appear:this.isMounted},{default:()=>this.show?s("div",{"aria-hidden":!0,class:[`${e}-drawer-mask`,this.showMask==="transparent"&&`${e}-drawer-mask--invisible`],onClick:this.handleMaskClick}):null}):null,s(_e,Object.assign({},this.$attrs,{class:[this.drawerClass,this.$attrs.class],style:[this.mergedBodyStyle,this.$attrs.style],blockScroll:this.blockScroll,contentStyle:this.contentStyle,contentClass:this.contentClass,placement:this.placement,scrollbarProps:this.scrollbarProps,show:this.show,displayDirective:this.displayDirective,nativeScrollbar:this.nativeScrollbar,onAfterEnter:this.onAfterEnter,onAfterLeave:this.onAfterLeave,trapFocus:this.trapFocus,autoFocus:this.autoFocus,resizable:this.resizable,maxHeight:this.maxHeight,minHeight:this.minHeight,maxWidth:this.maxWidth,minWidth:this.minWidth,showMask:this.showMask,onEsc:this.handleEsc,onClickoutside:this.handleOutsideClick}),this.$slots)),[[Oe,{zIndex:this.zIndex,enabled:this.show}]])}})}}),tt={title:String,headerClass:String,headerStyle:[Object,String],footerClass:String,footerStyle:[Object,String],bodyClass:String,bodyStyle:[Object,String],bodyContentClass:String,bodyContentStyle:[Object,String],nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,closable:Boolean},ct=A({name:"DrawerContent",props:tt,slots:Object,setup(){const e=G(U,null);e||Be("drawer-content","`n-drawer-content` must be placed inside `n-drawer`.");const{doUpdateShow:t}=e;function r(){t(!1)}return{handleCloseClick:r,mergedTheme:e.mergedThemeRef,mergedClsPrefix:e.mergedClsPrefixRef}},render(){const{title:e,mergedClsPrefix:t,nativeScrollbar:r,mergedTheme:m,bodyClass:u,bodyStyle:w,bodyContentClass:d,bodyContentStyle:h,headerClass:g,headerStyle:C,footerClass:x,footerStyle:B,scrollbarProps:k,closable:y,$slots:f}=this;return s("div",{role:"none",class:[`${t}-drawer-content`,r&&`${t}-drawer-content--native-scrollbar`]},f.header||e||y?s("div",{class:[`${t}-drawer-header`,g],style:C,role:"none"},s("div",{class:`${t}-drawer-header__main`,role:"heading","aria-level":"1"},f.header!==void 0?f.header():e),y&&s(xe,{onClick:this.handleCloseClick,clsPrefix:t,class:`${t}-drawer-header__close`,absolute:!0})):null,r?s("div",{class:[`${t}-drawer-body`,u],style:w,role:"none"},s("div",{class:[`${t}-drawer-body-content-wrapper`,d],style:h,role:"none"},f)):s(V,Object.assign({themeOverrides:m.peerOverrides.Scrollbar,theme:m.peers.Scrollbar},k,{class:`${t}-drawer-body`,contentClass:[`${t}-drawer-body-content-wrapper`,d],contentStyle:h}),f),f.footer?s("div",{class:[`${t}-drawer-footer`,x],style:B,role:"none"},f.footer()):null)}});export{ct as _,dt as a};
