#!/usr/bin/env bash
# restart.sh

# Load your zsh environment
# -L disables login shell behavior, -i forces interactive mode
# This makes sure all your ~/.zshrc settings (like aliases and PATH) are loaded
zsh -i -c 'dc down && dc up -d --build && dc logs -f api'
