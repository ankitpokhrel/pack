### Pack

A small utility to create archives while ignoring any hidden files/folders. You can additionally pass a config file with the list of patterns to ignore.

#### Install

Install the runnable binary to your `$GOPATH/bin`.

```sh
$ go install github.com/ankitpokhrel/pack
```

#### Usage

```sh
NAME:
   pack - Pack compresses file/folder ignoring any hidden and given files

USAGE:
   pack [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ignore value, --ig value  Get ignore list from given files
   --help, -h                  show help (default: false)
```

##### Example

Given a list `.myignorelist`
```sh
$ cat .myignorelist

vendor/
file.txt
*.png
```

The following command will create `destination.zip` file by ignoring all patterns mentioned in `.myignorelist`.
```sh
$ pack -ig .myignorelist /path/to/file-to-compress /part/to/destination.zip
```
