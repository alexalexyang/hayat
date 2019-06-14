window.onload = function() {
    var listRooms = document.getElementById('room');
    var inputRoom = document.getElementById('inputRoom');
    var chats = document.getElementById('chat');
    var socket = new WebSocket('ws://localhost:8000/clientlistws');

    socket.onopen = function(event) {
        console.log("Open")
    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ' + error);
    };


    submitter = function(roomid) {
        document.clientlistForm.inputRoom.value = roomid;
        document.getElementById('clientlistForm').submit();
    };

    socket.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        // console.log(event.data)

        for (let i = 0; i < msg.length; i++) {
            if (msg[i].beingserved == false) {
                var roomid = msg[i].roomid
                    // chats.innerHTML += `<ul id="chat_${roomid}" name="${roomid}" style="display: none;"><iframe id=${roomid} name="${roomid}"></iframe></br></ul>`;
                listRooms.innerHTML += `<li id=${roomid} onclick="submitter('${roomid}')"><a target="_blank" href="http://localhost:8000/clientprofile/${roomid}">${roomid}</a></li>`;
            } else {
                document.getElementById(msg[i].roomid).remove();
            };
        };
    };
};