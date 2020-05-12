#!/bin/sh

ALIAS="alias avia=/etc/avia/binary"
sudo mkdir /etc/avia
sudo curl -fsSL https://github.com/Joeverson/git-job/releases/download/0.0.1/binary --output /etc/avia/binary --silent || {
    echo "Error in download application"
    exit 1
}

sudo chmod -R +x /etc/avia/binary


if [[ -e ~/.zshrc ]]; then
    echo $ALIAS >> ~/.zshrc || {
        echo "Error in add alias in .zshrc"
        exit 1
    }
elif [[ -e ~/.bashrc ]]; then
    echo $ALIAS >> ~/.bashrc || {
        echo "Error in add alias in .bashrc"
        exit 1
    }
fi

echo "Done, close your terminal to effective actions."