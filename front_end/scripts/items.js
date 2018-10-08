var UseRouter = {
	insider: reuseInsiderTrading,
};


function useItem(item_uuid) {

	var msg = {
		'action': 'use',
		'o': {
			'uuid': item_uuid,
			'params': {}
		}
	};

	var callback = function(msg) {
		createInsiderArea(msg.msg.o);
	};

	doSend('item', msg, callback);

};


function checkUsedItems() {

	Object.values(vm_items.items).forEach(function(d) {
		if (d.used) {
			UseRouter[d.Type](d);
		}
	});
};


function reuseInsiderTrading(item) {
	createInsiderArea(item.target_prices);
};


function createInsiderArea(target_dict) {
	var stocks = []; 
	Object.keys(target_dict).forEach(function(d) {
		var stock = vm_stocks.stocks[d];
		console.log(stock);
		var insiderStock = {
			'name': stock.ticker_id,
			'target_price': target_dict[d],
			'current_price': stock.current_price,
		};

		stocks.push(insiderStock);		
	})

	Vue.set(vm_dash_tab, 'insiderStocks', stocks);

};



