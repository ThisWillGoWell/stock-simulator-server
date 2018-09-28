function setSVG(location) {
	var width = 500;
	var height = 300;
	d3.select(location).append('svg')
		.attr('width', width)
		.attr('height', height)//.append('g');
};

function addToLineGraph(location, dat) {
	var svg = d3.select(location).select('svg');
	var path = svg.append('path');

	var minTime = new Date('3000 Jan 1');
	var maxTime = new Date();
	var maxValue = 0;
	dat.forEach(function(d) {
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
	console.log(minTime)
	console.log(maxTime)
	console.log(dat);

	// // Time parsing
	// var parseTime = d3.timeParse("%B %d, %Y");
	// console.log(dat[0].time);
	// console.log(parseTime(dat[0].time));
	var scaleTime = d3.scaleTime()
		.domain([minTime, maxTime])
		.range([0, 100])
	var scaleValue = d3.scaleLinear()
		.domain([0, maxValue])
		.range([300, 0]);

	console.log(scaleValue(dat[30].value));

	var line = d3.line()
		.x(function(d) { return scaleTime(d.time); })
		.y(function(d) { return scaleValue(d.value); });

	// Adding 
	path.data([dat]).attr('d', line).attr('stroke', 'red').attr('stroke-width', '2px').attr('fill', 'none');

	console.log(line);
};