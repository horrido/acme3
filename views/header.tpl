<!DOCTYPE html>

<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8" />
<title>StarDust by Free Css Templates</title>
<meta name="keywords" content="" />
<meta name="description" content="" />
<link href="/static/css/default.css" rel="stylesheet" type="text/css" />
<link rel="stylesheet" href="/static/js/jquery-ui.min.css">
<script src="/static/js/external/jquery/jquery.js"></script>
<script src="/static/js/jquery-ui.min.js"></script>
<script>
$(function() {
	$( "input[type=submit], button" )
		.button()
	$("select").selectmenu();
});
</script>
</head>

<body>
<!-- start header -->
<div id="header-bg">
	<div id="header">
		<div align="right">{{if .InSession}}
		Welcome, {{.First}} [<a href="http://{{.domainname}}/user/logout">Logout</a>|<a href="http://{{.domainname}}/user/profile">Profile</a>]
		{{else}}
		[<a href="http://{{.domainname}}/user/login/home">Login</a>]
		{{end}}
		</div>
		<div id="logo">
			<h1><a href="#">StarDust<sup></sup></a></h1>
			<h2>Designed by FreeCSSTemplates</h2>
		</div>
		<div id="menu">
			<ul>
				<li class="active"><a href="http://{{.domainname}}/home">home</a></li>
				<li><a href="http://{{.domainname}}/photos">photos</a></li>
				<li><a href="http://{{.domainname}}/about">about</a></li>
				<li><a href="#">links</a></li>
				<li><a href="#">contact </a></li>
			</ul>
		</div>
	</div>
</div>
<!-- end header -->
<!-- start page -->
<div id="page">