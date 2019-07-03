let config = {
    "host": "hayat.notathoughtexperiment.me",
    "localhost": "localhost:8000",
    "http": "http://",
    "https": "https://",
    "ws": "ws://",
    "wss": "wss://"
}

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

function setPostPath(pathName) {
    let form = document.getElementsByClassName("form");
    if (window.location.hostname == "localhost") {
        form[0].action = `${httpProtocol}${host}${pathName}`;
    } else {
        form[0].action = `${httpProtocol}${host}${pathName}`;
    }
}

