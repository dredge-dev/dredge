#!/usr/bin/env bash

ARTIFACT_PATH="https://github.com/dredge-dev/dredge/releases/latest/download/"
ARTIFACT_NAME="drg"
INSTALL_PATH="/usr/local/bin/"

TTY_BOLD="\033[1m"
TTY_NORMAL="\033[0m"

echo -e "
 __   __   ___  __   __   ___ 
|  \ |__) |__  |  \\ / _\` |__  
|__/ |  \ |___ |__/ \\__> |___ 
Automates DevOps workflows

Installing the latest version of ${TTY_BOLD}${ARTIFACT_NAME}${TTY_NORMAL} in ${TTY_BOLD}${INSTALL_PATH}${TTY_NORMAL}
"

OS="$(uname)"
ARCH="$(uname -m)"

if [[ "$OS" == "Linux" && "$ARCH" == "x86_64" ]]; then
    FLAVOR="linux-amd64"
elif [[ "$OS" == "Linux" && "$ARCH" == "arm64" ]]; then
    FLAVOR="linux-arm64"
elif [[ "$OS" == "Darwin" && "$ARCH" == "x86_64" ]]; then
    FLAVOR="darwin-amd64"
elif [[ "$OS" == "Darwin" && "$ARCH" == "arm64" ]]; then
    FLAVOR="darwin-arm64"
else
    echo "Your system is not supported, we currently only support Linux and MacOS on x86_64 and amd64 processors."
    exit 1
fi

curl -fsSL "${ARTIFACT_PATH}${ARTIFACT_NAME}-${FLAVOR}" -o "${ARTIFACT_NAME}"
chmod +x $ARTIFACT_NAME

echo -e "Using ${TTY_BOLD}sudo${TTY_NORMAL} to install to $INSTALL_PATH, this might require your password."
sudo mv $ARTIFACT_NAME $INSTALL_PATH

echo
echo -e "${TTY_BOLD}Success!${TTY_NORMAL} You can now start using Dredge:"
echo
echo "  drg init          # initialize the project"
echo "  drg add release   # add your first resource"
echo
