<?php

// 创建 连接 发送消息 接收响应 关闭连接
$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
socket_connect($socket, '/tmp/keyword_match.sock');
$i = 0;
while ($i++ < 5) {
    $s = json_encode([
        'module' => 'database',
        'data1' => [
            'mode' => 'insert',
            'force' => false,
            'drives' => 'mysql',
            'isTransaction' =>  rand(10, 20) % 2 === 0 ? true : false,
            'sql' => 'sql' . $i,
            'params' => [0 => ['a' => 'b']]
        ]

    ]);
    $l = strlen($s);
    $msg = "length:{$l}\ncontent:{$s}";

    usleep(500000);
    echo 'send ', $i, PHP_EOL;
    socket_send($socket, $msg, strlen($msg), 0);
    $response = socket_read($socket, 1024);
    var_dump($response);
    usleep(100);
}


socket_close($socket);





