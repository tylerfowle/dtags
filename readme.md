# :bookmark: dtags
#### directory tagging and navigation utility.

Dtags makes tagging, and change directories fast and easy. Add arbitrary tags
to any directory, making it easy to jump back to those directories easy.

## Installation

Installation consists of 3 parts.
1. an alias in your bashrc/zshrc.
2. bash helper script that is used to launch commands and `cd` to tagged directories.
3. the guts, a go utility that does all the heavy lifting.


#### Download dtags and set permissions:
```
# Clone the repo
git clone https://github.com/tylerfowle/dtags
cd dtags

# Build both slices
GOARCH=arm64 GOOS=darwin go build -o dtags-arm64 main.go
GOARCH=amd64 GOOS=darwin go build -o dtags-amd64 main.go

# Combine into a universal binary
lipo -create -output dtags dtags-arm64 dtags-amd64

# Install & set permissions
sudo mv dtags /usr/local/bin/dtags
sudo chmod +x /usr/local/bin/dtags
```

#### Install bash helper script:

```
dtags install
```

#### Alias: add to your bashrc/zshrc:

```
alias dt=". ~/.config/dtags/dt"
```

## Usage
Command | Description
---     | ---
`<string>`         | returns a path from a tag, `cd`s you to the directory when called from bash helper script/alias
`add`              | add tag to current _working_ directory
`del`              | delete tag and associated path from database, dtags doesnt care were this command is ran from
`ls`               | prints all tags and their associated directories
`list`             | lists all tags on current working directory
`completion`       | returns a list of all tags in database.  (used for bash/zsh completion)
`install`          | install dt bash helper script

## Todos:
- [ ] add images/gifs to readme
- [x] make tags case insensitive
- [x] zsh completion
- [ ] bash completion
- [x] confirm overwrite?
- [x] add install script
- [ ] add uninstall script
- [x] add ability to manually enter path when adding tag
- [ ] ~~bash helper script - make dtags path easy to change~~
- [ ] add ability to add and delete multiple tags at once?
- [ ] ~~config file? alias name? install location? bucket name?~~
