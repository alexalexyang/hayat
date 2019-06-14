window.onload = function() {
    var listRooms = document.getElementById('room');
    var inputRoom = document.getElementById('inputRoom');
    var chats = document.getElementById('chat');
    var tabs = document.getElementById('tab');
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
        tabs.innerHTML += `<li><a onclick="channel('${roomid}')">Tab</a></li>`;
    };

    channel = function(roomid) {
        var frames = document.getElementsByClassName("ChannelView");
        var length = frames.length;
        for (var i = 0; i < length; i++) {
            if (frames[i].id == ("viewer" + roomid)) {
                frames[i].style.display = "inline";
            } else { frames[i].style.display = "none"; }
        }
    }

    socket.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        // console.log(event.data)

        for (let i = 0; i < msg.length; i++) {
            if (msg[i].beingserved == false) {
                var roomid = msg[i].roomid
                chats.innerHTML += `<ul id="menu">
                                        <iframe name="${roomid}" frameborder=0 id="viewer${roomid}" class="ChannelView" style="display:none"></iframe>
                                    </ul>`;

                listRooms.innerHTML += `<li id=${roomid} onclick="submitter('${roomid}')"><a target="${roomid}" href="http://localhost:8000/clientprofile/${roomid}">${roomid}</a></li>`;
            } else {
                document.getElementById(msg[i].roomid).remove();
            };
        };
    };
};