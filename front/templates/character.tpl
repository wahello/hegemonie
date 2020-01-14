{% include "header.tpl" %}

<div class="large">
    <ul>
        {% for c in Character.OwnerOf %}
        <li><a href="/game/land?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>
        {% endfor %}
        {% for c in Character.DeputyOf %}
        <li><a href="/game/land?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>
        {% endfor %}
    </ul>
</div>

{% include "footer.tpl" %}