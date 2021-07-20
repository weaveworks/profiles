import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

const FeatureList = [
  {
    title: 'What is Profiles?',
    description: (
      <>
      Profiles is a GitOps package management mechanism.
      <a href="/docs/intro">Read more</a>
      </>
    ),
  },
  {
    title: 'Powered by Flux',
    description: (
      <>
        <a href="https://fluxcd.io">Flux</a> is a leading CNCF project around GitOps automation.
        Weave GitOps builds on this foundation to create a highly effective GitOps runtime.
      </>
    ),
  },
  {
    title: 'Profiles CLI',
    description: (
      <>
      Profiles are installed and managed via the official CLI <code>pctl</code>. 
      Releases can be found <a href="https://github.com/weaveworks/pctl/releases">here</a>.  
      </>
    ),
  },
];

function Feature({title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md padding-vert--md">
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
