// Setting graph colors. TEMP: in future use class so brennan can manage the css
var graphColors = {
	'net_worth': 'red',
	'wallet': 'green',
};


function DrawPortfolioGraph(location, dat, id) {
	console.log(dat);
	var width = 500;
	var height = 300;
	var margin = {
		'top': '40px',
		'bottom': '40px',
		'left': '40px',
		'right': '40px',
	};

	var svg = d3.select(location).append('svg')
		.attr('width', width)
		.attr('height', height);
	
	let minTime = new Date('3000 Jan 1');
	let maxTime = new Date();
	let maxValue = 0;


	for (var line_key in dat) {

		dat[line_key].forEach(function(d) {
			// formatting time IMPORTANT!
			d.time = new Date(d.time.replace("T"," ").replace("Z", ""));
			if (maxValue < d.value) {
				maxValue = d.value;
			}
			if (minTime > d.time) {
				minTime = d.time;
			}
			if (maxTime < d.time) {
				maxTime = d.time;
			}
		});
	}

	// Creating graph scales
	let scaleTime = d3.scaleTime()
		.domain([minTime, maxTime])
		.range([0, width])
	let scaleValue = d3.scaleLinear()
		.domain([0, maxValue + (maxValue/10)])
		.range([height, 0]);

	for (line_key in dat) {
		console.log(line_key);
		let path = svg.append('path');

		let line = d3.line()
			.x(function(d) { return scaleTime(d.time); })
			.y(function(d) { return scaleValue(d.value); });

		// Adding line 
		path.data([dat[line_key]]).attr('d', line).attr('stroke', graphColors[line_key]).attr('stroke-width', '2px').attr('fill', 'none');
	}
};