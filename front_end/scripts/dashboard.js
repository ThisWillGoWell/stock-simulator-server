
// ADDED THIS BLOCK FOR AUTH - TELL JAKE
let authenticated = sessionStorage.getItem('authenticated');
let auth_uid = sessionStorage.getItem('uid');

var vm_portfolios, vm_ledger, vm_stocks, vm_users;

//let authenticated = sessionStorage.getItem('authenticated');


if(authenticated) {
	// Get saved data from sessionStorage
	$( document ).ready(function() {


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

		// 1. Data loads into global object
		// 2. Create Vue objects 
		//		-must be once all object are stored
		// 3. Load html
		//		-cant be done until vue objects have been created in case vue methods are called in the html
		var STOCKS = {};
		var vm_stocks = new Vue({
			  el: '#stock-list',
			  data: {
			    stocks: {},
			    sortBy: 'ticker_id',
			    sortDesc: 1,
			  },
			  methods: {
				    formatPrice: function(value) {
				        let val = (value/1).toFixed(2)/100
				        return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
				    },
				    // on column name clicks
				    sortCol: function(col) {
				    	// If sorting by selected column
				    	if (this.sortBy == col) {
				    		// Change sort direction
				    		// console.log(col);
				    		this.sortDesc = -this.sortDesc;
				    	} else {
				    		// Change sorted column
				    		this.sortBy = col;
				    		
				    	}
				    },
				},
				computed:{
					sortedStocks: function() {
			    		if (Object.keys(this.stocks).length === 0) {
				    		return [];
				    	} else {
				    	  	// Turn to array and sort 
					    	var stock_array = Object.values(vm_stocks.stocks).map(function(d){ return d; });

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
				}
			});
			

		console.log("------ STOCKS ------");
		console.log(vm_stocks.stocks);

		var LEDGER = {};
			var vm_ledger = new Vue({
			  el: '#ledger-list',
			  data: {
			    ledger: {},
			  }
			});
			console.log("------ LEDGER ------");
			console.log(vm_ledger.ledger);

		var PORTFOLIOS = {};
		var vm_portfolios = new Vue({
			el: '#portfolio-list',
			data: {
				portfolios: {},
			},
			computed: {
				//TODO IMPLEMENT
				portfolioStocks: function() {

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
		console.log("------ PORTFOLIOS ------");
		console.log(vm_portfolios.portfolios);

		$('#user-info-container .username-text').text(auth_uid)
		var USERS = {};
		var vm_users = new Vue({
		  el: '#user-info-container',
		  data: {
		    users: {},
		    currentUser: auth_uid,
		  },
		  methods: {
		  	getCurrentUser: function() {
		  		// Get userUUID of the person that is logged in
		  		var currentUser = sessionStorage.getItem('uuid');
		  		console.log()
		  		// Have they been added to the users object yet?
		  		if (this.users[currentUser]) {
		  			return this.users[currentUser].display_name;
		  		} else {
		  			return "";
		  		}
		  	}
		  }
		});
		console.log("----- USERS -----");
		console.log(vm_users.users);


		var getHighestStock = function(stocks) {
			stocks = Object.values(stocks).map((d) => d);
			var highestStock = stocks.reduce(function(a, b){ return a.current_price > b.current_price ? a : b });
			return highestStock;
		};
		var getMoverStock = function(stocks) {
			stocks = Object.values(stocks).map((d) => d);
			var mover = stocks.reduce((a, b) => a.change > b.change ? a : b);
			return mover;
		};

		var superlativeStocks = new Vue({
			el: '#stockSuperlatives',
			computed: {
				highestStock: function() {
					if (Object.values(vm_stocks.stocks).length === 0) {
						return "";
					} else {
						stocks = Object.values(vm_stocks.stocks).map((d) => d);
						var highestStock = stocks.reduce(function(a, b){ return a.current_price > b.current_price ? a : b });
						return highestStock.ticker_id;
					}
				},
				mostChange: function() {
					if (Object.values(vm_stocks.stocks).length === 0) {
						return "";
					} else {
						stocks = Object.values(vm_stocks.stocks).map((d) => d);
						var mover = stocks.reduce((a, b) => a.change > b.change ? a : b);
						return mover.ticker_id;
					}
				},
				mostShares: function() {
					if (Object.values(vm_stocks.stocks).length === 0) {
						return "";
					} else {
						stocks = Object.values(vm_stocks.stocks).map((d) => d);
						var mover = stocks.reduce((a, b) => a.open_shares > b.open_shares ? a : b);
						return mover.ticker_id;
					}
				},
			}
		});


	    var currUser = new Vue({
	    	el: '#dashboard--view',
	    	computed: {
	    		currUserPortfolio: function() {
	    			var currUser = sessionStorage.getItem('uuid');
	    			if (vm_users.users[currUser] !== undefined) {
		    			var currUserFolioUUID = USERS[currUser].portfolio_uuid;
	    				if (vm_portfolios.portfolios[currUserFolioUUID] !== undefined) {
			    			var folio = vm_portfolios.portfolios[currUserFolioUUID];
			    			return folio;
		    			}
		    		}
		    		return {};
		    	},
		    	currUserStocks: function() {
		    		var currUser = sessionStorage.getItem('uuid');
	    			if (vm_users.users[currUser] !== undefined) {
		    			
		    			// Current users portfolio uuid
		    			var portfolio_uuid = USERS[currUser].portfolio_uuid;
	    				
	    				// If objects are in ledger
	    				if (Object.keys(vm_ledger.ledger).length !== 0) {
	    					
	    					var ownedStocks = Object.values(LEDGER).filter((d) => d.portfolio_id === portfolio_uuid);
	    					
	    					// Augmenting owned stocks
	    					ownedStocks = ownedStocks.map(function(d) {
	    						d.stock_ticker = STOCKS[d.stock_id].ticker_id;
	    						d.stock_price = STOCKS[d.stock_id].current_price;
	    						d.stock_value = Number(d.stock_price * d.amount);

	    						// Formatting to dollars
	    						d.stock_price = formatPrice(d.stock_price);
	    						d.stock_value = formatPrice(d.stock_value);
	    						return d;
	    					})
	    					return ownedStocks;
	    				}
		    		} else {
		    			return [];
		    		}
		    	}
	    	}
	    });

	    var sidebarCurrUser = new Vue({
	    	el: '#stats--view',
	    	computed: {
	    		currUserPortfolio: function() {
	    			var currUser = sessionStorage.getItem('uuid');
	    			if (vm_users.users[currUser] === undefined) {
	    				return {};
	    			} else {
		    			var currUserFolioUUID = USERS[currUser].portfolio_uuid;
	    				if (vm_portfolios.portfolios[currUserFolioUUID] === undefined) {
	    					return {};
		    			} else {
			    			var folio = vm_portfolios.portfolios[currUserFolioUUID];
			    			return folio;
		    			}
		    		}
		    	}
	    	}
	    });

		
		$(document).scroll(function() {
			scrollVal = $(document).scrollTop();
		});

		function formatPrice(value) {
	        let val = value.toFixed(2)/100
	        return val.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
	    }

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

	    $('.account-settings-btn').click(function() {
	    	// console.log("clicked");
	        $('#top-bar--container .account-settings-menu--container').toggleClass('open');
	        
	    });

	    $('#account-settings-menu-close-btn').click(function() {
	    
	        $('#top-bar--container .account-settings-menu--container').toggleClass('open');
	        
	    });

	    $('.debug-btn').click(function() {
	    
	        $('#debug-module--container').toggleClass('visible');
	        
	    });

	    $('.calc-btn').click(function() {
	    
	        
	        
	    });



		// TODO: find a better spot for this

	    $('#stock-list').on('click', 'tr.clickable' , function (event) {
	    	// TODO: get all data elements
		    var ticker_id = this.getElementsByClassName('stock-ticker-id')[0].innerHTML;
	    	var stock = Object.values(vm_stocks.stocks).filter((d) => d.ticker_id === ticker_id)[0];
		    console.log(stock);
		    var current_price = formatPrice(stock.current_price);
		    
		    //TODO: update all data elements in the modal
		    $('#modal--container .modal-stock-id').html(ticker_id);
		    $('#modal--container .modal-stock-price').html('$' + current_price);

		    toggleModal();
        });



	    // $('#nav li button').click(function(this) {
	    
	    //    $(this).remove();
	        
	    // });

	    $('thead tr th').click(function(event) {
	    	
	    	// //let targetElem = this.find('.material-icon').first();
	    	// let toggleAsc = false;
	    	// let toggleDsc = false;
	    	
	    	if($(event.currentTarget).find('i').hasClass("shown")) {
	    		$(event.currentTarget).find('i').toggleClass("flipped");
	    		// console.log("is asc");
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
	    	let timestamp = formatDate12Hour(new Date($.now()));
	    	// let message_body = $('#chat-module--container textarea').val();
	    	let message_body = msg.msg.message_body;
	    	var currentUser = msg.msg.author;
	    	let temp_msg = {
		        author_uuid: currentUser,
		        author_display_name: vm_users.users[currentUser].display_name,
		        timestamp: timestamp,
		        body: message_body,
		    };
		    appendNewMessage(temp_msg, true);
	    }


	    $(document).keypress(function(e) {
	    	if($('#chat-module--container textarea').val()) {
			    if(e.which == 13) {

			    	let message_body = $('#chat-module--container textarea').val();
					
					var msg = {
						'action': 'chat',
						'msg': {
							'message_body': message_body,
						}	
					};

					doSend(JSON.stringify(msg))
			        
			        $('#chat-module--container textarea').val().replace(/\n/g, "");
			        $('#chat-module--container textarea').val('');
			        return false;

			    }
			}
		});


		/*  WEBSOCKETS */
		let externalServer = "mockstarket.com";
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
	    	if (sessionStorage.getItem('authenticated') !== null) {
	            var loginMessage = {
	                'action': 'renew',
	                'msg': {
	                    'token': sessionStorage.getItem('authenticated')
	                }
	            };
	            doSend(JSON.stringify(loginMessage));
	        }   
			onEvent("Connected");

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
	        // console.log(msg);
	    	var router = {
	    		'login': routeLogin,
	    		'object': routeObject,
	    		'update': routeUpdate,
	    		'alert': routeAlert,
	    		'chat': routeChat,
	    	};
	    	
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

	    var routeLogin = function(msg) {

	        // if success if true -> set cookie and forward to dashboard
	        console.log("login recieved");

	        if(!msg.msg.success) {
	            let err_msg = msg.msg.err;
	            console.log(err_msg);
	        }
	        
	    };

	    var routeObject = function(msg) {
			switch (msg.msg.type) {
				case 'portfolio':
				    PORTFOLIOS[msg.msg.uuid] = msg.msg.object;
			  		// Give the vue object reactivity with PORTFOLIOS
				    Vue.set(vm_portfolios.portfolios, msg.msg.uuid, PORTFOLIOS[msg.msg.uuid]);

				    break;

				case 'stock':
			  		// Add variables for stocks for vue module initialization 
			  		msg.msg.object.change = 0;
			  		// Add object
				    STOCKS[msg.msg.uuid] = msg.msg.object;
			  		// Give the vue object reactivity with STOCKS
			  		Vue.set(vm_stocks.stocks, msg.msg.uuid, STOCKS[msg.msg.uuid]);
				  	break;

				case 'ledger':
				    LEDGER[msg.msg.uuid] = msg.msg.object;
			  		// Give the vue object reactivity with LEDGER
				    Vue.set(vm_ledger.ledger, msg.msg.uuid, LEDGER[msg.msg.uuid]);
				    break;

				case 'user':
					USERS[msg.msg.uuid] = msg.msg.object;
			  		// Give the vue object reactivity with USERS
				    Vue.set(vm_users.users, msg.msg.uuid, USERS[msg.msg.uuid]);
				    // TODO remove
				    // Setting up current user portfolio	    				
				    // if (msg.msg.object.uuid == sessionStorage.getItem('uuid')) {
				    // 	createCurrentUser(msg.msg.object.portfolio_uuid);
				    // }
				    break;

			}
	    };

	    var routeUpdate = function(msg) {
	    	var updateRouter = {
	    		'stock': stockUpdate,
	    		'ledger': ledgerUpdate,
	    		'portfolio': portfolioUpdate,
	    		'user': userUpdate,
	    	};
	    	updateRouter[msg.msg.type](msg);
	    };

	    var stockUpdate = function(msg) {
			msg.msg.changes.forEach(function(changeObject){
				// Variables needed to update the stocks
				var targetUUID = msg.msg.uuid;
				var targetField = changeObject.field;
				var targetChange = changeObject.value;
				
				// if value to update is current price, calculate change
				if (targetField === "current_price") {
					// temp var for calculating price
					var currPrice = STOCKS[targetUUID][targetField];
					// Adding change amount
					STOCKS[targetUUID].change = Math.round((targetChange - currPrice) * 1000)/100000;
					// vm_stocks.stocks[targetUUID].change = Math.round((targetChange - currPrice) * 1000)/100000;
					
					// helper to color rows in the stock table 
					var targetElem = $("tr[uuid=\x22" + targetUUID + "\x22]");
					var targetChangeElem = $("tr[uuid=\x22" + targetUUID + "\x22] > td.stock-change");
					
					if ((targetChange - currPrice) > 0) {
						targetChangeElem.removeClass("falling");
						targetChangeElem.addClass("rising");
					} else {
						targetChangeElem.removeClass("rising");
						targetChangeElem.addClass("falling");
					}
				}

				// Adding new current price
				STOCKS[targetUUID][targetField] = targetChange;


			})
	    };

	    var ledgerUpdate = function(msg) {
			msg.msg.changes.forEach(function(changeObject){

				// Variables needed to update the ledger item
				var targetUUID = msg.msg.uuid;
				var targetField = changeObject.field;
				var targetChange = changeObject.value;

				// Update ledger item
				LEDGER[targetUUID][targetField] = targetChange;
			})
	    };

	    var portfolioUpdate = function(msg) {
			
			msg.msg.changes.forEach(function(changeObject){

				// Variables needed to update the ledger item
				var targetUUID = msg.msg.uuid;
				var targetField = msg.msg.changes[0].field;
				var targetChange = msg.msg.changes[0].value;

				// Update ledger item
				PORTFOLIOS[targetUUID][targetField] = targetChange;
			})
	    };

	    var userUpdate = function(msg) {

			msg.msg.changes.forEach(function(changeObject){
			
				// Variables needed to update the ledger item
				var targetUUID = msg.msg.uuid;
				var targetField = changeObject.field;
				var targetChange = changeObject.value;

				// Update ledger item
				USERS[targetUUID][targetField] = targetChange;
			})
	    }


	    var routeAlert = function(msg) {
	    	console.log(msg);
	    }

	    var routeChat = function(msg) {
	    	
	    	console.log("----- CHAT -----");
	    	console.log(msg);
	    	formatChatMessage(msg);
	    }


	    /* Sending trade requests */

	    document.getElementById("trade-request-submit").addEventListener("click", sendTradeOptions, false);
	    
	    function sendTradeOptions() {
	    	
	    	// Get request parameters
		    
		    var stockTickerId = $('#modal--container .modal-stock-id').html();
	    	var amount = 1; // TODO: get from UI

	    	//Get stockid from ticker
	    	var focusStock = Object.values(vm_stocks.stocks).filter(
	    		function(stock) {
	    			return stock.ticker_id === stockTickerId;
	    		})[0];
	    	// Creating message for the trade request
	    	var options = {
	    		'action': "trade",
	    		'msg': {
	    			'stock_id': focusStock.uuid,
	    			'amount': amount,
	    		}
	    	};
	    	// Sending through websocket
	    	console.log("SEND TRADE");
	    	doSend(JSON.stringify(options));

	    };

	   	/* End sending trade requests */




	    function onError(evt)
	    {
	    	console.log(evt);
	        // writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
	    }

	    function doSend(message)
	    {
	        webSocket.send(message);
	    }



	    $('.modal-card button').click(function() {

	        toggleModal();

	    });


	    function toggleModal() {
	    	$('#modal--container').toggleClass('open');
	        // console.log("modal show");	
	    }

	    var allViews = $('.view');
	    var dashboardView = $('#dashboard--view');
	    var businessView = $('#business--view');
	    var stocksView = $('#stocks--view');
	    var investorsView = $('#investors--view');
	    var futuresView = $('#futures--view');
	    var storeView = $('#store--view');

	    function renderContent(route) {
	    	switch (route) {
				case 'dashboard':
						allViews.removeClass('active');
						dashboardView.addClass('active');
				    	// console.log("show dashboard");
				    break;

				case 'business':
						allViews.removeClass('active');
						businessView.addClass('active');
			  			// console.log("show business");
				  	break;

				case 'stocks':
						allViews.removeClass('active');
						stocksView.addClass('active');
						// console.log("show stocks");
				    break;

				case 'investors':
						allViews.removeClass('active');
						investorsView.addClass('active');
						// console.log("show investors");
				    break;

				case 'futures':
						allViews.removeClass('active');
						futuresView.addClass('active');
						// console.log("show futures");
				    break;

				case 'perks':
						allViews.removeClass('active');
						storeView.addClass('active');
						// console.log("show perks");
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

