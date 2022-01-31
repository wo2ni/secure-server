#!/usr/bin/env bash

#-------anti-portscan----------;
#-------端口扫描攻击防御-------;
#-------anti-portscan----------;

#--Color code;
Cls="\033[0m"
White="\033[1;38m"
Red="\033[1;31m"
Green="\033[1;32m"
Yellow="\033[1;33m"
Blue="\033[1;34m"
Purple="\033[1;35m"
Cyan_blue="\033[1;36m"

#Error logo;
Error_show() {
    clear
    echo -e "${Red}"
    cat << "EOF"
      _____
     | ____|  _ __   _ __    ___    _ __
     |  _|   | '__| | '__|  / _ \  | '__|
     | |___  | |    | |    | (_) | | |
     |_____| |_|    |_|     \___/  |_|

EOF
     echo -e "${Cls}"
}

#Secure port logo;
Logo_show() {
    clear
    echo -e "${Cyan_blue}"
    cat << "EOF"
      _____                                             _
     / ____|                                           | |
    | (___   ___  ___ _   _ _ __ ___   _ __   ___  _ __| |_
     \___ \ / _ \/ __| | | | '__/ _ \ | '_ \ / _ \| '__| __|
     ____) |  __/ (__| |_| | | |  __/ | |_) | (_) | |  | |_
    |_____/ \___|\___|\__,_|_|  \___| | .__/ \___/|_|   \__|
                                      | |
                                      |_|
EOF
     echo -e "${Cls}"
}

#--运行环境检测;
Check_run() {
    if ! [[ $(command -v iptables |grep -o 'iptables') == 'iptables' ]] && [[ $(command -v ipset |grep -o 'ipset') == 'ipset' ]]; then
        __Error_show
        printf "${Red}Error!${Red},${Cls}运行此软件,您必须安装${Green}iptables${Cls}和${Cls}${Green}ipset${Cls}工具${Cls}\a\n"
    fi
}

#--1级主菜单;
Menu_0() {
    Logo_show
    printf "${Blue}0)${Purple}布置端口陷阱,抵御扫描工具;${Cls}\n"
    printf "${Blue}1)${Purple}关闭端口陷阱;${Cls}\n"
    printf "${Blue}2)${Purple}增加放行端口;${Cls}\n"
    printf "${Blue}3)${Purple}删除放行端口;${Cls}\n"
    printf "${Blue}4)${Purple}查看放行的端口;${Cls}\n"
    printf "${Blue}5)${Purple}查看禁用的IP;${Cls}\n"
    printf "${Blue}6)${Purple}放行禁用的IP;${Cls}\n"
    printf "${Blue}q)${Purple}退出软件;${Cls}\n"
    read -p "请输入:" _input_0
    case ${_input_0} in
        '0')
            Menu_1
            ;;
        '1')
            #--关闭端口陷阱;
            Stop_port_trap
            ;;
        '2')
            #--增加放行的端口;
            Allow_port
            ;;
        '3')
            #--删除放行端口;
            Deny_port
            ;;
        '4')
            #--查看放行的端口;
            ipset list pub-port-set
            ;;
        '5')
            ;;
        '6')
            ;;
        'q')
            clear
            exit 0
            ;;
        *)
            Error_show
            printf "${Red}错误的输入;\a\n${Cls}"
            sleep 1
            Menu_0
            ;;
    esac
}

#--2级主菜单(1);
Menu_1() {
    Logo_show
    printf "${Blue}0)${Purple}:选择网络接口;${Cls}\n"
    printf "${Blue}1)${Purple}:返回上层菜单;${Cls}\n"
    printf "${Blue}q)${Purple}退出软件;${Cls}\n"
    read -p "请输入:" _input_0
    case ${_input_0} in
        '0')
            #--获取并选择网络接口;
            Get_net_int
            Select_net_dev="${_select_dev}"
            readonly Select_net_dev
            Start_port_trap
            Allow_port
            ;;
        '1')
            #--返回上层菜单 -> Menu_0
            Menu_0
            ;;
        'q')
            clear
            exit 0
            ;;
        *)
            Error_show
            printf "${Red}错误的输入;\a\n${Cls}"
            sleep 1
            Menu_1
            ;;
    esac
}

