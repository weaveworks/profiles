(self.webpackChunkprofiles_dev=self.webpackChunkprofiles_dev||[]).push([[399],{3905:function(e,t,r){"use strict";r.d(t,{Zo:function(){return c},kt:function(){return d}});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var s=n.createContext({}),p=function(e){var t=n.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},c=function(e){var t=p(e.components);return n.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},f=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,a=e.originalType,s=e.parentName,c=l(e,["components","mdxType","originalType","parentName"]),f=p(r),d=o,m=f["".concat(s,".").concat(d)]||f[d]||u[d]||a;return r?n.createElement(m,i(i({ref:t},c),{},{components:r})):n.createElement(m,i({ref:t},c))}));function d(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var a=r.length,i=new Array(a);i[0]=f;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:o,i[1]=l;for(var p=2;p<a;p++)i[p]=r[p];return n.createElement.apply(null,i)}return n.createElement.apply(null,r)}f.displayName="MDXCreateElement"},42:function(e,t,r){"use strict";r.r(t),r.d(t,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return p},toc:function(){return c},default:function(){return f}});var n=r(2122),o=r(9756),a=(r(7294),r(3905)),i=["components"],l={sidebar_position:1},s="How a profile is structured",p={unversionedId:"author-docs/profile-structure",id:"author-docs/profile-structure",isDocsHomePage:!1,title:"How a profile is structured",description:"profile.yaml contents",source:"@site/docs/author-docs/profile-structure.md",sourceDirName:"author-docs",slug:"/author-docs/profile-structure",permalink:"/profiles/docs/author-docs/profile-structure",editUrl:"https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/docs/author-docs/profile-structure.md",version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"Install a profile",permalink:"/profiles/docs/tutorial-basics/install-a-profile"},next:{title:"Remote Helm Chart artifacts",permalink:"/profiles/docs/author-docs/remote-helm-chart"}},c=[{value:"<code>profile.yaml</code> contents",id:"profileyaml-contents",children:[]},{value:"Profile repo directories",id:"profile-repo-directories",children:[]}],u={toc:c};function f(e){var t=e.components,r=(0,o.Z)(e,i);return(0,a.kt)("wrapper",(0,n.Z)({},u,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"how-a-profile-is-structured"},"How a profile is structured"),(0,a.kt)("h2",{id:"profileyaml-contents"},(0,a.kt)("inlineCode",{parentName:"h2"},"profile.yaml")," contents"),(0,a.kt)("p",null,"A profile is defined in a single file which ",(0,a.kt)("strong",{parentName:"p"},"must")," be named ",(0,a.kt)("inlineCode",{parentName:"p"},"profile.yaml"),".\nThis file lives at the root of the profile directory."),(0,a.kt)("p",null,"The following fields are required:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-yaml"},"apiVersion: weave.works/v1alpha1\nkind: ProfileDefinition\nmetadata:\n  name: # the name of your profile\nspec:\n  description: # a brief description of what your profile installs\n")),(0,a.kt)("p",null,"These fields are optional:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-yaml"},"# ...\nspec:\n  # ...\n  maintainer: # the name(s) of the profile author\n  prerequisites:\n  - # a list of strings detailing things the profile needs to run.\n  - # this field is not processed at the moment, but will be soon.\n")),(0,a.kt)("p",null,"Finally, the ",(0,a.kt)("inlineCode",{parentName:"p"},"spec.artifacts")," lists all the components which the profile will install."),(0,a.kt)("p",null,"The following artifact types are supported:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/docs/author-docs/local-helm-chart"},"'Local' Helm Chart")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/docs/author-docs/remote-helm-chart"},"'Remote' Helm Chart")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/docs/author-docs/kustomize-and-raw"},"Raw Kubernetes yaml")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/docs/author-docs/kustomize-and-raw"},"Kustomize patch")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/docs/author-docs/nested-profiles"},"Another profile"))),(0,a.kt)("p",null,"Please refer to their dedicated docs pages for details on how to register different artifact\ntypes in a profile."),(0,a.kt)("h2",{id:"profile-repo-directories"},"Profile repo directories"),(0,a.kt)("p",null,"It will be assumed that everything contained within the same directory as a ",(0,a.kt)("inlineCode",{parentName:"p"},"profile.yaml"),"\nis related to that same profile."),(0,a.kt)("p",null,"A repository can contain multiple profiles when they are written in separate directories.\nFor example, the following structure shows a repo with three distinct profiles which\ncan be installed independently of each other:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-bash"},"org-profiles-repo/\n\u251c\u2500\u2500 logging-profile\n\u2502\xa0\xa0 \u2514\u2500\u2500 profile.yaml\n\u251c\u2500\u2500 observability-profile\n\u2502\xa0\xa0 \u2514\u2500\u2500 profile.yaml\n\u2514\u2500\u2500 our-awesome-apps\n    \u2514\u2500\u2500 profile.yaml\n")),(0,a.kt)("div",{className:"admonition admonition-tip alert alert--success"},(0,a.kt)("div",{parentName:"div",className:"admonition-heading"},(0,a.kt)("h5",{parentName:"div"},(0,a.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,a.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"12",height:"16",viewBox:"0 0 12 16"},(0,a.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M6.5 0C3.48 0 1 2.19 1 5c0 .92.55 2.25 1 3 1.34 2.25 1.78 2.78 2 4v1h5v-1c.22-1.22.66-1.75 2-4 .45-.75 1-2.08 1-3 0-2.81-2.48-5-5.5-5zm3.64 7.48c-.25.44-.47.8-.67 1.11-.86 1.41-1.25 2.06-1.45 3.23-.02.05-.02.11-.02.17H5c0-.06 0-.13-.02-.17-.2-1.17-.59-1.83-1.45-3.23-.2-.31-.42-.67-.67-1.11C2.44 6.78 2 5.65 2 5c0-2.2 2.02-4 4.5-4 1.22 0 2.36.42 3.22 1.19C10.55 2.94 11 3.94 11 5c0 .66-.44 1.78-.86 2.48zM4 14h5c-.23 1.14-1.3 2-2.5 2s-2.27-.86-2.5-2z"}))),"tip")),(0,a.kt)("div",{parentName:"div",className:"admonition-content"},(0,a.kt)("p",{parentName:"div"},"The name of each profile directory ",(0,a.kt)("strong",{parentName:"p"},"must")," match the name given in the ",(0,a.kt)("inlineCode",{parentName:"p"},"profile.yaml"),"\n",(0,a.kt)("inlineCode",{parentName:"p"},"metadata.name"),"."))),(0,a.kt)("p",null,"A repository can also contain just a single profile, with the ",(0,a.kt)("inlineCode",{parentName:"p"},"profile.yaml"),"\ndefined at the top level:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-bash"},"org-profiles-repo/\n\u2514\u2500\u2500 profile.yaml\n")),(0,a.kt)("p",null,"Profile directories can contain other objects related to various artifacts. These\nwill be demonstrated in subsequent pages."),(0,a.kt)("p",null,"Examples of profiles with various artifacts and configurations can be found ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/weaveworks/profiles-examples"},"here"),"."))}f.isMDXComponent=!0}}]);