<div id="content">
<h1>User Table</h1>
<h4><a href="http://{{.domainname}}/appadmin/add/{{.parms}}">Add a Record</a></h4>
<div>
<form method="POST">
<table>
<tr>
<td>
<select name="comparefield">
	<option value="id" selected="selected">Id</option>
	<option value="first">First</option>
	<option value="last">Last</option>
	<option value="email">Email</option>
	<option value="reg_key">Reg key</option>
	<option value="reg_date">Reg date</option>
</select>
<select name="compareop">
	<option value="__exact">=</option>
	<option value="__not__exact">!=</option>
	<option value="__lt">&lt;</option>
	<option value="__lte">&lt;=</option>
	<option value="__gt">&gt;</option>
	<option value="__gte" selected="selected">&gt;=</option>
	<option value="__contains">contains</option>
	<option value="__not__contains">not contains</option>
	<option value="__icontains">icontains</option>
	<option value="__not__icontains">not icontains</option>
</select>
</td>
<td> {{if .Errors.Compareval}}{{.Errors.Compareval}}{{end}}<input name="compareval" value="0"
	title="All search terms are compared as text strings. 'Reg date' has a date format of yyyy-mm-dd." /></td>
<td><input type="submit" value="Search" /></td>
</tr>
</table>
</form>
<br>
</div>
Total: {{.count}} records – {{.query}} ordered by {{.order}}
<table border="1" style="width:800px">
<tr>
<th style="width:100px"><a href="http://{{.domainname}}/appadmin/index/{{if eq .order "id"}}-{{end}}id!{{.offset}}!{{.query}}">Id</a></th>
<th style="width:100px"><a href="http://{{.domainname}}/appadmin/index/{{if eq .order "first"}}-{{end}}first!{{.offset}}!{{.query}}">First</a></th>
<th style="width:100px"><a href="http://{{.domainname}}/appadmin/index/{{if eq .order "last"}}-{{end}}last!{{.offset}}!{{.query}}">Last</a></th>
<th style="width:100px"><a href="http://{{.domainname}}/appadmin/index/{{if eq .order "email"}}-{{end}}email!{{.offset}}!{{.query}}">Email</a></th>
<th style="width:100px">Password</th>
<th style="width:200px">Reg key</th>
<th style="width:150px"><a href="http://{{.domainname}}/appadmin/index/{{if eq .order "reg_date"}}-{{end}}reg_date!{{.offset}}!{{.query}}">Reg date</a></th>
<th style="width:100px">Reset key</th>
</tr>
{{.Rows}}
</table>
{{if .ShowNav}}
<br>
<div id="progressbar"></div>pointer in data set
<div align="right">
<a href="http://{{.domainname}}/appadmin/index/{{.order}}!0!{{.query}}">&lt;&lt;Start</a>&nbsp;&nbsp;&nbsp;&nbsp;
{{if .showprev}}<a href="http://{{.domainname}}/appadmin/index/{{.order}}!{{.prev}}!{{.query}}">&lt;Prev</a>&nbsp;&nbsp;&nbsp;&nbsp;{{end}}
{{if .next}}<a href="http://{{.domainname}}/appadmin/index/{{.order}}!{{.next}}!{{.query}}">Next&gt;</a>&nbsp;&nbsp;&nbsp;&nbsp;{{end}}
<a href="http://{{.domainname}}/appadmin/index/{{.order}}!{{.end}}!{{.query}}">End&gt;&gt;</a>
</div>
{{end}}
</div>
<script>
$(function() {
$( document ).tooltip();
});

$(function() {
	$( "#progressbar" ).progressbar({
		value: {{.progress}}
	});
});
</script>
