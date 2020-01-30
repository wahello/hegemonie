{% include "header.tpl" %}

    <div>
        <h2>Armies</h2>
        <ul>
            {% for a in Land.Assets.Armies %}<li>{{a.Name}} (id {{a.Id}})</li>{% endfor %}
        </ul>
    </div>

{% include "footer.tpl" %}
