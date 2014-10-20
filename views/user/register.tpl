<div id="content">
<h1>Please register</h1>
<br>
{{if .flash.error}}
<h3>{{.flash.error}}</h3>
<br>
{{end}}
<form method="POST">
<table>
<tr>
    <td>First name: {{if .Errors.First}}{{.Errors.First}}{{end}}</td>
    <td><input name="first" type="text" value="{{.User.First}}" autofocus /></td>
</tr>
<tr>
    <td>Last name:</td>
    <td><input name="last" type="text" value="{{.User.Last}}" /></td>
</tr>
<tr>
    <td>Email address: {{if .Errors.Email}}{{.Errors.Email}}{{end}}</td>
    <td><input name="email" type="text" value="{{.User.Email}}" /></td>
</tr>
<tr>      
    <td>Password (must be at least 6 characters): {{if .Errors.Password}}{{.Errors.Password}}{{end}}</td>
    <td><input name="password" type="password" /></td>
</tr>
<tr>      
    <td>Confirm password: {{if .Errors.Confirm}}{{.Errors.Confirm}}{{end}}</td>
    <td><input name="password2" type="password" /></td>
</tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Register" /></td>
</tr>
</table>
</form>
</div>
