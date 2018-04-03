#!/bin/bash
PREFIX="sudo -k "
COMMAND="${PREFIX}find . -name \"*.pyc\" -o -name \"__pycache__\" -delete"
echo "Cleaning up stale Python bytecode ($COMMAND)..."
eval $COMMAND
