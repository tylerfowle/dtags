dtags
-----

Command Line Tool for Tagging Directories.


dtags allows you to:
* add and remove directory tags
* list all tags on current directory
* change directory (cd) to any tag


Prerequisites
-------------

dtags uses `fzf` for searching through an index of tagged directories.

install fzf here:
[junegunn/fzf](https://github.com/junegunn/fzf)



Installation
------------

### Using Git

Clone this repository and run
[install](https://github.com/tylerfowle/dtags/blob/master/install) script.

```sh
git clone --depth 1 https://github.com/tylerfowle/dtags.git ~/.dtags
~/.dtags/install
```

### Manual Download
Download this repo and run the
[install](https://github.com/tylerfowle/dtags/blob/master/install) script.

navigate to dtags directory
```
cd ./dtags
```
run install script
```
./install
```

Alias
-----

add the following to your (bash_profile|bashrc|zshrc)

```sh
# dtags - tag and cd to directories
alias t=". dtags"
```


Usage
-----

| command  | flag | description                    |
|----------|------|--------------------------------|
| `search` | null | default behavior. invoked by `. dtags`. uses `fzf` to search index of tags and paths  |
| `add`    | `-a` | add tag, takes argument        |
| `remove` | `-r` | remove tag, takes argument     |
| `list`   | `-l` | list tags on current directory |


Examples
--------

### adding and removing tags

![](https://raw.github.com/tylerfowle/i/master/dtags-addremove.gif)
