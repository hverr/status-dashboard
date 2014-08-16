<?php
header('Access-Control-Allow-Origin: http://dash.zelus.local');

class Memory {
    public $key = 'memory';

    public function getData($args=array()) {
        exec(
            "free -tmo | awk '{print $2 \",\" $3}'",
            $result
        );

        $a = explode(',', $result[1]);
        $total = intval($a[0]);
        $used = intval($a[1]);

        return array(
            'total' => $total,
            'used' => $used,
        );
    }
}

class Load {
    public $key = 'load';

    public function getData($args=array()) {
        exec("/bin/grep -c ^processor /proc/cpuinfo", $cpuinfo);
        $cores = intval($cpuinfo[0]);

        exec("awk '{print $1\",\"$2\",\"$3}' /proc/loadavg", $loadavg);
        $a = explode(',', $loadavg[0]);

        return array(
            'cores' => $cores,
            'load' => array_map(floatval, $a),
        );
    }
}

class UpTime {
    public $key = 'uptime';

    function getData($args=array()) {
        $uptime = explode('.', file_get_contents('/proc/uptime'));
        $secs = $uptime[0];

        $days = (int)($secs/60/60/24);
        $secs -= $days*60*60*24;

        $hours = (int)($secs/60/60);
        $secs -= $hours*60*60;

        $mins = (int)($secs/60);
        $secs -= $mins*60;

        return array(
            'days' => $days,
            'hours' => $hours,
            'minutes' => $mins,
            'seconds' => $secs,
        );
    }
}

$classes = ['Memory', 'Load', 'UpTime'];
$result = array();
foreach($classes as $class) {
    $mod = new $class();
    $data = $mod->getData();
    if($data) {
        $result[$mod->key] = $data;
    }
}

echo json_encode($result);

