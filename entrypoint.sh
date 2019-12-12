#!/bin/bash
set -e
if [ "$1" = 'run-json' -o "$1" = 'run' ]; then
    set -- './sandman' "$@"
    exec "$@"
else 
    echo "Usage: run-json or run"
fi