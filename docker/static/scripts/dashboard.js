$( document ).ready(function() {



	// // Entrance Anims 
	// $('#login-container').addClass('entrance-anim ');

	// var scrollVal = 0;

	// $( window ).on('resize', function(){
	// 	screenWidth = $(window).width();
	//     if (screenWidth > 1018) {
		
	// 	} 
		
	// });

	var sampleMessages = [
		{
	        id: 0,
	        author: 'Matty Ice',
	        timestamp:'11:38am',
	        body:"Hi bb gurl. @Lisa",
	    },
	    {
	        id: 1,
	        author: 'Lisa',
	        timestamp:'11:41am',
	        body:"Matt I told you not to talk dirty to me in this chat. Save it for the DM's when they are finally implemented.",
	    },
	    {
	        id: 2,
	        author: 'Matty Ice',
	        timestamp:'11:44am',
	        body:"Ohh srry bb. I nvr meant to hurt u 💖",
	    },
	    {
	        id: 3,
	        author: 'Andys Woody',
	        timestamp:'11:46am',
	        body:"Lisa, would you like to model for a new Rustangelo painting I'm working on?",
	    },
	    {
	        id: 4,
	        author: 'Lisa',
	        timestamp:'11:51am',
	        body:"Absolutely! want me to come over to your place? xD",
	    },
	];

	$(document).scroll(function() {

		scrollVal = $(document).scrollTop();

	    //console.log("SCROLL: "+scrollVal);


	});

	function formatDate12Hour(date) {
	  let hours = date.getHours();
	  let minutes = date.getMinutes();
	  let ampm = hours >= 12 ? 'pm' : 'am';
	  hours = hours % 12;
	  hours = hours ? hours : 12; // the hour '0' should be '12'
	  minutes = minutes < 10 ? '0'+minutes : minutes;
	  let strTime = hours + ':' + minutes + ' ' + ampm;
	  return strTime;
	}

	var chat_feed = $('#chat-module--container .chat-message--list');
	var debug_feed = $('#debug-module--container .debug-message--list');

	function appendNewMessage(msg, fromMe){

		let isMe = "";
		if (fromMe) {
			isMe = "is-me";
		}
		let msg_text = msg.body;
		let msg_author = msg.author;
		let msg_timestamp = formatDate12Hour(new Date($.now()));
		let msg_template = '<li>'+
				'				<div class="msg-username '+ isMe +'">'+ msg_author +' <span class="msg-timestamp">'+ msg_timestamp +'</span></div>'+
				'				<div class="msg-text">'+ msg_text +'</div>'+
				'			</li>';

		chat_feed.append(msg_template);
		chat_feed.animate({scrollTop: chat_feed.prop("scrollHeight")}, $('#chat-module--container .chat-message--list').height());

	}

	function appendNewServerMessage(msg){

		
		
		let msg_template = '<li>'+			
				'				<div class="msg-text">'+ msg +'</div>'+
				'			</li>';

		debug_feed.append(msg_template);
		debug_feed.animate({scrollTop: chat_feed.prop("scrollHeight")}, $('#chat-module--container .chat-message--list').height());

	}

	

    var i=0;

	setInterval(function() {

		if (i == sampleMessages.length) {
			i = 0;
			chat_feed.empty();
		}

	    appendNewMessage(sampleMessages[i], false);

	    i++;

	}, 5500);

	$('.chat-title-bar button').click(function() {
    
        $('#chat-module--container').toggleClass('closed');
        $('#chat-text-input').focus();
    });

    $('.debug-title-bar button').click(function() {
    
        $('#debug-module--container').toggleClass('closed');
        //$('#debug-text-input').focus();
    });

    $('#top-bar--container .account-settings-btn').click(function() {
    
        $('#top-bar--container .account-settings-menu--container').toggleClass('open');
        
    });

    $('#account-settings-menu-close-btn').click(function() {
    
        $('#top-bar--container .account-settings-menu--container').toggleClass('open');
        
    });

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

	var wsUri = "ws://" + window.location.host + "/ws";
    var output;
    var webSocket;

    function init()
    {
        output = debug_feed;
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
        doSend('{"action": "login", "value": {"username": "username", "password":"password"}}');

        setTimeout(function () {
            //doSend('{"action": "trade", "value": {"stock_ticker": "CHUNT", "exchange_id":"US", "amount":10}}');
        }, 1000);

        setTimeout(function () {

            doSend('{"action": "chat", "value":{"message_body":"hey matty ice"}}')

        }, 1000);



        //doSend('{"op":"subscribe","type":"alert", "system":"irRemote"}');
    }

    function onClose(evt)
    {
        onEvent("Disconnected");
    }

    function onEvent(message){
        appendNewServerMessage(message);
    }

    function onMessage(evt)
    {
        appendNewServerMessage(evt.data);

    }

    function onSend(message)
    {
        appendNewServerMessage(message);

    }

    function onError(evt)
    {
        appendNewServerMessage(evt.data);
    }

    function doSend(message)
    {
        onSend(message)
        webSocket.send(message);
    }

   



    window.addEventListener("load", init, false);

	

});