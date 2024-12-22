#!/bin/bash

# 현재 셸 타입을 확인하는 함수
get_shell_type() {
    ps -p $$ -o cmd= | awk '{print $1}'
}

shell_type=$(get_shell_type)

# 셸에 따라 적절한 설정 파일을 수정
if [[ "$shell_type" == "bash" ]]; then
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
elif [[ "$shell_type" == "zsh" ]]; then
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.zshrc
    source ~/.zshrc
else
    echo "Unsupported shell: $shell_type"
fi

