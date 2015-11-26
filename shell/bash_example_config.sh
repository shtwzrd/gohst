#!/bin/sh

# Example rc file for setting up `gohst` on a non-bash/non-zsh shell
# Author: warreq
# based on Glyf's preexec (glyf.livejournal.com/63106.html)

# This configuration is based around emulating the preexec and
# precmd command hooks available in zsh.

# Though portable, this implementation requires certain features
# not explicitly defined in the POSIX standard.

# The shell must have:
#   a) support for the DEBUG signal in conjunction with the trap command
#   b) dynamic evaluation of commands from within the PS1 variable

# Shells lacking this functionality will unfortunately have to settle
# for a more basic integration with `gohst`

# duplicate sourcing guard
if [ "$__prexec_defined" = "true" ]; then
    return 0
fi
__preexec_defined="true"

# contains(string, substring)
#
# Returns 0 if the specified string contains the specified substring,
# otherwise returns 1.
contains() {
    string="$1"
    substring="$2"
    if test "${string#*$substring}" != "$string"
    then
        return 0    # $substring is in $string
    else
        return 1    # $substring is not in $string
    fi
}

__gohst_precmd_hook() {
    r="$1"
    cmd='gohst -u man -d gohst.herokuapp.com log -f result'
    cmd="$cmd $r"
    cmd="$cmd &"
    sh -c "$cmd"
}

# precmd is a hook executed every time the command prompt is drawn.
# Think of it as a set of commands you trigger at the end of every
# command.
#
# Because we have to basically recreate the $PROMPT_COMMAND from
# Bash, you'll need to place any functions you use for setting the
# PS1 prompt here. DO NOT modify $PS1 directly.
precmd () {
    PRECMD_GUARD=true
    __gohst_precmd_hook $1

    # Your precmd hooks and prompt-setters HERE
    PS1_COMMAND $1

    PRECMD_GUARD=false
}

# an example precmd command for configuring your prompt
PS1_COMMAND() {
    if [ "$1" = "0" ]; then
        echo "$(tput setaf 2):) $(tput sgr0)"
    else
        echo "$(tput setaf 1);( $(tput sgr0)"
    fi
}

# capture the state of the shell as it is right before the execution
# of a command, and log it into gohst's history index
__gohst_preexec_hook() {
    __pwd="$(pwd)"
    __user="$(whoami)"
    __shell="$SHELL"
    __host="$(cat /etc/hostname)"
    __cmd="$@"
    gohst -u test -d gohst.herokuapp.com log context $__user $__host $__shell $__pwd "$__cmd"
}

# preexec is invoked right before the execution of every command,
# apart from those that get thrown out by the preexec_invoke_filter.
preexec () {
    __gohst_preexec_hook "$1"

    # Your preexec hooks HERE
}

# preexec_invoke_filter is the function called by the trap, and
# serves as a guard to prevent the triggering of preexec on
# commands that were not issued interactively
preexec_invoke_filter () {
    [ "$COMP_LINE" != "" ] && return  # do nothing if completing

    # For Bash: don't cause a preexec for $PROMPT_COMMAND
    [ -n "$BASH_COMMAND" -a "$BASH_COMMAND" = "$PROMPT_COMMAND" ] && return

    __this_command="$(fc -lr | sed -n 1p | cut -f 2)"
    preexec "$__this_command"
}

# source_preexec_and_precmd installs precmd into the prompt
# and activates a DEBUG trap to create the preexec functionality.
# we pass the last exit code ($?) to precmd from within PS1 so
# that it will be recalculated every time
source_preexec_and_precmd() {
    export PS1='$(precmd $?)'
    trap 'preexec_invoke_filter' DEBUG
}

source_preexec_and_precmd
