let socket = new WebSocket(`${wsProtocol}${host}/dashboardws`);
let listRooms = document.getElementById('room');
let chats = document.getElementById('chats');
let clientprofiles = document.getElementById('clientprofiles');
let tabs = document.getElementById('tabs');
document.clientlistForm.inputRoom.action = `${httpProtocol}${host}/dashboard`

socket.onopen = function (event) {
    console.log("Open")
};

socket.onerror = function (error) {
    console.log('WebSocket error: ' + error);
};

submitter = function (roomid, username) {
    console.log(document.clientlistForm.inputRoom.action)
    document.clientlistForm.inputRoom.value = roomid;
    document.clientlistForm.deleteRoom.value = "";
    console.log(document.clientlistForm.inputRoom.value)
    document.getElementById('clientlistForm').submit();
    insertChats(roomid, username);
};

insertChats = function (roomid, username) {
    chats.insertAdjacentHTML('beforeend', `<iframe id="tabbedchat-${roomid}" class="tabbedchat" style="display:none" src="${httpProtocol}${host}/chatclient/${roomid}"></iframe>`);
    clientprofiles.insertAdjacentHTML('beforeend', `<iframe id="tabbedprofile-${roomid}" class="tabbedchat" style="display:none" src="${httpProtocol}${host}/clientprofile/${roomid}"></iframe>`);
    tabs.innerHTML += `<button class="tablinks" onclick="showTab(event, '${roomid}')" id="tablink-${roomid}">${username} <img class="closechat" src="static/icons8-close-window-50.png" onclick="closeChat('${roomid}')"></button>`;
}

showTab = function (evt, roomid) {
    var i, tabbedchat, tablinks;
    tabbedchat = document.getElementsByClassName("tabbedchat");
    for (i = 0; i < tabbedchat.length; i++) {
        tabbedchat[i].style.display = "none";
    }
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    document.getElementById("tabbedchat-" + roomid).style.display = "flex";
    document.getElementById("tabbedprofile-" + roomid).style.display = "flex";
    evt.currentTarget.className += " active";
}

closeChat = function (roomid) {
    let confirmClose = confirm("Are you sure you want to close this chat?")
    if (confirmClose == true) {
        document.getElementById("tabbedchat-" + roomid).remove();
        document.getElementById("tablink-" + roomid).remove();

        // Send signal down to delete room.
        document.clientlistForm.deleteRoom.value = roomid;
        document.clientlistForm.inputRoom.value = "";
        document.getElementById('clientlistForm').submit();
    }
}

socket.onmessage = function (event) {
    let msg = JSON.parse(event.data);
    console.log(msg)

    if (msg == null) {
        return
    }

    for (let i = 0; i < msg.length; i++) {
        let roomid = msg[i].roomid;
        let username = msg[i].username;
        let consultantName = msg[i].servedby;

        // Show all clients that consultant is already serving.
        if (consultantName.length > 0) {
            console.log("Triggering.")
            insertChats(roomid, username);
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