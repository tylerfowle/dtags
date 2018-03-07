# dtags
add the following bash alias to your bashrc/zshrc
```
alias dt=". ~/go/src/github.com/tylerfowle/dtags/dt"
```


## commands

`add` : add tag to current _working_ directory

`del` : delete tag and associated path from database, dtags doesnt care were this command is ran from

`list` : lists all tags associated with current directory, is backwards compatible with bash version

`[string]` : returns a path from a tag, `cd`s you to the directory when called from bash helper script/alias

`all` , `more`, or 'no arg': prints all tags and their associated directories


#### issues/todos:
- [ ] case sensitive tags?
- [ ] prompt for overwrite of tag?
- [ ] change all commands to be flags?
- [ ] bash completion
- [x] zsh completion

- [ ] add fuzzy searching when called with no arguments
- [ ] add install script
- [ ] add uninstall script
- [ ] add ability to manually enter path when adding tag

#### bash helper:
- [ ] make dtags path easy to change

#### rework go:
- [x] move database open/close to init function
- [ ] cleanup functions

