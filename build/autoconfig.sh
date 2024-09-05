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


network_state=$(virsh net-list --all | awk '/default/ {print $2}')

if [ "$network_state" != "active" ]; then
    virsh net-autostart default
    virsh net-start default
else
    echo "The default network is already active."
fi
