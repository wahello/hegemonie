{% include "header.tpl" %}

<div>
    <ul>
        {% for c in Cities %}
        <li><a href="/game/land/overview?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>
        {% endfor %}
    </ul>
</div>

<div>
    {% for msg in Log %}{{msg}}<br/>
    {% endfor %}
</div>

{% include "footer.tpl" %}