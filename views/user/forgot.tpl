<div id="content">
<h1>Reset password</h1>
<br>
{{if .flash.error}}
<h3>{{.flash.error}}</h3>
<br>
{{end}}
<form method="POST">
<table>
<tr>
    <td>Email address: {{if .Errors.email}}{{.Errors.email}}{{end}}</td>
    <td><input name="email" type="text" autofocus /></td>
</tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Request reset" /></td>
</tr>
</table>
</form>
</div>
