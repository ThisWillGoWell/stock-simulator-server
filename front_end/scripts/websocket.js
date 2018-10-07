// ADDED THIS BLOCK FOR AUTH - TELL JAKE
var token = sessionStorage.getItem('token');
var auth_uuid = sessionStorage.getItem('uuid');

var REQUESTS = {};
var REQUEST_ID = 1;

if (token) {
	/*  WEBSOCKETS */
	var externalServer = "mockstarket.com";
	var localServer = window.location.host;
	var wsUri = "wss://"+ externalServer + "/ws";
	var port = location.port
	if(port == "8000"){
		wsUri = "ws://localhost:8000/ws"
	}

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
	                'token': sessionStorage.getItem('token')
	        };
	        doSend('connect', loginMessage);
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
			try {
				router[msg.action](msg);
			} catch (err) {
				console.log(msg);
				console.error(err);
			}
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

	function doSend(action, msg, callback) {
		if (callback === undefined) {
			var callback = function(msg) {
				console.log("no callback supplied")
				console.log(msg)
			}
		}

		REQUESTS[REQUEST_ID] = callback;

		var message = {
			'action': action,
			'msg': msg,
			'request_id': REQUEST_ID.toString()
		};
		
		REQUEST_ID++;
	    webSocket.send(JSON.stringify(message));
	};

	function registerRoute(route, callback) {
		router[route] = callback;
	};

		
	registerRoute("response", function(msg) {
		try {
			REQUESTS[msg.request_id](msg);
		} catch (err) {
			console.error(err);
			console.log("no request_id key for " + JSON.stringify(msg));
			console.log(REQUESTS);
			console.log(REQUEST_ID);
		}
		console.log(msg);
		delete REQUESTS[msg.request_id];
	});

	init();

} else {

	window.location.href = "/login.html";

}