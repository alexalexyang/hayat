let httpProtocol
let wsProtocol
let host

if (window.location.hostname == "localhost") {
    wsProtocol = config["ws"]
    httpProtocol = config["http"]
    host = config["localhost"]
} else {
    wsProtocol = config["wss"]
    httpProtocol = config["https"]
    host = config["host"]
}

window.onload = function () {
    let socket = new WebSocket(`${wsProtocol}${host}/dashboardws`);
    let listRooms = document.getElementById('room');
    let inputRoom = document.getElementById('inputRoom');
    let chats = document.getElementById('chat');
    let tabs = document.getElementById('tab');

    socket.onopen = function (event) {
        console.log("Open")
    };

    socket.onerror = function (error) {
        console.log('WebSocket error: ' + error);
    };

    submitter = function (roomid, username) {
        document.clientlistForm.inputRoom.value = roomid;
        document.getElementById('clientlistForm').submit();
        chats.insertAdjacentHTML('beforeend', `<iframe name="frame-${roomid}" id="viewer${roomid}" class="ChannelView" style="display:none" src="${httpProtocol}${host}/clientprofile/${roomid}"></iframe>`);
        tabs.innerHTML += `<li><a onclick="channel('${roomid}')">${username}</a></li>`;
    };

    channel = function (roomid) {
        let frames = document.getElementsByClassName("ChannelView");
        let length = frames.length;
        for (let i = 0; i < length; i++) {
            if (frames[i].id == ("viewer" + roomid)) {
                frames[i].style.display = "inline";
            } else { frames[i].style.display = "none"; }
        }
    }

    socket.onmessage = function (event) {
        let msg = JSON.parse(event.data);

        for (let i = 0; i < msg.length; i++) {
            if (msg[i].beingserved == false) {
                let roomid = msg[i].roomid
                let username = msg[i].username
                listRooms.innerHTML += `<li id=${roomid} onclick="submitter('${roomid}', '${username}')">${username}</li>`;
            } else {
                document.getElementById(msg[i].roomid).remove();
            };
        };
    };
};