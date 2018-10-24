var root = document.body;
var geolocationInfo = '';
var ipToSearchFor = 'IP address';

var searchIP = function() {
	m.request({
		method: "GET",
		url: "ip2geo",
		data: {i: ipToSearchFor},
		withCredentials: true,
	})
	.then(function(data) {
		console.log(data);
		geolocationInfo = data;
	});
}

var App = {
	view: function() {
		return m("main", [
			m("div", {class: 'ipsearch'}, 
				m("input[type=text][placeholder=IP 8.8.8.8]", {
					oninput: m.withAttr("value", function(value) {ipToSearchFor = value; searchIP(ipToSearchFor);}),
					value: ipToSearchFor
				})
			),
			m("div", {class: 'result'}, m.trust(formatLocation(geolocationInfo))),
		]);
	}
}

m.mount(root, App);

function formatLocation(geolocation) {
	formatted = '';
	var skipField = 'This parameter is unavailable for selected data file. Please upgrade the data file.';
	var invalidIP = 'Invalid IP address.';
	if (geolocation['CountryShort'] == invalidIP) {
		return invalidIP;
	}
	for (var property in geolocation) {
		if (geolocation.hasOwnProperty(property) && geolocation[property] != skipField && property != 'Elevation') {
			formatted += '</br>' + property + ': ' + geolocation[property];
		}
	}
	return formatted;
}