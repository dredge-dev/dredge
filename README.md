# Dredge - Toil less, code more

Dredge automates software development workflows.

## Dredgefile

Define a Dredgefile to put your development workflows into code:

```
workflows:
- name: hello
  description: Say hello
  steps:
  - shell:
      cmd: echo hello
```

## Installation

### Linux

```
curl -L https://github.com/dredge-dev/dredge/releases/latest/download/drg-linux-amd64 > drg
chmod +x drg
mv drg /usr/local/bin/
```

### macOS

```
curl -L https://github.com/dredge-dev/dredge/releases/latest/download/drg-darwin-amd64 > drg
chmod +x drg
mv drg /usr/local/bin/
```

### Windows

```
Invoke-WebRequest -Uri https://github.com/dredge-dev/dredge/releases/latest/download/drg-windows-amd64 -OutFile drg
```

### ARM

ARM binaries are available in [the GitHub releases](https://github.com/dredge-dev/dredge/releases).

## Try it out

Create a Dredgefile as shown above and run Dredge:

```
drg
```

## License

Dredge is licensed under the Apache 2.0 License.
