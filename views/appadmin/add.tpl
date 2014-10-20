<div id="content">
<h1>Add</h1>
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
<tr><td>Id:</td><td><input name="-" value="{{.User.Id}}" size="10" readonly /></td></tr>
<tr><td>First: {{if .Errors.First}}{{.Errors.First}}{{end}}</td><td><input name="first" value="{{.User.First}}" size="30" /></td></tr>
<tr><td>Last:</td><td><input name="last" value="{{.User.Last}}" size="30" /></td></tr>
<tr><td>Email:  {{if .Errors.Email}}{{.Errors.Email}}{{end}}</td><td><input name="email" value="{{.User.Email}}" size="30" /></td></tr>
<tr><td>Password: {{if .Errors.Password}}{{.Errors.Password}}{{end}}</td><td><input name="password" value="{{.User.Password}}" size="30" /></td></tr>
<tr><td>Reg key</td><td><input name="reg_key" value="{{.User.Reg_key}}" size="50" readonly /></td></tr>
<tr><td>Reg date</td><td><input name="reg_date" value="{{.User.Reg_date}}" size="50" readonly /></td></tr>
<tr><td>Reset key</td><td><input name="reset_key" value="{{.User.Reset_key}}" size="50" readonly /></td></tr>
<tr><td>&nbsp;</td></tr>
<tr>
    <td>&nbsp;</td><td><input type="submit" value="Add" /></td>
</tr>
</table>
</form>
<h4><a href="http://{{.domainname}}/appadmin/index/{{.parms}}">Return to Index</a></h4>
</div>