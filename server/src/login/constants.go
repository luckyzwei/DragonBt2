package main

var ModName_Login = "Login"
var ModName_ServerList = "ServerList"

var CMD_Login_Guest = "login_guest"
var CMD_Login_Sdk = "login_sdk"
var CMD_Login_Ios = "login_ios"
var CMD_ServerList = "serverlist"

const PARAM_TIMESTAMP = "ts"
const PARAM_SIGNATURE = "sig"

var RET_OK = 0
var RET_PARAM_ERROR = 1
var RET_PARAM_ERROR_SIGNATURE = 2
var RET_LOGIN_ERROR_PASSWORD = 3
var RET_LOGIN_ERROR_SDK = 4
var RET_LOGIN_ERROR_NETWORKBUSY_SDK = 5
var RET_LOGIN_ERROR_SYSTEMBUSY_SDK = 6
var RET_LOGIN_ERROR_TOKEN_INVALID_SDK = 7
var RET_LOGIN_ERROR_TOKEN_TIMEOUT_SDK = 8

var RET_LOGIN_ERROR_IOS_APPID_SDK = 8
var RET_LOGIN_ERROR_IOS_SIGNATURE_SDK = 8
var RET_LOGIN_ERROR_IOS_TOKEN_SDK = 8
var RET_LOGIN_ERROR_IOS_TOKEN_TIMEOUT_SDK = 8
