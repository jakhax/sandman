#!/bin/bash
set -e
if [ "$1" = 'run_json' -o "$1" = 'run' ]; then
    set -- './sandman' "$@"
fi
exec "$@"