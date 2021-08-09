(self.webpackChunkprofiles_dev=self.webpackChunkprofiles_dev||[]).push([[3171],{3905:function(e,t,r){"use strict";r.d(t,{Zo:function(){return u},kt:function(){return f}});var o=r(7294);function n(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);t&&(o=o.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,o)}return r}function l(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){n(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function i(e,t){if(null==e)return{};var r,o,n=function(e,t){if(null==e)return{};var r,o,n={},a=Object.keys(e);for(o=0;o<a.length;o++)r=a[o],t.indexOf(r)>=0||(n[r]=e[r]);return n}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(o=0;o<a.length;o++)r=a[o],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(n[r]=e[r])}return n}var p=o.createContext({}),s=function(e){var t=o.useContext(p),r=t;return e&&(r="function"==typeof e?e(t):l(l({},t),e)),r},u=function(e){var t=s(e.components);return o.createElement(p.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return o.createElement(o.Fragment,{},t)}},m=o.forwardRef((function(e,t){var r=e.components,n=e.mdxType,a=e.originalType,p=e.parentName,u=i(e,["components","mdxType","originalType","parentName"]),m=s(r),f=n,d=m["".concat(p,".").concat(f)]||m[f]||c[f]||a;return r?o.createElement(d,l(l({ref:t},u),{},{components:r})):o.createElement(d,l({ref:t},u))}));function f(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var a=r.length,l=new Array(a);l[0]=m;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i.mdxType="string"==typeof e?e:n,l[1]=i;for(var s=2;s<a;s++)l[s]=r[s];return o.createElement.apply(null,l)}return o.createElement.apply(null,r)}m.displayName="MDXCreateElement"},5904:function(e,t,r){"use strict";r.r(t),r.d(t,{frontMatter:function(){return i},contentTitle:function(){return p},metadata:function(){return s},toc:function(){return u},default:function(){return m}});var o=r(2122),n=r(9756),a=(r(7294),r(3905)),l=["components"],i={slug:"mlops-profile",title:"Install MLOps Profile with Weave GitOps",author:"Chanwit Kaewkasi",author_title:"DX Engineer, WeaveWorks",tags:["mlops","profile","gitops"]},p=void 0,s={permalink:"/blog/mlops-profile",editUrl:"https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/blog/blog/2021-08-04-MLOps-Profile-with-Weave-GitOps.md",source:"@site/blog/2021-08-04-MLOps-Profile-with-Weave-GitOps.md",title:"Install MLOps Profile with Weave GitOps",description:"\u201cProfiles\u201d is a GitOps-native package management system, which allows you to apply a Profile to add a set of capabilities to your GitOps-enabled cluster.  Here\u2019s a quick tutorial on how to provision a cluster, install Weave GitOps and apply the MLOps profile to enable Kubeflow Pipeline for your cluster. To try this setup, we strongly recommend that you use Kind as a local cluster on your laptop.",date:"2021-08-04T00:00:00.000Z",formattedDate:"August 4, 2021",tags:[{label:"mlops",permalink:"/blog/tags/mlops"},{label:"profile",permalink:"/blog/tags/profile"},{label:"gitops",permalink:"/blog/tags/gitops"}],readingTime:1.815,truncated:!0},u=[],c={toc:u};function m(e){var t=e.components,i=(0,n.Z)(e,l);return(0,a.kt)("wrapper",(0,o.Z)({},c,i,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("p",null,"\u201c",(0,a.kt)("em",{parentName:"p"},"Profiles"),"\u201d is a GitOps-native package management system, which allows you to apply a Profile to add a set of capabilities to your GitOps-enabled cluster.  Here\u2019s a quick tutorial on how to provision a cluster, install Weave GitOps and apply the MLOps profile to enable Kubeflow Pipeline for your cluster. To try this setup, we strongly recommend that you use ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/kubernetes-sigs/kind"},"Kind")," as a local cluster on your laptop."),(0,a.kt)("p",null,"First, you need Weave GitOps Core. Please download it from ",(0,a.kt)("a",{parentName:"p",href:"https://docs.gitops.weave.works/docs/installation/"},"here")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},'curl -L "https://github.com/weaveworks/weave-gitops/releases/download/v0.2.2/wego-$(uname)-$(uname -m)" -o wego\nchmod +x wego\nsudo mv ./wego /usr/local/bin/wego\nwego version\n')),(0,a.kt)("p",null,"Second, you also need Profiles CLI. You can download the CLI from ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/weaveworks/pctl/releases"},"here")),(0,a.kt)("p",null,"The following are links of the Profile CLI for both Linux and macOS."),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"https://github.com/weaveworks/pctl/releases/download/v0.2.0/pctl_linux_amd64.tar.gz"},"Profiles CLI for Linux")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"https://github.com/weaveworks/pctl/releases/download/v0.2.0/pctl_darwin_amd64.tar.gz"},"Profiles CLI for macOS"))),(0,a.kt)("h4",{id:"heres-steps-to-install-mlops-profile-with-weave-gitops"},"Here\u2019s steps to install MLOps Profile with Weave GitOps"),(0,a.kt)("ol",null,(0,a.kt)("li",{parentName:"ol"},"Create cluster")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"kind create cluster\n")),(0,a.kt)("ol",{start:2},(0,a.kt)("li",{parentName:"ol"},"Create a new GitOps repository on GitHub.")),(0,a.kt)("p",null,"For example, my GitOps repo is ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/chanwit/test-mlops-profile"},"test-mlops-profile")),(0,a.kt)("p",null,"Obtain your GitHub token, and export it as an environment variable along with your user and repo names. Please change ",(0,a.kt)("inlineCode",{parentName:"p"},"GITHUB_USER")," to your GitHub user name."),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"export GITHUB_TOKEN=gh0_1234567890\nexport GITHUB_USER=chanwit\nexport GITHUB_REPO=test-mlops-profile\n")),(0,a.kt)("ol",{start:3},(0,a.kt)("li",{parentName:"ol"},"Check that you already installed both Weave GitOps Core and Profiles CLI, then we install GitOps runtime to our cluster using the following command.")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"wego gitops install\n")),(0,a.kt)("ol",{start:5},(0,a.kt)("li",{parentName:"ol"},"After that, we install the ProfileManager controller using the pctl install command.\nPlease note that we need to tell it that our GitOps runtime is running in the wego-system namespace.")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"pctl install --flux-namespace=wego-system\n")),(0,a.kt)("ol",{start:6},(0,a.kt)("li",{parentName:"ol"},"Then, we clone the GitOps repo and setup our Weave GitOps automation for it.")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"git clone git@github.com:$GITHUB_USER/$GITHUB_REPO\ncd $GITHUB_REPO\nwego app add .\n")),(0,a.kt)("ol",{start:7},(0,a.kt)("li",{parentName:"ol"},(0,a.kt)("p",{parentName:"li"},"Weave GitOps will create a PR for you to merge. After you merge it, GitOps automation will kick off.")),(0,a.kt)("li",{parentName:"ol"},(0,a.kt)("p",{parentName:"li"},"Finally, we are able to install the MLOps profile using the following commands."))),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"kubectl create ns kubeflow\n\npctl add --profile-url=https://github.com/chanwit/mlops-profile --git-repository=wego-system/test-mlops-profile --name=mlops-profile --namespace=kubeflow\n\ngit add .\ngit commit -am \u201cinstall mlops profile\u201d\ngit push origin main\n")),(0,a.kt)("p",null,"Here\u2019s the pod list after you have the MLOps profile installed."),(0,a.kt)("p",null,(0,a.kt)("img",{alt:"MLOps profile",src:r(4768).Z})))}m.isMDXComponent=!0},4768:function(e,t,r){"use strict";t.Z=r.p+"assets/images/mlops-5118515c7e9d096b6859baf5ed4d23a6.png"}}]);