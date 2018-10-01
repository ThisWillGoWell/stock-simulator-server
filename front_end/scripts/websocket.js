/*  WEBSOCKETS */
var externalServer = "mockstarket.com";
var localServer = window.location.host;
var wsUri = "wss://"+ externalServer + "/ws";
var output;
var webSocket;

var router = {};


function init() {
    webSocket = new WebSocket(wsUri);
    webSocket.onopen = function(evt) { onOpen(evt) };
    webSocket.onclose = function(evt) { onClose(evt) };
    webSocket.onmessage = function(evt) { onMessage(evt) };
    webSocket.onerror = function(evt) { onError(evt) };
};

function onOpen(evt) {
	if (sessionStorage.getItem('token') !== null) {
        var loginMessage = {
            'action': 'connect',
            'msg': {
                'token': sessionStorage.getItem('token')
            }
        };
        doSend(JSON.stringify(loginMessage));
    }   
	onEvent("Connected");
};

function onClose(evt) {
    onEvent("Disconnected");
};

function onEvent(message){
    console.log(message);
};

function onMessage(evt) {
    var msg = JSON.parse(evt.data);
    // console.log(msg);
	
	if (msg.action) {
		router[msg.action](msg);
	} else {
		if (msg.type == "error") {
			console.log("ERROR")
			console.log(msg);
		}
		console.log("No message action");
		console.log(console.log(msg));	
	}
};

function onError(evt) {
	console.log(evt);
    // writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
};

function doSend(message) {
    webSocket.send(message);
};

function registerRoute(route, callback) {
	router[route] = callback;
};

