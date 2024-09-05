shell_type=$0

if [[ "$shell_type" == "-bash" ]]; then
    echo export PATH=$PATH:/usr/local/go/bin >> ~/.bashrc
    echo "$shell_type"
    source ~/.bashrc
elif [[ $"shell_type" == "-zsh" ]]; then
    echo export PATH=$PATH:/usr/local/go/bin >> ~/.zshrc
    echo "$shell_type"
    source ~/.zshrc
else 
    echo "hey"
fi