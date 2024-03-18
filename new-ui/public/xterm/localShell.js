
function connect(){
	
	
	url = "ws://"+document.location.host+"/ws/localshell"
	console.log(url);
	let term = new Terminal({
		"cursorBlink":true,
	});
	if (window["WebSocket"]) {
		term.open(document.getElementById("terminal"));
		term.fit();
		// term.toggleFullScreen(true);
		term.on('data', function (data) {
			const encodedData = btoa(data);
    		conn.send(encodedData);
		});

		conn = new WebSocket(url);
		conn.onopen = function(e) {
			conn.send(btoa("sh && clear \n"))
			
		};
		conn.onmessage = function(event) {			
			const decodedData = atob(event.data);
			term.write(decodedData);
		};
		conn.onclose = function(event) {
			if (event.wasClean) {
				console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
			} else {
				console.log('[close] Connection died');
				term.writeln("")
			}
			term.write('Connection Reset By Peer! Try Refresh.');
		};
		conn.onerror = function(error) {
			console.log('[error] Connection error');
			term.write("error: "+error.message);
			term.destroy();
		};
	} else {
		var item = document.getElementById("terminal");
		item.innerHTML = "<h2>Your browser does not support WebSockets.</h2>";
	}
}