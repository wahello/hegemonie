{% include "header.tpl" %}

    <div>
        <h2>Buildings</h2>
        <ul>
            {% for b in Land.Assets.Buildings %}<li>{{b.Type.Name}} (id {{b.Id}})</li>{% endfor %}
        </ul>
    </div>

    <div>
        <h2>Construction</h2>
        <form action="/action/city/build" method="post">
            <select name="bid">
                {% for b in Land.Evol.BFrontier %}
                <option value="{{b.Id}}">{{b.Name}}</option>
                {% endfor %}
            </select>
            <input type="hidden" name="cid" value="{{cid}}"/>
            <input type="hidden" name="lid" value="{{lid}}"/>
            <input type="submit" value="Start!"/>
        </form>
    </div>

{% include "footer.tpl" %}
