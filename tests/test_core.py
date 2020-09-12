# -*- coding:utf-8 -*-
import sys
import onetimepass as otp  # pip3 install onetimepass
from urllib import parse, request, error
import demjson  # pip3 install demjson
import xxtea  # pip3 install xxtea-py cffi
import base64
import re
import datetime
import hashlib
import os
import json
from M2Crypto import BIO, RSA  # dnf install python3-m2crypto.x86_64 -y


def getjsonfiledata(encrypt: "检查是否已经获取密钥对" = True):
    """读入配置文件 testconfig.json ，请先配置它，并先执行 test_gettotptoken.py 。"""

    # tlog("读入配置文件 ...")
    f = open("testconfig.json", 'r')
    lines = f.read()
    f.close()
    jsonfiledata = demjson.decode(lines)
    if jsonfiledata["apiver"] == "" or jsonfiledata["url"] == "":
        terr("错误： 'testconfig.json' 配置不完全。")
        exit()
    if encrypt and (jsonfiledata["publickey"] == "" or jsonfiledata["privateKey"] == ""):
        terr("错误： 需要一个初始的密钥对。")
        exit()
    return jsonfiledata


def rsaEncrypt(public_key: "公钥", message: "要加密的信息", showAllInfo=True):
    """RSA 加密"""
    bio = BIO.MemoryBuffer(public_key)
    rsa_pub = RSA.load_pub_key_bio(bio)
    buffer = None
    while message:
        input = message[:245]
        if showAllInfo:
            tlog("正在加密分段 ...")
            tlog(input)
        snidata = rsa_pub.public_encrypt(input, RSA.pkcs1_padding)
        if buffer == None:
            buffer = snidata
        else:
            buffer = buffer+snidata
        message = message[245:]
    ctxt64_pri = base64.b64encode(buffer)
    return ctxt64_pri


def rsaDecrypt(private_key: "私钥", message: "要解密的信息", showAllInfo=True):
    """RSA 解密"""
    bio = BIO.MemoryBuffer(private_key)
    rsa_pri = RSA.load_key_bio(bio)
    buffer = None
    while message:
        input = message[:512]
        if showAllInfo:
            tlog("正在解密分段 ...")
        snidata = rsa_pri.private_decrypt(input, RSA.pkcs1_padding)
        if showAllInfo:
            tlog(snidata)
        if buffer == None:
            buffer = snidata
        else:
            buffer = buffer+snidata
        message = message[512:]
    return buffer


def postarray_p(postUrl: "提交到指定的URL", jsonDataArr: "提交的数据数组", showAllInfo=True):
    """[明文传输]向服务器提交内容并显示返回内容，明文操作"""

    jsonfiledata = getjsonfiledata(False)
    apiverAppidSecret = [jsonfiledata["apiver"], jsonfiledata["apptoken"]]

    if (showAllInfo):
        tlog("传输模式：明文")
        tlog("准备输入的数据 ...")
    tlog(postUrl)
    tlog(jsonDataArr)
    if (showAllInfo):
        tlog("读取 testconfig.json ...")
    totptoken = jsonfiledata["apptoken"]
    if (showAllInfo):
        tlog("插入固定提交信息 ...")
    jsonDataArr["apptoken"] = totptoken
    jsonDataArr["apiver"] = apiverAppidSecret[0]
    postMod = parse.urlencode(jsonDataArr).encode(encoding='utf-8')
    if (showAllInfo):
        tlog(demjson.encode(jsonDataArr))
        tlog("↑ 发送请求:")
        tlog(postMod.decode())
    postReq = request.Request(url=postUrl, data=postMod)
    try:
        postRes = request.urlopen(postReq)
    except error.HTTPError as e:
        terr("错误：HTTP 连接遇到问题！")
        tlog(e)
        tlog("使用 cURL 获取原始数据 ...")
        curlcmd = 'curl -X POST -d "'+postMod.decode()+'" "'+postUrl+'"'
        tlog(curlcmd)
        output = os.popen(curlcmd)
        tlog(output.read())
        sys.exit(1)
    except error.URLError as e:
        terr("错误：网址不正确！")
        tlog(e)
        sys.exit(1)
    postRes = postRes.read()
    postRes = postRes.decode(encoding='utf-8')
    if (showAllInfo):
        tlog("↓ 收到数据:")
        tlog(postRes)
        tlog("JSON 解析 ...")
    try:
        dataarr = demjson.decode(postRes)
    except:
        terr("错误：解密失败。")
        tlog("原始内容：")
        tlog(postRes)
        sys.exit()
    tlog(dataarr)
    tok("完成。")
    return dataarr


