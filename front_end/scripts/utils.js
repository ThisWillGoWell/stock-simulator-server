function formatPrice(value) {
	let val = (value/100).toFixed(2).toString();
	val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
	return val;

};

function abbrevPrice(value) {
	// TODO if value is greater than something abbreviate
	// Handle negative values for formatting change
	if (value < 1000000) {
		let val = (value/100).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val;
	} else if (value < 1000000000) {
		let val = ((value/100)/1000).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val + "K";
	} else if (value < 100000000000) {
		let val = ((value/100)/1000000).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val + "M";
	} else if (value < 100000000000000) {
		let val = ((value/100)/1000000000).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val + "B";
	} else if (value < 100000000000000000) {
		let val = ((value/100)/1000000000000).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val + "T";
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