
<% 
String result = (String)request.getAttribute("result");
String path = request.getContextPath();
%>

<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="icon" href="./favicon.ico">

    <title>Chain A Observer</title>

    <!-- Bootstrap core CSS -->
    <link href="<%=path%>/css/dist/css/bootstrap.min.css" rel="stylesheet">
	<link href="<%=path%>/css/dist/dashboard.css" rel="stylesheet">
  </head>
  <body class="text-center">
  <div class="container">
  <body>
    <nav class="navbar navbar-dark fixed-top bg-dark flex-md-nowrap p-0 shadow">
      <a class="navbar-brand col-sm-3 col-md-2 mr-0" href="#">Cross-Blockchain Message Transfer System ---- Chain A</a>
      <ul class="navbar-nav px-3">
        <li class="nav-item text-nowrap">
          <a class="nav-link" href="index">Back to Index</a>
        </li>
      </ul>
    </nav>

    <div class="container-fluid">
      <div class="row">
        <nav class="col-md-2 d-none d-md-block bg-light sidebar">
          <div class="sidebar-sticky">
            <ul class="nav flex-column">
              <li class="nav-item">
                <a class="nav-link active" href="#">
                  <span data-feather="home"></span>
                  Dashboard <span class="sr-only">(current)</span>
                </a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#">
                  <span data-feather="file"></span>
                 Message Hisory
                </a>
              </li>
              
        </nav>

        <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-4">
          

          <h2>HISTORY MESSAGES of Chain A</h2>
          <div class="table-responsive">
            <table class="table table-striped table-sm" id="table">
              <thead>
            <tr>
			<th>Sender</th>
			<th>Receiver</th>
			<th>Source network</th>
			<th>Destination network</th>
			<th>Content</th>
		    </tr>
              </thead>
              <tbody>
                
              </tbody>
            </table>
          </div>
          <button class="btn btn-lg btn-primary btn-block" type="button" onclick="refresh()" >Refresh</button>
        </main>
      </div>
    </div>
    </body>
</html>




<script>
    var a = document.createElement("a"); 
    document.body.appendChild(a); 
    a.style = "display: none"; 
    
    function refresh(){
    	location.reload();
    }
    
    window.onload = listResult();
	window.onload = function socket() {
		var ws = new WebSocket('ws://10.6.71.96:8081');
		ws.onopen = function() {
		};
		ws.onmessage = function(event) {
			if(event.data == "\"redundant\""){
				alert("You receive a redundant transaction from other bridge, it will be discarded.");
			} else if (event.data == "\"Hash\""){
				alert("You receive a transaction from other bridge, but its \"Hash\" has been tampered!");
			} else if (event.data == "\"Content\""){
				alert("You receive a transaction from other bridge, but its \"Content\" has been tampered!");
			} else if (event.data == "\"TransactionID\""){
				alert("You receive a transaction from other bridge, but its \"TransactionID\" has been tampered!");
			} else if (event.data == "\"SenderID\""){
				alert("You receive a transaction from other bridge, but its \"SenderID	\" has been tampered!");
			} else if (event.data == "\"ReceiverID\""){
				alert("You receive a transaction from other bridge, but its \"ReceiverID\" has been tampered!");
			} else if (event.data == "\"SourceID\""){
				alert("You receive a transaction from other bridge, but its \"SourceID\" has been tampered!");
			} else if (event.data == "\"DestinationID\""){
				alert("You receive a transaction from other bridge, but its \"DestinationID\" has been tampered!");
			} else if (event.data == "\"signature\""){
				alert("You receive a transaction from other bridge, but its \"signature\" is failed to be verified!");
			} else{
				alert("You receive a new valid transaction. You can refresh the page to review it.");
			}			
		};
	}
	
	function listResult() {
		var dataArray = <%=result%>;
		var table = document.getElementById("table");
		table.getElementsByTagName("tbody").innerHTML = "";
		for (var i = 0; i < dataArray.length; i++) {
			var data = dataArray[i];
			var row = table.insertRow(table.rows.length);
			var c1 = row.insertCell(0);
			c1.innerHTML = data.Record.header.senderID;
			var c2 = row.insertCell(1);
			c2.innerHTML = data.Record.header.receiverID;
			var c3 = row.insertCell(2);
			c3.innerHTML = data.Record.header.sourceID;
			var c4 = row.insertCell(3);
			c4.innerHTML = data.Record.header.destinationID;
			var c5 = row.insertCell(4);
			c5.innerHTML = data.Record.content;
		}	
	}
	
	function base64ToArrayBuffer(base64) {
	    var binary_string = window.atob(base64);
	    var len = binary_string.length;
	    var bytes = new Uint8Array(len);
	    for (var i = 0; i < len; i++) {
	        bytes[i] = binary_string.charCodeAt(i);
	    }
	    return bytes.buffer;
	}
	
</script>