def postarray(postUrl: "提交到指定的URL", jsonDataArr: "提交的数据数组", showAllInfo=True, publicKey: "服务器公钥" = None, privateKey: "客户端私钥" = None):
    """[加密传输]向服务器提交内容并显示返回内容，自动处理加密解密"""
    jsonfiledata = getjsonfiledata(True)
    if (showAllInfo):
        tlog("传输模式：加密")
        tlog(postUrl)
    if (showAllInfo):
        tlog("读取 testconfig.json ...")
    if publicKey == None:
        publicKey = jsonfiledata["publickey"]
    if privateKey == None:
        privateKey = jsonfiledata["privateKey"]
    if (showAllInfo):
        tlog("插入固定提交信息 ...")
    if jsonfiledata["apptoken"]:
        jsonDataArr["apptoken"] = jsonfiledata["apptoken"]
    jsonDataArr["apiver"] = jsonfiledata["apiver"]
    if (showAllInfo):
        tlog("准备输入的数据 ...")
    jsondata = json.dumps(jsonDataArr)
    jsondata = str.encode(jsondata)
    if (showAllInfo):
        tlog(jsondata)
    if (showAllInfo):
        tlog("正在加密数据 ...")
    publicKey = str.encode(publicKey)
    postData = {
        'd': rsaEncrypt(publicKey, jsondata, showAllInfo)
    }
    postMod = parse.urlencode(postData).encode(encoding='utf-8')
    if (showAllInfo):
        tlog("↑ 发送请求:")
    if (showAllInfo):
        tlog(jsonDataArr)
    postReq = request.Request(url=postUrl, data=postMod)
    postRes = request.urlopen(postReq)
    postRes = postRes.read()
    if (showAllInfo):
        tlog("↓ 收到数据:")
    postRes = postRes.decode(encoding='utf-8')
    if (showAllInfo):
        tlog(postRes)
    if (showAllInfo):
        tlog("还原 JSON ...")
    postRes = postRes.replace('-', '+').replace('_', '/')
    mod4 = len(postRes) % 4
    if mod4:
        postRes += "===="[0:4-mod4]
    postRes = bytes(postRes, encoding="utf8")
    tlog(postRes.decode())
    if (showAllInfo):
        tlog("解密数据 ...")
    try:
        postRes = base64.b64decode(postRes)
        postRes = rsaDecrypt(privateKey, postRes, showAllInfo)
    except:
        terr("解密不成功。")
        quit()
    if (showAllInfo):
        tlog("检查返回的数据 ...")
    if (showAllInfo):
        tlog(str(postRes, encoding="utf-8"))
    try:
        resArr = json.loads(postRes)
    except:
        terr("返回数据错误。")
        quit()
    if resArr['code'] != 1000000 and resArr['code'] != 1000100:
        terr("返回状态码错误。")
        quit()
    tok("网络操作完成。")
    return resArr


def tlog(loginfo: "信息内容", end='\n'):
    """输出前面带时间的信息"""
    nowtime = datetime.datetime.now().strftime('[%Y-%m-%d %H:%M:%S.%f]')
    print("\033[35m", end='')
    print(nowtime, end='\033[0m ')
    print(loginfo, end=end)


def terr(loginfo: "信息内容"):
    """输出错误"""
    tlog("\033[31m"+loginfo+"\033[0m")


def tok(loginfo: "信息内容"):
    """输出正确"""
    tlog("\033[32m"+loginfo+"\033[0m")


def title(loginfo: "信息内容"):
    """输出标题"""
    tlog("\033[1m"+loginfo.center(40, '=')+"\033[0m")


def instr(alertinfo: "提示用户要输入的内容", isint=False):
    """接收用户输入"""
    tlog("\033[1m"+alertinfo+"\033[4m", '')
    userinput = input()
    print("\033[0m", end='')
    if isint:
        return int(userinput)
    return userinput
