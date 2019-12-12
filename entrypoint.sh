#!/bin/bash
set -e
if [ "$1" = 'run_json' -o "$1" = 'run' ]; then
    set -- './sandman' "$@"
    exec "$@"
else 
    echo "Usage: run_json or run"
fi