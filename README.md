# res2
A resouce manager for arbitrary content.

## Getting Started

- Create a YAML file called `res2.yaml` with a `files` property whose value is a mapping of file names to HTTP or HTTPS URLs.
- Run `res2` from within the directory with the `res2.yaml` file.

## Example Manifest

```yaml
files:
    .editorconfig: https://raw.githubusercontent.com/wareismymind/peer/main/.editorconfig
    .megalinter: https://raw.githubusercontent.com/wareismymind/peer/main/.mega-linter.yml
```
