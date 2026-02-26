#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HARNESS="$ROOT_DIR/test/harness.sh"
MUXTREE_BIN="${MUXTREE_BIN:-$ROOT_DIR/muxtree}"
MXT_BIN="${MXT_BIN:-$ROOT_DIR/mxt}"

if [[ ! -x "$MUXTREE_BIN" ]]; then
  echo "muxtree binary not found at $MUXTREE_BIN" >&2
  exit 1
fi

if [[ ! -x "$MXT_BIN" ]]; then
  echo "mxt binary not found at $MXT_BIN" >&2
  exit 1
fi

if ! command -v git >/dev/null 2>&1; then
  echo "git is required for this runner" >&2
  exit 1
fi

if ! command -v tmux >/dev/null 2>&1; then
  echo "tmux is required for this runner" >&2
  exit 1
fi

TMP_ROOT="$(mktemp -d -t mxt-feature-suite.XXXXXX)"
CONFIG_DIR="$TMP_ROOT/config"
WORKTREE_DIR="$TMP_ROOT/worktrees"
TEST_REPO="$TMP_ROOT/testrepo"
INIT_INPUT="$TMP_ROOT/init-input.txt"

cleanup() {
  rm -rf "$TMP_ROOT"
}
trap cleanup EXIT

mkdir -p "$CONFIG_DIR" "$WORKTREE_DIR" "$TEST_REPO"

init_repo() {
  pushd "$TEST_REPO" >/dev/null
  git init -b main >/dev/null
  git config user.email "test@example.com"
  git config user.name "Test User"
  printf "README" > README.md
  printf "claude" > CLAUDE.md
  printf "env" > .env
  printf "test" > test.md
  git add . >/dev/null
  git commit -m "Initial commit" >/dev/null
  git branch develop >/dev/null
  popd >/dev/null
}

write_config() {
  local terminal="$1"
  local copy_files="$2"
  local pre_session_cmd="$3"
  local tmux_layout="$4"
  cat > "$CONFIG_DIR/config" <<EOF
worktree_dir=$WORKTREE_DIR
terminal=$terminal
copy_files=$copy_files
pre_session_cmd=$pre_session_cmd
tmux_layout=$tmux_layout
EOF
}

run_harness() {
  local muxtree_path="$1"
  local mxt_path="$2"
  shift 2
  MUXTREE_CONFIG_DIR="$CONFIG_DIR" MUXTREE_PATH="$muxtree_path" MXT_PATH="$mxt_path" "$HARNESS" "$@"
}

run_init_check() {
  local muxtree_out="$TMP_ROOT/init.muxtree.out"
  local muxtree_err="$TMP_ROOT/init.muxtree.err"
  local mxt_out="$TMP_ROOT/init.mxt.out"
  local mxt_err="$TMP_ROOT/init.mxt.err"
  local muxtree_norm="$TMP_ROOT/init.muxtree.norm"
  local mxt_norm="$TMP_ROOT/init.mxt.norm"

  set +e
  MUXTREE_CONFIG_DIR="$CONFIG_DIR" "$INIT_MUXTREE_WRAPPER" init >"$muxtree_out" 2>"$muxtree_err"
  local muxtree_code=$?
  MUXTREE_CONFIG_DIR="$CONFIG_DIR" "$INIT_MXT_WRAPPER" init >"$mxt_out" 2>"$mxt_err"
  local mxt_code=$?
  set -e

  if [[ $muxtree_code -ne $mxt_code ]]; then
    echo "init exit codes differ: muxtree=$muxtree_code mxt=$mxt_code" >&2
    exit 1
  fi

  sed -E 's/^# Generated on .*/# Generated on <timestamp>/' "$muxtree_out" > "$muxtree_norm"
  sed -E 's/^# Generated on .*/# Generated on <timestamp>/' "$mxt_out" > "$mxt_norm"

  if ! diff -u "$muxtree_norm" "$mxt_norm" >/dev/null; then
    echo "init stdout differs after timestamp normalization" >&2
    diff -u "$muxtree_norm" "$mxt_norm" || true
    exit 1
  fi

  if ! diff -u "$muxtree_err" "$mxt_err" >/dev/null; then
    echo "init stderr differs" >&2
    diff -u "$muxtree_err" "$mxt_err" || true
    exit 1
  fi

  echo "✓ init output matches (timestamp normalized)"
}

make_init_wrappers() {
  local muxtree_wrapper="$TMP_ROOT/muxtree-init-wrapper.sh"
  local mxt_wrapper="$TMP_ROOT/mxt-init-wrapper.sh"

  cat > "$INIT_INPUT" <<EOF
$WORKTREE_DIR
terminal
README.md


EOF

  cat > "$muxtree_wrapper" <<EOF
#!/usr/bin/env bash
export MUXTREE_CONFIG_DIR="$CONFIG_DIR"
rm -f "$CONFIG_DIR/config"
exec < "$INIT_INPUT"
exec "$MUXTREE_BIN" "\$@"
EOF

  cat > "$mxt_wrapper" <<EOF
#!/usr/bin/env bash
export MUXTREE_CONFIG_DIR="$CONFIG_DIR"
rm -f "$CONFIG_DIR/config"
exec < "$INIT_INPUT"
exec "$MXT_BIN" "\$@"
EOF

  chmod +x "$muxtree_wrapper" "$mxt_wrapper"
  INIT_MUXTREE_WRAPPER="$muxtree_wrapper"
  INIT_MXT_WRAPPER="$mxt_wrapper"
}

