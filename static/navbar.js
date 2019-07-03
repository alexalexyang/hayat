let httpProtocol
let host

if (window.location.hostname == "localhost") {
    httpProtocol = config["http"]
    host = config["localhost"]
} else {
    httpProtocol = config["https"]
    host = config["host"]
}

let menulinks = document.getElementsByClassName("menulinks");

for (let i = 0; i < menulinks.length; i++) {
    splitLink = menulinks[i].getAttribute('href').split('#');
    menulinks[i].href = `${httpProtocol}${host}` + splitLink[1]
}