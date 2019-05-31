window.onload = function() {
    var listRooms = document.getElementById('room');
    var inputRoom = document.getElementById('inputRoom');
    var socket = new WebSocket('ws://localhost:8000/clientlistws');

    socket.onopen = function(event) {
        console.log("Open")
    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ' + error);
    };

    submitter = function(roomid) {
        console.log("submitting: " + roomid);
        document.clientlistForm.inputRoom.value = roomid;
        document.getElementById('clientlistForm').submit();
    };

    socket.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        // console.log(event.data)

        for (let i = 0; i < msg.length; i++) {
            if (msg[i].beingserved == false) {
                var roomid = msg[i].roomid
                listRooms.innerHTML += `<li id=${roomid} onclick="submitter('${roomid}')"><a target="_blank" href="http://localhost:8000/chatclient/${roomid}">${roomid}</a></li>`;
            } else {
                document.getElementById(msg[i].roomid).remove();
            };
        };
    };
};