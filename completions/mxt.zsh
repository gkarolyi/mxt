#compdef mxt

_mxt_managed_branches() {
    local worktree_dir repo_name wt_base
    worktree_dir="$HOME/worktrees"
    local config_dir="${MXT_CONFIG_DIR:-$HOME/.mxt}"
    local config_file="$config_dir/config"

    # Load worktree_dir from global config
    if [[ -f "$config_file" ]]; then
        local val
        val=$(grep -E '^worktree_dir=' "$config_file" 2>/dev/null | head -1 | cut -d= -f2-)
        val="${val#"${val%%[![:space:]]*}"}"
        val="${val%"${val##*[![:space:]]}"}"
        [[ -n "$val" ]] && worktree_dir="$val"
    fi

    # Override with project-local config
    local repo_root
    repo_root=$(git rev-parse --show-toplevel 2>/dev/null) || return
    if [[ -f "$repo_root/.mxt" ]]; then
        local val
        val=$(grep -E '^worktree_dir=' "$repo_root/.mxt" 2>/dev/null | head -1 | cut -d= -f2-)
        val="${val#"${val%%[![:space:]]*}"}"
        val="${val%"${val##*[![:space:]]}"}"
        [[ -n "$val" ]] && worktree_dir="$val"
    fi

    worktree_dir="${worktree_dir/#\~/$HOME}"
    repo_name=$(basename "$repo_root")
    wt_base="$worktree_dir/$repo_name"

    local branches=()
    while IFS= read -r line; do
        local wt_dir branch
        wt_dir=$(echo "$line" | awk '{print $1}')
        branch=$(echo "$line" | sed -n 's/.*\[\(.*\)\].*/\1/p')
        [[ -n "$branch" ]] || continue
        [[ "$wt_dir" == "$wt_base"/* ]] || continue
        branches+=("$branch")
    done < <(git worktree list 2>/dev/null)
    echo "${branches[@]}"
}

_mxt_git_branches() {
    git branch -a --format='%(refname:short)' 2>/dev/null
}

_mxt() {
    local -a commands session_actions
    commands=(
        'init:Set up mxt config'
        'config:Show current config'
        'new:Create worktree + tmux session'
        'list:List worktrees and session status'
        'ls:List worktrees and session status'
        'delete:Delete worktree and branch'
        'rm:Delete worktree and branch'
        'sessions:Manage tmux sessions'
        's:Manage tmux sessions'
        'help:Show help message'
        'version:Print version number'
    )
    session_actions=(
        'open:Create session and open terminal'
        'launch:Create session and open terminal'
        'start:Create session and open terminal'
        'close:Kill tmux session'
        'kill:Kill tmux session'
        'stop:Kill tmux session'
        'relaunch:Close and reopen session'
        'restart:Close and reopen session'
        'attach:Attach to session'
    )

    local curcontext="$curcontext" state line
    _arguments -C \
        '1:command:->command' \
        '*::arg:->args'

    case $state in
        command)
            _describe -t commands 'mxt command' commands
            ;;
        args)
            local cmd="${line[1]}"
            case "$cmd" in
                init)
                    _arguments \
                        '(-l --local)'{-l,--local}'[Create project-local config]' \
                        '--reinit[Overwrite existing config without prompting]'
                    ;;
                config|list|ls|help|version)
                    ;;
                new)
                    _arguments \
                        '1:branch:' \
                        '--from[Base branch]:branch:($(_mxt_git_branches))' \
                        '--run[Auto-run command in agent window]:command:(claude codex)' \
                        '--bg[Create session without opening terminal]'
                    ;;
                delete|rm)
                    _arguments \
                        '1:branch:($(_mxt_managed_branches))' \
                        '(-f --force)'{-f,--force}'[Skip confirmation]'
                    ;;
                sessions|s)
                    _arguments -C \
                        '1:action:->action' \
                        '*::arg:->session_args'

                    case $state in
                        action)
                            _describe -t session_actions 'session action' session_actions
                            ;;
                        session_args)
                            local action="${line[1]}"
                            case "$action" in
                                open|launch|start)
                                    _arguments \
                                        '1:branch:($(_mxt_managed_branches))' \
                                        '--run[Auto-run command]:command:(claude codex)' \
                                        '--bg[Create without opening terminal]'
                                    ;;
                                close|kill|stop)
                                    _arguments \
                                        '1:branch:($(_mxt_managed_branches))'
                                    ;;
                                relaunch|restart)
                                    _arguments \
                                        '1:branch:($(_mxt_managed_branches))' \
                                        '--run[Auto-run command]:command:(claude codex)' \
                                        '--bg[Create without opening terminal]'
                                    ;;
                                attach)
                                    _arguments \
                                        '1:branch:($(_mxt_managed_branches))' \
                                        '2:window:(dev agent)'
                                    ;;
                            esac
                            ;;
                    esac
                    ;;
            esac
            ;;
    esac
}

# Register the completion function.
# When sourced directly, compdef handles registration.
# When autoloaded from fpath (as _mxt), #compdef at the top handles it.
compdef _mxt mxt 2>/dev/null
