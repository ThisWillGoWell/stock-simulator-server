// Vue object for notifications
var vm_notify = new Vue({
  data: {
    notes: {},
  },
});

// Register route for websocket messages for notifications
registerRoute("notification", function(msg) {
	// Add notification to the list
	Vue.set(vm_notify.notes, msg.msg.uuid, msg.msg);

// TODO: do not route if the message has been seen
	var seen = msg.msg.seen;
	if (!seen) {
		// Execute notification type
		routeNote[msg.msg.type](msg.msg);
	}

});


var routeNote = {
	trade: notifyTrade,
	send_money: notifyTransfer,
	recieve_money: notifyTransfer,
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

	// sendAck(msg.uuid);

};

function notifyTransfer(msg) {
	var color, message;
	var success = msg.notification.success;

	// Getting usernames
	var receiver = msg.notification.receiver;
	receiver = vm_users.users[receiver].display_name;

	if (isRecent(msg.time, 30000)) {
		// If trade was a success
		if (success) {
			// Getting amount 
			var amount = msg.notification.amount;
			message = "Sucessful tranfer of " + formatPrice(amount) + " to " + receiver + ".";
			color = "#1abc9c";

		} else {
			message = "Tranfer to " + receiver + " failed.";
			color = "#f44336";
		}

		notifyTopBar(message, color, success);
	}

	// sendAck(msg.uuid);

};

function notifyTrade(msg) {

	var color, message;
	var success = msg.notification.success;
	
	// Notify only if the trade was within the last minute
	if (isRecent(msg.time, 30000)) {
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

	}

	// sendAck(msg.uuid);

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



