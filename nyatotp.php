<?php
require_once "src/nyacore.class.php";
require_once 'src/nyatotp.class.php';
$argv = count($_POST) > 0 ? $_POST : $_GET;
if (isset($argv["n"]) && isset($argv["s"])) {
    $nyatotp = new nyatotp();
    $nyatotp->newdevicetotp($argv["n"],$argv["s"]);
}
?>