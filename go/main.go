package main

import (
    "log"
    "fmt"
    "strings"
    "os"
    "io/ioutil"
    "time"
    //"os/exec"
    "./utils/systemd"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//读取TOKEN到file中,再利用ioutil将file直接读取到[]byte中,这是性能最好最快的方法;
func Read_Token()  (string) {
    f, err := os.Open("/etc/token")
    if err != nil {
        fmt.Println("read file fail", err)
        return ""
    }
    defer f.Close()

    fd, err := ioutil.ReadAll(f)
    if err != nil {
        fmt.Println("read to token fail", err)
        return ""
    }

    return string(fd)
}

//Get当前系统时间;
func Get_Time_Now()(string) {
    /*当前时间的字符串,2006-01-02 15:04:05 golang的诞生时间,固定写法;*/
    timeStr := time.Now().Format("2006-01-02 15:04:05")
    return timeStr
}

//定义内连键盘菜单;update.CallbackQuery.Data ,返回参数;
var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
        //tgbotapi.NewInlineKeyboardButtonSwitch(test,Switch),
        //tgbotapi.NewInlineKeyboardButtonData("0:开启/关闭SSH服务","0"),
        tgbotapi.NewInlineKeyboardButtonData("0:开启/关闭SSH服务","0"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("1:查看服务器负载情况", "1"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("2:禁止所有用户登录服务器(包括root)", "2"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("3:踢掉所有登录服务器用户(包括root)", "3"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("4:查看服务器运行时间", "4"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("5:查看服务器所有定时任务", "5"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("6:查看所有具有攻击嫌疑的IP", "6"),
    ),
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonURL("播放丹丹大姐姐最美广场舞", "https://www.youtube.com/watch?v=SkskZGoVFvI&list=PLRkymQjXyCD6SnxH4WztCPaG1zMTnNwmx&index=24"),
    ),
)

//检测ssh状态;
func Kiss_Ssh() {
    fmt.Println("执行SSH检测SSH服务状态;")
    //ssh_status := '执行SSH检测程序'
    //return ssh_status
}


func main() {
    bot, err := tgbotapi.NewBotAPI(strings.Replace(Read_Token(), "\n", "", -1))      //取得Token;
    if err != nil {
        log.Panic(err)
    }

    //bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    //循环遍历每个更新;
    for update := range updates {
        //检查我们是否收到消息更新;
        if update.Message != nil {
            //从给定的聊天 ID 构造一条新消息并包含;
            //我们收到的文本;

            /*当输入其他任何消息,机器执行复读机;*/
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
            //fmt.Println("执行复读机")
            //如果消息已打开，请添加我们的数字键盘的副本;
            switch update.Message.Text {
            //输入Open指令: "open":
            case "open":
                msg.ReplyMarkup = numericKeyboard
            }
            //发送这消息;
            if _, err = bot.Send(msg); err != nil {
                panic(err)
            }
        } else if update.CallbackQuery != nil {
            //响应回调查询，告诉 Telegram 向用户展示;
            //收到数据的消息;
            //收到消息响应后,闪烁回复消息,然后消失;
            //callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

            callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

            if _, err := bot.Request(callback); err != nil {
                panic(err)
            }
            
            //内链键盘数据返回;
            switch update.CallbackQuery.Data {
            case "0":
                /*开启或关闭SSH服务;*/;
                chatID := update.CallbackQuery.Message.Chat.ID

                fmt.Println("/*开启或关闭SSH服务;*/")
                fmt.Println(Get_Time_Now)

                //Get ssh服务信息;
                ssh_status := systemd.Auto_Control_Server("sshd.service")
                //Telegram机器人发送ssh服务信息;
                bot.Send(tgbotapi.NewMessage(chatID, Get_Time_Now()+"  "+ssh_status))

                //Get 当前ssh服务配置文件;
                ssh_conf := systemd.Get_Server_Conf("/etc/ssh/sshd_config")     //Get ssh server config;
                //Telegram机器人发送当前ssh服务配置;
                bot.Send(tgbotapi.NewMessage(chatID, "当前ssh服务配置:\n"+ssh_conf))

                //bot.Send(tgbotapi.NewMessage(chatID,"/*开启或关闭SSH服务*/"))
                //Kiss_Ssh()
            case "1":
                /*查看服务器负载情况;*/
                fmt.Println("/*查看服务器负载情况;*/")
            case "2":
                /*禁止所有用户登录服务器(包括root)*/
                fmt.Println("/*禁止所有用户登录服务器(包括root)*/")
            case "3":
                /*踢掉所有登录服务器用户(包括root)*/
                fmt.Println("/*踢掉所有登录服务器用户(包括root)*/")
            case "4":
                /*查看服务器运行时间*/
                fmt.Println("/*查看服务器运行时间*/")
            case "5":
                /*查看服务器所有定时任务*/
                fmt.Println("/*查看服务器所有定时任务*/")
            case "6":
                /*查看所有具有攻击嫌疑的IP*/
                fmt.Println("/*查看所有具有攻击嫌疑的IP*/")
            }

            //收到数据的消息,输出的数据消息(复读机);
            //msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "这是一个测试!!")
            /*msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
            if _, err := bot.Send(msg); err != nil {
                panic(err)
            }*/
        }
    }
}
