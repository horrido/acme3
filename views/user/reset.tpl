<div id="content">
<h1>Reset password</h1>
<br>
{{if .flash.error}}
<h3>{{.flash.error}}</h3>
<br>
{{end}}{{if .flash.notice}}
<h3>{{.flash.notice}}</h3>
<br>
{{end}}
<form method="POST">
<table>
<tr>      
    <td>New password (must be at least 6 characters): {{if .Errors.password}}{{.Errors.password}}{{end}}</td>
    <td><input name="password" type="password" /></td>
</tr>
<tr>      
    <td>Confirm new password: {{if .Errors.password2}}{{.Errors.password2}}{{end}}</td>
    <td><input name="password2" type="password" /></td>
</tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Reset password" /></td>
</tr>
</table>
</form>
</div>