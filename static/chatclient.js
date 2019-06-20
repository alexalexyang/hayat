// var request = new XMLHttpRequest();
// request.open('GET', '../static/config.json');
// request.responseType = 'json';
// request.send();

var protocol
var host

// if (window.location.hostname == "localhost") {
//     request.onload = function() {
//         protocol = request.response[0].ws
//         host = request.response[0].localhost
//         console.log(protocol)
//         console.log(host)
//     }
// } else {
//     request.onload = function() {
//         protocol = request.response[0].ws
//         host = request.response[0].host
//         console.log(protocol)
//         console.log(host)
//     }
// }

if (window.location.hostname == "localhost") {
    protocol = config["ws"]
    host = config["localhost"]
    console.log(protocol)
    console.log(host)
} else {
    protocol = config["ws"]
    host = config["host"]
    console.log(protocol)
    console.log(host)
}


var pathname = document.location.pathname;
var form = document.getElementById('form-msg');
var txtMsg = document.getElementById('msg');
var listMsgs = document.getElementById('msgs');
var socketStatus = document.getElementById('status');
var btnClose = document.getElementById('close');
// Creating a new WebSocket connection.
// console.log(`${protocol}${host}/chatclientws/` + pathname.split("/")[2])
// var socket = new WebSocket(`${protocol}${host}/chatclientws/` + pathname.split("/")[2]);

console.log(`${protocol}${host}/chatclientws/` + pathname.split("/")[2])
var socket = new WebSocket(`${protocol}${host}/chatclientws/` + pathname.split("/")[2]);
socket.onopen = function(event) {
    // socketStatus.innerHTML = 'Connected to: ' + event.currentTarget.url;
    socketStatus.innerHTML = ''
    socketStatus.className = 'open';
};

socket.onerror = function(error) {
    console.log('WebSocket error: ' + error.message);
};

form.onsubmit = function(e) {
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

socket.onmessage = function(event) {
    var msg = JSON.parse(event.data);
    // console.log(event)
    listMsgs.innerHTML += '<li class="received">' + msg.username + ": " + msg.message + '</li>';
};

socket.onclose = function(event) {
    socketStatus.innerHTML = 'Disconnected from the WebSocket.';
    socketStatus.className = 'closed';
};