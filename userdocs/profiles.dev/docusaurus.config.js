const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'Profiles',
  tagline: 'GitOps native package manager',
  url: 'https://docs.profiles.dev/',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon_150px.png',
  organizationName: 'weaveworks', // Usually your GitHub org/user name.
  projectName: 'profiles', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Profiles',
      logo: {
        alt: 'something cute coming soon',
        src: 'img/weave-logo.png',
      },
      items: [
        {
          to: '/docs/tutorial-basics/setup',
          position: 'left',
          label: 'Getting started',
        },
        {to: '/docs/intro', label: 'Docs', position: 'left'},
        {to: '/docs/installer-docs/installing-via-gitops', label: 'Installing Profiles', position: 'left'},
        {to: '/blog', label: 'Blog', position: 'left'},
        {
          href: 'https://github.com/weaveworks/profiles',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/tutorial-basics/setup',
            },
            {
              label: 'Developer Docs: Profile Author',
              to: '/docs/author-docs/profile-structure',
            },
            {
              label: 'Developer Docs: Profile User',
              to: '/docs/installer-docs/installing-via-gitops',
            },
            {
              label: 'Developer Docs: Catalog Manager',
              to: '/docs/catalog-docs/add-profiles',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Slack',
              href: 'https://slack.weave.works/',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/weaveworks',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'FAQ',
              to: '/docs/faq',
            },
            {
              label: 'Contributing',
              href: '/docs/contributing',
            },
            {
              label: 'Report a problem with Profiles',
              href: 'https://github.com/weaveworks/profiles/issues',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Weaveworks, Inc. Built with Docusaurus.`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl:
            'https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/weaveworks/profiles/edit/main/userdocs/profiles.dev/blog',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
