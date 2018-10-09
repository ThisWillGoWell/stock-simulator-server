// Vue object for notifications
var vm_notify = new Vue({
  data: {
    notes: {},
  },
});

var RED = "#f44336";
var GREEN = "#1abc9c";
var BLUE = "blue";


var routeNote = {
	trade: notifyTrade,
	send_money: notifyTransfer,
	receive_money: notifyTransfer,
	new_item: notifyNewItem,
};

function sendAck(note_id, callback) {
	var msg = {
		'uuid': note_id
	};
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

function notifyTransfer(msg) {
	var color, message;
	var success = msg.notification.success;

	// Getting usernames
	var receiver = msg.notification.receiver;
	receiver = vm_users.users[receiver].display_name;

	// If trade was a success
	if (success) {
		// Getting amount 
		var amount = msg.notification.amount;
		message = "Sucessful tranfer of " + formatPrice(amount) + " to " + receiver + ".";
		color = GREEN;

	} else {
		message = "Tranfer to " + receiver + " failed.";
		color = RED;
	}

	notifyTopBar(message, color, success);

	sendAck(msg.uuid);

};

function notifyTrade(msg) {

	var color, message;
	var success = msg.notification.success;
	
	// If trade was a success
	if (success) {
		var amount = Number(msg.notification.amount);
		var tradeType = "";
		if (amount < 0) {
			tradeType = 'sell';
			amount *= -1;
		} else {
			tradeType = 'buy';
		}

		var stock_item = vm_stocks.stocks[msg.notification.stock];

		if (tradeType === 'sell') {
			message = "Successful sale of " + amount + " " + stock_item.ticker_id + ".";
		} else if (tradeType === 'buy') {
			message = "Successful purchase of " + amount + " " + stock_ticker + "."; 
		} else {
			console.error("tradeType not set by server message");
		}

		color = "#1abc9c";

	} else {

		color = "#f44336";
		message = "Trade failed.";
	}

	notifyTopBar(message, color, success);

	sendAck(msg.uuid);

};

function notifyTopBar(message, color, success) {

	// Set text as the message
	d3.select('#notification-module--container span').html(message);

	// Set color and start motion
	d3.select('#notification-module--container')
		.style('background-color', color)
		.transition().duration(300)
		.style('opacity', 1).style('top', '0px');
		

	// Hide and move back up
	d3.select('#notification-module--container')
		.transition().delay(3500).duration(300)
		.style('opacity', 0).style('top', '-60px');

};


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

	//Set up the element
	mover
		.style('opacity', 0)
		.style('left', ((parent_right - el_width - 5) + 'px'))
		.transition().duration(2000)
		.style('opacity', 1);

	mover
		.transition().delay(1800).ease(d3.easeLinear).duration(15000)
		.style('left', (parent_left + 5) + 'px');

	mover
		.transition().delay(17000).duration(1000)
		.style('opacity', 0);

	// emove
	mover
		.transition().delay(18000).duration(0)
		.remove();

	// .remove();
};

// Let the page get going and then scroll
setInterval(function() {
	var stock = Object.values(vm_stocks.stocks)[Math.floor(Math.random()*Object.values(vm_stocks.stocks).length)];
	scrollTopBar(stock.ticker_id + ": " + formatPrice(stock.current_price));

}, 3000);
