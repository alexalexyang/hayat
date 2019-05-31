window.onload = function() {
    var listMsgs = document.getElementById('msgs');
    var socket = new WebSocket('ws://localhost:8000/clientlistws');

    socket.onopen = function(event) {
        console.log("Open")

    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ' + error);
    };

    socket.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        console.log(event.data)

        for (let i = 0; i < msg.length; i++) {
            listMsgs.innerHTML += '<li class=' + msg[i].roomid + '>' + msg[i].roomid + '</li>';
        }

        // listMsgs.innerHTML += '<li class=' + msg.roomid + '>' + msg.roomid + '</li>';
    };
};