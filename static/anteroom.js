console.log("Able to access .js")

var request = new XMLHttpRequest();
request.open('GET', '../static/config.json');
request.responseType = 'json';
request.send();
request.onload = function() {
    host = request.response[0].host
    console.log(host)
    document.anteroomForm.action = `${host}/anteroom`
}