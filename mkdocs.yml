site_name: PhoneInfoga
site_url: https://sundowndev.github.io/phoneinfoga
repo_name: 'sundowndev/phoneinfoga'
repo_url: 'https://github.com/sundowndev/phoneinfoga'
site_description: 'Advanced information gathering & OSINT tool for phone numbers.'
site_author: 'Sundowndev'
copyright: 'PhoneInfoga was developed by sundowndev and is licensed under GPL-3.0.'
nav:
  - 'Home': index.md
  - 'Getting Started':
    - 'Installation': getting-started/install.md
    - 'Usage': getting-started/usage.md
    - 'Scanners': getting-started/scanners.md
    - 'Go module usage': getting-started/go-module-usage.md
  - 'Resources':
      - 'Formatting phone numbers': resources/formatting.md
      - 'Additional resources': resources/additional-resources.md
  - 'Contribute': contribute.md
theme:
  name: material
  logo: './images/logo_white.svg'
  favicon: './images/logo.svg'
  palette:
    - scheme: default
      primary: blue grey
      accent: indigo
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    - scheme: slate
      primary: blue grey
      accent: indigo
      toggle:
        icon: material/brightness-4
        name: Switch to light mode
  features:
    - content.tabs.link
    - navigation.instant
    - navigation.sections
    - navigation.tabs
extra:
  social:
    - icon: fontawesome/brands/github-alt
      link: 'https://github.com/sundowndev/phoneinfoga'

# Extensions
markdown_extensions:
  - markdown.extensions.admonition
  - pymdownx.superfences
  - pymdownx.tabbed:
      alternate_style: true
  - attr_list
  - admonition
  - pymdownx.details
plugins:
  - search
  - minify:
      minify_html: true
  - redirects:
      redirect_maps:
        'install.md': 'getting-started/install.md'
        'usage.md': 'getting-started/usage.md'
        'scanners.md': 'getting-started/scanners.md'
        'go-module-usage.md': 'getting-started/go-module-usage.md'
        'formatting.md': 'resources/formatting.md'
        'resources.md': 'resources/additional-resources.md'
