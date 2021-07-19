import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

const FeatureList = [
  {
    title: 'What are Profiles?',
    Svg: require('../../static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
      Profiles is a GitOps package management mechanism. 
      <p><a href="/docs/intro">Read more</a></p>
      </>
    ),
  },
  {
    title: 'Powered by Flux',
    Svg: require('../../static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        <a href="https://fluxcd.io">Flux</a> is a leading CNCF project around GitOps automation.
        Weave GitOps builds on this foundation to create a highly effective GitOps runtime.
      </>
    ),
  },
  {
    title: 'The Profiles CLI',
    Svg: require('../../static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
      Profiles are installed and managed via the official CLI <code>pctl</code>. 
      Releases can be found <a href="https://github.com/weaveworks/pctl/releases">here</a>.  
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} alt={title} />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
