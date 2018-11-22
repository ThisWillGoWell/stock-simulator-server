// Vue object for notifications
var vm_notify = new Vue({
  data: {
    notes: {},
  },
});

var RED = "#f44336";
var GREEN = "#1abc9c";
var BLUE = "blue";

var notificationList = $('#notification-list');

var routeNote = {
	trade: notifyTrade,
	send_money: notifySentMoney,
	receive_money: notifyReceiveMoney,
	new_item: notifyNewItem,
};

function sendAck(note_id, callback) {
	var msg = {
		'uuid': note_id
	};
	// TODO: make a callback maybe? 
	doSend('ack', msg);
};

function notifyNewItem(msg) {
	var color = "blue";
	var item = msg.notification.item_type;

	var message = "Bought the " + item + " item";
	var success = true;

	notifyTopBar(message, color, success);

	sendAck(msg.uuid);
};

function notifySentMoney(msg) {
	// Getting usernames
	var receiver = vm_portfolios.portfolios[msg.notification.receiver].name;

	var amount = msg.notification.amount;
	var message = "Successful transfer of $" + formatPrice(amount) + " to " + receiver + ".";
	var color = GREEN;

	console.log(message)
	notifyTopBar(message, color, true);

	sendAck(msg.uuid);
};

function notifyReceiveMoney(msg) {
	// Getting giver name
	var sender = vm_portfolios.portfolios[msg.notification.sender].name;

	var amount = msg.notification.amount;
	var message = "You received $" + formatPrice(amount) + " from " + sender + ".";
	var color = GREEN;
	console.log(message)
	notifyTopBar(message, color, true);

	sendAck(msg.uuid);
};

function notifyTrade(msg) {

	var message;
	var success = msg.notification.success;
	
	// If trade was a success
	if (success) {
		// Vars needed to form note
		var amount = Number(msg.notification.amount);
		var stock_item = vm_stocks.stocks[msg.notification.stock];
		
		if (amount < 0) {
			amount *= -1;
			message = "Successful sale of " + amount + " " + stock_item.ticker_id + " stocks.";
		} else {
			message = "Successful purchase of " + amount + " " + stock_item.ticker_id + " stocks."; 
		}
		notifyTopBar(message, GREEN, success);

	} else {
		message = "Trade failed.";
		notifyTopBar(message, RED, success);
	}

	sendAck(msg.uuid);

};



// function createNotificationHTML(message, color, success, icon) {
	
// 	var notificationObj = {
// 		'message': message,
// 		'color': color,
// 		'success': success,
// 		'icon': icon,
// 	};

// 	var notificationHTML = '<li class="notification-title-bar on-screen" style="background-color: '+ color +';">'+
// 							'	<span>'+ message +'</span>'+
// 							'	<i class="material-icons">'+
// 							'        '+ icon +''+
// 							'	</i>'+
// 							'</li>';
// 	console.log("NOTIFICATION CREATED");
// 	return notificationHTML;
	
// }
var activeNotificationIndex = 0;
var numActiveNotifications = 0;
var activeNotificationList = [];
var timerStart = false;
var timerWaitTime = 2500;

function notifyTopBar(message, color, success) {
	
	timerStart = true;
	activeNotificationIndex++;

	var notificationHTML = '<li class="notification-title-bar" style="background-color: '+ color +';" stackId="'+ activeNotificationIndex +'">'+
							'	<span>'+ message +'</span>'+
							'	<i class="material-icons">'+
							'        done'+
							'	</i>'+
							'</li>';


	$("#notification-list").prepend(notificationHTML);
	
	activeNotificationList.push(activeNotificationIndex);

	setTimeout(function() {

		$('#notification-list li:first').addClass("on-screen");
	
	}, 100);
	
	setTimeout(function() {
		$('#notification-list li:last').removeClass("on-screen");	
		setTimeout(function() {
			$('#notification-list li:last').remove();
		}, 2700);
	
	}, timerWaitTime);
	console.log("NOTIFICATION ADDED");
	

};

// setInterval(function() {
// 	if(timerStart && activeNotificationIndex > 0) {
// 		$('#notification-list li:last').removeClass("on-screen");
	
// 		setTimeout(function() {
			
// 			$('#notification-list li:last').remove();
		
// 		}, 200);
// 		console.log("notification removed");
// 		activeNotificationIndex--;
// 	} else {
// 		timerStart = false;
// 	}
	
	
	
	
	

// }, 1000);



function startStockScroll() {

}


function scrollTopBar(message) {
	
	var mover = d3.select('#scroll-module--container').append('span');
	
	// Set text as the message
	mover.html(message.trim());

	// Getting start and end values
	var parent_bounds = d3.select('#scroll-module--container').node().getBoundingClientRect();
	var parent_left = parent_bounds.left;
	var parent_right = parent_bounds.right;
	var parent_width = parent_bounds.width;

	// Getting element width
	var bounds = mover.node().getBoundingClientRect();
	var el_width = bounds.width;

	//Set up the element and fade in
	mover
		.style('opacity', 0)
		.style('left', ((parent_right - el_width - 5) + 'px'))
		.transition().duration(2000)
		.style('opacity', 1);
	
	// set motion
	mover
		.transition().delay(1800).ease(d3.easeLinear).duration(15000)
		.style('left', (parent_left + 5) + 'px');

	// Fade out
	mover
		.transition().delay(17000).duration(1000)
		.style('opacity', 0);

	// Remove once stopped and faded
	mover
		.transition().delay(18000).duration(0)
		.remove();

};

// Let the page get going and then scroll
setTimeout(function() {
	if (getConfigSetting('ticker')) {
		$("#scroll-module--container").addClass('raiseTicker')
		setInterval(function() {
			var stock = Object.values(vm_stocks.stocks)[Math.floor(Math.random()*Object.values(vm_stocks.stocks).length)];
			scrollTopBar(stock.ticker_id + ": " + formatPrice(stock.current_price));
		
		}, 3000);
	}
}, 5000)
