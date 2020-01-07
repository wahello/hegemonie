{% include "header.tpl" %}

<p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>

<section class="col2">
    <div>
        <h2>Characters</h2>
        <ul>
            {% for c in User.Characters %}
            <li><a href="/game/character?cid={{c.Id}}">{{c.Name}}</a></li>
            {% endfor %}
        </ul>
    </div>
    <div>
        <h2>Admin</h2>
        <form action="/action/produce" method="post"><input type="submit" value="Produce"/></form>
        <form action="/action/move" method="post"><input type="submit" value="Movement"/></form>
    </div>
</section>

{% include "footer.tpl" %}