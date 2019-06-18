var request = new XMLHttpRequest();
request.open('GET', '../static/config.json');
request.responseType = 'json';
request.send();

var httpProtocol
var wsProtocol
var host

if (window.location.hostname == "localhost") {
    request.onload = function() {
        httpProtocol = request.response[0].http
        wsProtocol = request.response[0].ws
        host = request.response[0].localhost

        document.clientlistForm.action = `${httpProtocol}${host}/clientlist`
    }
} else {
    request.onload = function() {
        httpProtocol = request.response[0].http
        wsProtocol = request.response[0].ws
        host = request.response[0].host

        document.clientlistForm.action = `${httpProtocol}${host}/clientlist`
    }
}

window.onload = function() {
    var socket = new WebSocket(`${wsProtocol}${host}/clientlistws`);
    var listRooms = document.getElementById('room');
    var inputRoom = document.getElementById('inputRoom');
    var chats = document.getElementById('chat');
    var tabs = document.getElementById('tab');

    socket.onopen = function(event) {
        console.log("Open")
    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ' + error);
    };

    submitter = function(roomid, username) {
        console.log("once or twice")
        document.clientlistForm.inputRoom.value = roomid;
        document.getElementById('clientlistForm').submit();
        chats.insertAdjacentHTML('beforeend', `<iframe name="frame-${roomid}" id="viewer${roomid}" class="ChannelView" style="display:none" src="http://localhost:8000/clientprofile/${roomid}"></iframe>`);
        tabs.innerHTML += `<li><a onclick="channel('${roomid}')">${username}</a></li>`;
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
        console.log(event.data)

        for (let i = 0; i < msg.length; i++) {
            if (msg[i].beingserved == false) {
                var roomid = msg[i].roomid
                var username = msg[i].username
                listRooms.innerHTML += `<li id=${roomid} onclick="submitter('${roomid}', '${username}')">${username}</li>`;
            } else {
                document.getElementById(msg[i].roomid).remove();
            };
        };
    };
};