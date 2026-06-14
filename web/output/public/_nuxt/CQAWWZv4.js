import{a as j,m as Q}from"./6B1W4Q65.js";import{f as U}from"./CiG8jYj-.js";import{i as X,N as Z,c as ee}from"./QpI9WcJO.js";import{e as oe,f as C,g as x,h as I,an as b,j as w,k as a,be as ne,bf as re,l as P,ak as se,ao as ie,u as z,n as S,al as p,ad as te,ax as le,bg as ae,bh as ce,bi as de,bj as ue,m as y,v as ge,aD as fe,c0 as ve,F as he,ap as O,aE as me,Y as pe}from"./DmHPR5lg.js";import{N as be}from"./DO8alW5h.js";const Ce={margin:"0 0 8px 0",padding:"10px 20px",maxWidth:"720px",minWidth:"420px",iconMargin:"0 10px 0 0",closeMargin:"0 0 0 10px",closeSize:"20px",closeIconSize:"16px",iconSize:"20px",fontSize:"14px"};function xe(t){const{textColor2:o,closeIconColor:l,closeIconColorHover:n,closeIconColorPressed:i,infoColor:c,successColor:u,errorColor:g,warningColor:r,popoverColor:e,boxShadow2:s,primaryColor:d,lineHeight:v,borderRadius:f,closeColorHover:h,closeColorPressed:m}=t;return Object.assign(Object.assign({},Ce),{closeBorderRadius:f,textColor:o,textColorInfo:o,textColorSuccess:o,textColorError:o,textColorWarning:o,textColorLoading:o,color:e,colorInfo:e,colorSuccess:e,colorError:e,colorWarning:e,colorLoading:e,boxShadow:s,boxShadowInfo:s,boxShadowSuccess:s,boxShadowError:s,boxShadowWarning:s,boxShadowLoading:s,iconColor:o,iconColorInfo:c,iconColorSuccess:u,iconColorWarning:r,iconColorError:g,iconColorLoading:d,closeColorHover:h,closeColorPressed:m,closeIconColor:l,closeIconColorHover:n,closeIconColorPressed:i,closeColorHoverInfo:h,closeColorPressedInfo:m,closeIconColorInfo:l,closeIconColorHoverInfo:n,closeIconColorPressedInfo:i,closeColorHoverSuccess:h,closeColorPressedSuccess:m,closeIconColorSuccess:l,closeIconColorHoverSuccess:n,closeIconColorPressedSuccess:i,closeColorHoverError:h,closeColorPressedError:m,closeIconColorError:l,closeIconColorHoverError:n,closeIconColorPressedError:i,closeColorHoverWarning:h,closeColorPressedWarning:m,closeIconColorWarning:l,closeIconColorHoverWarning:n,closeIconColorPressedWarning:i,closeColorHoverLoading:h,closeColorPressedLoading:m,closeIconColorLoading:l,closeIconColorHoverLoading:n,closeIconColorPressedLoading:i,loadingColor:d,lineHeight:v,borderRadius:f})}const Ie={common:oe,self:xe},H={icon:Function,type:{type:String,default:"info"},content:[String,Number,Function],showIcon:{type:Boolean,default:!0},closable:Boolean,keepAliveOnHover:Boolean,onClose:Function,onMouseenter:Function,onMouseleave:Function},ye=C([x("message-wrapper",`
 margin: var(--n-margin);
 z-index: 0;
 transform-origin: top center;
 display: flex;
 `,[U({overflow:"visible",originalTransition:"transform .3s var(--n-bezier)",enterToProps:{transform:"scale(1)"},leaveToProps:{transform:"scale(0.85)"}})]),x("message",`
 box-sizing: border-box;
 display: flex;
 align-items: center;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 transform .3s var(--n-bezier),
 margin-bottom .3s var(--n-bezier);
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 flex-wrap: nowrap;
 overflow: hidden;
 max-width: var(--n-max-width);
 color: var(--n-text-color);
 background-color: var(--n-color);
 box-shadow: var(--n-box-shadow);
 `,[I("content",`
 display: inline-block;
 line-height: var(--n-line-height);
 font-size: var(--n-font-size);
 `),I("icon",`
 position: relative;
 margin: var(--n-icon-margin);
 height: var(--n-icon-size);
 width: var(--n-icon-size);
 font-size: var(--n-icon-size);
 flex-shrink: 0;
 `,[["default","info","success","warning","error","loading"].map(t=>b(`${t}-type`,[C("> *",`
 color: var(--n-icon-color-${t});
 transition: color .3s var(--n-bezier);
 `)])),C("> *",`
 position: absolute;
 left: 0;
 top: 0;
 right: 0;
 bottom: 0;
 `,[X()])]),I("close",`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 flex-shrink: 0;
 `,[C("&:hover",`
 color: var(--n-close-icon-color-hover);
 `),C("&:active",`
 color: var(--n-close-icon-color-pressed);
 `)])]),x("message-container",`
 z-index: 6000;
 position: fixed;
 height: 0;
 overflow: visible;
 display: flex;
 flex-direction: column;
 align-items: center;
 `,[b("top",`
 top: 12px;
 left: 0;
 right: 0;
 `),b("top-left",`
 top: 12px;
 left: 12px;
 right: 0;
 align-items: flex-start;
 `),b("top-right",`
 top: 12px;
 left: 0;
 right: 12px;
 align-items: flex-end;
 `),b("bottom",`
 bottom: 4px;
 left: 0;
 right: 0;
 justify-content: flex-end;
 `),b("bottom-left",`
 bottom: 4px;
 left: 12px;
 right: 0;
 justify-content: flex-end;
 align-items: flex-start;
 `),b("bottom-right",`
 bottom: 4px;
 left: 0;
 right: 12px;
 justify-content: flex-end;
 align-items: flex-end;
 `)])]),we={info:()=>a(ue,null),success:()=>a(de,null),warning:()=>a(ce,null),error:()=>a(ae,null),default:()=>null},Se=w({name:"Message",props:Object.assign(Object.assign({},H),{render:Function}),setup(t){const{inlineThemeDisabled:o,mergedRtlRef:l}=P(t),{props:n,mergedClsPrefixRef:i}=se(j),c=ie("Message",l,i),u=z("Message","-message",ye,Ie,n,i),g=S(()=>{const{type:e}=t,{common:{cubicBezierEaseInOut:s},self:{padding:d,margin:v,maxWidth:f,iconMargin:h,closeMargin:m,closeSize:L,iconSize:M,fontSize:k,lineHeight:A,borderRadius:E,iconColorInfo:R,iconColorSuccess:_,iconColorWarning:T,iconColorError:$,iconColorLoading:W,closeIconSize:F,closeBorderRadius:N,[p("textColor",e)]:B,[p("boxShadow",e)]:K,[p("color",e)]:V,[p("closeColorHover",e)]:D,[p("closeColorPressed",e)]:q,[p("closeIconColor",e)]:Y,[p("closeIconColorPressed",e)]:G,[p("closeIconColorHover",e)]:J}}=u.value;return{"--n-bezier":s,"--n-margin":v,"--n-padding":d,"--n-max-width":f,"--n-font-size":k,"--n-icon-margin":h,"--n-icon-size":M,"--n-close-icon-size":F,"--n-close-border-radius":N,"--n-close-size":L,"--n-close-margin":m,"--n-text-color":B,"--n-color":V,"--n-box-shadow":K,"--n-icon-color-info":R,"--n-icon-color-success":_,"--n-icon-color-warning":T,"--n-icon-color-error":$,"--n-icon-color-loading":W,"--n-close-color-hover":D,"--n-close-color-pressed":q,"--n-close-icon-color":Y,"--n-close-icon-color-pressed":G,"--n-close-icon-color-hover":J,"--n-line-height":A,"--n-border-radius":E}}),r=o?te("message",S(()=>t.type[0]),g,{}):void 0;return{mergedClsPrefix:i,rtlEnabled:c,messageProviderProps:n,handleClose(){var e;(e=t.onClose)===null||e===void 0||e.call(t)},cssVars:o?void 0:g,themeClass:r==null?void 0:r.themeClass,onRender:r==null?void 0:r.onRender,placement:n.placement}},render(){const{render:t,type:o,closable:l,content:n,mergedClsPrefix:i,cssVars:c,themeClass:u,onRender:g,icon:r,handleClose:e,showIcon:s}=this;g==null||g();let d;return a("div",{class:[`${i}-message-wrapper`,u],onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave,style:[{alignItems:this.placement.startsWith("top")?"flex-start":"flex-end"},c]},t?t(this.$props):a("div",{class:[`${i}-message ${i}-message--${o}-type`,this.rtlEnabled&&`${i}-message--rtl`]},(d=Oe(r,o,i))&&s?a("div",{class:`${i}-message__icon ${i}-message__icon--${o}-type`},a(Z,null,{default:()=>d})):null,a("div",{class:`${i}-message__content`},ne(n)),l?a(re,{clsPrefix:i,class:`${i}-message__close`,onClick:e,absolute:!0}):null))}});function Oe(t,o,l){if(typeof t=="function")return t();{const n=o==="loading"?a(ee,{clsPrefix:l,strokeWidth:24,scale:.85}):we[o]();return n?a(le,{clsPrefix:l,key:o},{default:()=>n}):null}}const je=w({name:"MessageEnvironment",props:Object.assign(Object.assign({},H),{duration:{type:Number,default:3e3},onAfterLeave:Function,onLeave:Function,internalKey:{type:String,required:!0},onInternalAfterLeave:Function,onHide:Function,onAfterHide:Function}),setup(t){let o=null;const l=y(!0);ge(()=>{n()});function n(){const{duration:s}=t;s&&(o=window.setTimeout(u,s))}function i(s){s.currentTarget===s.target&&o!==null&&(window.clearTimeout(o),o=null)}function c(s){s.currentTarget===s.target&&n()}function u(){const{onHide:s}=t;l.value=!1,o&&(window.clearTimeout(o),o=null),s&&s()}function g(){const{onClose:s}=t;s&&s(),u()}function r(){const{onAfterLeave:s,onInternalAfterLeave:d,onAfterHide:v,internalKey:f}=t;s&&s(),d&&d(f),v&&v()}function e(){u()}return{show:l,hide:u,handleClose:g,handleAfterLeave:r,handleMouseleave:c,handleMouseenter:i,deactivate:e}},render(){return a(be,{appear:!0,onAfterLeave:this.handleAfterLeave,onLeave:this.onLeave},{default:()=>[this.show?a(Se,{content:this.content,type:this.type,icon:this.icon,showIcon:this.showIcon,closable:this.closable,onClose:this.handleClose,onMouseenter:this.keepAliveOnHover?this.handleMouseenter:void 0,onMouseleave:this.keepAliveOnHover?this.handleMouseleave:void 0}):null]})}}),Pe=Object.assign(Object.assign({},z.props),{to:[String,Object],duration:{type:Number,default:3e3},keepAliveOnHover:Boolean,max:Number,placement:{type:String,default:"top"},closable:Boolean,containerClass:String,containerStyle:[String,Object]}),Ae=w({name:"MessageProvider",props:Pe,setup(t){const{mergedClsPrefixRef:o}=P(t),l=y([]),n=y({}),i={create(r,e){return c(r,Object.assign({type:"default"},e))},info(r,e){return c(r,Object.assign(Object.assign({},e),{type:"info"}))},success(r,e){return c(r,Object.assign(Object.assign({},e),{type:"success"}))},warning(r,e){return c(r,Object.assign(Object.assign({},e),{type:"warning"}))},error(r,e){return c(r,Object.assign(Object.assign({},e),{type:"error"}))},loading(r,e){return c(r,Object.assign(Object.assign({},e),{type:"loading"}))},destroyAll:g};O(j,{props:t,mergedClsPrefixRef:o}),O(Q,i);function c(r,e){const s=me(),d=pe(Object.assign(Object.assign({},e),{content:r,key:s,destroy:()=>{var f;(f=n.value[s])===null||f===void 0||f.hide()}})),{max:v}=t;return v&&l.value.length>=v&&l.value.shift(),l.value.push(d),d}function u(r){l.value.splice(l.value.findIndex(e=>e.key===r),1),delete n.value[r]}function g(){Object.values(n.value).forEach(r=>{r.hide()})}return Object.assign({mergedClsPrefix:o,messageRefs:n,messageList:l,handleAfterLeave:u},i)},render(){var t,o,l;return a(he,null,(o=(t=this.$slots).default)===null||o===void 0?void 0:o.call(t),this.messageList.length?a(fe,{to:(l=this.to)!==null&&l!==void 0?l:"body"},a("div",{class:[`${this.mergedClsPrefix}-message-container`,`${this.mergedClsPrefix}-message-container--${this.placement}`,this.containerClass],key:"message-container",style:this.containerStyle},this.messageList.map(n=>a(je,Object.assign({ref:i=>{i&&(this.messageRefs[n.key]=i)},internalKey:n.key,onInternalAfterLeave:this.handleAfterLeave},ve(n,["destroy"],void 0),{duration:n.duration===void 0?this.duration:n.duration,keepAliveOnHover:n.keepAliveOnHover===void 0?this.keepAliveOnHover:n.keepAliveOnHover,closable:n.closable===void 0?this.closable:n.closable}))))):null)}});export{Ae as _};
