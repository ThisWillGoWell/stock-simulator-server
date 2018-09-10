$( document ).ready(function() {

	let sampleMessages = [
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


	let vm_nav = new Vue({
		el: '#nav',
		methods: {
			nav: function (event) {
				console.log("Click on " + event.toElement.innerHTML);
		    }
		}
	});


	let vm_stocks = new Vue({
	  el: '#stock-list',
	  data: {
	    stocks: {},
	    sortBy: 'ticker_id',
	    sortDesc: 1,
	  },
	  methods: {
		    formatPrice: function(value) {
		        let val = (value/1).toFixed(2)
		        return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".")
		    },
		    // on column name clicks
		    sortCol(col) {
		    	// If sorting by selected column
		    	if (this.sortBy == col) {
		    		// Change sort direction
		    		vm_stocks.sortDesc = -vm_stocks.sortDesc;
		    	} else {
		    		// Change sorted column
		    		vm_stocks.sortBy = col;
		    		
		    	}
		    },
		    // html uses v-for stock in sortStocks(stocks)
		    sortStocks(data) {
		    	if (Object.keys(data).length === 0) {
		    		return data;
		    	}
	    	  	// Turn to array and sort 
		    	var stock_array = Object.keys(data).map(function(key) {
		    		return data[key];
		    	})
		    	// Sorting array
		    	stock_array = stock_array.sort(function(a,b) {
		    		if (a[vm_stocks.sortBy] > b[vm_stocks.sortBy]) {
		    			return -vm_stocks.sortDesc;
		    		} 
		    		if (a[vm_stocks.sortBy] < b[vm_stocks.sortBy]) {
		    			return vm_stocks.sortDesc;
		    		}
		    		return 0;
		    	})
		    	return stock_array;
		    }
		}
	});

	let vm_ledger = new Vue({
	  el: '#ledger',
	  data: {
	    ledger: {}
	  }
	});

	Vue.component('stock-portfolio', {
		props: ['portfoliouuid'],
		data: function() {
			// Turn to array and sort 
	    	var port_array = Object.keys(vm_ledger.ledger).map(function(key) {
	    		return vm_ledger.ledger[key];
	    	})
	    	console.log(port_array);
	    	console.log(this.portfoliouuid);
			return port_array.filter(function(d) {
				return d.uuid == this.portfoliouuid;
			});
		},
		template: `
			<div>
			</div>
		`

	});

	let vm_portfolios = new Vue({
		el: '#portfolios',
		data: {
			portfolios: {},
		}
	});

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

	let chat_feed = $('#chat-module--container .chat-message--list');
	let debug_feed = $('#debug-module--container .debug-message--list');

	/* TODO send a chat message to server */
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
				'				<div class="msg-text">'+ formatted_msg +'</div>'+
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
	let externalServer = "bookingsgolf.com:8000";
	let localServer = window.location.host;
	let wsUri = "wss://"+ externalServer + "/ws";
    let output;
    let webSocket;

    function init()
    {
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
    };

    function onMessage(evt)
    {
        var msg = JSON.parse(evt.data);

    	var router = {
    		'login': routeLogin,
    		'object': routeObject,
    		'update': routeUpdate,
    		'alert': alertUpdate,
    	}
    	router[msg.action](msg);
    };

    var routeLogin = function(msg) {
		console.log("login recieved");
    };

    var routeObject = function(msg) {
		switch (msg.msg.type) {
			case 'portfolio':
			    Vue.set(vm_portfolios.portfolios, msg.msg.uuid, msg.msg.object);
			    break;

			case 'stock':
		  		// Add variables for stocks for vue module initialization 
		  		msg.msg.object.change = 0;
		  		// New - cannot add to vm_stocks.stocks directly (https://vuejs.org/v2/guide/reactivity.html in Change Detection Caveats section)
		  		Vue.set(vm_stocks.stocks, msg.msg.uuid, msg.msg.object);
			  	break;

			case 'ledger':
				/* Add owner names to portfolio uuid */
			    Vue.set(vm_ledger.ledger, msg.msg.uuid, msg.msg.object);
			    break;
		}
    };
    // TODO remove later debugging purposes only
    let logDataOnce = true;

    var routeUpdate = function(msg) {
    	/* ledgers or portfolios. ledgers can build the portfolio object */
    	if (logDataOnce) {
    		console.log(vm_stocks.stocks);
    		console.log(vm_ledger.ledger);
    		console.log(vm_portfolios.portfolios);
    		logDataOnce = false;
    	}
		switch (msg.msg.type) {
			case 'stock':
				// Variables needed to update the stocks
				var targetUUID = msg.msg.uuid;
				var targetField = msg.msg.changes[0].field;
				var targetChange = msg.msg.changes[0].value;

				// temp var for calculating price
				var currPrice = vm_stocks.stocks[targetUUID][targetField];
				// Adding change amount
				vm_stocks.stocks[targetUUID].change = Math.round((targetChange - currPrice) * 1000)/1000;
				// Adding new current price
				vm_stocks.stocks[targetUUID][targetField] = targetChange;

				// helper to color rows in the stock table 
				var targetElem = $("tr[uuid=\x22" + targetUUID + "\x22]");
				targetElem.addClass("updated");
			    break;

			case 'ledger':
				// Variables needed to update the ledger item
				var targetUUID = msg.msg.uuid;
				var targetField = msg.msg.changes[0].field;
				var targetChange = msg.msg.changes[0].value;

				// Update ledger item
				vm_ledger.ledger[targetUUID][targetField] = targetChange;

			    break;
		}

		if (msg.msg.type == "stock") {
	
		}
    };

    var alertUpdate = function(msg) {
    	console.log(msg);
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


    

	

});