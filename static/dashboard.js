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
        tabs.innerHTML += `<button class="tablinks" onclick="showTab(event, '${roomid}')" id="tablink-${roomid}">${username} <a onclick="closeChat('${roomid}')">x</a></button>`;
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
        document.getElementById("tabcontent-" + roomid).remove();
        document.getElementById("tablink-" + roomid).remove();
        // Send signal down to delete room.
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

            if (consultantName.length > 0) {
                chats.insertAdjacentHTML('beforeend', `<iframe id="tabcontent-${roomid}" class="tabcontent" style="display:none" src="${httpProtocol}${host}/chatclient/${roomid}"></iframe>`);
                tabs.innerHTML += `<button class="tablinks" onclick="showTab(event, '${roomid}')" id="tablink-${roomid}">${username} <a onclick="closeChat('${roomid}')">x</a></button>`;
                document.getElementById(`tablink-${roomid}`).click();
            }

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
};