make_stateful_wrappers() {
  local muxtree_wrapper="$TMP_ROOT/muxtree-stateful-wrapper.sh"
  local mxt_wrapper="$TMP_ROOT/mxt-stateful-wrapper.sh"

  local wrapper_body
  wrapper_body=$(cat <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
CONFIG_DIR="__CONFIG_DIR__"
CONFIG_FILE="$CONFIG_DIR/config"
WORKTREE_DIR="$(grep '^worktree_dir=' "$CONFIG_FILE" | head -n 1 | cut -d= -f2-)"
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
REPO_NAME="$(basename "$REPO_ROOT")"
COMMAND="${1:-}"
ACTION=""
BRANCH=""

sanitize_branch_path() {
  echo "$1" | sed -E 's/[^A-Za-z0-9._-]+/-/g; s/^-+//'
}

sanitize_session_branch() {
  echo "$1" | sed -E 's/[^A-Za-z0-9_-]+/-/g; s/^-+//'
}

set_branch_context() {
  local branch="$1"
  BRANCH="$branch"
  local branch_path
  branch_path="$(sanitize_branch_path "$branch")"
  local branch_session
  branch_session="$(sanitize_session_branch "$branch")"
  WORKTREE_PATH="$WORKTREE_DIR/$REPO_NAME/$branch_path"
  SESSION_NAME="${REPO_NAME}_${branch_session}"
}

cleanup_state() {
  if [[ -z "$BRANCH" ]]; then
    return
  fi
  if tmux has-session -t "$SESSION_NAME" >/dev/null 2>/dev/null; then
    tmux kill-session -t "$SESSION_NAME" >/dev/null 2>/dev/null || true
  fi
  if [[ -d "$WORKTREE_PATH" ]]; then
    git worktree remove "$WORKTREE_PATH" --force >/dev/null 2>/dev/null || rm -rf "$WORKTREE_PATH" >/dev/null 2>/dev/null || true
    git worktree prune >/dev/null 2>/dev/null || true
  fi
  if git show-ref --verify --quiet "refs/heads/$BRANCH"; then
    git branch -D "$BRANCH" >/dev/null 2>/dev/null || true
  fi
}

ensure_worktree() {
  if [[ -z "$BRANCH" ]]; then
    return
  fi
  if [[ ! -d "$WORKTREE_PATH" ]]; then
    git worktree add -b "$BRANCH" "$WORKTREE_PATH" main >/dev/null 2>/dev/null || true
  fi
}

case "$COMMAND" in
  new)
    set_branch_context "${2:-}"
    cleanup_state
    ;;
  delete)
    set_branch_context "${2:-}"
    cleanup_state
    ensure_worktree
    ;;
  sessions)
    ACTION="${2:-}"
    set_branch_context "${3:-}"
    cleanup_state
    ensure_worktree
    case "$ACTION" in
      close|kill|stop|relaunch|restart)
        tmux new-session -d -s "$SESSION_NAME" -c "$WORKTREE_PATH" >/dev/null 2>/dev/null || true
        ;;
    esac
    ;;
  *)
    ;;
esac

export MUXTREE_CONFIG_DIR="$CONFIG_DIR"
exec "__BIN_PATH__" "$@"
EOF
  )

  wrapper_body="${wrapper_body//__CONFIG_DIR__/$CONFIG_DIR}"

  printf "%s" "${wrapper_body//__BIN_PATH__/$MUXTREE_BIN}" > "$muxtree_wrapper"
  printf "%s" "${wrapper_body//__BIN_PATH__/$MXT_BIN}" > "$mxt_wrapper"

  chmod +x "$muxtree_wrapper" "$mxt_wrapper"
  STATEFUL_MUXTREE_WRAPPER="$muxtree_wrapper"
  STATEFUL_MXT_WRAPPER="$mxt_wrapper"
}

confirm() {
  local prompt="$1"
  read -r -p "$prompt [y/N] " response
  if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo "Aborted by user" >&2
    exit 1
  fi
}

manual_compare() {
  local title="$1"
  local muxtree_cmd="$2"
  local mxt_cmd="$3"
  local instructions="$4"

  echo ""
  echo "MANUAL CHECK: $title"
  echo "$instructions"
  confirm "Run muxtree command now?"
  set +e
  eval "$muxtree_cmd"
  local muxtree_status=$?
  set -e
  echo "muxtree exit code: $muxtree_status"
  confirm "Did muxtree behavior match expectations?"

  confirm "Run mxt command now?"
  set +e
  eval "$mxt_cmd"
  local mxt_status=$?
  set -e
  echo "mxt exit code: $mxt_status"
  confirm "Did mxt behavior match muxtree?"
}

