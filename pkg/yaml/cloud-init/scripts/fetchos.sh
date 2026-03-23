ASCII_ART="::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::
::                                  ::
::        __          __    ____    ::
::   ____/ /___  ____/ /___/ / /__  ::
::  / __  / __ \/ __  / __  / / _ \ ::
:: / /_/ / /_/ / /_/ / /_/ / /  __/ ::
:: \__,_/\____/\__,_/\__,_/\___/  ::
::                                  ::
::                                  ::
::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::"

display_info() {

    echo "$ASCII_ART"


    local user=$(whoami)
    local hostname=$(hostname)
    local os=$(grep PRETTY_NAME /etc/os-release | cut -d'"' -f2)
    local kernel=$(uname -r)
    local uptime=$(uptime -p | sed 's/up //')
    local cpu=$(lscpu | grep 'Model name:' | sed 's/Model name:[ \t]*//')
    local mem_info=$(free -h | awk 'NR==2{printf "%s / %s", $3, $2}')
    local disk_info=$(df -h / | awk 'NR==2{printf "%s / %s (%s)", $3, $2, $5}')


    printf "  %-34s\n" "User: $user@$hostname"
    printf "  %-34s\n" "OS: $os"
    printf "  %-34s\n" "Kernel: $kernel"
    printf "  %-34s\n" "Uptime: $uptime"
    printf "  %-34s\n" "CPU: $cpu"
    printf "  %-34s\n" "Memory: $mem_info"
    printf "  %-34s\n" "Disk(/): $disk_info"
}


display_info
