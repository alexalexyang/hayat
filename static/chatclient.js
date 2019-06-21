var protocol
var host

if (window.location.hostname == "localhost") {
    protocol = config["ws"]
    host = config["localhost"]
} else {
    protocol = config["ws"]
    host = config["host"]
}


var pathname = document.location.pathname.split("/")[2];
var form = document.getElementById('form-msg');
var txtMsg = document.getElementById('msg');
var listMsgs = document.getElementById('msgs');
var socketStatus = document.getElementById('status');
var btnClose = document.getElementById('close');

console.log(`${protocol}${host}/chatclientws/${pathname}`)
var socket = new WebSocket(`${protocol}${host}/chatclientws/${pathname}`);
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
    var msg = JSON.stringify(myObj);

    // Sending the msg via WebSocket.
    socket.send(msg);

    // Cleaning up the field after sending.
    txtMsg.value = '';

    return false;
};

socket.onmessage = function (event) {
    var msg = JSON.parse(event.data);
    // console.log(event)
    listMsgs.innerHTML += '<li class="received">' + msg.username + ": " + msg.message + '</li>';
};

socket.onclose = function (event) {
    socketStatus.innerHTML = 'Disconnected from the WebSocket.';
    socketStatus.className = 'closed';
};