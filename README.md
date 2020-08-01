## Pack

A small utility to create archives while ignoring any hidden or unnecessary files and folders. It uses git style pattern matching.

### Installation

Install the runnable binary to your `$GOPATH/bin`.

```sh
$ go install github.com/ankitpokhrel/pack
```

### Usage

```sh
NAME:
   pack - Pack create archives while ignoring any hidden or unnecessary files and folders

USAGE:
   pack [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ignore value, --ig value  Ignore list from given files
   --help, -h                  show help (default: false)
```

### Example

Given a `.gitignore` and `.ignoremetoo` file:
```sh
$ cat .gitignore

vendor/
*.swp
*~

$ cat .ignoremetoo

file.txt
*.png
```

The following command will create `destination.zip` file by ignoring all patterns mentioned in `.gitignore` and `.ignoremetoo`.
```sh
$ pack -ig .gitignore -ig .ignoremetoo /path/to/file-to-compress /part/to/destination.zip
```
