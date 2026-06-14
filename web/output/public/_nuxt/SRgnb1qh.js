import{i as K,u as se,N as de,c as ce}from"./QpI9WcJO.js";import{e as ue,aA as he,g as A,h as t,f as T,an as l,i as U,j as be,k as a,l as fe,u as I,m as j,aq as ge,n as z,ad as ve,al as m,bp as L,bX as s}from"./DmHPR5lg.js";import{i as D,a as p,c as N}from"./CoaUF789.js";import{u as me}from"./Bjh6Msyv.js";const pe={buttonHeightSmall:"14px",buttonHeightMedium:"18px",buttonHeightLarge:"22px",buttonWidthSmall:"14px",buttonWidthMedium:"18px",buttonWidthLarge:"22px",buttonWidthPressedSmall:"20px",buttonWidthPressedMedium:"24px",buttonWidthPressedLarge:"28px",railHeightSmall:"18px",railHeightMedium:"22px",railHeightLarge:"26px",railWidthSmall:"32px",railWidthMedium:"40px",railWidthLarge:"48px"};function we(e){const{primaryColor:d,opacityDisabled:f,borderRadius:n,textColor3:r}=e;return Object.assign(Object.assign({},pe),{iconColor:r,textColor:"white",loadingColor:d,opacityDisabled:f,railColor:"rgba(0, 0, 0, .14)",railColorActive:d,buttonBoxShadow:"0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)",buttonColor:"#FFF",railBorderRadiusSmall:n,railBorderRadiusMedium:n,railBorderRadiusLarge:n,buttonBorderRadiusSmall:n,buttonBorderRadiusMedium:n,buttonBorderRadiusLarge:n,boxShadowFocus:`0 0 0 2px ${he(d,{alpha:.2})}`})}const xe={common:ue,self:we},ye=A("switch",`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[t("children-placeholder",`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),t("rail-placeholder",`
 display: flex;
 flex-wrap: none;
 `),t("button-placeholder",`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),A("base-loading",`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[K({left:"50%",top:"50%",originalTransform:"translateX(-50%) translateY(-50%)"})]),t("checked, unchecked",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),t("checked",`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),t("unchecked",`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),T("&:focus",[t("rail",`
 box-shadow: var(--n-box-shadow-focus);
 `)]),l("round",[t("rail","border-radius: calc(var(--n-rail-height) / 2);",[t("button","border-radius: calc(var(--n-button-height) / 2);")])]),U("disabled",[U("icon",[l("rubber-band",[l("pressed",[t("rail",[t("button","max-width: var(--n-button-width-pressed);")])]),t("rail",[T("&:active",[t("button","max-width: var(--n-button-width-pressed);")])]),l("active",[l("pressed",[t("rail",[t("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])]),t("rail",[T("&:active",[t("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])])])])])]),l("active",[t("rail",[t("button","left: calc(100% - var(--n-button-width) - var(--n-offset))")])]),t("rail",`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[t("button-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[K()]),t("button",`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),l("active",[t("rail","background-color: var(--n-rail-color-active);")]),l("loading",[t("rail",`
 cursor: wait;
 `)]),l("disabled",[t("rail",`
 cursor: not-allowed;
 opacity: .5;
 `)])]),ke=Object.assign(Object.assign({},I.props),{size:{type:String,default:"medium"},value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},onChange:[Function,Array]});let B;const $e=be({name:"Switch",props:ke,slots:Object,setup(e){B===void 0&&(typeof CSS<"u"?typeof CSS.supports<"u"?B=CSS.supports("width","max(1px)"):B=!1:B=!0);const{mergedClsPrefixRef:d,inlineThemeDisabled:f}=fe(e),n=I("Switch","-switch",ye,xe,e,d),r=se(e),{mergedSizeRef:y,mergedDisabledRef:g}=r,k=j(e.defaultValue),R=ge(e,"value"),v=me(R,k),S=z(()=>v.value===e.checkedValue),w=j(!1),o=j(!1),c=z(()=>{const{railStyle:i}=e;if(i)return i({focused:o.value,checked:S.value})});function u(i){const{"onUpdate:value":$,onChange:_,onUpdateValue:V}=e,{nTriggerFormInput:F,nTriggerFormChange:W}=r;$&&N($,i),V&&N(V,i),_&&N(_,i),k.value=i,F(),W()}function E(){const{nTriggerFormFocus:i}=r;i()}function X(){const{nTriggerFormBlur:i}=r;i()}function Y(){e.loading||g.value||(v.value!==e.checkedValue?u(e.checkedValue):u(e.uncheckedValue))}function q(){o.value=!0,E()}function G(){o.value=!1,X(),w.value=!1}function J(i){e.loading||g.value||i.key===" "&&(v.value!==e.checkedValue?u(e.checkedValue):u(e.uncheckedValue),w.value=!1)}function Q(i){e.loading||g.value||i.key===" "&&(i.preventDefault(),w.value=!0)}const O=z(()=>{const{value:i}=y,{self:{opacityDisabled:$,railColor:_,railColorActive:V,buttonBoxShadow:F,buttonColor:W,boxShadowFocus:Z,loadingColor:ee,textColor:te,iconColor:ie,[m("buttonHeight",i)]:h,[m("buttonWidth",i)]:oe,[m("buttonWidthPressed",i)]:ae,[m("railHeight",i)]:b,[m("railWidth",i)]:C,[m("railBorderRadius",i)]:ne,[m("buttonBorderRadius",i)]:re},common:{cubicBezierEaseInOut:le}}=n.value;let M,H,P;return B?(M=`calc((${b} - ${h}) / 2)`,H=`max(${b}, ${h})`,P=`max(${C}, calc(${C} + ${h} - ${b}))`):(M=L((s(b)-s(h))/2),H=L(Math.max(s(b),s(h))),P=s(b)>s(h)?C:L(s(C)+s(h)-s(b))),{"--n-bezier":le,"--n-button-border-radius":re,"--n-button-box-shadow":F,"--n-button-color":W,"--n-button-width":oe,"--n-button-width-pressed":ae,"--n-button-height":h,"--n-height":H,"--n-offset":M,"--n-opacity-disabled":$,"--n-rail-border-radius":ne,"--n-rail-color":_,"--n-rail-color-active":V,"--n-rail-height":b,"--n-rail-width":C,"--n-width":P,"--n-box-shadow-focus":Z,"--n-loading-color":ee,"--n-text-color":te,"--n-icon-color":ie}}),x=f?ve("switch",z(()=>y.value[0]),O,e):void 0;return{handleClick:Y,handleBlur:G,handleFocus:q,handleKeyup:J,handleKeydown:Q,mergedRailStyle:c,pressed:w,mergedClsPrefix:d,mergedValue:v,checked:S,mergedDisabled:g,cssVars:f?void 0:O,themeClass:x==null?void 0:x.themeClass,onRender:x==null?void 0:x.onRender}},render(){const{mergedClsPrefix:e,mergedDisabled:d,checked:f,mergedRailStyle:n,onRender:r,$slots:y}=this;r==null||r();const{checked:g,unchecked:k,icon:R,"checked-icon":v,"unchecked-icon":S}=y,w=!(D(R)&&D(v)&&D(S));return a("div",{role:"switch","aria-checked":f,class:[`${e}-switch`,this.themeClass,w&&`${e}-switch--icon`,f&&`${e}-switch--active`,d&&`${e}-switch--disabled`,this.round&&`${e}-switch--round`,this.loading&&`${e}-switch--loading`,this.pressed&&`${e}-switch--pressed`,this.rubberBand&&`${e}-switch--rubber-band`],tabindex:this.mergedDisabled?void 0:0,style:this.cssVars,onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},a("div",{class:`${e}-switch__rail`,"aria-hidden":"true",style:n},p(g,o=>p(k,c=>o||c?a("div",{"aria-hidden":!0,class:`${e}-switch__children-placeholder`},a("div",{class:`${e}-switch__rail-placeholder`},a("div",{class:`${e}-switch__button-placeholder`}),o),a("div",{class:`${e}-switch__rail-placeholder`},a("div",{class:`${e}-switch__button-placeholder`}),c)):null)),a("div",{class:`${e}-switch__button`},p(R,o=>p(v,c=>p(S,u=>a(de,null,{default:()=>this.loading?a(ce,{key:"loading",clsPrefix:e,strokeWidth:20}):this.checked&&(c||o)?a("div",{class:`${e}-switch__button-icon`,key:c?"checked-icon":"icon"},c||o):!this.checked&&(u||o)?a("div",{class:`${e}-switch__button-icon`,key:u?"unchecked-icon":"icon"},u||o):null})))),p(g,o=>o&&a("div",{key:"checked",class:`${e}-switch__checked`},o)),p(k,o=>o&&a("div",{key:"unchecked",class:`${e}-switch__unchecked`},o)))))}});export{$e as _};
