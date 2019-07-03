let pathname = document.location.pathname.split("/")[2];
let form = document.getElementById('form-msg');
let txtMsg = document.getElementById('msg');
let listMsgs = document.getElementById('msgs');
let socketStatus = document.getElementById('status');
let btnClose = document.getElementById('close');

// console.log(`${wsProtocol}${host}/chatclientws/${pathname}`)
let socket = new WebSocket(`${wsProtocol}${host}/chatclientws/${pathname}`);
socket.onopen = function (event) {
    // socketStatus.innerHTML = 'Connected to: ' + event.currentTarget.url;
    socketStatus.innerHTML = ''
    socketStatus.className = 'open';
};

socket.onerror = function (error) {
    console.log('WebSocket error: ' + error.message);
};

form.onsubmit = function (e) {
    e.preventDefault();
    const myObj = {
        message: txtMsg.value
    };

    // Recovering the message of the textarea.
    let msg = JSON.stringify(myObj);

    // Sending the msg via WebSocket.
    socket.send(msg);

    // Cleaning up the field after sending.
    txtMsg.value = '';

    return false;
};

socket.onmessage = function (event) {
    let msg = JSON.parse(event.data);
    // console.log(msg)
    for (i = 0; i < msg.length; i++) {
        switch (msg[i].type) {
            case "open":
                listMsgs.innerHTML += `<li class="received"><i>${msg[i].username} has connected.</i></li>`;
                break
            case "close":
                listMsgs.innerHTML += `<li class="received"><i>${msg[i].username} has disconnected.</i></li>`;
                break
            default:
                listMsgs.innerHTML += `<li class="received">${msg[i].username}: ${msg[i].message}</li>`;
        }
    }
};

socket.onclose = function (event) {
    socketStatus.innerHTML = 'Disconnected from the WebSocket.';
    socketStatus.className = 'closed';
};