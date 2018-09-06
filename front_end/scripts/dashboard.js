
// ADDED THIS BLOCK FOR AUTH - TELL JAKE
let authenticated = sessionStorage.getItem('authenticated');
let auth_obj = $.parseJSON(sessionStorage.getItem('auth_obj'));
let auth_uid = auth_obj.uid;
let auth_pw = auth_obj.pw;
//let authenticated = sessionStorage.getItem('authenticated');


if(authenticated) {
	// Get saved data from sessionStorage
	$( document ).ready(function() {

		let sampleMessages = [
			{
		        id: 0,
		        author_display_name: 'Matty Ice',
		        timestamp:'11:38am',
		        body:"Hi bb gurl. @Lisa",
		    },
		    {
		        id: 1,
		        author_display_name: 'Lisa',
		        timestamp:'11:41am',
		        body:"Matt I told you not to talk dirty to me in this chat. Save it for the DM's when they are finally implemented.",
		    },
		    {
		        id: 2,
		        author_display_name: 'Matty Ice',
		        timestamp:'11:44am',
		        body:"Ohh srry bb. I nvr meant to hurt u 💖",
		    },
		    {
		        id: 3,
		        author_display_name: 'Andys Woody',
		        timestamp:'11:46am',
		        body:"Lisa, would you like to model for a new Rustangelo painting I'm working on?",
		    },
		    {
		        id: 4,
		        author_display_name: 'Lisa',
		        timestamp:'11:51am',
		        body:"Absolutely! want me to come over to your place? xD",
		    },
		];


		let vm_nav = new Vue({
			el: '#nav',
			methods: {
				nav: function (event) {
					//console.log(event.currentTarget.getAttribute('data-route'));
					route = event.currentTarget.getAttribute('data-route');
					// console.log("Click on " + route);
					
					renderContent(route);
			    
			    }
			}
		});

		
		let vm_popout_menu = new Vue({
			el: '#btn-logout',
			methods: {
				logout: function (event) {
					// delete cookie
					// Get saved data from sessionStorage
					console.log("logout");
					sessionStorage.removeItem('authenticated');
					sessionStorage.removeItem('auth_obj');
					// send back to index
					window.location.href = "/";
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
			        let val = (value/1).toFixed(2)/100
			        return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".")
			    },
			    // on column name clicks
			    sortCol(col) {
			    	// If sorting by selected column
			    	if (this.sortBy == col) {
			    		// Change sort direction
			    		console.log(col);
			    		vm_stocks.sortDesc = -vm_stocks.sortDesc;
			    	} else {
			    		// Change sorted column
			    		vm_stocks.sortBy = col;
			    		
			    	}
			    },
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
		  el: '#ledger-list',
		  data: {
		    ledger: {}
		  }
		});

		let vm_portfolios = new Vue({
			el: '#portfolio-list',
			data: {
				portfolios: {},
			},
			methods: {
				getPortfolioStocks: function(portfolioUUID) {
					// List of all ledger items
					var stocks = Object.keys(vm_ledger.ledger).map(function(key){
						return vm_ledger.ledger[key];
					});
					// Ledger items of interest
					stocks = stocks.filter(function(d) {
						return d.portfolio_id == portfolioUUID;
					});
					// Grabbing additional info from stock objects
					stocks = stocks.map(function(d) {
						d.ticker_id = vm_stocks.stocks[d.stock_id].ticker_id;
						d.stock_name = vm_stocks.stocks[d.stock_id].name;
						d.current_price = vm_stocks.stocks[d.stock_id].current_price;
						return d;
					});
					return stocks;
				}
			}
		});

		let vm_users = new Vue({
		  el: '#user-info-container',
		  data: {
		    activeUsers: {},
		  },
		  methods: {
		  	getCurrentUser : function() {
		  		// Get userUUID of the person that is logged in
		  		var currentUser = sessionStorage.getItem('authenticated');
		  		// Have they been added to the users object yet?
		  		if (this.activeUsers[currentUser]) {
		  			return this.activeUsers[currentUser].display_name;
		  		} else {
		  			return "";
		  		}
		  	}
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

		function appendNewMessage(msg, fromMe){

			let isMe = "";
			if (fromMe) {
				isMe = "is-me";
			}
			let msg_text = msg.body;
			let msg_author_display_name = msg.author_display_name;
			let msg_author_uuid = msg.author_uuid;
			let msg_timestamp = formatDate12Hour(new Date($.now()));
			let msg_template = '<li '+ msg_author_uuid +'>'+
					'				<div class="msg-username '+ isMe +'">'+ msg_author_display_name +' <span class="msg-timestamp">'+ msg_timestamp +'</span></div>'+
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

	    $('thead tr th').click(function(event) {
	    	
	    	// //let targetElem = this.find('.material-icon').first();
	    	// let toggleAsc = false;
	    	// let toggleDsc = false;
	    	
	    	if($(event.currentTarget).find('i').hasClass("shown")) {
	    		$(event.currentTarget).find('i').toggleClass("flipped");
	    		console.log("is asc");
	    	} else {
	    		$('thead tr th i').removeClass("shown");
	    		$(event.currentTarget).find('i').addClass("shown");
	    	}

	    	// if($(event.currentTarget).find('i').hasClass("asc")) {
	    	// 	toggleDsc = true;
	    	// 	console.log("is dsc");
	    	// }

	    	
	    	
	    	// if(toggleAsc) {
	    		
	    	// 	$(event.currentTarget).find('i').addClass("asc");
	    	// } else {
	    	// 	$(event.currentTarget).find('i').addClass("dsc");
	    	// }

	    	// if(toggleDsc) {
	    	// 	$(event.currentTarget).find('i').addClass("dsc");
	    	// } else {
	    	// 	$(event.currentTarget).find('i').addClass("asc");
	    	// }
	        
	        //$('#debug-module--container').toggleClass('visible');
	        
	    });

	    function formatChatMessage(msg) {
			// console.log("from me");
	  //   	console.log(msg);
	    	let timestamp = formatDate12Hour(new Date($.now()));
	    	// let message_body = $('#chat-module--container textarea').val();
	    	let message_body = msg.msg.message_body;
	    	var currentUser = msg.msg.author;
	    	let temp_msg = {
		        author_uuid: currentUser,
		        author_display_name: vm_users.activeUsers[currentUser].display_name,
		        timestamp: timestamp,
		        body: message_body,
		    };
		    appendNewMessage(temp_msg, true);
	    }


	    $(document).keypress(function(e) {
	    	if($('#chat-module--container textarea').val()) {
			    if(e.which == 13) {
			    	let message_body = $('#chat-module--container textarea').val();
				    doSend('{"action": "chat", "msg": {"message_body": "'+message_body+'"}}');
			        
			        $('#chat-module--container textarea').val().replace(/\n/g, "");
			        $('#chat-module--container textarea').val('');
			        return false;

			    }
			}
		});


		/*  WEBSOCKETS */
		let externalServer = "bookingsgolf.com:8000";
		let localServer = window.location.host;
		let wsUri = "ws://"+ externalServer + "/ws";
	    let output;
	    let webSocket;

	    function init()
	    {
	        testWebSocket();

	        console.log(vm_users.activeUsers);
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
	        doSend('{"action": "login", "msg": {"username": "'+ auth_uid +'", "password":"'+ auth_pw +'"}}');
	        // if(authenticated) {

	        // }
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
	        console.log(msg);
	    	var router = {
	    		'login': routeLogin,
	    		'object': routeObject,
	    		'update': routeUpdate,
	    		'alert': alertUpdate,
	    		// TODO chatUpdate
	    		'chat': chatUpdate,
	    	};

	    	
    		if (msg.action) {
    			router[msg.action](msg);
    		} else {
    			console.log("No message action");
    			
    		}
    		
    		
    	
	    };

	    var routeLogin = function(msg) {

	        // if success if true -> set cookie and forward to dashboard
	        console.log(msg.msg.success);
	        console.log("login recieved");

	        if(msg.msg.success) {
	            // Save data to sessionStorage
	            sessionStorage.setItem("authenticated", msg.msg.uuid);
	            // window.location.href = "/dashboard.html";

	        } else {
	            let err_msg = msg.msg.err;
	            console.log(err_msg);
	            // $('.login-err').text("Username or password is incorrect");
	            //$('.login-err').text(err_msg);
	        }   
	        
	        console.log(msg.msg.uuid);
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

				case 'user':
					/* Add owner names to portfolio uuid */
				    Vue.set(vm_users.activeUsers, msg.msg.uuid, msg.msg.object);
				    break;

			}
	    };
	    // TODO remove later debugging purposes only
	    let logDataOnce = true;

	    var routeUpdate = function(msg) {
	    	/* ledgers or portfolios. ledgers can build the portfolio object */
	    	if (logDataOnce) {
	    		console.log("------ STOCKS ------");
	    		console.log(vm_stocks.stocks);
	    		console.log("------ LEDGER ------");
	    		console.log(vm_ledger.ledger);
	    		console.log("------ PORTFOLIOS ------");
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
					var targetChangeElem = $("tr[uuid=\x22" + targetUUID + "\x22] > td.stock-change");

					if((targetChange - currPrice) > 0) {
						targetChangeElem.removeClass("falling");
						targetChangeElem.addClass("rising");
					} else {
						targetChangeElem.removeClass("rising");
						targetChangeElem.addClass("falling");
					}
					// targetElem.addClass("updated");
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

	    var chatUpdate = function(msg) {
	    	
	    	console.log("----- CHAT -----");
	    	console.log(msg);
	    	formatChatMessage(msg);
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

	    

	    $('.modal-card button').click(function() {
	    
	        $('#modal--container').toggleClass('open');
	        
	    });

	    

	    function toggleModal() {
	    	$('#modal--container').toggleClass('open');
	        console.log("modal show");	
	    }

	    function renderContent(route) {
	    	switch (route) {
				case 'dashboard':
				    	console.log("show dashboard");
				    break;

				case 'business':
			  			console.log("show business");
				  	break;

				case 'stocks':
						console.log("show stocks");
				    break;

				case 'investors':
						console.log("show investors");
				    break;

				case 'futures':
						console.log("show futures");
				    break;

				case 'perks':
						console.log("show perks");
				    break;
			}
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

	} else {

	window.location.href = "/";

}

