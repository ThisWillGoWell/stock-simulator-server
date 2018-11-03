function isRecent(time, since) {
	let noteTime = new Date(time);
	let currTime = new Date();
	let timeSince = currTime - noteTime;
	return (timeSince < since);
}

function formatPrice(value) {
	let val = (value/100).toFixed(2).toString();
	val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
	return val;

};

function formatPercent(value) {
	return Math.abs(value);
};

function findPercentChange(newPrice, oldPrice) {
    if (newPrice > oldPrice) {
      return ((newPrice - oldPrice)/oldPrice * 100).toFixed(2);
    } else if (newPrice < oldPrice) {
      return ((oldPrice - newPrice)/oldPrice * -100).toFixed(2);
    } else {
      return (0).toFixed(2);
    }
  };

function abbrevPrice(value) {
	// TODO if value is greater than something abbreviate
	// Handle negative values for formatting change
	if (value < 1000000) {
		let val = (value/100).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val;
	} else if (value < 10000000) {
		let val = ((value/100)/1000).toFixed(2).toString();
		val = val.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		return val + "K";
	} else if (value < 1000000000) {
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

// get selected values from a selectize
function getSelectized(loc) {
	var $select = $(loc);
	var selectize = $select[0].selectize;
	var fields = $.map(selectize.items, function(val) {
		return selectize.options[val].value;
	});
	return fields;
}

var dateFormat = function () {
	var token = /d{1,4}|m{1,4}|yy(?:yy)?|([HhMsTt])\1?|[LloSZ]|"[^"]*"|'[^']*'/g,
			timezone = /\b(?:[PMCEA][SDP]T|(?:Pacific|Mountain|Central|Eastern|Atlantic) (?:Standard|Daylight|Prevailing) Time|(?:GMT|UTC)(?:[-+]\d{4})?)\b/g,
			timezoneClip = /[^-+\dA-Z]/g,
			pad = function (val, len) {
					val = String(val);
					len = len || 2;
					while (val.length < len) val = "0" + val;
					return val;
			};

	// Regexes and supporting functions are cached through closure
	return function (date, mask, utc) {
			var dF = dateFormat;

			// You can't provide utc if you skip other args (use the "UTC:" mask prefix)
			if (arguments.length == 1 && Object.prototype.toString.call(date) == "[object String]" && !/\d/.test(date)) {
					mask = date;
					date = undefined;
			}

			// Passing date through Date applies Date.parse, if necessary
			date = date ? new Date(date) : new Date;
			if (isNaN(date)) throw SyntaxError("invalid date");

			mask = String(dF.masks[mask] || mask || dF.masks["default"]);

			// Allow setting the utc argument via the mask
			if (mask.slice(0, 4) == "UTC:") {
					mask = mask.slice(4);
					utc = true;
			}

			var _ = utc ? "getUTC" : "get",
					d = date[_ + "Date"](),
					D = date[_ + "Day"](),
					m = date[_ + "Month"](),
					y = date[_ + "FullYear"](),
					H = date[_ + "Hours"](),
					M = date[_ + "Minutes"](),
					s = date[_ + "Seconds"](),
					L = date[_ + "Milliseconds"](),
					o = utc ? 0 : date.getTimezoneOffset(),
					flags = {
							d:    d,
							dd:   pad(d),
							ddd:  dF.i18n.dayNames[D],
							dddd: dF.i18n.dayNames[D + 7],
							m:    m + 1,
							mm:   pad(m + 1),
							mmm:  dF.i18n.monthNames[m],
							mmmm: dF.i18n.monthNames[m + 12],
							yy:   String(y).slice(2),
							yyyy: y,
							h:    H % 12 || 12,
							hh:   pad(H % 12 || 12),
							H:    H,
							HH:   pad(H),
							M:    M,
							MM:   pad(M),
							s:    s,
							ss:   pad(s),
							l:    pad(L, 3),
							L:    pad(L > 99 ? Math.round(L / 10) : L),
							t:    H < 12 ? "a"  : "p",
							tt:   H < 12 ? "am" : "pm",
							T:    H < 12 ? "A"  : "P",
							TT:   H < 12 ? "AM" : "PM",
							Z:    utc ? "UTC" : (String(date).match(timezone) || [""]).pop().replace(timezoneClip, ""),
							o:    (o > 0 ? "-" : "+") + pad(Math.floor(Math.abs(o) / 60) * 100 + Math.abs(o) % 60, 4),
							S:    ["th", "st", "nd", "rd"][d % 10 > 3 ? 0 : (d % 100 - d % 10 != 10) * d % 10]
					};

			return mask.replace(token, function ($0) {
					return $0 in flags ? flags[$0] : $0.slice(1, $0.length - 1);
			});
	};
}();

function pathTween(d1, precision) {
	return function() {
	  var path0 = this,
		  path1 = path0.cloneNode(),
		  n0 = path0.getTotalLength(),
		  n1 = (path1.setAttribute("d", d1), path1).getTotalLength();
  
	  // Uniform sampling of distance based on specified precision.
	  var distances = [0], i = 0, dt = precision / Math.max(n0, n1);
	  while ((i += dt) < 1) distances.push(i);
	  distances.push(1);
  
	  // Compute point-interpolators at each distance.
	  var points = distances.map(function(t) {
		var p0 = path0.getPointAtLength(t * n0),
			p1 = path1.getPointAtLength(t * n1);
		return d3.interpolate([p0.x, p0.y], [p1.x, p1.y]);
	  });
  
	  return function(t) {
		return t < 1 ? "M" + points.map(function(p) { return p(t); }).join("L") : d1;
	  };
	};
  }