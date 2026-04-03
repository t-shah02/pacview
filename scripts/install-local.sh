#!/usr/bin/env bash
# Build pacview, install to ~/.local/bin, and add a shell alias to your profile.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_NAME="pacview"
INSTALL_DIR="${HOME}/.local/bin"
INSTALLED_BIN="${INSTALL_DIR}/${BIN_NAME}"
MARKER_BEGIN="# >>> pacview install-local (managed block)"
MARKER_END="# <<< pacview install-local"

die() {
	echo "install-local: $*" >&2
	exit 1
}

detect_profile() {
	case "${SHELL##*/}" in
	zsh)
		echo "${ZDOTDIR:-${HOME}}/.zshrc"
		;;
	bash)
		if [[ -f "${HOME}/.bashrc" ]]; then
			echo "${HOME}/.bashrc"
		elif [[ -f "${HOME}/.bash_profile" ]]; then
			echo "${HOME}/.bash_profile"
		else
			echo "${HOME}/.bashrc"
		fi
		;;
	*)
		echo "${HOME}/.profile"
		;;
	esac
}

ensure_profile_block() {
	local profile="$1"
	local block
	block="$(printf '%s\n' \
		"${MARKER_BEGIN}" \
		"alias ${BIN_NAME}=\"\${HOME}/.local/bin/${BIN_NAME}\"" \
		"${MARKER_END}")"

	if [[ ! -f "$profile" ]]; then
		touch "$profile" || die "cannot create profile: $profile"
	fi

	if grep -qF "${MARKER_BEGIN}" "$profile" 2>/dev/null; then
		echo "Profile already contains pacview block: $profile"
		return 0
	fi

	printf '\n%s\n' "$block" >>"$profile" || die "cannot append to $profile"
	echo "Added pacview alias to $profile"
}

cd "$REPO_ROOT" || die "cannot cd to repo root: $REPO_ROOT"

command -v make >/dev/null 2>&1 || die "make is not installed"
make build

[[ -f "${REPO_ROOT}/bin/${BIN_NAME}" ]] || die "build did not produce bin/${BIN_NAME}"

mkdir -p "$INSTALL_DIR"
cp -f "${REPO_ROOT}/bin/${BIN_NAME}" "$INSTALLED_BIN"
chmod +x "$INSTALLED_BIN"
echo "Installed ${INSTALLED_BIN}"

PROFILE="$(detect_profile)"
ensure_profile_block "$PROFILE"

echo
echo "Done. Run: source \"${PROFILE}\""
echo "Or open a new terminal, then run: ${BIN_NAME}"
