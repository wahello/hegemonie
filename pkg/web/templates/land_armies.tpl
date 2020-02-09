{% include "header.tpl" %}
{% if Land.Assets.Units %}
<div><h2>Create an Army</h2>
    <p>Cancel the Army and give both its freight and its troops to the local City.
    The action only works if there is a City on the local position of the Army.</p>
    <form action="/action/army/make" method="post">
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <select name="uid">{% for u in Land.Assets.Units %}
            <option value="{{u.Id}}">{{u.Id}} / {{u.Type.Name}}</option>{% endfor %}
        </select>
        <input type="submit" value="Army!"/>
    </form>
</div>
{% endif %}

<div><h2>Armies</h2>
    <ul>{% for a in Land.Assets.Armies %}
        <li>
            <a href="/game/army?cid={{Character.Id}}&lid={{Land.Id}}&aid={{a.Id}}">{{a.Name}}</a>
        </li>{% endfor %}
    </ul>
</div>

<div><h2>Pack a Caravan</h2>
    <p>Create an army around a pile of Resources.</p>
    <form action="/action/army/caravan" method="post">
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <ul>
            <li>Resource 0: <input type="text" name="r0" value="0"/> (max {{Land.Stock.Actual.R0}})</li>
            <li>Resource 1: <input type="text" name="r1" value="0"/> (max {{Land.Stock.Actual.R1}})</li>
            <li>Resource 2: <input type="text" name="r2" value="0"/> (max {{Land.Stock.Actual.R2}})</li>
            <li>Resource 3: <input type="text" name="r3" value="0"/> (max {{Land.Stock.Actual.R3}})</li>
            <li>Resource 4: <input type="text" name="r4" value="0"/> (max {{Land.Stock.Actual.R4}})</li>
            <li>Resource 5: <input type="text" name="r5" value="0"/> (max {{Land.Stock.Actual.R5}})</li>
        </ul>
        <input type="submit" value="Caravan!"/>
    </form>
</div>
{% include "footer.tpl" %}
