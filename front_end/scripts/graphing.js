const TICKS = 5;


function formatData(data) {
	// Setting local time			
	Object.values(data).forEach(function(d) {
		d.forEach(function(i) {
			i.time = new Date(i.time);
		})
	})
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
		console.log(i)
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
		data.data.push(points);//[msg.msg.message.field] = points;

		// Make note the data is available
		responses.push(msg.request_id);
	};
  
	// Send message
	doSend("query", msg, callback);
  
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
    
    var g = svg.append("g");
	
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
		.range([margin.right, width - margin.left])
	var scaleValue = d3.scaleLinear()
		.domain([minValue - (minValue/10), maxValue + (maxValue/10)])
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
        .call(brush);
    
    // Used when finding which point to tooltip
    var bisectTime = d3.bisector(function(d) {
            return d.time
        }).left;

	// Y grid
	function yGrid() {
		return d3.axisLeft(scaleValue).ticks(TICKS);
	}

	var labels = [];

	for (line_key in dat) {
		// const timeFormat = d3.timeFormat("%I:%M %p")
		
		// var minY = Number.POSITIVE_INFINITY;
		// var maxY = Number.NEGATIVE_INFINITY;
		// var minX, maxX;
		// dat[line_key].forEach(function(d) {
		// 	if (minY > d.value) {
		// 		minY = d.value;
		// 		minX = d.time;
		// 	}
		// 	if (maxY < d.value) {
		// 		maxY = d.value;
		// 		maxX = d.time;
		// 	}
		// });

		// //Add annotations
		// var newLabels = [{
		// 	data: { date: minX, value: formatPrice(minY) },
		// 		x: scaleTime(minX),
		// 		y: scaleValue(minY).toFixed(2),
		// 		dx: 10,
		// 		dy: 10
		// 	}, {
		// 	data: { date: maxX, value: formatPrice(maxY) },
		// 		x: scaleTime(maxX),
		// 		y:scaleValue(maxY).toFixed(2),	
		// 		dx: 10,
		// 		dy: 10
		// }].map(function (l) {
		// 	l.note = Object.assign({}, l.note, {
		// 		title: "$" + l.data.value,
		// 		label: "" + timeFormat(l.data.date)
		// 	});
		// 	l.connector = {
		// 		end: "dot",
		// 	};

		// 	return l;
		// });

        // newLabels.forEach(d => labels.push(d));
        
        // Sorting the data
        // dat[line_key].sort(function(a,b) {
        //     return a > b;
        // });

		let path = g.append('path').attr('class','graph-line');

		// Adding line 
		path.data([dat[line_key]]).attr('d', line).attr('stroke', 'black').attr('stroke-width', '2px').attr('fill', 'none');
	}

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

	// Add gridlines
	// add the Y gridlines
	g.append("g")			
		.attr("class", "graph-grid")
		.call(yGrid()
			.tickSize(-width)
			.tickFormat("")
        )
        
    // Adding tooltip
    var ttip = svg.append('g')
        .attr('class', 'graph-tooltip')
        .style('display', 'none');

    // Adding hover tooltip layer
    svg.append('rect')
        .attr('class', 'hover-overlay')
        .attr('width', width)
        .attr('height', height)
        .on('mouseover', function() { ttip.style('display', null); })
        .on('mousemove', function() {
            console.log(d3.mouse(this))
            var xVal = scaleTime.invert(d3.mouse(this)[0]);
            Object.values(dat).forEach(function(d) {
                    var i = bisectTime(d, xVal, 1, d.length - 1 );
                    var d0 = d[i - 1]; 
                    var d1 = d[i];
                    var dat = xVal - d0.time > d1.time - xVal ? d1 : d0;
                    console.log(dat);
                });
        })
        .on('mouseout', function() { ttip.style('display', 'none'); });

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

	// var annotate = d3.annotation().annotations(labels).editMode(true)
	// 	.type(d3.annotationCallout)
	// 	.accessors({ 
	// 		x: d => scaleTime(d.date), 
	// 		y: d => scaleValue(d.value).toFixed(2)
	// 	  })
	// 	.accessorsInverse({
	// 		date: d => scaleTime.invert(d.x),
	// 		value: d => scaleValue.invert(d.y).toFixed(2) 
	// 	});

	// // Adding annotations for extreme values
	// svg.append('g')
	// 	.attr('class', 'graph-annotation')
	// 	.call(annotate);

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