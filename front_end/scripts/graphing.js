const TICKS = 5;
const COLOR_PALETTE = [
	"#EF5350",
	"#AB47BC",
	"#5C6BC0",
	"#29B6F6",
	"#66BB6A",
	"#FFCA28",
	"#FF7043",
	"#D4E157",
];





function formatData(data) {
	// Setting local time			
	Object.values(data).forEach(function(d) {
		d.forEach(function(i) {
			i.time = new Date(i.time);
		});
		
		// Sorting by time
		d.sort(function(a,b) {
			if (a.time > b.time) {
				return 1;
			} else if (a.time < b.time) {
				return -1;
			} else {
				return 0;
			}
		});
	});


	// if networth add points in to make a step graph
	return data;
};

// Get graph data
function queryDrawGraph(location, uuids, fields, append = false) {
	if (uuids.length !== fields.length) {
		console.error("In getGraphData(): fields and uuids are not the same length");
	}

	var data = {
	  data: [],
	  tags: {},
	};
	
	var stillWaiting = true;
	var responses = [];
	var requests = [];

	uuids.forEach(function(d, i) {
		queryDB(uuids[i], fields[i], requests, responses, data);
	});
  
	var drawGraphOnceDone = null;
	drawGraphOnceDone = function() {
	  if (requests.every(r => responses.indexOf(r) > -1)) {
		stillWaiting = false;
	  }
  
	  if (!stillWaiting) {
		// draw graph once all the data is back
		DrawLineGraph(location, data, append = append);
	  } else {
		setTimeout(drawGraphOnceDone, 100);
	  }
	};
  
	setTimeout(drawGraphOnceDone, 100);

}


// Store graphing data
function queryDB(uuid, field, requests, responses, data, num_points = 1000, length = "6h", use_cache = true) {
	var msg = {
		uuid: uuid,
		field: field,
		num_points: num_points,
		use_cache: use_cache,
		length: length
	}; 
  
	// Store request on front end
	requests.push(REQUEST_ID.toString());
	var callback = function(msg) {
		console.log("GRAPH msg HERE")
		console.log(msg)
		// Pull out the data and format it
		var points = msg.msg.points;
		points = points.map(function(d) {
			return { time: d[0], value: d[1] };
		});

		// Store the data
		let key = msg.msg.message.field + "-" + msg.msg.message.uuid;
		data.data[key] = points;//.push(points);//[msg.msg.message.field] = points;

		// Make note the data is available
		responses.push(msg.request_id);
	};
  
	// Send message
	doSend("query", msg, callback);
  
}

// Takes a legend lable key and cleans up uuids
function cleanLegendLabel(label) {
	var parts = label.split('-');
	var field = parts[0];
	var uuid = parts[1];
	var label;

	// According to the field get the object label
	switch(field) {
		case 'current_price':
			label = vm_stocks.stocks[uuid].ticker_id;
			break;
		case 'net_worth':
			label = vm_portfolios.portfolios[uuid].name;
			label += "'s networth";
			break;
		case 'wallet':
			label = vm_portfolios.portfolios[uuid].name;
			label += "'s wallet";
			break;
	}
	return label;
}


