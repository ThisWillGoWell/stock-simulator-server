// ADDED THIS BLOCK FOR AUTH - TELL JAKE
var token = sessionStorage.getItem('token');
var auth_uuid = sessionStorage.getItem('uuid');

var REQUESTS = {};
var REQUEST_ID = 1;

if (token) {
	wsUri = "wss://mockstarket.com/api/ws";

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
		notifyTopBar("DISCONNECTED FROM WS", RED, false);
		onEvent("Disconnected");
	};

	function onEvent(message){
	    // console.log(message);
	};

	function onMessage(evt) {
	    var msg = JSON.parse(evt.data);
		
		if (msg.action) {
			try {
				router[msg.action](msg);
			} catch (err) {
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
		delete REQUESTS[msg.request_id];
	});

} else {

	window.location.href = "/login.html";

}