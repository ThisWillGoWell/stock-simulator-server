var formatPrice = function(value) {
	// TODO if value is greater than something abbreviate
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
}