// TODO: tags for d3 plotting(title labels etc) sent with dat object in an serparate property
//			tags can pass the type of data being sent through so more data structuring can be done here like min an maxs 
function DrawLineGraph(location, data, id, append) {
	console.log(data)

	// Pulling out data, use tags to change data if need
	var dat = formatData(data.data);
	console.log(dat)
	
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

	if (!append) {
		d3.select(location).selectAll('svg').remove();
	}

	var svg = d3.select(location).append('svg')
		.attr('width', width)
		.attr('height', height)
    
    var g = svg.append("g").attr("class", "line-area");
	
	var minTime = new Date('3000 Jan 1');
	var maxTime = new Date('1999 Jan 1');
	var maxValue = Number.NEGATIVE_INFINITY;
	var minValue = Number.POSITIVE_INFINITY;
	// var minValueTime, maxValueTime;
	for (var line_key in dat) {

		dat[line_key].forEach(function(d) {
			if (minValue > d.value) {
				minValue = d.value;
				minValueTime = d.time;
			}
			if (maxValue < d.value) {
				maxValue = d.value;
				maxValueTime = d.time;
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
	var scaleTime = d3.scaleTime()
		.domain([minTime, maxTime])
		.range([margin.left, width - margin.right])
	var scaleValue = d3.scaleLinear()
		.domain([minValue - (maxValue/10), maxValue + (maxValue/10)])
		.range([height  - margin.top, margin.bottom]);

    // Creating the line logic
    var line = d3.line()
        .x(function(d) { return scaleTime(d.time); })
        .y(function(d) { return scaleValue(d.value); });

        
    function waiting() {
        waitingTimeout = null;
    }
    var delay = 500, waitingTimeout;

    // When brushing stops
    function brushEnd() {
        var event = d3.event.selection;
        console.log(event)
        if (!event) {
            if (!waitingTimeout) return waitingTimeout = setTimeout(waiting, delay);
            scaleTime.domain([minTime, maxTime]);
        } else {
            console.log(event)
            scaleTime.domain([event[0], event[1]].map(scaleTime.invert, scaleTime));
            svg.select('.graph-brush').call(brush.move, null);
        }
        brushZoom();
    };
    
    function brushZoom() {
        // Create transition
        var transition = svg.transition().duration(500);
        
        svg.select("#x-axis").transition(transition).call(xAxisCall);
        // svg.select("#y-axis").transition(transition).call(yAxisCall);
        
        svg.selectAll(".graph-line").attr('d', line);
    };

    var brush = d3.brushX().on('end', brushEnd);

    // Add brush
    svg.append("g")
        .attr('class', 'graph-brush')
		.on('mouseover', function() { toolTips.style('display', null); })
		.on('mousemove', function() {
			var mouseX = d3.mouse(this)[0];
			var mouseY = d3.mouse(this)[1];
			var xVal = scaleTime.invert(mouseX);
			Object.keys(dat).forEach(function(key) {
					// Get index of where a 'new' point would fit 
					var i = bisectTime(dat[key], xVal, 1, dat[key].length - 1);
					// Find points on either side
					var d0 = dat[key][i - 1]; 
					var d1 = dat[key][i];
					// Compare which is closer
					var tipPoint = xVal - d0.time > d1.time - xVal ? d1 : d0;

					d3.select('#' + key).attr('transform', 'translate(' + scaleTime(tipPoint.time) + ',' + scaleValue(tipPoint.value) + ')');
					d3.select('#legend-' + key).html(cleanLegendLabel(key) + ': $' + formatPrice(tipPoint.value));
				});
			
			// Get legend size
			var w = legendParent.select('div').node().getBoundingClientRect().width;
			var h = legendParent.select('div').node().getBoundingClientRect().height;

			// orientate the legend correctly
			if (scaleTime(mouseX) > scaleTime(width/2)) {
				if (scaleValue(mouseY) > scaleValue(height/2)) {
					legendParent.attr('transform', 'translate(' + (mouseX - w - 30) + ',' + (mouseY + 15) + ')');
				} else {
					legendParent.attr('transform', 'translate(' + (mouseX - w - 30) + ',' + (mouseY - h - 30) + ')');
				}
			} else {
				if (scaleValue(mouseY) > scaleValue(height/2)) {
					legendParent.attr('transform', 'translate(' + (mouseX + 15) + ',' + (mouseY + 15) + ')');
				} else {
					legendParent.attr('transform', 'translate(' + (mouseX + 15) + ',' + (mouseY - h - 30) + ')');
				}
			}
		})
		.on('mouseout', function() { 
			// Get max of each graph
			toolTips.style('display', 'none'); })
        .call(brush);
    
    // Used when finding which point to tooltip
    var bisectTime = d3.bisector(d => d.time).left;

	// Creating graph legend
	var legendParent = svg.append('g')
		// .attr('class', 'graph-legend')
		.style('pointer-events', 'none')
		
	var legend = legendParent.append('foreignObject')
		.append('xhtml:div')
		.attr('class', 'graph-legend');
	
	
	// Y grid
	function yGrid() {
		return d3.axisLeft(scaleValue).ticks(TICKS);
	}

	var labels = [];
	var i = 0;

	for (line_key in dat) {
		// Creating space for each line graph
		let path = g.append('path').attr('class','graph-line');

		// Adding line 
		path.data([dat[line_key]])
			.attr('d', line)
			.attr('stroke', COLOR_PALETTE[i])
			.attr('stroke-width', '2px')
			.attr('fill', 'none');

		console.log(line_key);
		// Adding tooltip for each line
		let ttip = svg.append('g')
			.attr('class', 'graph-tooltip')
			.attr('id', line_key)
			
		ttip.append("circle")
			.attr('r', 4)
			.style('fill', COLOR_PALETTE[i]);

		legend.append('span')
			.attr('id', 'legend-' + line_key)
			.attr('class', 'legend-item')
			.style('color', COLOR_PALETTE[i]);

		i++;
	}

	// Selecting all tooltips
	var toolTips = d3.selectAll('.graph-legend');


	// Creating x axis
	var xAxisCall = d3.axisBottom(scaleTime)
		// %a for day of the week
		.tickFormat(d3.timeFormat("%I:%M%p"))
		.ticks(TICKS);
    
    // Creating y axis 
	var yAxisCall = d3.axisLeft(scaleValue)
		.tickFormat(function(d) {
			return "$" + abbrevPrice(d);
		})
		.ticks(TICKS);

	// Add axis
	g.append('g')
		.attr('id', 'x-axis')
		.attr('transform', 'translate(0, '+ (height - (margin.top-5)) +')')
		.style('font-size', 12)
		.call(xAxisCall)
	    .selectAll("text")	
	        .attr('font-size', '10px');

	g.append('g')
		.attr('id', 'y-axis')
		.attr('transform', 'translate(' + (margin.left-5) +', ' + '0' + ')')
		.style('font-size', 12)
		.call(yAxisCall);

	// adding horizontal gridlines
	g.append("g")
		.attr("class", "graph-grid")
		.call(yGrid()
			.tickSize(-width)
			.tickFormat(""));

       



	// Add graph title
	if (tags) {
		if (tags.title) {
			g.append('text').text(tags.title)
				.attr('class', 'stockGraph-title')
				.attr('font-size', '20px')
				.attr('text-anchor', 'middle')
				.attr('transform', 'translate(' + (width/2) + ', 40)');
		}
	}
	console.log(labels)
};
	

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


$('#select-timeframe').selectize({
    create: true,
    sortField: 'text'
});