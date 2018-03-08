# dtags

### convenient directory tagging and navigation
dtags allows you to tag any directory, view your tags, and `cd` to those tags

## Installation

* [go install](#go-install)
* [github install](#github-install)
* [direct download](#direct-download-install)

#### Alias:
add the following bash alias to your bashrc/zshrc
```
alias dt=". ~/go/src/github.com/tylerfowle/dtags/dt"
```

#### Go Install:


```
go get tylerfowle/dtags
```


#### Github Install:
```
git clone git@github.com:tylerfowle/dtags.git
```


#### Direct Download Install:


## Usage
Command | Description
---     | ---
`add`                         | add tag to current _working_ directory
`del`                         | delete tag and associated path from database, dtags doesnt care were this command is ran from
`list`                        | lists all tags associated with current directory, is backwards compatible with bash version
`all`,`more`,`<no arg>`       | prints all tags and their associated directories
`<string>`                    | returns a path from a tag, `cd`s you to the directory when called from bash helper script/alias


## Todos:
- [ ] case sensitive tags?
- [ ] prompt for overwrite of tag?
- [ ] change all commands to be flags?
- [ ] bash completion
- [x] zsh completion
- [ ] legacy database conversion?
- [ ] add fuzzy searching when called with no arguments
- [ ] add install script
- [ ] add uninstall script
- [ ] add ability to manually enter path when adding tag
- [ ] bash helper script - make dtags path easy to change