cleanup_branch() {
  local repo_root="$1"
  local repo_name="$2"
  local branch="$3"
  local branch_path
  branch_path="$(echo "$branch" | sed -E 's/[^A-Za-z0-9._-]+/-/g; s/^-+//')"
  local session_branch
  session_branch="$(echo "$branch" | sed -E 's/[^A-Za-z0-9_-]+/-/g; s/^-+//')"
  local session_name="${repo_name}_${session_branch}"
  local wt_path="$WORKTREE_DIR/$repo_name/$branch_path"

  tmux kill-session -t "$session_name" >/dev/null 2>/dev/null || true
  if [[ -d "$wt_path" ]]; then
    git -C "$repo_root" worktree remove "$wt_path" --force >/dev/null 2>/dev/null || rm -rf "$wt_path" >/dev/null 2>/dev/null || true
    git -C "$repo_root" worktree prune >/dev/null 2>/dev/null || true
  fi
  if git -C "$repo_root" show-ref --verify --quiet "refs/heads/$branch"; then
    git -C "$repo_root" branch -D "$branch" >/dev/null 2>/dev/null || true
  fi
}

init_repo
make_init_wrappers
make_stateful_wrappers

pushd "$TEST_REPO" >/dev/null
REPO_NAME="$(basename "$TEST_REPO")"

write_config "terminal" "README.md,missing.txt" "echo \"Setup complete\"" ""

echo "== Running non-interactive harness checks =="
run_harness "$MUXTREE_BIN" "$MXT_BIN" help
run_harness "$MUXTREE_BIN" "$MXT_BIN" version
run_init_check

NO_CONFIG_DIR="$TMP_ROOT/no-config"
mkdir -p "$NO_CONFIG_DIR"
MUXTREE_CONFIG_DIR="$NO_CONFIG_DIR" MUXTREE_PATH="$MUXTREE_BIN" MXT_PATH="$MXT_BIN" "$HARNESS" config

write_config "terminal" "README.md" "" ""
run_harness "$MUXTREE_BIN" "$MXT_BIN" config

cat > "$TEST_REPO/.muxtree" <<EOF
copy_files=README.md,CLAUDE.md
pre_session_cmd=
tmux_layout=
EOF
run_harness "$MUXTREE_BIN" "$MXT_BIN" config
rm -f "$TEST_REPO/.muxtree"

write_config "terminal" "README.md,missing.txt" "echo \"Setup complete\"" ""
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" new feature-basic --bg
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" new feature-run --run claude --bg
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" new feature-from --from develop --bg

write_config "terminal" "README.md" "" "dev:hx|lazygit,agent:"
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" new feature-layout --bg

write_config "terminal" "README.md" "false" ""
manual_compare "Pre-session command failure" \
  "MUXTREE_CONFIG_DIR=$CONFIG_DIR $MUXTREE_BIN new feature-fail --bg" \
  "MUXTREE_CONFIG_DIR=$CONFIG_DIR $MXT_BIN new feature-fail --bg" \
  "When prompted, answer 'n' to abort. Verify warning and abort message match muxtree." 
cleanup_branch "$TEST_REPO" "$REPO_NAME" "feature-fail"

write_config "terminal" "README.md" "" ""
rm -rf "$WORKTREE_DIR/$REPO_NAME"
run_harness "$MUXTREE_BIN" "$MXT_BIN" list

git worktree add -b feature-list "$WORKTREE_DIR/$REPO_NAME/feature-list" main >/dev/null
run_harness "$MUXTREE_BIN" "$MXT_BIN" list
cleanup_branch "$TEST_REPO" "$REPO_NAME" "feature-list"

write_config "terminal" "README.md" "" ""
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" sessions open feature-session --bg
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" sessions close feature-session
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" sessions relaunch feature-session --bg

write_config "terminal" "README.md" "" ""
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" delete feature-delete --force

write_config "terminal" "README.md" "" ""
run_harness "$STATEFUL_MUXTREE_WRAPPER" "$STATEFUL_MXT_WRAPPER" new feature-attach --bg
manual_compare "sessions attach" \
  "MUXTREE_CONFIG_DIR=$CONFIG_DIR $MUXTREE_BIN sessions attach feature-attach" \
  "MUXTREE_CONFIG_DIR=$CONFIG_DIR $MXT_BIN sessions attach feature-attach" \
  "This will attach to tmux in the current terminal. Detach with Ctrl-b d, then confirm behavior matches." 
cleanup_branch "$TEST_REPO" "$REPO_NAME" "feature-attach"

for terminal in terminal iterm2 ghostty current; do
  write_config "$terminal" "README.md" "" ""
  manual_compare "Terminal integration ($terminal)" \
    "MUXTREE_CONFIG_DIR=$CONFIG_DIR $MUXTREE_BIN new terminal-$terminal" \
    "MUXTREE_CONFIG_DIR=$CONFIG_DIR $MXT_BIN new terminal-$terminal" \
    "Verify terminal launching/attach behavior. Detach from tmux if needed before continuing." 
  cleanup_branch "$TEST_REPO" "$REPO_NAME" "terminal-$terminal"
done

popd >/dev/null

echo ""
echo "Feature spec run complete."
