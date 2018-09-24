
// ADDED THIS BLOCK FOR AUTH - TELL JAKE
let authenticated = sessionStorage.getItem('authenticated');
let auth_uid = sessionStorage.getItem('uid');


//let authenticated = sessionStorage.getItem('authenticated');

if(authenticated) {
	
	// Get saved data from sessionStorage
	$( document ).ready(function() {


		/* Highest level Vue data object */
		var vm_stocks = new Vue({
			data: {
				stocks: {}
			}
		});

		var vm_ledger = new Vue({
			data: {
				ledger: {},
			}
		});
		
		var vm_portfolios = new Vue({
			data: {
				portfolios: {},
			}
		});		
		
		var vm_users = new Vue({
			data: {
			  users: {},
			  currentUser: auth_uid,
			},
			methods: {
		  		getCurrentUser: function() {
		  			// Get userUUID of the person that is logged in
		  			var currentUser = sessionStorage.getItem('uuid');
		  			// Have they been added to the users object yet?
		  			if (vm_users.users[currentUser]) {
		  				return vm_users.users[currentUser].display_name;
		  			} else {
		  				return "";
		  			}
		  		}
		  	},
		});
		console.log("----- USERS -----");
		console.log(vm_users.users);
		console.log("------ STOCKS ------");
		console.log(vm_stocks.stocks);
		console.log("------ LEDGER ------");
		console.log(vm_ledger.ledger);
		console.log("------ PORTFOLIOS ------");
		console.log(vm_portfolios.portfolios);



		/* Vues that are used to display data */

		// Vue for sidebar navigation
		let vm_nav = new Vue({
			el: '#nav',
			methods: {
				nav: function (event) {

					let route = event.currentTarget.getAttribute('data-route');
					
					renderContent(route);
			    }
			}
		});

		// Vue for any sidebar data
	    var sidebarCurrUser = new Vue({
	    	el: '#stats--view',
	    	methods: {
	    		toPrice: formatPrice,
	    	},
	    	computed: {
	    		currUserPortfolio: function() {
	    			var currUserUUID = sessionStorage.getItem('uuid');
	    			if (vm_users.users[currUserUUID] !== undefined) {
		    			var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
	    				if (vm_portfolios.portfolios[currUserFolioUUID] !== undefined) {
			    			var folio = vm_portfolios.portfolios[currUserFolioUUID];
							folio.investments = folio.net_worth - folio.wallet;
			    			return folio;
		    			}
		    		}
		    		return {};
		    	}
	    	}
		});
		
	    // Vue for username top right
	    let displayName = new Vue({
	    	el: "#user-info-container",
			computed: {
				userDisplayName: function() {
	    			var currUserUUID = sessionStorage.getItem('uuid');
	    			if (vm_users.users[currUserUUID] !== undefined) {
		    			return vm_users.users[currUserUUID].display_name;
		    		}
		    		return "";
				}
			}
	    });

		// Vue for all options data 
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

		// Vue for all dashboard data
	    var vm_dash_tab = new Vue({
	    	el: '#dashboard--view',
			data: {
			  sortBy: 'amount',
			  sortDesc: 1,
			},
	    	methods: {
	    		toPrice: formatPrice,
	    		// on column name clicks
			    sortCol: function(col) {
					// If sorting by selected column
			    	if (vm_dash_tab.sortBy == col) {
						// Change sort direction
			    		// console.log(col);
			    		vm_dash_tab.sortDesc = -vm_dash_tab.sortDesc;
			    	} else {
						// Change sorted column
			    		vm_dash_tab.sortBy = col;
			    	}
			    }
	    	},
	    	computed: {
				currUserPortfolio: function() {
					var currUserUUID = sessionStorage.getItem('uuid');
	    			if (vm_users.users[currUserUUID] !== undefined) {
						var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
	    				if (vm_portfolios.portfolios[currUserFolioUUID] !== undefined) {
							var folio = vm_portfolios.portfolios[currUserFolioUUID];
							folio.investments = folio.net_worth - folio.wallet;
			    			return folio;
		    			}
		    		}
		    		return {};
		    	},
				currUserStocks: function() {
					var currUserUUID = sessionStorage.getItem('uuid');
					if (vm_users.users[currUserUUID] !== undefined) {
						
						// Current users portfolio uuid
						var portfolio_uuid = vm_users.users[currUserUUID].portfolio_uuid;
						
						// If objects are in ledger
						if (Object.keys(vm_ledger.ledger).length !== 0) {
							
							var ownedStocks = Object.values(vm_ledger.ledger).filter((d) => d.portfolio_id === portfolio_uuid);
							
							// Remove stocks that user owns 0 of
							ownedStocks = ownedStocks.filter(d => d.amount !== 0);
	
							// Augmenting owned stocks
							ownedStocks = ownedStocks.map(function(d) {
								d.stock_ticker = vm_stocks.stocks[d.stock_id].ticker_id;
								d.stock_price = vm_stocks.stocks[d.stock_id].current_price;
								d.stock_value = Number(d.stock_price) * Number(d.amount);
	
								return d;
							})

					    	// Sorting array
					    	ownedStocks = ownedStocks.sort(function(a,b) {
					    		if (a[vm_dash_tab.sortBy] > b[vm_dash_tab.sortBy]) {
					    			return -vm_dash_tab.sortDesc;
					    		} 
					    		if (a[vm_dash_tab.sortBy] < b[vm_dash_tab.sortBy]) {

					    			return vm_dash_tab.sortDesc;
					    		}
					    		return 0;
					    	});
							return ownedStocks;
						}
					}
					return [];
				}
	    	}
	    });

		// Vue for all stocks tab data 
		var vm_stocks_tab = new Vue({
			el: '#stocks--view',
			data: {
			  sortBy: 'ticker_id',
			  sortDesc: 1,
			},
			methods: {
				toPrice: formatPrice,
			    // on column name clicks
			    sortCol: function(col) {
					// If sorting by selected column
			    	if (vm_stocks_tab.sortBy == col) {
						// Change sort direction
			    		// console.log(col);
			    		vm_stocks_tab.sortDesc = -vm_stocks_tab.sortDesc;
			    	} else {
						// Change sorted column
			    		vm_stocks_tab.sortBy = col;
			    		
			    	}
			    },
			},
			computed:{
				sortedStocks: function() {
		    		if (Object.keys(vm_stocks.stocks).length !== 0) {
			    	  	// Turn to array and sort 
						var stock_array = Object.values(vm_stocks.stocks);

				    	// Sorting array
				    	stock_array = stock_array.sort(function(a,b) {
				    		if (a[vm_stocks_tab.sortBy] > b[vm_stocks_tab.sortBy]) {
				    			return -vm_stocks_tab.sortDesc;
				    		}
				    		if (a[vm_stocks_tab.sortBy] < b[vm_stocks_tab.sortBy]) {
				    			return vm_stocks_tab.sortDesc;
				    		}
				    		return 0;
				    	})
				    	return stock_array;
					}
					return [];
				},
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

		// Vue for all investors tab data
		var vm_investors_tab = new Vue({
			el: '#investors--view',
			methods: {
				toPrice: formatPrice,
			},
			computed: {
				investors: function() {
					var investors = Object.values(vm_portfolios.portfolios);
					// List of all ledger items
					var ledgerItems = Object.values(vm_ledger.ledger);

					investors.map(function(d) {
						// Augment investor data
						d.name = vm_users.users[d.user_uuid].display_name;
						// Get all stocks
						d.stocks = ledgerItems.filter(l => (l.portfolio_id === d.uuid) & (l.amount !== 0)); // ledgers can have amount == 0
						// Augment stock data
						d.stocks = d.stocks.map(function(d) {
							d.ticker_id = vm_stocks.stocks[d.stock_id].ticker_id;
							d.stock_name = vm_stocks.stocks[d.stock_id].name;
							d.current_price = vm_stocks.stocks[d.stock_id].current_price;
							d.value = d.current_price * d.amount;
							return d;
						});
						return d;
					})
					return investors;
				},
			}
		});

		Vue.component('investor-card', {
			computed: {
				investor: function() {
					return vm_investors_tab.investors.filter(d => d.user_uuid === this.user_uuid)[0];
				}
			},
			props: ['user_uuid'],
			template: "<div>{{ investor }}</div>"
		});


		var vm_chat = new Vue({
			el: '#chat-module--container',
			data: {
				showingChat: false,
				unreadMessages: false,

			},
			methods: {
				toggleChat: function() {
					this.showingChat = !this.showingChat;
	        		$('#chat-module--container').toggleClass('closed');
	        		$('#chat-text-input').focus();
				}
			},
			watch: {
				unreadMessages: function() {
					// make css changes here to show a notification for unread messages
					if (this.unreadMessages) {
						console.log("unread messages");
					} else {
						console.log("all messages read");
					}
				}

			},
		});


		$(document).scroll(function() {
			scrollVal = $(document).scrollTop();
		});


		let chat_feed = $('#chat-module--container .chat-message--list');
		let debug_feed = $('#debug-module--container .debug-message--list');

		var cleanInput = (input) => {
			return $('<div/>').text(input).html();
		}

		function appendNewMessage(msg, fromMe){
			// if your chat is closed, add notification
			if (!vm_chat.showingChat) {
				vm_chat.unreadMessages = true;
			}

			let isMe = "";
			if (fromMe) {
				isMe = "is-me";
			}
			let msg_text = cleanInput(msg.body);
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


		// NOTE: I DONT THINK THIS IS BEING USED AT ALL
		// function appendNewServerMessage(msg){
			
		// 	let msg_template = '<li>'+			
		// 			'				<div class="msg-text">'+ formatted_msg +'</div>'+
		// 			'			</li>';

		// 	debug_feed.append(msg_template);
		// 	debug_feed.animate({scrollTop: chat_feed.prop("scrollHeight")}, $('#chat-module--container .chat-message--list').height());

		// }


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




		// TODO: find a better spot for this

	    $('table').on('click', 'tr.clickable' , function (event) {
	    	
	    	var ticker_id = $(this).find('.stock-ticker-id').attr('tid');
	    	
	    	console.log("TID: "+ticker_id);
	    	// TODO: get all data elements
		    // var ticker_id = this.getElementsByClassName('stock-ticker-id')[0].innerHTML;

	    	var stock = Object.values(vm_stocks.stocks).filter(d => d.ticker_id === ticker_id)[0];

	    	// Set show modal to true
	    	buySellModal.showModal = true;
	    	buySellModal.stock_uuid = stock.uuid;
		    
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

	        console.log("login recieved");

	        if(!msg.msg.success) {
	            let err_msg = msg.msg.err;
	            console.log(err_msg);
	        }
	        
	    };

	    var routeObject = function(msg) {
			switch (msg.msg.type) {
				case 'portfolio':
					console.log(msg.msg.object)
				    Vue.set(vm_portfolios.portfolios, msg.msg.uuid, msg.msg.object);
				    break;

				case 'stock':
			  		// Add variables for stocks for vue module initialization 
			  		msg.msg.object.change = 0;
			  		Vue.set(vm_stocks.stocks, msg.msg.uuid, msg.msg.object);
				  	break;
					  
				case 'ledger':
			  		Vue.set(vm_ledger.ledger, msg.msg.uuid, msg.msg.object);
				  	break;

				case 'user':
				    Vue.set(vm_users.users, msg.msg.uuid, msg.msg.object);
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
			var targetUUID = msg.msg.uuid;
			msg.msg.changes.forEach(function(changeObject){
				// Variables needed to update the stocks
				var targetField = changeObject.field;
				var targetChange = changeObject.value;
				
				// if value to update is current price, calculate change
				if (targetField === "current_price") {
					// temp var for calculating price
					var currPrice = vm_stocks.stocks[targetUUID][targetField];
					// Adding change amount
					vm_stocks.stocks[targetUUID].change = targetChange - currPrice;
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
				vm_stocks.stocks[targetUUID][targetField] = targetChange;

			})
	    };

	    var ledgerUpdate = function(msg) {

			var targetUUID = msg.msg.uuid;
			msg.msg.changes.forEach(function(changeObject){
				// Variables needed to update the ledger item
				var targetField = changeObject.field;
				var targetChange = changeObject.value;

				// Update ledger item
				vm_ledger.ledger[targetUUID][targetField] = targetChange;
			})
	    };

	    var portfolioUpdate = function(msg) {

			var targetUUID = msg.msg.uuid;
			msg.msg.changes.forEach(function(changeObject){
				// Variables needed to update the ledger item
				var targetField = changeObject.field;
				var targetChange = changeObject.value;

				// Update ledger item
				vm_portfolios.portfolios[targetUUID][targetField] = targetChange;
			})
	    };

	    var userUpdate = function(msg) {

			var targetUUID = msg.msg.uuid;
			msg.msg.changes.forEach(function(changeObject){
				// Variables needed to update the ledger item
				var targetField = changeObject.field;
				var targetChange = changeObject.value;

				// Update ledger item
				vm_users.users[targetUUID][targetField] = targetChange;
			})
	    };


	    var routeAlert = function(msg) {
	    	console.log(msg);
	    }

	    var routeChat = function(msg) {
	    	
	    	console.log("----- CHAT -----");
	    	console.log(msg);
	    	formatChatMessage(msg);
	    }


	    /* Sending trade requests */

	    function sendTrade() {
	    	
	    	// Creating message for the trade request
	    	var options = {
	    		'action': "trade",
	    		'msg': {
	    			'stock_id': buySellModal.stock_uuid,
	    			'amount': buySellModal.buySellAmount,
	    		}
	    	};
	    	// Sending through websocket
	    	console.log("SEND TRADE");
	    	doSend(JSON.stringify(options));

	    	// Reset buy sell amount
	    	buySellModal.buySellAmount = 0;

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


	    function toggleModal() {
	    	$('#modal--container').toggleClass('open');
	        // console.log("modal show");	
	    }


	    // Vue object for the buy and sell modal
	    var buySellModal = new Vue({
	    	el: '#modal--container',
	    	data: {
	    		showModal: false,
	    		buySellAmount: 0,
	    		isBuying: true,
	    		stock_uuid: 'OSRS',
	    	},
	    	methods: {
	    		toPrice: formatPrice,
	    		addAmount: function(amt) {
	    			buySellModal.buySellAmount += amt; 
	    		},
	    		clearAmount: function() {
	    			buySellModal.buySellAmount = 0;
	    		},
	    		determineMax: function() {
	    			if (buySellModal.isBuying) {
	    				buySellModal.buySellAmount = buySellModal.stock.open_shares; 
	    			} else {
	    				//determine current users holdings
	    				let stock = vm_dash_tab.currUserStocks.filter(d => d.stock_id === buySellModal.stock_uuid)[0];
	    				buySellModal.buySellAmount = stock.amount;
	    			}
	    		},
	    		setIsBuying: function(bool) {
	    			// Change buying or selling
	    			buySellModal.isBuying = bool;

	    			// Set styling
	    			if (buySellModal.isBuying) {
			        	$('#calc-btn-buy').addClass("fill");
			        	$('#calc-btn-sell').removeClass("fill");
	    			} else {
			        	$('#calc-btn-sell').addClass("fill");
			        	$('#calc-btn-buy').removeClass("fill");
	    			}
	    		},
	    		submitTrade: function() {
	    			// Change amount depending on buy/sell
	    			if (!buySellModal.isBuying) {
	    				buySellModal.buySellAmount *= -1;
	    			}
	    			sendTrade();
	    		},
	    		closeModal: function(){
	    			toggleModal();
	    			buySellModal.buySellAmount = 0;
	    			buySellModal.showModal = false;
	    			buySellModal.isBuying = true;
	    		}
	    	},
	    	computed: {
	    		stock: function() {

	    			var clickedStock = Object.values(vm_stocks.stocks).filter(d => d.uuid === buySellModal.stock_uuid)[0];
	    			return clickedStock;
	    		}
	    	}
	    });


	    var allViews = $('.view');
	    var dashboardView = $('#dashboard--view');
	    var businessView = $('#business--view');
	    var stocksView = $('#stocks--view');
	    var investorsView = $('#investors--view');
	    var futuresView = $('#futures--view');
	    var storeView = $('#store--view');
	    var currentViewName = $('#current-view');

	    function renderContent(route) {
	    	switch (route) {
				case 'dashboard':
						allViews.removeClass('active');
						dashboardView.addClass('active');
						currentViewName[0].innerHTML = "Dashboard";
				    	// console.log("show dashboard");
				    break;

				case 'business':
						allViews.removeClass('active');
						businessView.addClass('active');
						console.log(currentViewName)
						currentViewName[0].innerHTML = "Business";
			  			// console.log("show business");
				  	break;

				case 'stocks':
						allViews.removeClass('active');
						stocksView.addClass('active');
						currentViewName[0].innerHTML = "Stocks";
						// console.log("show stocks");
				    break;

				case 'investors':
						allViews.removeClass('active');
						investorsView.addClass('active');
						currentViewName[0].innerHTML = "Investors";
						// console.log("show investors");
				    break;

				case 'futures':
						allViews.removeClass('active');
						futuresView.addClass('active');
						currentViewName[0].innerHTML = "Futures";
						// console.log("show futures");
				    break;

				case 'perks':
						allViews.removeClass('active');
						storeView.addClass('active');
						currentViewName[0].innerHTML = "Perks";
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

