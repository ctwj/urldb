import{e as se,aA as be,j as H,k as u,l as N,m as P,n as A,ap as ue,aq as U,ah as he,f as x,g as n,ai as fe,aj as ke,an as D,h as S,ak as ve,u as K,ao as me,ad as xe,aE as ge,ay as pe,al as j}from"./DmHPR5lg.js";import{u as E,i as Ce,N as ye,a as ze}from"./QpI9WcJO.js";import{u as L}from"./Bjh6Msyv.js";import{c as i,a as we}from"./CoaUF789.js";const Re={sizeSmall:"14px",sizeMedium:"16px",sizeLarge:"18px",labelPadding:"0 8px",labelFontWeight:"400"};function Se(o){const{baseColor:c,inputColorDisabled:h,cardColor:g,modalColor:w,popoverColor:k,textColorDisabled:p,borderColor:t,primaryColor:v,textColor2:r,fontSizeSmall:R,fontSizeMedium:d,fontSizeLarge:a,borderRadiusSmall:m,lineHeight:f}=o;return Object.assign(Object.assign({},Re),{labelLineHeight:f,fontSizeSmall:R,fontSizeMedium:d,fontSizeLarge:a,borderRadius:m,color:c,colorChecked:v,colorDisabled:h,colorDisabledChecked:h,colorTableHeader:g,colorTableHeaderModal:w,colorTableHeaderPopover:k,checkMarkColor:c,checkMarkColorDisabled:p,checkMarkColorDisabledChecked:p,border:`1px solid ${t}`,borderDisabled:`1px solid ${t}`,borderDisabledChecked:`1px solid ${t}`,borderChecked:`1px solid ${v}`,borderFocus:`1px solid ${v}`,boxShadowFocus:`0 0 0 2px ${be(v,{alpha:.3})}`,textColor:r,textColorDisabled:p})}const Te={name:"Checkbox",common:se,self:Se},O=he("n-checkbox-group"),De={min:Number,max:Number,size:String,value:Array,defaultValue:{type:Array,default:null},disabled:{type:Boolean,default:void 0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onChange:[Function,Array]},Ue=H({name:"CheckboxGroup",props:De,setup(o){const{mergedClsPrefixRef:c}=N(o),h=E(o),{mergedSizeRef:g,mergedDisabledRef:w}=h,k=P(o.defaultValue),p=A(()=>o.value),t=L(p,k),v=A(()=>{var d;return((d=t.value)===null||d===void 0?void 0:d.length)||0}),r=A(()=>Array.isArray(t.value)?new Set(t.value):new Set);function R(d,a){const{nTriggerFormInput:m,nTriggerFormChange:f}=h,{onChange:l,"onUpdate:value":C,onUpdateValue:y}=o;if(Array.isArray(t.value)){const s=Array.from(t.value),_=s.findIndex(B=>B===a);d?~_||(s.push(a),y&&i(y,s,{actionType:"check",value:a}),C&&i(C,s,{actionType:"check",value:a}),m(),f(),k.value=s,l&&i(l,s)):~_&&(s.splice(_,1),y&&i(y,s,{actionType:"uncheck",value:a}),C&&i(C,s,{actionType:"uncheck",value:a}),l&&i(l,s),k.value=s,m(),f())}else d?(y&&i(y,[a],{actionType:"check",value:a}),C&&i(C,[a],{actionType:"check",value:a}),l&&i(l,[a]),k.value=[a],m(),f()):(y&&i(y,[],{actionType:"uncheck",value:a}),C&&i(C,[],{actionType:"uncheck",value:a}),l&&i(l,[]),k.value=[],m(),f())}return ue(O,{checkedCountRef:v,maxRef:U(o,"max"),minRef:U(o,"min"),valueSetRef:r,disabledRef:w,mergedSizeRef:g,toggleCheckbox:R}),{mergedClsPrefix:c}},render(){return u("div",{class:`${this.mergedClsPrefix}-checkbox-group`,role:"group"},this.$slots)}}),$e=()=>u("svg",{viewBox:"0 0 64 64",class:"check-icon"},u("path",{d:"M50.42,16.76L22.34,39.45l-8.1-11.46c-1.12-1.58-3.3-1.96-4.88-0.84c-1.58,1.12-1.95,3.3-0.84,4.88l10.26,14.51  c0.56,0.79,1.42,1.31,2.38,1.45c0.16,0.02,0.32,0.03,0.48,0.03c0.8,0,1.57-0.27,2.2-0.78l30.99-25.03c1.5-1.21,1.74-3.42,0.52-4.92  C54.13,15.78,51.93,15.55,50.42,16.76z"})),Me=()=>u("svg",{viewBox:"0 0 100 100",class:"line-icon"},u("path",{d:"M80.2,55.5H21.4c-2.8,0-5.1-2.5-5.1-5.5l0,0c0-3,2.3-5.5,5.1-5.5h58.7c2.8,0,5.1,2.5,5.1,5.5l0,0C85.2,53.1,82.9,55.5,80.2,55.5z"})),Fe=x([n("checkbox",`
 font-size: var(--n-font-size);
 outline: none;
 cursor: pointer;
 display: inline-flex;
 flex-wrap: nowrap;
 align-items: flex-start;
 word-break: break-word;
 line-height: var(--n-size);
 --n-merged-color-table: var(--n-color-table);
 `,[D("show-label","line-height: var(--n-label-line-height);"),x("&:hover",[n("checkbox-box",[S("border","border: var(--n-border-checked);")])]),x("&:focus:not(:active)",[n("checkbox-box",[S("border",`
 border: var(--n-border-focus);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),D("inside-table",[n("checkbox-box",`
 background-color: var(--n-merged-color-table);
 `)]),D("checked",[n("checkbox-box",`
 background-color: var(--n-color-checked);
 `,[n("checkbox-icon",[x(".check-icon",`
 opacity: 1;
 transform: scale(1);
 `)])])]),D("indeterminate",[n("checkbox-box",[n("checkbox-icon",[x(".check-icon",`
 opacity: 0;
 transform: scale(.5);
 `),x(".line-icon",`
 opacity: 1;
 transform: scale(1);
 `)])])]),D("checked, indeterminate",[x("&:focus:not(:active)",[n("checkbox-box",[S("border",`
 border: var(--n-border-checked);
 box-shadow: var(--n-box-shadow-focus);
 `)])]),n("checkbox-box",`
 background-color: var(--n-color-checked);
 border-left: 0;
 border-top: 0;
 `,[S("border",{border:"var(--n-border-checked)"})])]),D("disabled",{cursor:"not-allowed"},[D("checked",[n("checkbox-box",`
 background-color: var(--n-color-disabled-checked);
 `,[S("border",{border:"var(--n-border-disabled-checked)"}),n("checkbox-icon",[x(".check-icon, .line-icon",{fill:"var(--n-check-mark-color-disabled-checked)"})])])]),n("checkbox-box",`
 background-color: var(--n-color-disabled);
 `,[S("border",`
 border: var(--n-border-disabled);
 `),n("checkbox-icon",[x(".check-icon, .line-icon",`
 fill: var(--n-check-mark-color-disabled);
 `)])]),S("label",`
 color: var(--n-text-color-disabled);
 `)]),n("checkbox-box-wrapper",`
 position: relative;
 width: var(--n-size);
 flex-shrink: 0;
 flex-grow: 0;
 user-select: none;
 -webkit-user-select: none;
 `),n("checkbox-box",`
 position: absolute;
 left: 0;
 top: 50%;
 transform: translateY(-50%);
 height: var(--n-size);
 width: var(--n-size);
 display: inline-block;
 box-sizing: border-box;
 border-radius: var(--n-border-radius);
 background-color: var(--n-color);
 transition: background-color 0.3s var(--n-bezier);
 `,[S("border",`
 transition:
 border-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border: var(--n-border);
 `),n("checkbox-icon",`
 display: flex;
 align-items: center;
 justify-content: center;
 position: absolute;
 left: 1px;
 right: 1px;
 top: 1px;
 bottom: 1px;
 `,[x(".check-icon, .line-icon",`
 width: 100%;
 fill: var(--n-check-mark-color);
 opacity: 0;
 transform: scale(0.5);
 transform-origin: center;
 transition:
 fill 0.3s var(--n-bezier),
 transform 0.3s var(--n-bezier),
 opacity 0.3s var(--n-bezier),
 border-color 0.3s var(--n-bezier);
 `),Ce({left:"1px",top:"1px"})])]),S("label",`
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 padding: var(--n-label-padding);
 font-weight: var(--n-label-font-weight);
 `,[x("&:empty",{display:"none"})])]),fe(n("checkbox",`
 --n-merged-color-table: var(--n-color-table-modal);
 `)),ke(n("checkbox",`
 --n-merged-color-table: var(--n-color-table-popover);
 `))]),_e=Object.assign(Object.assign({},K.props),{size:String,checked:{type:[Boolean,String,Number],default:void 0},defaultChecked:{type:[Boolean,String,Number],default:!1},value:[String,Number],disabled:{type:Boolean,default:void 0},indeterminate:Boolean,label:String,focusable:{type:Boolean,default:!0},checkedValue:{type:[Boolean,String,Number],default:!0},uncheckedValue:{type:[Boolean,String,Number],default:!1},"onUpdate:checked":[Function,Array],onUpdateChecked:[Function,Array],privateInsideTable:Boolean,onChange:[Function,Array]}),je=H({name:"Checkbox",props:_e,setup(o){const c=ve(O,null),h=P(null),{mergedClsPrefixRef:g,inlineThemeDisabled:w,mergedRtlRef:k}=N(o),p=P(o.defaultChecked),t=U(o,"checked"),v=L(t,p),r=ze(()=>{if(c){const e=c.valueSetRef.value;return e&&o.value!==void 0?e.has(o.value):!1}else return v.value===o.checkedValue}),R=E(o,{mergedSize(e){const{size:z}=o;if(z!==void 0)return z;if(c){const{value:b}=c.mergedSizeRef;if(b!==void 0)return b}if(e){const{mergedSize:b}=e;if(b!==void 0)return b.value}return"medium"},mergedDisabled(e){const{disabled:z}=o;if(z!==void 0)return z;if(c){if(c.disabledRef.value)return!0;const{maxRef:{value:b},checkedCountRef:T}=c;if(b!==void 0&&T.value>=b&&!r.value)return!0;const{minRef:{value:M}}=c;if(M!==void 0&&T.value<=M&&r.value)return!0}return e?e.disabled.value:!1}}),{mergedDisabledRef:d,mergedSizeRef:a}=R,m=K("Checkbox","-checkbox",Fe,Te,o,g);function f(e){if(c&&o.value!==void 0)c.toggleCheckbox(!r.value,o.value);else{const{onChange:z,"onUpdate:checked":b,onUpdateChecked:T}=o,{nTriggerFormInput:M,nTriggerFormChange:I}=R,F=r.value?o.uncheckedValue:o.checkedValue;b&&i(b,F,e),T&&i(T,F,e),z&&i(z,F,e),M(),I(),p.value=F}}function l(e){d.value||f(e)}function C(e){if(!d.value)switch(e.key){case" ":case"Enter":f(e)}}function y(e){switch(e.key){case" ":e.preventDefault()}}const s={focus:()=>{var e;(e=h.value)===null||e===void 0||e.focus()},blur:()=>{var e;(e=h.value)===null||e===void 0||e.blur()}},_=me("Checkbox",k,g),B=A(()=>{const{value:e}=a,{common:{cubicBezierEaseInOut:z},self:{borderRadius:b,color:T,colorChecked:M,colorDisabled:I,colorTableHeader:F,colorTableHeaderModal:V,colorTableHeaderPopover:G,checkMarkColor:W,checkMarkColorDisabled:q,border:Y,borderFocus:J,borderDisabled:Q,borderChecked:X,boxShadowFocus:Z,textColor:ee,textColorDisabled:oe,checkMarkColorDisabledChecked:re,colorDisabledChecked:ae,borderDisabledChecked:ne,labelPadding:ce,labelLineHeight:le,labelFontWeight:ie,[j("fontSize",e)]:te,[j("size",e)]:de}}=m.value;return{"--n-label-line-height":le,"--n-label-font-weight":ie,"--n-size":de,"--n-bezier":z,"--n-border-radius":b,"--n-border":Y,"--n-border-checked":X,"--n-border-focus":J,"--n-border-disabled":Q,"--n-border-disabled-checked":ne,"--n-box-shadow-focus":Z,"--n-color":T,"--n-color-checked":M,"--n-color-table":F,"--n-color-table-modal":V,"--n-color-table-popover":G,"--n-color-disabled":I,"--n-color-disabled-checked":ae,"--n-text-color":ee,"--n-text-color-disabled":oe,"--n-check-mark-color":W,"--n-check-mark-color-disabled":q,"--n-check-mark-color-disabled-checked":re,"--n-font-size":te,"--n-label-padding":ce}}),$=w?xe("checkbox",A(()=>a.value[0]),B,o):void 0;return Object.assign(R,s,{rtlEnabled:_,selfRef:h,mergedClsPrefix:g,mergedDisabled:d,renderedChecked:r,mergedTheme:m,labelId:ge(),handleClick:l,handleKeyUp:C,handleKeyDown:y,cssVars:w?void 0:B,themeClass:$==null?void 0:$.themeClass,onRender:$==null?void 0:$.onRender})},render(){var o;const{$slots:c,renderedChecked:h,mergedDisabled:g,indeterminate:w,privateInsideTable:k,cssVars:p,labelId:t,label:v,mergedClsPrefix:r,focusable:R,handleKeyUp:d,handleKeyDown:a,handleClick:m}=this;(o=this.onRender)===null||o===void 0||o.call(this);const f=we(c.default,l=>v||l?u("span",{class:`${r}-checkbox__label`,id:t},v||l):null);return u("div",{ref:"selfRef",class:[`${r}-checkbox`,this.themeClass,this.rtlEnabled&&`${r}-checkbox--rtl`,h&&`${r}-checkbox--checked`,g&&`${r}-checkbox--disabled`,w&&`${r}-checkbox--indeterminate`,k&&`${r}-checkbox--inside-table`,f&&`${r}-checkbox--show-label`],tabindex:g||!R?void 0:0,role:"checkbox","aria-checked":w?"mixed":h,"aria-labelledby":t,style:p,onKeyup:d,onKeydown:a,onClick:m,onMousedown:()=>{pe("selectstart",window,l=>{l.preventDefault()},{once:!0})}},u("div",{class:`${r}-checkbox-box-wrapper`}," ",u("div",{class:`${r}-checkbox-box`},u(ye,null,{default:()=>this.indeterminate?u("div",{key:"indeterminate",class:`${r}-checkbox-icon`},Me()):u("div",{key:"check",class:`${r}-checkbox-icon`},$e())}),u("div",{class:`${r}-checkbox-box__border`}))),f)}});export{Ue as N,je as _,Te as c};
