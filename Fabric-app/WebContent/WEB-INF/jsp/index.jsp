<%
String path = request.getContextPath();
String basePath = request.getScheme()+"://"+request.getServerName()+":"+request.getServerPort()+path+"/";
String socketPath = "ws://" + request.getServerName() + ":8081";
%>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="Colorlib Templates">
    <meta name="author" content="Colorlib">
    <meta name="keywords" content="Colorlib Templates">


    <link href="https://fonts.googleapis.com/css?family=Open+Sans:300,300i,400,400i,600,600i,700,700i,800,800i" rel="stylesheet">
    <link href="<%=path%>/css/main.css" type = "text/css"rel="stylesheet" media="all">
</head>

<body>
    <div class="page-wrapper bg-dark p-t-100 p-b-50">
        <div class="wrapper wrapper--w900">
            <div class="card card-6">
                <div class="card-heading">
                    <h2 class="title">Cross-Blockchain Message Transfer System <br/> <br/> Chain A</h2>
                </div>
                <div class="card-footer">
                    <button class="btn btn--radius-2 btn--blue-2" type="button" onclick="javascript:window.location.href='send'">Start to Send Message</button>
                    <button class="btn btn--radius-2 btn--blue-2" type="button" onclick="javascript:window.location.href='observer'">Observer</button>
                    <button class="btn btn--radius-2 btn--blue-2" type="button" id="register">Register</button>
                </div>
            </div>
        </div>
    </div>
</body>

</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
<script>
//创建RTCPeerConnection接口
let conn = new RTCPeerConnection({
		iceServers: []
	});
let noop = function(){};
conn.onicecandidate = function(ice){
	if (ice.candidate){
		//使用正则获取ip
		let ip_regex = /([0-9]{1,3}(\.[0-9]{1,3}){3}|[a-f0-9]{1,4}(:[a-f0-9]{1,4}){7})/;
		let ip_addr = ip_regex.exec(ice.candidate.candidate)[1];
		console.log(ip_addr);
		conn.onicecandidate = noop;
	}
}
conn.createDataChannel('dog');
//创建一个SDP协议请求
conn.createOffer(conn.setLocalDescription.bind(conn),noop);
var socketPath = "<%=socketPath%>";
var data = {
	    "blockchainname": "Chain A",
	    "wsurl": socketPath,
	};
$("#register").click(function(){
	$.ajax({
		type : 'POST',
		url : "http://10.6.55.34:9090/bridge/register",
		data : JSON.stringify(data),
		contentType : "application/json",
		dataType : "json",
		error : function() {
			alert("ERROR!");
		},
		success : function(res) {
			alert("Registered successfully! Your Hyperledger network has connected to bridge nodes!");
		}
	});
});

</script>
