let menulinks = document.getElementsByClassName("menulinks");

for (let i = 0; i < menulinks.length; i++) {
    splitLink = menulinks[i].getAttribute('href').split('#');
    menulinks[i].href = `${httpProtocol}${host}` + splitLink[1]
}