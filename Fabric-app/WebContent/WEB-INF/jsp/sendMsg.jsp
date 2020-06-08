<%
String path = request.getContextPath();
String basePath = request.getScheme()+"://"+request.getServerName()+":"+request.getServerPort()+path+"/";
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
                <form method="POST" action="sendMessage" enctype="multipart/form-data">
                <div class="card-body">
                   
                        <div class="form-row">
                            <div class="name">sender's ID</div>
                            <div class="value">
                                <input class="input--style-6" type="text" name="senderID"/>
                            </div>
                        </div>
                        <div class="form-row">
                            <div class="name">receiver's ID</div>
                            <div class="value">
                                <input class="input--style-6" type="text" name="receiverID"/>
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="name">Source NetID</div>
                            <div class="value">
                                <div class="input-group">
                                    <input class="input--style-6" type="text" name="sourceID" placeholder="Source Chain NetID"/>
                                </div>
                            </div>
                        </div>
                        <div class="form-row">
                            <div class="name">Destination NetID</div>
                            <div class="value">
                                <div class="input-group">
                                    <input class="input--style-6" type="text" name="destinationID" placeholder="Target Chain NetID"/>
                                </div>
                            </div>
                        </div>
                        <div class="form-row">
                            <div class="name">Message</div>
                            <div class="value">
                                <div class="input-group">
                                    <textarea class="textarea--style-6" name="message" placeholder="Message sent to the Target Chain"></textarea>
                                </div>
                            </div>
                        </div>
                        <div class="form-row">
                            <div class="name">File</div>
                            <div class="value">
                                <div class="input-group">
                                    File to upload: <input type="file" name="file">
                                </div>
                            </div>
                        </div>
                        
                </div>
                <div class="card-footer">
                    <button class="btn btn--radius-2 btn--blue-2" type="submit">Send Message</button>
                    <a class="btn btn--radius-2 btn--blue-2" href="index">Back</a>
                </div>
                </form>
                
            </div>
        </div>
    </div>

</body>

</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
<script>

</script>
