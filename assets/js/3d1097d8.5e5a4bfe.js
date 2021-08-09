(self.webpackChunkprofiles_dev=self.webpackChunkprofiles_dev||[]).push([[9584],{3905:function(e,r,t){"use strict";t.d(r,{Zo:function(){return p},kt:function(){return d}});var n=t(7294);function o(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function i(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function a(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?i(Object(t),!0).forEach((function(r){o(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function l(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},i=Object.keys(e);for(n=0;n<i.length;n++)t=i[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)t=i[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var s=n.createContext({}),c=function(e){var r=n.useContext(s),t=r;return e&&(t="function"==typeof e?e(r):a(a({},r),e)),t},p=function(e){var r=c(e.components);return n.createElement(s.Provider,{value:r},e.children)},u={inlineCode:"code",wrapper:function(e){var r=e.children;return n.createElement(n.Fragment,{},r)}},f=n.forwardRef((function(e,r){var t=e.components,o=e.mdxType,i=e.originalType,s=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),f=c(t),d=o,g=f["".concat(s,".").concat(d)]||f[d]||u[d]||i;return t?n.createElement(g,a(a({ref:r},p),{},{components:t})):n.createElement(g,a({ref:r},p))}));function d(e,r){var t=arguments,o=r&&r.mdxType;if("string"==typeof e||o){var i=t.length,a=new Array(i);a[0]=f;var l={};for(var s in r)hasOwnProperty.call(r,s)&&(l[s]=r[s]);l.originalType=e,l.mdxType="string"==typeof e?e:o,a[1]=l;for(var c=2;c<i;c++)a[c]=t[c];return n.createElement.apply(null,a)}return n.createElement.apply(null,t)}f.displayName="MDXCreateElement"},3645:function(e,r,t){"use strict";t.r(r),t.d(r,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return c},toc:function(){return p},default:function(){return f}});var n=t(2122),o=t(9756),i=(t(7294),t(3905)),a=["components"],l={sidebar_position:5},s="Upgrading profiles",c={unversionedId:"installer-docs/upgrading-profiles",id:"installer-docs/upgrading-profiles",isDocsHomePage:!1,title:"Upgrading profiles",description:"When newer versions of a profile are available you will be able to discover them by",source:"@site/docs/installer-docs/upgrading-profiles.md",sourceDirName:"installer-docs",slug:"/installer-docs/upgrading-profiles",permalink:"/docs/installer-docs/upgrading-profiles",editUrl:"https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/docs/installer-docs/upgrading-profiles.md",version:"current",sidebarPosition:5,frontMatter:{sidebar_position:5},sidebar:"tutorialSidebar",previous:{title:"Listing installed profiles",permalink:"/docs/installer-docs/listing-installed-profiles"},next:{title:"Removing profiles",permalink:"/docs/installer-docs/removing-profiles"}},p=[],u={toc:p};function f(e){var r=e.components,t=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,n.Z)({},u,t,{components:r,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"upgrading-profiles"},"Upgrading profiles"),(0,i.kt)("p",null,"When newer versions of a profile are available you will be able to discover them by\n",(0,i.kt)("a",{parentName:"p",href:"/docs/installer-docs/listing-installed-profiles#listing-installed-profiles"},"listing installed profiles"),".\nOnce you know which version you want to upgrade to, run the following:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# The first argument is the relative or aboslute path to the local installation directory\n#\xa0and the second argument is the version to upgrade to.\npctl upgrade ~/workspace/demo-profile/ v0.0.2\n")),(0,i.kt)("p",null,"This will then perform an upgrade of your local installation. You can also pass in the ",(0,i.kt)("inlineCode",{parentName:"p"},"--create-pr")," flag to automatically create a PR\n. Pctl uses a 3-way merge behind the scenes to perform the upgrade. If you have made local modifications to\nyour installation that conflict with changes in the upgrades you will get merge conflicts, and will have to manually resolve them."))}f.isMDXComponent=!0}}]);