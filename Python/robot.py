#!/usr/bin/env python
# -*- coding:utf-8 -*-
# pylint: disable=C0116,W0613

#--Telegram机器人;

import logging
import subprocess
import os
import re
import shutil
import datetime
import time
import socket

from telegram import InlineKeyboardButton, InlineKeyboardMarkup, Update
from telegram.ext import Updater, CommandHandler, CallbackQueryHandler, CallbackContext
from time import sleep

#输入日志;
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.INFO
)
logger = logging.getLogger(__name__)

#获取当前时间;
def Get_Time() -> str:
    return time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())

#读取TOKEN;
def Read_Token() -> str:
    with open('/etc/token', 'r') as f:
        TOKEN = (f.read().strip('\n'))
        return TOKEN

#/start  开启内链键盘;
def start(update: Update, context: CallbackContext) -> None:
    keyboard = [
        #--内连键盘;
        [InlineKeyboardButton("1:开启/关闭SSH服务", callback_data='1'),],
        [InlineKeyboardButton("2:查看服务器负载情况", callback_data='2'),],
        [InlineKeyboardButton("3:禁止所有用户登录服务器(包括root)", callback_data='3'),],
        [InlineKeyboardButton("4:踢掉所有登录服务器用户(包括root)", callback_data='4'),],
        [InlineKeyboardButton("5:查看服务器运行时间", callback_data='5'),],
        [InlineKeyboardButton("6:查看服务器当前所有进程", callback_data='6'),],
        [InlineKeyboardButton("7:查看服务器所有定时任务", callback_data='7'),],
        [InlineKeyboardButton("8:查看具有攻击嫌疑的IP", callback_data='8'),],
        [InlineKeyboardButton("9:播放丹丹姐最美广场舞", callback_data='9'),],
    ]

    reply_markup = InlineKeyboardMarkup(keyboard)

    #update.message.reply_text('Hello Sir,好久不见,你依旧那般年轻,英俊且潇洒.')
    update.message.reply_text('请选择:', reply_markup=reply_markup)

#读取ssh服务配置文件,匹配端口号,通过Systemd获取ssh启动状态;
def Read_ssh_conf() -> str:
    global ssh_conf_file        #ssh配置;
    ssh_conf_file = str(subprocess.check_output("grep -Ev '#|^$' /etc/ssh/sshd_config",shell=True).decode("utf-8"))
    ssh_port = re.search("Port +(6[0-5]{2}[0-3][0-5]|[1-5]\d{4}|[1-9]\d{1,3}|[0-9])",ssh_conf_file).group(0)
    port_int = re.sub("\D", "",ssh_port)
    ssh_status = str(subprocess.check_output("systemctl status sshd |grep -Eo 'running'",shell=True).decode("utf-8"))

    if not ssh_status.strip('\n') == 'running':
        print("ssh服务没有运行!!!")
        ssh_status = '!!SSH服务没有运行!!'

    return (ssh_port,port_int,ssh_status)

def button(update: Update, context: CallbackContext) -> None:
    '''解析 CallbackQuery 并更新消息文本'''
    query = update.callback_query

    #即使不需要通知用户,也需要回答回调查询;
    query.answer()

    if query.data == '1':
        print("开启/关闭SSH服务")
        #query.edit_message_text(text=f"Selected option: {query.data}")
        out_port,ssh_port,ssh_status = Read_ssh_conf()             #获取ssh端口,和配置;
        #--SSH状态返回;
        query.edit_message_text(
                Get_Time()+"\n-----SSH端口:-----\n"+
                out_port+"\n\n"+
                Get_Time()+"\n-----SSH配置-----:\n"+
                ssh_conf_file+"\n"+
                Get_Time()+"\n-----SSH服务当前状态-----:\n"+
                Get_Time()+'    '+ssh_status
                )
        context.bot.send_message(chat_id=update.effective_chat.id, text="主人!!夜已深,为了你的颜值,要注意休息哦")
        #query.edit_message_text(ssh_msg)
    elif query.data == '2':
        print("use 查看服务器负载情况")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '3':
        print("use 禁止所有用户登录服务器(包括root)")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '4':
        print("use 踢掉所有登录服务器用户(包括root)")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '5':
        print("use 查看服务器运行时间")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '6':
        print("use 查看服务器当前所有进程")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '7':
        print("use 查看服务器所有定时任务")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '8':
        print("use 查看具有攻击嫌疑的IP")
        query.edit_message_text(text=f"Selected option: {query.data}")
    elif query.data == '9':
        print("use 播放丹丹姐最美广场舞")
        query.edit_message_text(text=f"Selected option: {query.data}")


#/help    取得帮助;
def help_command(update: Update, context: CallbackContext) -> None:
    """Displays info on how to use the bot."""
    update.message.reply_text("Use /start to test this bot.")


def main() -> None:
    """Run the bot."""
    #--取得Token;
    updater = Updater(token = Read_Token())

    updater.dispatcher.add_handler(CommandHandler('start', start))
    updater.dispatcher.add_handler(CallbackQueryHandler(button))

    updater.dispatcher.add_handler(CommandHandler('help', help_command))

    # Start the Bot
    updater.start_polling()

    # Run the bot until the user presses Ctrl-C or the process receives SIGINT,
    # SIGTERM or SIGABRT
    updater.idle()


if __name__ == '__main__':
    main()
