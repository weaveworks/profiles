(self.webpackChunkprofiles_dev=self.webpackChunkprofiles_dev||[]).push([[5622],{3905:function(e,t,r){"use strict";r.d(t,{Zo:function(){return u},kt:function(){return m}});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var s=n.createContext({}),p=function(e){var t=n.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},u=function(e){var t=p(e.components);return n.createElement(s.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},f=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,a=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),f=p(r),m=o,y=f["".concat(s,".").concat(m)]||f[m]||c[m]||a;return r?n.createElement(y,i(i({ref:t},u),{},{components:r})):n.createElement(y,i({ref:t},u))}));function m(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var a=r.length,i=new Array(a);i[0]=f;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:o,i[1]=l;for(var p=2;p<a;p++)i[p]=r[p];return n.createElement.apply(null,i)}return n.createElement.apply(null,r)}f.displayName="MDXCreateElement"},1524:function(e,t,r){"use strict";r.r(t),r.d(t,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return p},assets:function(){return u},toc:function(){return c},default:function(){return m}});var n=r(2122),o=r(9756),a=(r(7294),r(3905)),i=["components"],l={slug:"mlops-profile",title:"Install MLOps Profile with Weave GitOps",author:"Chanwit Kaewkasi",author_title:"DX Engineer, WeaveWorks",tags:["mlops","profile","gitops"]},s=void 0,p={permalink:"/blog/mlops-profile",editUrl:"https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/blog/blog/2021-08-04-MLOps-Profile-with-Weave-GitOps.md",source:"@site/blog/2021-08-04-MLOps-Profile-with-Weave-GitOps.md",title:"Install MLOps Profile with Weave GitOps",description:"\u201cProfiles\u201d is a GitOps-native package management system, which allows you to apply a Profile to add a set of capabilities to your GitOps-enabled cluster.  Here\u2019s a quick tutorial on how to provision a cluster, install Weave GitOps and apply the MLOps profile to enable Kubeflow Pipeline for your cluster. To try this setup, we strongly recommend that you use Kind as a local cluster on your laptop.",date:"2021-08-04T00:00:00.000Z",formattedDate:"August 4, 2021",tags:[{label:"mlops",permalink:"/blog/tags/mlops"},{label:"profile",permalink:"/blog/tags/profile"},{label:"gitops",permalink:"/blog/tags/gitops"}],readingTime:1.815,truncated:!0,authors:[{name:"Chanwit Kaewkasi",title:"DX Engineer, WeaveWorks"}]},u={authorsImageUrls:[void 0]},c=[],f={toc:c};function m(e){var t=e.components,r=(0,o.Z)(e,i);return(0,a.kt)("wrapper",(0,n.Z)({},f,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("p",null,"\u201c",(0,a.kt)("em",{parentName:"p"},"Profiles"),"\u201d is a GitOps-native package management system, which allows you to apply a Profile to add a set of capabilities to your GitOps-enabled cluster.  Here\u2019s a quick tutorial on how to provision a cluster, install Weave GitOps and apply the MLOps profile to enable Kubeflow Pipeline for your cluster. To try this setup, we strongly recommend that you use ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/kubernetes-sigs/kind"},"Kind")," as a local cluster on your laptop."))}m.isMDXComponent=!0}}]);