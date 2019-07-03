let socket = new WebSocket(`${wsProtocol}${host}/dashboardws`);
let listRooms = document.getElementById('room');
let chats = document.getElementById('chats');
let tabs = document.getElementById('tabs');

socket.onopen = function (event) {
    console.log("Open")
};

socket.onerror = function (error) {
    console.log('WebSocket error: ' + error);
};

submitter = function (roomid, username) {
    document.clientlistForm.inputRoom.value = roomid;
    document.getElementById('clientlistForm').submit();
    chats.insertAdjacentHTML('beforeend', `<iframe id="tabcontent-${roomid}" class="tabcontent" style="display:none" src="${httpProtocol}${host}/clientprofile/${roomid}"></iframe>`);
    tabs.innerHTML += `<button class="tablinks" onclick="showTab(event, '${roomid}')" id="tablink-${roomid}">${username} <img class="closechat" src="static/icons8-close-window-50.png" onclick="closeChat('${roomid}')"></button>`;
};

showTab = function (evt, roomid) {
    var i, tabcontent, tablinks;
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    document.getElementById("tabcontent-" + roomid).style.display = "flex";
    evt.currentTarget.className += " active";
}

closeChat = function (roomid) {
    let confirmClose = confirm("Are you sure you want to close this chat?")
    if (confirmClose == true) {
        document.getElementById("tabcontent-" + roomid).remove();
        document.getElementById("tablink-" + roomid).remove();

        // Send signal down to delete room.
        document.clientlistForm.deleteRoom.value = roomid;
        document.getElementById('clientlistForm').submit();
    }
}

socket.onmessage = function (event) {
    let msg = JSON.parse(event.data);

    if (msg == null) {
        return
    }

    for (let i = 0; i < msg.length; i++) {
        let roomid = msg[i].roomid;
        let username = msg[i].username;
        let consultantName = msg[i].servedby;

        // Show all clients that consultant is already serving.
        if (consultantName.length > 0) {
            chats.insertAdjacentHTML('beforeend', `<iframe id="tabcontent-${roomid}" class="tabcontent" style="display:none" src="${httpProtocol}${host}/chatclient/${roomid}"></iframe>`);
            tabs.innerHTML += `<button class="tablinks" onclick="showTab(event, '${roomid}')" id="tablink-${roomid}">${username} <img class="closechat" src="static/icons8-close-window-50.png" onclick="closeChat('${roomid}')"></button>`;
            document.getElementById(`tablink-${roomid}`).click();
        }

        // Show all clients of the organisation currently waiting to be served.
        if (msg[i].beingserved == false) {
            listRooms.innerHTML += `<li id=${roomid} onclick="submitter('${roomid}', '${username}')">${username}</li>`;
        } else {
            let elem = document.getElementById(msg[i].roomid)
            if (elem != null) {
                elem.remove();
            }
        }
    };

};
