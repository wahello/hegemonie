{% include "header.tpl" %}

<div class="large">
    <ul>
        {% for c in Cities %}
        <li><a href="/game/land/overview?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>
        {% endfor %}
    </ul>
</div>

{% include "footer.tpl" %}