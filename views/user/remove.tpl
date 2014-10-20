<div id="content">
<h1>Remove User Account</h1>
<br>
{{if .flash.error}}
<h3>{{.flash.error}}</h3>
<br>
{{end}}
<p><font size="3">Caution: all related transactions and data will also be removed. Are you sure?</font></p>
<form method="POST">
<table>
<tr>      
    <td>Current password: {{if .Errors.current}}{{.Errors.current}}{{end}}</td>
    <td><input name="current" type="password" /></td>
</tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Remove" /></td><td><a href="http://localhost:8080/home">Cancel</a></td>
</tr>
</table>
</form>
</div>
