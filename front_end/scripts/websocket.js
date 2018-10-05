// ADDED THIS BLOCK FOR AUTH - TELL JAKE
var token = sessionStorage.getItem('token');
var auth_uuid = sessionStorage.getItem('uuid');

if(token) {
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

	function doSend(action, msg, request_id) {
		if (request_id === undefined) {
			var message = {
				'action': action,
				'msg': msg,
			};	
		} else {
			var message = {
				'action': action,
				'msg': msg,
				'request_id': request_id
			};
		}
	    webSocket.send(JSON.stringify(message));
	};

	function registerRoute(route, callback) {
		router[route] = callback;
	};+

	init();

} else {

	window.location.href = "/login.html";

}