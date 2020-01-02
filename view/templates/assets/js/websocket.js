var visitTxt = document.getElementById("visits");
var trafficTable = document.getElementById("trafficTable");
var data;
w = new WebSocket("wss://" + HOST + "/websocket");
w.onopen = function () {
	console.log("Websocket connection enstablished");
};

w.onclose = function () {
	console.log("Websocket Disconnected form Server")
};
w.onmessage = function (message) {
	console.log(message.data);
	data=JSON.parse(message.data);
	renderPage();
	
};

function renderPage() {
visitTxt.innerHTML=data.VisitCount;
var tableRow = `<tr><th scope="row">`+data.TrafficData.time+ `</th><td>` + data.TrafficData.IpAddr+`</td>
					<td>`+data.TrafficData.Location+`</td>
					<td>`+data.TrafficData.path+` </td>
					<td>`;
if (data.TrafficData.SecurityCheck == "Pass"){
tableRow += ` <div class="badge badge-pill badge-success">PASS</div>`
} else if (data.TrafficData.SecurityCheck == "Warning") {
tableRow += `<div class="badge badge-pill badge-warning">WARNING</div>`
}else{
	tableRow += `<div class="badge badge-pill badge-danger">BLOCKED</div>`
}
  tableRow += `</td></tr>`;
trafficTable.insertAdjacentHTML('afterbegin', tableRow);
}


