{% include "header.tpl" %}

    <div>
        <h2>Defence</h2>
        <ul>
            {% for u in Land.Assets.Units %}<li>{{u.Type.Name}} (id {{u.Id}})</li>{% endfor %}
        </ul>
    </div>

    <div>
        <h2>Train</h2>
        <form action="/action/city/train" method="post">
            <select name="uid">
                {% for b in Land.Evol.UFrontier %}
                <option value="{{b.Id}}">{{b.Name}}</option>
                {% endfor %}
            </select>
            <input type="hidden" name="cid" value="{{cid}}"/>
            <input type="hidden" name="lid" value="{{lid}}"/>
            <input type="submit" value="Start!"/>
        </form>
    </div>

{% include "footer.tpl" %}
