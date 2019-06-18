var page = "/anteroom"

var request = new XMLHttpRequest();
request.open('GET', '../static/config.json');
request.responseType = 'json';
request.send();

if (window.location.hostname == "localhost") {
    request.onload = function() {
        protocol = request.response[0].http
        host = request.response[0].localhost
        document.anteroomForm.action = `${protocol}${host}${page}`
    }
} else {
    request.onload = function() {
        protocol = request.response[0].http
        host = request.response[0].host
        document.anteroomForm.action = `${protocol}${host}${page}`
    }
}