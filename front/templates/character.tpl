{% include "header.tpl" %}

<p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>

<section class="col2">
    <div>
        <ul>
            {% for c in Character.OwnerOf %}
            <li><a href="/game/land?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>
            {% endfor %}
            {% for c in Character.DeputyOf %}
            <li><a href="/game/land?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>
            {% endfor %}
        </ul>
    </div>
</section>

{% include "footer.tpl" %}