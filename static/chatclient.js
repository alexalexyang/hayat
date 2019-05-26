window.onload = function() {

    // get the references of the page elements.
    var form = document.getElementById('form-msg');
    var txtMsg = document.getElementById('msg');
    var listMsgs = document.getElementById('msgs');
    var socketStatus = document.getElementById('status');
    var btnClose = document.getElementById('close');

    // Creating a new WebSocket connection.
    var socket = new WebSocket('ws://localhost:8000/chatclientws');

    socket.onopen = function(event) {
        socketStatus.innerHTML = 'Connected to: ' + event.currentTarget.url;
        socketStatus.className = 'open';

    };

    socket.onerror = function(error) {
        console.log('WebSocket error: ' + error);
    };

    form.onsubmit = function(e) {
        e.preventDefault();


        const myObj = {
            username: 'Skip',
            message: txtMsg.value
        };


        // Recovering the message of the textarea.
        var msg = JSON.stringify(myObj);

        // Sending the msg via WebSocket.
        socket.send(msg);

        // Adding the msg in a list of sent messages.
        listMsgs.innerHTML += '<li class="sent"><span>Sent:</span>' + msg + '</li>';

        // Cleaning up the field after sending.
        txtMsg.value = '';

        return false;
    };

    socket.onmessage = function(event) {
        var msg = JSON.parse(event.data);
        console.log(event)
        console.log(msg["username"])
        listMsgs.innerHTML += '<li class="received"><span>Received:</span>' + msg.message + '</li>';
    };

    socket.onclose = function(event) {
        socketStatus.innerHTML = 'Disconnected from the WebSocket.';
        socketStatus.className = 'closed';
    };
};