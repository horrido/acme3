<div id="content">
<h1>Your Profile</h1>
<br>
{{if .flash.error}}
<h3>{{.flash.error}}</h3>
<br>
{{end}}
{{if .flash.notice}}
<h3>{{.flash.notice}}</h3>
<br>
{{end}}
<form method="POST">
<table>
<tr>
    <td>First name: {{if .Errors.First}}{{.Errors.First}}{{end}}</td>
    <td><input name="first" type="text" value="{{.User.First}}" /></td>
</tr>
<tr>
    <td>Last name:</td>
    <td><input name="last" type="text" value="{{.User.Last}}"/></td>
</tr>
<tr>
    <td>Email address: {{if .Errors.Email}}{{.Errors.Email}}{{end}}</td>
    <td><input name="email" type="text" value="{{.User.Email}}"/></td>
</tr>
<tr>      
    <td>Current password: {{if .Errors.Current}}{{.Errors.Current}}{{end}}</td>
    <td><input name="current" type="password" /></td>
</tr>
<tr>
<td>Optional:</td>
</tr>
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
    <td>&nbsp;</td><td><input type="submit" value="Update" /></td>
</tr>
</table>
<a href="http://localhost:8080/user/remove">Remove account</a>
</form>
</div>