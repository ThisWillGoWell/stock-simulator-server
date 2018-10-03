function notify(message, success) {

	// Set text as the message
	d3.select('#notification-module--container span').html(message);

	// Set notification color
	if (success) {
		var color = "#1abc9c";
	} else {
		var color = "#f44336";
	}

	// Set color and start motion
	d3.select('#notification-module--container')
		.style('background-color', color)
		.transition().duration(300)
		.style('opacity', 1).style('top', '0px');
		

	// Hide and move back up
	d3.select('#notification-module--container')
		.transition().delay(5000).duration(300)
		.style('opacity', 0).style('top', '-60px');

};