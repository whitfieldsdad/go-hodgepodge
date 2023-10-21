# hodgepodge

This repository contains a hodgepodge of Go code that is used in various projects.

> 👷 🚧: this project is experimental, doesn't have a stable API, and is under active development.

## Features

The following features are available:

<img src="./docs/diagrams/features.png" width="300"/>

| Name                          | Functional area | Description                                                     | Notes                                                                                                      |
| ----------------------------- | --------------- | --------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| Collect environment variables | Data collection | Collect environment variables as sequences of key/value pairs.  |                                                                                                            |
| Collect file metadata         | Data collection | Collect file metadata (e.g. size, traits, MACb timestamps)      |                                                                                                            |
| Collect system metadata       | Data collection |                                                                 |                                                                                                            |
| Collect system metrics        | Data collection | Measure disk usage.                                             |                                                                                                            |
| Discover network location     | Data collection | Hostname, primary network interface.                            |                                                                                                            |
| Encrypt/decrypt files         | Data collection | Encrypt/decrypt files using symmetric or asymmetric encryption. | RSA is used for asymmetric encryption. AES-256-GCM or ChaCha20-Poly1305 are used for symmetric encryption. |
| Hash files                    | Data collection | Hash files using a variety of common hash functions.            | Supported hash functions are MD5, SHA-1, SHA-256, SHA-512.                                                 |
| List disks                    | Data collection |                                                                 |                                                                                                            |
| List network connections      | Data collection |                                                                 |                                                                                                            |
| List network interfaces       | Data collection |                                                                 |                                                                                                            |
| List running processes        | Data collection |                                                                 |                                                                                                            |
| Search for files              | Data collection | Search for files by path or filename.                           |                                                                                                            |
| Execute shell commands        | Code execution  | Execute shell commands and capture their output.                | Supports executing commands in a sh, bash, Powershell, or Command Prompt session.                          |

## Examples

### Listing processes

```shell
go run main.go processes list | jq
```

```json
{
  "name": "docker",
  "working_directory": "/mnt/c/Users/tyler/AppData/Roaming/Docker",
  "pid": 151,
  "ppid": 150,
  "file": {
    "path": "/mnt/wsl/docker-desktop/cli-tools/usr/bin/docker",
    "filename": "docker",
    "directory": "/mnt/wsl/docker-desktop/cli-tools/usr/bin",
    "extension": "",
    "size": 56389632,
    "timestamps": {
      "modify_time": "2023-01-19T08:05:49-05:00",
      "access_time": "2023-01-19T08:05:49-05:00",
      "change_time": "2023-01-19T08:06:15-05:00",
      "birth_time": null
    },
    "traits": {
      "is_directory": false,
      "is_regular_file": true,
      "is_symbolic_link": false,
      "is_socket": false,
      "is_hard_link": false,
      "is_named_pipe": false,
      "is_block_device": false,
      "is_character_device": false,
      "is_hidden": false
    },
    "hashes": {
      "md5": "8b39ff078adbd2a209104b8163934d3c",
      "sha1": "a49f2f33afd484064637282d14f6d71b59d43a16",
      "sha256": "f1896376f8d504e4e5143d3695298d884bf0a4eb30395382b310bfd3af277951",
      "sha512": "887431ac3d37c932b31247f2aa6b973b024cb61e63eb23aa1077f7be0e0fb6075fe48225b7a45c0e66220854aeec52deabde9979fde640a7809f8f9df3747f9f",
      "xxh64": "a09722851a6e6353"
    }
  },
  "command_line": "docker serve --address unix:///root/.docker/run/docker-cli-api.sock",
  "argv": [
    "docker",
    "serve",
    "--address",
    "unix:///root/.docker/run/docker-cli-api.sock"
  ],
  "argc": 4,
  "start_time": "2023-10-16T19:17:52.26-04:00"
}
...
```
