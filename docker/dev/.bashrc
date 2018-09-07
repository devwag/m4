# ~/.bashrc: executed by bash(1) for non-login shells.

# ls aliases
export LS_OPTIONS='--color=auto'
eval "`dircolors`"
alias ls='ls $LS_OPTIONS'
alias ll='ls $LS_OPTIONS -alF'
alias l='ls $LS_OPTIONS -lA'

export PATH="$PATH:/usr/local/go/bin"
export GOPATH="/root/go"
