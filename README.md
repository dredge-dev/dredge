# Dredge - Automates software DevOps workflows

Dredge is a open-source tool that seamlessly integrates with your DevOps tools to streamline and standardize your development and operations workflows, helping your team to work more efficiently and effectively.

[![cast](https://asciinema.org/a/564048.png)](https://asciinema.org/a/564048)

## Installation

Install the `drg` command line tool and initialize the project.

```bash
curl https://dredge.dev/install.sh | bash
drg init
```

## Dredgefile

A Dredgefile contains resources and workflows. See [the docs](https://dredge.dev/docs/dredgefile/) for more information.

```
resources:
  release:
    - provider: github-releases
workflows:
- name: hello
  description: Say hello
  steps:
  - shell:
      cmd: echo hello
```

## License

Dredge is licensed under the Apache 2.0 License.
