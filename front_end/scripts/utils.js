function formatPrice(value) {
	// TODO if value is greater than something abbreviate
	// Handle negative values for formatting change
	if (value < 1000000) {
		let val = (value/100).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val;
	} else if (value < 100000000) {
		let val = ((value/100)/1000).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val + "K";
	} else {
		let val = (value/100).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val;
	}
};


function formatDate12Hour(date) {
  	let hours = date.getHours();
	let minutes = date.getMinutes();
	let ampm = hours >= 12 ? 'pm' : 'am';
	hours = hours % 12;
	hours = hours ? hours : 12; // the hour '0' should be '12'
	minutes = minutes < 10 ? '0'+minutes : minutes;
	let strTime = hours + ':' + minutes + ' ' + ampm;
	return strTime;
};


function getHighestStock(stocks) {
	stocks = Object.values(stocks).map((d) => d);
	var highestStock = stocks.reduce(function(a, b){ return a.current_price > b.current_price ? a : b });
	return highestStock;
};

function getMoverStock(stocks) {
	stocks = Object.values(stocks).map((d) => d);
	var mover = stocks.reduce((a, b) => a.change > b.change ? a : b);
	return mover;
};