{% include "header.tpl" %}

<h1>Status of {{Land.Name}}</h1>

<h2>Resources</h2>
<ul>
    <li>Stock: {% for r in Land.Stock %}{{r}},{% endfor %}</li>
    <li>Production: {% for r in Land.Production %}{{r}},{% endfor %}</li>
</ul>

<h2>Buildings</h2>
<ul>
    {% for b in Land.Buildings %}
    <li>{{b.Type.Name}} (id {{b.Id}})</li>
    {% endfor %}
</ul>

<h2>Troops defending</h2>
<ul>
    {% for u in Land.Units %}
    <li>{{u.Type.Name}} (id {{u.Id}})</li>
    {% endfor %}
</ul>

<h2>Stocks</h2>
<table>
    <tr>
        <td>Base Capacity</td>
        {% for r in Land.Stock.BaseProduction %}<td>{{r}}</td>{% endfor %}
    </tr>

    <th><td>Mult</td></th>
    <tr>
        <td>Buildings</td>
        {% for r in Land.Stock.Buildings.Mult %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Knowledge</td>
        {% for r in Land.Stock.Knowledge.Mult %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Troops</td>
        {% for r in Land.Stock.Troops.Mult %}<td>{{r}}</td>{% endfor %}
    </tr>

    <th><td>Bonus</td></th>
    <tr>
        <td>Buildings</td>
        {% for r in Land.Stock.Buildings.Plus %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Knowledge</td>
        {% for r in Land.Stock.Knowledge.Plus %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Troops</td>
        {% for r in Land.Stock.Troops.Plus %}<td>{{r}}</td>{% endfor %}
    </tr>

    <th><td>Total</td></th>
    <tr>
        <td>Actual Capacity</td>
        {% for r in Land.Stock.ActualProduction %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Usage</td>
        {% for r in Land.Stock.Usage %}<td>{{r}}</td>{% endfor %}
    </tr>
</table>

<h2>Resources</h2>
<table>
    <tr>
        <td>Base Production</td>
        {% for r in Land.Production.BaseProduction %}<td>{{r}}</td>{% endfor %}
    </tr>

    <th><td>Multipliers</td></th>
    <tr>
        <td>Buildings</td>
        {% for r in Land.Production.Buildings.Mult %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Knowledge</td>
        {% for r in Land.Production.Knowledge.Mult %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Troops</td>
        {% for r in Land.Production.Troops.Mult %}<td>{{r}}</td>{% endfor %}
    </tr>

    <th><td>Bonus</td></th>
    <tr>
        <td>Buildings</td>
        {% for r in Land.Production.Buildings.Plus %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Knowledge</td>
        {% for r in Land.Production.Knowledge.Plus %}<td>{{r}}</td>{% endfor %}
    </tr>
    <tr>
        <td>Troops</td>
        {% for r in Land.Production.Troops.Plus %}<td>{{r}}</td>{% endfor %}
    </tr>

    <th><td>Total</td></th>
    <tr>
        <td>Actual Production</td>
        {% for r in Land.Production.ActualProduction %}<td>{{r}}</td>{% endfor %}
    </tr>
</table>

<h2>Production</h2>

{% include "footer.tpl" %}