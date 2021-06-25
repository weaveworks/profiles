import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

const FeatureList = [
  {
    title: 'GitOps Native package manager',
    Svg: require('../../static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
      some cool text here
      </>
    ),
  },
  {
    title: 'something else',
    Svg: require('../../static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
      more cool text
      </>
    ),
  },
  {
    title: 'hmmmm',
    Svg: require('../../static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
      something about flux maybe. don't forget to also change the logos
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
