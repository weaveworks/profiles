(self.webpackChunkprofiles_dev=self.webpackChunkprofiles_dev||[]).push([[8018],{3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return p},kt:function(){return m}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=r.createContext({}),u=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},p=function(e){var t=u(e.components);return r.createElement(s.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,s=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),d=u(n),m=a,h=d["".concat(s,".").concat(m)]||d[m]||c[m]||i;return n?r.createElement(h,o(o({ref:t},p),{},{components:n})):r.createElement(h,o({ref:t},p))}));function m(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,o=new Array(i);o[0]=d;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:a,o[1]=l;for(var u=2;u<i;u++)o[u]=n[u];return r.createElement.apply(null,o)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},9945:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return u},toc:function(){return p},default:function(){return d}});var r=n(2122),a=n(9756),i=(n(7294),n(3905)),o=["components"],l={sidebar_position:1},s="Environment setup",u={unversionedId:"tutorial-basics/setup",id:"tutorial-basics/setup",isDocsHomePage:!1,title:"Environment setup",description:"This tutorial assumes you have some knowledge of the concept of GitOps and are comfortable using",source:"@site/docs/tutorial-basics/setup.md",sourceDirName:"tutorial-basics",slug:"/tutorial-basics/setup",permalink:"/docs/tutorial-basics/setup",editUrl:"https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/docs/tutorial-basics/setup.md",version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"Introduction",permalink:"/docs/intro"},next:{title:"Write a profile",permalink:"/docs/tutorial-basics/create-a-profile"}},p=[{value:"Prerequisites",id:"prerequisites",children:[{value:"Kubernetes cluster",id:"kubernetes-cluster",children:[]},{value:"Profiles CLI",id:"profiles-cli",children:[]},{value:"Profiles CRDs and Flux CRDs",id:"profiles-crds-and-flux-crds",children:[]},{value:"A GitHub repo, synced to Flux",id:"a-github-repo-synced-to-flux",children:[]},{value:"Personal Access Token",id:"personal-access-token",children:[]}]},{value:"Get started!",id:"get-started",children:[]}],c={toc:p};function d(e){var t=e.components,n=(0,a.Z)(e,o);return(0,i.kt)("wrapper",(0,r.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"environment-setup"},"Environment setup"),(0,i.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,i.kt)("div",{parentName:"div",className:"admonition-heading"},(0,i.kt)("h5",{parentName:"div"},(0,i.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,i.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,i.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"Assumed knowledge")),(0,i.kt)("div",{parentName:"div",className:"admonition-content"},(0,i.kt)("p",{parentName:"div"},"This tutorial assumes you have some knowledge of the concept of GitOps and are comfortable using\n",(0,i.kt)("a",{parentName:"p",href:"https://fluxcd.io/"},"Flux"),"."),(0,i.kt)("p",{parentName:"div"},"Please refer to the ",(0,i.kt)("a",{parentName:"p",href:"/docs/intro"},"Introduction")," to read about the core concepts of Profiles."))),(0,i.kt)("p",null,"In this tutorial you will create and install a simple profile onto your Kubernetes cluster using various GitOps tools."),(0,i.kt)("p",null,(0,i.kt)("em",{parentName:"p"},"If you are only interested in ",(0,i.kt)("strong",{parentName:"em"},"installing")," profiles, not writing them, please skip ahead to the relevant section\nonce you have set up your environment.")),(0,i.kt)("p",null,(0,i.kt)("em",{parentName:"p"},"If you are only interested in ",(0,i.kt)("strong",{parentName:"em"},"writing")," profiles, not installing them, you can skip the environment\nsetup steps.")),(0,i.kt)("hr",null),(0,i.kt)("h2",{id:"prerequisites"},"Prerequisites"),(0,i.kt)("p",null,"In order to install profiles, you need to have the following set up:"),(0,i.kt)("h3",{id:"kubernetes-cluster"},"Kubernetes cluster"),(0,i.kt)("p",null,"For local testing, we recommend using ",(0,i.kt)("a",{parentName:"p",href:"https://kind.sigs.k8s.io/docs/user/quick-start/"},"kind"),".\nThe cluster must be version 1.16 or newer."),(0,i.kt)("h3",{id:"profiles-cli"},"Profiles CLI"),(0,i.kt)("p",null,"Profiles are installed and managed via the official CLI: ",(0,i.kt)("inlineCode",{parentName:"p"},"pctl"),".\nReleases can be found ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/weaveworks/pctl/releases"},"here"),".\n",(0,i.kt)("inlineCode",{parentName:"p"},"pctl")," binaries are not backwards compatible, and we recommended keeping your local\nversion regularly updated."),(0,i.kt)("h3",{id:"profiles-crds-and-flux-crds"},"Profiles CRDs and Flux CRDs"),(0,i.kt)("p",null,"Profiles relies on Flux to deploy artifacts to your cluster, this means that at a minimum\nyou much have the following Flux CRDs and associated controllers installed:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"helmreleases.helm.toolkit.fluxcd.io")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"gitrepositories.source.toolkit.fluxcd.io")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"helmrepositories.source.toolkit.fluxcd.io")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"kustomizations.kustomize.toolkit.fluxcd.io"))),(0,i.kt)("p",null,"You can install everything by running Flux's ",(0,i.kt)("a",{parentName:"p",href:"https://fluxcd.io/docs/cmd/flux_install/"},"install command"),":"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"flux install\n")),(0,i.kt)("p",null,"Or to install specific components:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},'flux install --components="source-controller,kustomize-controller,helm-controller"\n')),(0,i.kt)("p",null,"Next install the Profiles CRD, with:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"pctl install\n")),(0,i.kt)("p",null,"Note: This will install the latest version of the profiles CRD, which may not always be stable."),(0,i.kt)("p",null,"To specify a ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/weaveworks/profiles/releases"},"specific version"),", use the ",(0,i.kt)("inlineCode",{parentName:"p"},"--version")," flag."),(0,i.kt)("h3",{id:"a-github-repo-synced-to-flux"},"A GitHub repo, synced to Flux"),(0,i.kt)("p",null,"This tutorial will require a GitHub account. (More git providers will be added in the future.)"),(0,i.kt)("p",null,"The repo can be public or private (note: you will not be asked to push any sensitive information) and must\nbe linked to the Flux instance running in your cluster."),(0,i.kt)("p",null,"You can do this by running ",(0,i.kt)("a",{parentName:"p",href:"https://fluxcd.io/docs/installation/#github-and-github-enterprise"},(0,i.kt)("inlineCode",{parentName:"a"},"flux bootstrap github"))," with the relevant arguments."),(0,i.kt)("div",{className:"admonition admonition-caution alert alert--warning"},(0,i.kt)("div",{parentName:"div",className:"admonition-heading"},(0,i.kt)("h5",{parentName:"div"},(0,i.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,i.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"16",height:"16",viewBox:"0 0 16 16"},(0,i.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M8.893 1.5c-.183-.31-.52-.5-.887-.5s-.703.19-.886.5L.138 13.499a.98.98 0 0 0 0 1.001c.193.31.53.501.886.501h13.964c.367 0 .704-.19.877-.5a1.03 1.03 0 0 0 .01-1.002L8.893 1.5zm.133 11.497H6.987v-2.003h2.039v2.003zm0-3.004H6.987V5.987h2.039v4.006z"}))),"Private repos")),(0,i.kt)("div",{parentName:"div",className:"admonition-content"},(0,i.kt)("p",{parentName:"div"},"If you choose to use a private repo, please ensure that your local git environment is set\nup correctly for the rest of the tutorial."))),(0,i.kt)("h3",{id:"personal-access-token"},"Personal Access Token"),(0,i.kt)("p",null,"The profile will be installed in a GitOps way, therefore ",(0,i.kt)("inlineCode",{parentName:"p"},"pctl")," will push all manifests to your cluster git repo.\nCreate a ",(0,i.kt)("a",{parentName:"p",href:"https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line"},"personal access token")," (check all permissions under repo)\non your GitHub account and export it:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"export GIT_TOKEN=<your token>\n")),(0,i.kt)("h2",{id:"get-started"},"Get started!"),(0,i.kt)("p",null,"Check you have everything on this list and go back if something is missing."),(0,i.kt)("p",null," \u2705 ","\xa0","\xa0"," ",(0,i.kt)("a",{parentName:"p",href:"#kubernetes-cluster"},"Cluster")),(0,i.kt)("p",null," \u2705 ","\xa0","\xa0"," ",(0,i.kt)("a",{parentName:"p",href:"#pctl-the-profiles-cli"},"Pctl binary")),(0,i.kt)("p",null," \u2705 ","\xa0","\xa0"," ",(0,i.kt)("a",{parentName:"p",href:"#profiles-crds-and-flux-crds"},"Profiles CRDs and Flux CRDs")),(0,i.kt)("p",null," \u2705 ","\xa0","\xa0"," ",(0,i.kt)("a",{parentName:"p",href:"#a-github-repo-synced-to-flux"},"GitHub repo")),(0,i.kt)("p",null," \u2705 ","\xa0","\xa0"," ",(0,i.kt)("a",{parentName:"p",href:"#personal-access-token"},"GitHub token")),(0,i.kt)("p",null,"Once you have completed the prerequisites installation you can start writing a profile!"))}d.isMDXComponent=!0}}]);