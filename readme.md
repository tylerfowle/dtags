# dtags


add the following bash alias to your bashrc/zshrc
```
alias dt=". ~/go/src/github.com/tylerfowle/dtags/dt"
```


## commands

`add` : add tag to current directory
`del` : delete tag and associated path from database
`` : returns a path from a tag
`tags` : lists all tags associated with current directory

`list`: prints all tags and their associated directories





#### issues:
case sensitive tags?
change all commands to be flags?
add fuzzy searching when called with no arguments
add install script
add uninstall script
add ability to manually enter path when adding tag

bash helper:
- make dtags path easy to change

rework go:
- move database open/close to on function
- cleanup functions
