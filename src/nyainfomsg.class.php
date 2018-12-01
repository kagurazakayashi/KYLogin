<?php
class nyainfomsg {
    public $imsg = array(
        /*
        ABBCCDD
        A: 1成功 2失败
        BB: 模块，例如「安全类」
        CC: 错误类型
        DD: 详细错误
        */
        // A=1 : 操作成功执行
        // A=1/BB=00 : 通用成功类型
        // A=1/BB=00/CC=00 : 通用成功
        // A=1/BB=00/CC=00/DD=00 :
        1000000 => '执行成功。',
        // A=1/BB=01 : 数据库类
        // A=1/BB=01/CC=00 : 数据库相关
        // A=1/BB=01/CC=00/DD=00 :
        1010000 => 'SQL语句成功执行。',
        // A=1/BB=01/CC=00/DD=01 :
        1010001 => 'SQL语句成功执行，但没有查询到数据。',
        // A=2 : 操作出现问题
        // A=2/BB=00 : 通用
        // A=2/BB=00/CC=00 : 未知问题
        // A=2/BB=00/CC=00/DD=00 :
        2000000 => '出现未知错误。',
        // A=2/BB=00/CC=01 : 参数相关
        // A=2/BB=00/CC=01/DD=00 :
        2000100 => '没有参数。',
        // A=2/BB=00/CC=01/DD=01 :
        2000101 => '需要更多参数。',
        // A=2/BB=00/CC=01/DD=02 :
        2000102 => '参数无效。',
        // A=2/BB=00/CC=01/DD=03 :
        2000103 => '不支持的提交方式。',
        // A=2/BB=01 : 数据库类
        // A=2/BB=01/CC=01 : MySQL 数据库连接
        // A=2/BB=01/CC=01/DD=00 :
        2010100 => '未能连接到数据库。',
        // A=2/BB=01/CC=01/DD=01 :
        2010101 => '数据库错误。',
        // A=2/BB=01/CC=01/DD=02 :
        2010102 => '数据库未能返回正确的数据。',
        // A=2/BB=01/CC=01/DD=03 :
        2010103 => '缺少数据库配置。',
        // A=2/BB=01/CC=02 : Redis 数据库连接
        // A=2/BB=01/CC=02/DD=00 :
        2010200 => 'Redis 数据库插件初始化失败',
        // A=2/BB=01/CC=02/DD=01 :
        2010201 => 'Redis 数据库连接失败',
        // A=2/BB=02 : 安全类
        // A=2/BB=02/CC=01 : 字符串规则检查
        // A=2/BB=02/CC=01/DD=00 :
        2020100 => '无效字符串。',
        // A=2/BB=02/CC=01/DD=01 :
        2020101 => '字符格式不正确。',
        // A=2/BB=02/CC=01/DD=02 :
        2020102 => 'SQL语句不正确。',
        // A=2/BB=02/CC=01/DD=03 :
        2020103 => '不应包含HTML代码。',
        // A=2/BB=02/CC=02 : 字符串格式检查
        // A=2/BB=02/CC=02/DD=01 :
        2020201 => '不是有效的电子邮件地址。',
        // A=2/BB=02/CC=02/DD=02 :
        2020202 => '不是有效的 IPv4 地址。',
        // A=2/BB=02/CC=02/DD=03 :
        2020203 => '不是有效的 IPv6 地址。',
        // A=2/BB=02/CC=02/DD=04 :
        2020204 => '不是有效的整数数字。',
        // A=2/BB=02/CC=02/DD=05 :
        2020205 => '不是有效的中国电话号码。',
        // A=2/BB=02/CC=03 : 合规性检查
        // A=2/BB=02/CC=03/DD=00 :
        2020300 => '包含违禁词汇。',
        // A=2/BB=02/CC=04 : 加密通信和IP验证
        // A=2/BB=02/CC=04/DD=00 :
        2020400 => '不正确的参数',
        // A=2/BB=02/CC=04/DD=01 :
        2020401 => '无效的 app_id 或 app_secret。',
        // A=2/BB=02/CC=04/DD=02 :
        2020402 => '无法验证IP地址',
        // A=2/BB=02/CC=04/DD=03 :
        2020403 => 'IP地址处于封禁状态',
        // A=2/BB=02/CC=04/DD=04 :
        2020404 => '写入历史记录失败',
        // A=2/BB=02/CC=04/DD=05 :
        2020405 => '重置加密传输失败',
        // A=2/BB=02/CC=04/DD=06 :
        2020406 => '创建加密过程失败',
        // A=2/BB=02/CC=04/DD=07 :
        2020407 => '访问过于频繁',
        // A=2/BB=02/CC=04/DD=08 :
        2020408 => '用于解密的参数不正确',
        // A=2/BB=02/CC=04/DD=09 :
        2020409 => 'apptoken不正确',
        // A=2/BB=02/CC=04/DD=10 :
        2020410 => 'json解析失败',
        // A=2/BB=02/CC=04/DD=11 :
        2020411 => '加密json解析失败',
        // A=2/BB=02/CC=05 : 应用验证
        // A=2/BB=02/CC=05/DD=00 :
        2020500 => '此应用不可用',
    );
    /**
     * @description: 创建异常信息提示JSON
     * @param Int code 错误代码
     * @param Bool showmsg 是否显示错误信息（否则直接无输出）
     * @param String str 附加错误信息
     * @return String 异常信息提示JSON
     */
    function m($code = -1,$showmsg = true,$str = "") {
        if (!$showmsg) return null;
        return json_encode(array(
            "code" => $code,
            "msg" => $this->imsg[$code],
            "info" => $str
        ));
    }
    /**
     * @description: 返回信息的同时，抛出403错误，结束程序
     * @param Int code 错误代码
     * @param Bool showmsg 是否显示错误信息（否则直接403）
     */
    function http403($code=null,$showmsg=true) {
        // header('HTTP/1.1 403 Forbidden');
        if ($code && $showmsg) {
            global $nlcore;
            $json = $nlcore->msg->m($code);
            header('Content-Type:application/json;charset=utf-8');
            echo $json;
        }
        die();
    }
    function __destruct() {
        $this->imsg = null;
        unset($this->imsg);
    }
}
?>