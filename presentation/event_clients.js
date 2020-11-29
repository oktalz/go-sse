// Create a new HTML5 EventSource
var source = new EventSource('http://127.0.0.1:'+port+'/api/sse?sub=status,time');

// Create a callback for when a new message is received.
source.onmessage = function(e) {

    //console.log(e.data);
    // Append the `data` attribute of the message to the DOM.
    var data = JSON.parse(e.data)
    if (data.type == "time"){
        document.getElementById("time").innerHTML = data.time;
    }
    if (data.type == "status"){
        document.getElementById("status").innerHTML = data.status;
    }
};
