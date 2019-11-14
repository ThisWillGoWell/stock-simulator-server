// ADDED THIS BLOCK FOR AUTH - TELL JAKE
var token = sessionStorage.getItem('token');
var auth_uuid = sessionStorage.getItem('uuid');

var REQUESTS = {};
var REQUEST_ID = 1;

if (token) {
	/*  WEBSOCKETS */
	// var externalServer = "mockstarket.com/api/ws";
	// if(cururl.contains("dev")) {
	// 	externalServer = "dev.mockstarket.com/api/ws";
	// }
	//
	// let host = window.location.hostname;
	// let port = window.location.port;
	//
	// host = " wss://mockstarket.com/api/ws";
	// if(host.includes("localhost")){
	// 	host = "dev.mockstarket.com"
	// } else if(host.includes("home.")){
	// 	host = host.replace("home.", "")
	// }
	//
	// host = "mockstarket.com";
	wsUri = "wss://mockstarket.com/api/ws";

	if( window.location.host.includes("localhost") ){
		if(window.location.port === "8080"){
			// force use localhost
			wsUri = "ws://localhost:8000/api/ws"
		}else if(window.location.port === "8081"){
			// force use dev
			wsUri = "wss://dev.mockstarket.com/api/ws"
		} else if (window.location.port === "8082"){
			// force use prod
			wsUri = "wss://mockstarket.com/api/ws"
		}
	}else if(window.location.host.includes("dev")){
		url = "wss://dev.mockstarket.com"
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
		notifyTopBar("DISCONNECTED FROM WS", RED, false);
		onEvent("Disconnected");
		// Uncomment underneath line for return to login on disconnect
		// window.location.href = "/login.html";
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
				//console.error(err);
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
		delete REQUESTS[msg.request_id];
	});

} else {

	window.location.href = "/login.html";

}