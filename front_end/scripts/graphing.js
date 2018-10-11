
function formatData(data) {

	// if networth add points in to make a step graph
	return data.data;
	
};

// TODO: tags for d3 plotting(title labels etc) sent with dat object in an serparate property
//			tags can pass the type of data being sent through so more data structuring can be done here like min an maxs 
function DrawLineGraph(location, data, id) {
	// Pulling out data, use tags to change data if need
	var dat = formatData(data);
	var tags = data.tags;

	// logging remove later
	console.log("DATA");
	console.log(dat);
	console.log("TAGS");
	console.log(data.tags);
	// logging remove later

	var width = 700;
	var height = 500;
	var margin = {
		'top': 60,
		'bottom': 60,
		'left': 60,
		'right': 60,
	};
	console.log(location);
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
		let path = svg.append('path');

		let line = d3.line()
			.x(function(d) { return scaleTime(d.time); })
			.y(function(d) { return scaleValue(d.value); });

		// Adding line 
		path.data([dat[line_key]]).attr('d', line).attr('stroke', 'black').attr('stroke-width', '2px').attr('fill', 'none');
	}

	// Creating axis
	var xAxisCall = d3.axisBottom()
		.tickFormat(d3.timeFormat("%a %I:%M%p"));
	xAxisCall.scale(scaleTime);
	var yAxisCall = d3.axisLeft()
		.tickFormat(function(d) {
			return "$" + abbrevPrice(d);
		});
	yAxisCall.scale(scaleValue);

	svg.append('g')
		.attr('id', 'x-axis')
		.attr('transform', 'translate(0, '+ (height - (margin.top-5)) +')')
		.style('font-size', 12)
		.call(xAxisCall)
	    .selectAll("text")	
	        .style("text-anchor", "end")
	        .attr("transform", "rotate(-35)")
	        .attr('font-size', '10px');

	svg.append('g')
		.attr('id', 'y-axis')
		.attr('transform', 'translate(' + (margin.left-5) +', ' + '0' + ')')
		.style('font-size', 12)
		.call(yAxisCall);

	// Add graph title
	if (tags) {

		if (tags.title) {
			svg.append('text').text(tags.title)
				.attr('class', 'stockGraph-title')
				.attr('font-size', '20px')
				.attr('text-anchor', 'middle')
				.attr('transform', 'translate(' + (width/2) + ', 40)');
		}
	}

	// // Adding axis labels
	// d3.select('#x-axis').append("text")
	// 	.attr('transform', 'translate(' + (margin.left/2) + ', ' + (margin.top/2) + ')')
	// 	.attr('fill', 'black')
	// 	.attr('font-size', 14)
	// 	.attr('font-weight', 'bold')
	// 	.text('Time');

	// d3.select('#y-axis').append('text')
	// 	.attr('transform', 'translate(0, ' + (margin.top/2) + ')')
	// 	.attr('fill', 'black')
	// 	.attr('font-size', 15)
	// 	.attr('font-weight', 'bold')
	// 	.text('$');

};

$('#select-timeframe').selectize({
    create: true,
    sortField: 'text'
});
