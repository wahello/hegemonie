{% include "header.tpl" %}

<h1>Profile of {{User.Meta.Name}}</h1>

<p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>

<h2>Characters</h2>
<ul>{% for c in User.Characters %}
    <li><a href="/game/character?cid={{c.Id}}">{{c.Name}}</a></li>{% endfor %}
</ul>

<h2>Admin</h2>
<p>Logged as {{User.Meta.Name}}.</p>
<p>Your email is {{User.Meta.Email}}</p>
<form action="/action/logout" method="post"><input type="submit" value="Log Out"/></form>
<form action="/action/produce" method="post"><input type="submit" value="Produce"/></form>
<form action="/action/move" method="post"><input type="submit" value="Movement"/></form>

{% include "footer.tpl" %}