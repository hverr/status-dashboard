#!/bin/bash

##############################################################################
# Usage: ./test_release.sh executable1 [executable2 ...]                     #
#                                                                            #
# This script makes sure version information is available in an executable.  #
# Furthermore it makes sure the executable was not built in a dirty          #
# worktree.                                                                  #
##############################################################################

set -e

run_tests() {
    local arg="$1"
    local version="$( "$arg" -version )"
    local status=$?

    if [ $status -ne 0 ]; then
        echo "$arg does not have any version information. (got return code 1)"
        return 1
    fi

    echo "$version" | grep -i 'No version info' > /dev/null
    if [ $? -eq 0 ]; then
        echo "$arg does not have any version information. (parsed output)"
        return 1
    fi

    local commit=$(echo "$version" | awk '{print $2}')
    git rev-parse --verify "$commit" >/dev/null 2>/dev/null
    if [ $? -ne 0 ]; then
        echo "$arg does not have a valid commit version."
        return 1
    fi

    echo "$version" | grep -i 'dirty' > /dev/null
    if [ $? -eq 0 ]; then
        echo "$arg is built in a dirty worktree."
        return 1
    fi

    return 0
}

test_executable() {
    if ! run_tests "$1"; then
        echo "Tests dit not pass for $1, aborting." >&2
        exit 1
    fi
}

if [ $# -lt 1 ]; then
    echo "Usage: $0 executable1 [executable2 ...]"
    echo ""
    echo "Check if the executable is properly compiled."
    exit 1
fi

while [ $# -ge 1 ]; do
    test_executable "$1"
    shift
done
