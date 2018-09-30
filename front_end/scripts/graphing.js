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
		'top': 40,
		'bottom': 40,
		'left': 40,
		'right': 40,
	};

	var svg = d3.select(location).append('svg')
		.attr('width', width)
		.attr('height', height)
		.append("g");
	
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
		.range([margin.right, width - margin.left])
	let scaleValue = d3.scaleLinear()
		.domain([0, maxValue + (maxValue/10)])
		.range([height  - margin.top, margin.bottom]);


	for (line_key in dat) {
		console.log(line_key);
		let path = svg.append('path');

		let line = d3.line()
			.x(function(d) { return scaleTime(d.time); })
			.y(function(d) { return scaleValue(d.value); });

		// Adding line 
		path.data([dat[line_key]]).attr('d', line).attr('stroke', graphColors[line_key]).attr('stroke-width', '2px').attr('fill', 'none');
	}

	// Creating axis
	var xAxisCall = d3.axisBottom()
		// .tickFormat(d3.timeFormat("%Y-%m-%d"));
	xAxisCall.scale(scaleTime);
	var yAxisCall = d3.axisLeft();
	yAxisCall.scale(scaleValue);

	svg.append('g')
		.attr('id', 'x-axis')
		.attr('transform', 'translate(0, '+ (height - margin.top) +')')
		.style('font-size', 12)
		.call(xAxisCall);
	svg.append('g')
		.attr('id', 'y-axis')
		.attr('transform', 'translate(' + margin.left +', ' + '0' + ')')
		.style('font-size', 12)
		.call(yAxisCall);

	// Adding axis labels
	d3.select('#x-axis').append("text")
		.attr('transform', 'translate(' + margin.left + ', ' + margin.top + ')')
		.attr('fill', 'black')
		.attr('font-size', 14)
		.attr('font-weight', 'bold')
		.text('Time');

	d3.select('#y-axis').append('text')
		.attr('transform', 'translate(50, -5)')
		.attr('fill', 'black')
		.attr('font-size', 15)
		.attr('font-weight', 'bold')
		.text('$');

};