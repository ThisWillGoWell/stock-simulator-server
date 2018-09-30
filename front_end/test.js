
loginUrl = window.location.host

websocketUrl =  "ws://" + window.location.host + '/ws'
if(location.protocol == "https") {
    websocketUrl =  "wss://" + window.location.host + '/ws'
}
webSocket = null;
objects = {};
requests = {};
requestId = 0;


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
    webSocket.onmessage = function(evt) {
        var msg = JSON.parse(evt.data);
        onMessageRouter[msg](msg.action);
    };
    webSocket.onopen = function(evt) {
        let token = getToken('Will', 'pass');
        SendWs(newConnectMsg(token))
    };
}

function SendWs(msg){
    webSocket.send(JSON.stringify(msg))
}

const routeObject = function(msg) {
    if (!(msg.msg.type in objects)){
        objects[msg.msg.type] = {}
    }
    objects[msg.msg.type][msg.msg.uuid] = msg.msg.object
};

const routeUpdate = function(msg) {
    for (change of msg.msg.changes){
        objects[msg.msg.type][msg.msg.uuid][change.field] = change.value
    }
};

const responseRouter = function(msg){
    if (msg.request_id in requests){
        requests[msg.request_id](msg);
        delete responses[msg.request_id];
    }
};

var onMessageRouter = {
    'update': routeUpdate,
    'object': routeObject,
    'response': responseRouter,
};

function makeRequest(msg, callback){
    let r = requestId;
    responses[r] = callback;
    SendWs(msg);
    return r;
}

function newBaseMsg(action, msg, requestId=''){
    return {
        'action': action,
        'msg': msg,
        'requestId': requestId
    }
}

function newConnectMsg(token){
    return newBaseMsg('connect', {'token':token})
}

function stockOrderMsg(stockId, amount ){
    return newBaseMsg('order', {'amount': amount, 'stock_id':stockId})
}

