<?php
    require 'yaloginUserInfo.php';
    require 'yaloginGlobal.php';
    require 'yaloginStatus.php';
    require 'yaloginSQLC.php';

    class yaloginInformation
    {
        private $user;
        private $sqlset;
        private $ysqlc;

        function init()
        {
            $this->user = new YaloginUserInfo();
            $this->ysqlc = new yaloginSQLC();
            $this->sqlset = $this->ysqlc->sqlset;
        }

        /*获取指定的用户资料
        $table 为空使用 $db_user_table
        table 使用别名
        column 列名(逗号分隔)
        $db_safetable
        $db_safecolumn
        */
        function getInformation() {
             $db = "";
             $tablename = "";
             $column = isset($_POST["column"]) ? $_POST["column"] : null;

             if (isset($_POST["db"]) && $_POST["db"] != "") {
                 $db = $this->aliasconv($_POST["db"],1);
                 if ($db == null) {
                     return 13005;
                 }
             } else {
                 $db = $this->sqlset->db_name;
             }

             if (isset($_POST["table"]) && $_POST["table"] != "") {
                 $tablename = $this->aliasconv($_POST["table"],2);
                 if ($tablename == null) {
                     return 13002;
                 }
             } else {
                 $tablename = $this->sqlset->db_user_table;
             }

             if ($column == null) {
                 return 13001;
             }
             $tablename = isset($this->sqlset->db_tablealias[$table]) ? $this->sqlset->db_tablealias[$table] : null;
             if ($tablename == null) {
                 return 13002;
             }
             if (in_array($table,$this->sqlset->db_safetable) == true) {
                 return 13003; //包含禁止查询表
             }
             $columnarr = explode(",",$column);
             $columnarrintersect = array_intersect($columnarr,$this->sqlset->db_safecolumn);
             if ($columnarrintersect != null && count($columnarrintersect) > 0) {
                 return 13004; //包含禁止查询列
             }
             $status = new YaloginStatus();
             $status->init();
             $statusarr = $status->loginuser();
             if ($statusarr["autologinby"] == "fail") {
                 return 90901;
             }
             $userhash = $statusarr["userhash"];
             return $this->subsql($tablename,$columnarr,$table,$db,$userhash);
        }

        function subsql($tablename,$columnarr,$table,$db,$userhash) {
            $sqlstr = "SELECT `";
            $columns = implode('`,`',$columnarr);
            sqlstr = sqlstr.columns."` FROM `".$db."`.`".$table."` WHERE `hash` = '".$userhash."';";
            $result_array = $this->ysqlc->sqlc($sqlcmd,true,false);
            return $result_array;
        }

        //从别名提取名字 mode 1=数据库 2=表 3=列（暂不支持）
        function aliasconv($name,$mode) {
            $alias = null;
            if ($mode == 1) {
                $alias = $this->sqlset->db_dbalias;
            } else if ($mode == 2) {
                $alias = $this->sqlset->db_tablealias;
            } else {
                return null;
            }
            $resultname = isset($alias[$name]) ? $alias[$name] : null;
            return $resultname;
        }

        function echohtml($result_array) {
            $html = '<!doctype html><html><head><meta charset="utf-8"><title>用户信息查看</title></head><body><table border="1"><tbody>';
            while(list($key,$val)= each($result_array)) { 
            	$html = $html.'<tr><th scope="row">'.$key.'</th><td>'.$val.'</td></tr>';
            }
            return $html."</tbody></table></body></html>";
        }
    }
    
    

?>