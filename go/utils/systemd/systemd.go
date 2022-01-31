package systemd

import (
    "fmt"
    "os/exec"
    "strings"
)

//自动那个开启/关闭某服务;
func Auto_Control_Server(server string)(string) {
    shell := "systemctl status "
    server_name := server
    grep := "|grep -Eo 'running'"

    server_status, err := exec.Command("bash", "-c", shell+server_name+grep).Output()

    //将获取到的ssh服务状态转换成字符串,并去掉换行符;
    status := string(server_status)
    status = strings.Replace(status, "\n", "", -1)

    if err != nil {
        //fmt.Println("Error")
        //fmt.Println("Error")
        Start_Server(server_name)
        status = server+" Running"
    }


    if status != "running" {
        //fmt.Println("这个服务没有启动,即将启动该服务!!")
        Start_Server(server_name)
    } else {
        //fmt.Println("这个服务已经启动,即将关闭该服务!!")
        status = server+" Closed"
        Stop_Server(server_name)
    }
    return status
}

//Get 某个服务状态;
func Get_Server_Status(server string)(string) {
    shell := "systemctl status "
    server_name := server
    grep := "|grep -Eo 'running'"

    server_status, err := exec.Command("bash", "-c", shell+server_name+grep).Output()

    if err != nil {
        fmt.Printf("Failed to execute command: %s\n", shell, server_name, grep)
    }

    //将获取到的ssh服务状态转换成字符串,并去掉换行符;
    status := string(server_status)
    status = strings.Replace(status, "\n", "", -1)

    if status != "running" {
        fmt.Println("这个服务没有启动!!")
        status = "!!NOT running!!;"
    }
    return status
}

//Get 某个服务配置文件;
func Get_Server_Conf(conf_file string)(string) {
    grep := "grep -Ev ';|#|^$'  "
    file := conf_file

    conf, err := exec.Command("bash", "-c", grep+file).Output()

    if err != nil {
        fmt.Printf("Failed to execute command: %s", grep,file)
    }

    server_conf := string(conf)
    //ssh_conf = strings.Replace(ssh_conf, "\n", "", -1)

    return server_conf
}

//Start 某个服务;
func Start_Server(server_start string) {
    shell_start := "systemctl start "
    server_name := server_start

    exec.Command("bash", "-c", shell_start+server_name).Output()
}

//Stop 某个服务;
func Stop_Server(server_stop string) {
    shell := "systemctl stop "
    server_name := server_stop

    exec.Command("bash", "-c", shell+server_name).Output()
}
