package install

import (
	"os"
	"path/filepath"
	"strings"
)

const script = `#!/usr/bin/env bash
dtags={{BIN_PATH}}
declare -a arr=("" "add" "del" "tags" "list" "ls" "completion", "bash-script")
found=0

# check if arg is a single arg command in arr array
for i in "${arr[@]}"; do
    if [[ $i == $1 ]]; then
        found=1
        ${dtags} $@
        false 1
    fi
done

if [[ -d "$(${dtags} $1)" && ${found} -ne 1 ]]; then
    cd "$(${dtags} $1)"
elif [[ ${found} -ne 1 ]]; then
    echo "no directory found for tag [$1]"
    false 1
fi`

func BashHelper() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	s := strings.Replace(script, "{{BIN_PATH}}", filepath.Dir(ex), 1)
	return s, nil
}
