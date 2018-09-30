
loginUrl = window.location.host

websocketUrl =  "ws://" + window.location.host + '/ws'
if(location.protocol == "https") {
    websocketUrl =  "wss://" + window.location.host + '/ws'
}
webSocket = null
objects = {}

function basicAuthEncode(user, password) {
    var token = user + ":" + password;

    // Should i be encoding this value????? does it matter???
    // Base64 Encoding -> btoa
    var hash = btoa(token);

    return "Basic " + hash;
};

function createUser(user, password, nickname) {
    const Http = new XMLHttpRequest();
    Http.open("PUT", url+"/create", false);
    Http.setRequestHeader("Authorization", basicAuthEncode(user, password));
    Http.setRequestHeader("DisplayName", nickname);
    Http.send();

    if (Http.status !== 200) {
        console.error(Http.responseText);
        return null;
    } else {
        sessionStorage.setItem('token', Http.responseText);
        window.location.href = "/";
        return  Http.responseText;
    }
};

function getToken(user, password) {
    const Http = new XMLHttpRequest();
    Http.open("GET", url+"/token", false);
    Http.setRequestHeader("Authorization", basicAuthEncode(user, password));
    Http.send();

    if (Http.status !== 200) {
        console.error(Http.responseText);
        return null;
    } else {
        sessionStorage.setItem('token', Http.responseText);
        window.location.href = "/";
        return  Http.responseText;
    }
};

function connectWebsocket(){
    webSocket =  new WebSocket(wsUri)
    webSocket.onmessage = function(evt) { onMessage(evt) };
}

var routeObject = function(msg) {
    if (!(msg.msg.type in objects)){
        objects[msg.msg.type] = {}
    }
    objects[msg.msg.type][msg.msg.uuid] = msg.msg.object
};

var routeUpdate = function(msg) {
    for (change of msg.msg.changes){
        objects[msg.msg.type][msg.msg.uuid][change.feld] = change.value
    }
}

var onMessageRouter = {
    'update': routeUpdate,
    'object': routeObject,
}

onMessage(evt) {
    var msg = JSON.parse(evt.data);
    onMessageRouter[msg]
}

