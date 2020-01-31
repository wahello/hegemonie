{% include "header.tpl" %}
<div><h2>Defence</h2>
    <ul>{% for u in Land.Assets.Units %}
        <li>{{u.Type.Name}} (id {{u.Id}})</li>{% endfor %}
    </ul>
</div>
<div><h2>Armies</h2>
    <ul>{% for a in Land.Assets.Armies %}
        <li>{{a.Name}} (id {{a.Id}})</li>{% endfor %}
    </ul>
</div>
<div><h2>Buildings</h2>
    <ul>{% for b in Land.Assets.Buildings %}
        <li>{{b.Type.Name}} (id {{b.Id}})</li>{% endfor %}
    </ul>
</div>
<div><h2>Knowledge</h2>
    <ul>{% for k in Land.Assets.Knowledges %}
        <li>{{k.Type.Name}} (id {{k.Id}})</li>{% endfor %}
    </ul>
</div>
{% include "footer.tpl" %}
