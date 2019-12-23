{% include "header.tpl" %}

<h1>The cities managed by {{Character.Meta.Name}}</h1>

<ul>{% for c in Character.OwnerOf %}
    <li><a href="/game/land?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>{% endfor %}
</ul>
<ul>{% for c in Character.DeputyOf %}
    <li><a href="/game/land?cid={{Character.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>{% endfor %}
</ul>

{% include "footer.tpl" %}
