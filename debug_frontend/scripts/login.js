$( document ).ready(function() {


	console.log("login js");
	

    $('.debug-btn').click(function() {
    
        $('#debug-module--container').toggleClass('visible');
        
    });

    



    $(document).keypress(function(e) {
    	if($('#chat-module--container textarea').val()) {
		    if(e.which == 13) {
		    	let timestamp = formatDate12Hour(new Date($.now()));
		    	let message_body = $('#chat-module--container textarea').val();
		    	let temp_msg = {
			        id: 0,
			        author: 'Sebio',
			        timestamp: timestamp,
			        body: message_body,
			    };
		        appendNewMessage(temp_msg, true);
		        $('#chat-module--container textarea').val().replace(/\n/g, "");
		        $('#chat-module--container textarea').val('');
		        return false;

		    }
		}
	});


	/*  WEBSOCKETS */

	var wsUri = "ws://"+ window.location.host + "/ws";
    var output;
    var webSocket;

    function init()
    {
        output = document.getElementById("output");
        testWebSocket();
    }

    function testWebSocket()
    {
        webSocket = new WebSocket(wsUri);
        webSocket.onopen = function(evt) { onOpen(evt) };
        webSocket.onclose = function(evt) { onClose(evt) };
        webSocket.onmessage = function(evt) { onMessage(evt) };
        webSocket.onerror = function(evt) { onError(evt) };
    }

    function onOpen(evt)
    {
        onEvent("Connected");
        //doSend('{"container_type": "register", "register_action": "register", "device_type": "test", "device_name":"' + window.prompt("device_name","test")  + '"}');
        doSend('{"action": "login", "msg": {"username": "Will", "password":"pass"}}');

        //doSend('{"op":"subscribe","type":"alert", "system":"irRemote"}');
    }

    function onClose(evt)
    {
        onEvent("Disconnected");
    }

    function onEvent(message){
        // writeToScreen('<span style="color: darkorange;">'+ message+'</span>')
        console.log(message);
    }

    function syntaxHighlight(json) {
        if (typeof json != 'string') {
            json = JSON.stringify(json, undefined, 4);
        }
        json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
            var cls = 'number';
            if (/^"/.test(match)) {
                if (/:$/.test(match)) {
                    cls = 'key';
                } else {
                    cls = 'string';
                }
            } else if (/true|false/.test(match)) {
                cls = 'boolean';
            } else if (/null/.test(match)) {
                cls = 'null';
            }
            return '<span class="' + cls + '">' + match + '</span>';
        });
    }

    function putScren(inp) {
        // document.body.appendChild(document.createElement('pre')).innerHTML = inp;
    }

    function syntaxHighlight(json) {
        // https://stackoverflow.com/questions/4810841/how-can-i-pretty-print-json-using-javascript
        json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
            var cls = 'number';
            if (/^"/.test(match)) {
                if (/:$/.test(match)) {
                    cls = 'key';
                } else {
                    cls = 'string';
                }
            } else if (/true|false/.test(match)) {
                cls = 'boolean';
            } else if (/null/.test(match)) {
                cls = 'null';
            }
            return '<span class="' + cls + '">'  + match + '</span>';
        });
    }

    function onMessage(evt)
    {
        //writeToScreen(syntaxHighlight(evt.data));
        console.log(evt.data);
        var str = JSON.stringify(JSON.parse(evt.data), undefined, 4);
        //putScren(syntaxHighlight(str));
    }

    function onSend(message)
    {
        // writeToScreen('<span style="color: lightblue;">SEND: ' + message +'</span>');

    }

    function onError(evt)
    {
        // writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
    }

    function doSend(message)
    {
        onSend(message)
        webSocket.send(message);
    }

    // function writeToScreen(message)
    // {
    //     var pre = document.createElement("p");
    //     //pre.style.wordWrap = "break-word";
    //     pre.innerHTML = message;
    //     output.appendChild(pre);
    // }



    window.addEventListener("load", init, false);

    function attemptLogin() {
    	
    	let input_uid = $('#login-uid').val();
		let input_pw = $('#login-pw').val();
		let auth_msg = {};

    	if(input_uid != '' && input_pw != '') {
    		auth_msg = {
    					action: "login",
				        value: {
				        	"username": input_uid, 
				        	"password": input_pw
				        }
				    };
    	}



    	console.log(auth_msg);

    	try {
    		onEvent("Connected");
	  		doSend(JSON.stringify(auth_msg));
		}
		catch(error) {
	  		console.error(error);
		  	$('.login-err').text("Username or password is incorrect");
		}
		
    }

   



    window.addEventListener("load", init, false);


    $('#input-login-submit').click(function() {
    
    	attemptLogin();
        console.log("login clicked");
        
    });

	

});