#compdef dt
#
# zsh completion for dtags
#
# Recommended installation:
#
# copy this script to `~/.zsh/_dtags`
# and then add the following to your ~/.zshrc file:
#
# fpath=(~/.zsh $fpath)

_dtags_all_tags() {
  all_tags=(`dtags completion`)
}

local -a _1st_arguments
_1st_arguments=(
'add:add tag to current directory'
'del:delete any tag from database'
'list:list all tags associated with current directory'
'all:list all tags and directories'
'more:list all tags and directories'
)

local -a all_tags

case "$words[0]" in
  *)
    _dtags_all_tags
    _wanted all_tags expl 'all tags' compadd -a all_tags
    ;;
esac
