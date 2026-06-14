import{g as a,i as d,an as e,f as o,h as s,j as f,k as p,l as g,aG as b,ao as m,ap as h}from"./DmHPR5lg.js";import{a as $}from"./DO8alW5h.js";const t="0!important",u="-1px!important";function i(r){return e(`${r}-type`,[o("& +",[a("button",{},[e(`${r}-type`,[s("border",{borderLeftWidth:t}),s("state-border",{left:u})])])])])}function n(r){return e(`${r}-type`,[o("& +",[a("button",[e(`${r}-type`,[s("border",{borderTopWidth:t}),s("state-border",{top:u})])])])])}const y=a("button-group",`
 flex-wrap: nowrap;
 display: inline-flex;
 position: relative;
`,[d("vertical",{flexDirection:"row"},[d("rtl",[a("button",[o("&:first-child:not(:last-child)",`
 margin-right: ${t};
 border-top-right-radius: ${t};
 border-bottom-right-radius: ${t};
 `),o("&:last-child:not(:first-child)",`
 margin-left: ${t};
 border-top-left-radius: ${t};
 border-bottom-left-radius: ${t};
 `),o("&:not(:first-child):not(:last-child)",`
 margin-left: ${t};
 margin-right: ${t};
 border-radius: ${t};
 `),i("default"),e("ghost",[i("primary"),i("info"),i("success"),i("warning"),i("error")])])])]),e("vertical",{flexDirection:"column"},[a("button",[o("&:first-child:not(:last-child)",`
 margin-bottom: ${t};
 margin-left: ${t};
 margin-right: ${t};
 border-bottom-left-radius: ${t};
 border-bottom-right-radius: ${t};
 `),o("&:last-child:not(:first-child)",`
 margin-top: ${t};
 margin-left: ${t};
 margin-right: ${t};
 border-top-left-radius: ${t};
 border-top-right-radius: ${t};
 `),o("&:not(:first-child):not(:last-child)",`
 margin: ${t};
 border-radius: ${t};
 `),n("default"),e("ghost",[n("primary"),n("info"),n("success"),n("warning"),n("error")])])])]),x={size:{type:String,default:void 0},vertical:Boolean},w=f({name:"ButtonGroup",props:x,setup(r){const{mergedClsPrefixRef:l,mergedRtlRef:c}=g(r);return b("-button-group",y,l),h($,r),{rtlEnabled:m("ButtonGroup",c,l),mergedClsPrefix:l}},render(){const{mergedClsPrefix:r}=this;return p("div",{class:[`${r}-button-group`,this.rtlEnabled&&`${r}-button-group--rtl`,this.vertical&&`${r}-button-group--vertical`],role:"group"},this.$slots)}});export{w as _};
