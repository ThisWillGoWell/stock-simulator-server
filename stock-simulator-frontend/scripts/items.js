// var UseRouter = {
// 	insider: reuseInsiderTrading,
// };


var vm_items = new Vue({
	data: {
	  items: {}
	},
	computed: {
		userItems: function() {
			return Object.values(vm_items.items).forEach(function(i) {
				console.log(i)
				var item = {
					name: i.name,
					config_id: i.config,
					uuid: i.uuid,
					portfolio_uuid: i.portfolio_uuid,
				};
				switch(i.config) {
					case 'personal_broker':
						item.duration = prettifyItemDuration(i.duration);
						item.desc = [
							['Purchase Fee', '$0'],
							['Sale Fee', '$0']
						]
						break;
				}
				console.log(i)
				console.log(item)
				return item;
			})
	  	}
	}
  
});

function prettifyItemDuration(dur) {
	var split = dur.indexOf('h');
	var hours = dur.substring(0, split);
	dur = dur.substring(split+1);
	split = dur.indexOf('m');
	var minutes = dur.substring(0, split);
	dur = dur.substring(split+1);
	split = dur.indexOf('s');
	var seconds = dur.substring(0, split);
	
	var ret = "";
	if (hours > 0) {
		if (hours == 1) {
			ret += (hours + " Hour ")
		} else {
			ret += (hours + " Hours ")
		}
	}
	if (minutes > 0) {
		if (minutes == 1) {
			ret += (minutes + " Minute ")
		} else {
			ret += (minutes + " Minutes ")
		}
	}
	if (seconds > 0) {
		if (seconds == 1) {
			ret += (seconds + " Second")
		} else {
			ret += (seconds + " Seconds")
		}
	}
	return ret;
}

function useItem(item_uuid) {

	var msg = {
		'action': 'use',
		'o': {
			'uuid': item_uuid,
			'params': {}
		}
	};

	var callback = function(msg) {
		if (msg.msg.o.success) {
			console.log(msg)
			notifyTopBar("Successfully used!")
		} else {
			var message = msg.msg.o.err;
			notifyTopBar(message, RED, msg.msg.o.success);
		}
	};

	doSend('item', msg, callback);

};


// Vue.component('personal-broker', {
// 	props: {

// 		desc:
// 	}
// })

// function checkUsedItems() {

// 	Object.values(vm_items.items).forEach(function(d) {
// 		if (d.used) {
// 			UseRouter[d.type](d);
// 		}
// 	});
// };


// function reuseInsiderTrading(item) {
// 	createInsiderArea(item.target_prices);
// };


// function createInsiderArea(target_dict) {
// 	var stocks = []; 
// 	Object.keys(target_dict).forEach(function(d) {
// 		console.log(d);
// 		var stock = vm_stocks.stocks[d];
// 		var insiderStock = {
// 			'name': stock.ticker_id,
// 			'target_price': target_dict[d],
// 			'current_price': stock.current_price,
// 		};

// 		stocks.push(insiderStock);		
// 	})

// 	Vue.set(vm_dash_tab, 'insiderStocks', stocks);

// };