#Select Network Intfaces;
Select_net_int() {
    net_num=0
    for i in /sys/class/net/*
    do
        array="$(echo ${i}|cut -d / -f 5 |grep -Ev 'lo')"
        net_dev[${net_num}]="${array}"
        ((net_num++))
    done
    printf "${White}本机网卡成员:${Cls}${Blue}${net_dev[*]}${Cls}\n"
    printf "${White}本机网卡数量:$((${#net_dev[@]} - 1))${Cls}\n"

    printf "${Green}请输入网卡接口:${Cls}" ;read _select_dev
    if [ ! ${_select_dev} ]; then
        Select_net_int
    else
        for (( x=${_select_dev};x>=0;x-- ))
        do
            if echo ${_select_dev} |grep "${net_dev[x]}" ;then
                #Net_Dev="${_select_dev}"
                #readonly Net_Dev
                return 0
            else
                printf "${Red}错误的网络接口;\n${Cls}"
                return 2
            fi
        done
    fi
}

#Get User Services Network Intfaces;
Get_net_int() {
    while ((1))
    do
        if Select_net_int == 0 ; then
            echo Pass
            #readonly _select_dev
            break
        else
            printf "${Red}FAIL:当前计算机没有${_select_dev}网络接口!${Cls}\n"
            #Select_net_int
            #continue
        fi
    done
}

#--布置端口陷阱;
Start_port_trap() {
    #如果有IP连接未开放端口,该IP将进入扫描者名单,过期时间,IP_DENY_SECOND秒;(默认:30S)
    #如果该IP连续连接未开放端口,过期时间不复位,但包计数器会累计,如果累计超过PORT_SCAN_MAX,该IP将无法连接任何端口,直到过期.
    #
    IP_DENY_SECOND=30
    PORT_SCAN_MAX=3

    ipset create pub-port-set bitmap:port range 0-65535
    IP_SET_MAX=$((100 * 1024 * 1024 / 8 / 60 * $IP_DENY_SECOND))
    ipset create scanner-ip-set hash:ip timeout $IP_DENY_SECOND maxelem $IP_SET_MAX counters
    iptables -N trap-scan
    iptables -A trap-scan -m set --match-set scanner-ip-set src -j DROP
    iptables -A trap-scan -j SET --add-set scanner-ip-set src
    iptables -A trap-scan -j DROP
    iptables -i $Select_net_dev -A INPUT -p tcp --syn -m set ! --match-set pub-port-set dst -j trap-scan
    iptables -i $Select_net_dev -A INPUT -p tcp --syn -m set ! --update-counters --match-set scanner-ip-set src --packets-gt $PORT_SCAN_MAX  -j DROP
    iptables -i $Select_net_dev -A INPUT -p tcp ! --syn -m conntrack ! --ctstate ESTABLISHED,RELATED -j DROP
}

#--关闭端口陷阱;
Stop_port_trap() {
    #TODO: 只删除我们创建的规则,保留原有的;
    while :
    do
        printf "${Red}只删除我们创建的规则,保留原有的规则吗?${Cls}\n"
        read -p "(Yy保留|nN清空所有规则)" _int_cls
        case ${_int_cls} in
            Y|y)
                echo '只删除我们所创建的规则;'
                # iptables -L | grep -P "trap-scan|scanner"
                #iptables -F
                iptables -F trap-scan
                iptables -X trap-scan
                ipset destroy scanner-ip-set
                ipset destroy pub-port-set
                break
                ;;
            N|n)
                echo '删除所有的规则'
                iptables -F
                break
                ;;
            *)
                echo '错误的输入'
                ;;
        esac
    done
}

#--放行某个端口;
Allow_port() {
    while :
    do
        printf "请输入端口:(1~65535)" ;read _int_port
        case $_int_port in
            [1-9] | [1-9][0-9] | [1-9][0-9][0-9] | [1-9][0-9][0-9][0-9] | [1-5][0-9][0-9][0-9][0-9] | 6[0-4][0-9][0-9][0-9] | 65[0-4][0-9][0-9] | 655[0-3][0-5])
                break
                ;;
            *)
                printf "错误的端口输入:端口范围(1~65535)"
                ;;
        esac
    done
    ipset add pub-port-set ${_int_port}
    printf "${Green}${_int_port}${Cls} -- ${Green} is Allow${Cls}\n"
}

#--删除某个端口;
Deny_port() {
    while :
    do
        printf "请输入端口:(1~65535)" ;read _int_port
        case $_int_port in
            [1-9] | [1-9][0-9] | [1-9][0-9][0-9] | [1-9][0-9][0-9][0-9] | [1-5][0-9][0-9][0-9][0-9] | 6[0-4][0-9][0-9][0-9] | 65[0-4][0-9][0-9] | 655[0-3][0-5])
                break
                ;;
            *)
                printf "错误的端口输入:端口范围(1~65535)"
                ;;
        esac
    done
    ipset del pub-port-set ${_int_port}
    printf "${Green}${_int_port}${Cls} -- ${Green} is Allow${Cls}\n"
}

main ()
{
    Check_run
    Menu_0
}

main
