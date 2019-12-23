{% include "header.tpl" %}
<section class="col2">
    <div>
        <h1>Status of {{Land.Name}}</h1>
        <h2>Troops defending</h2>
        <ul>
            {% for u in Land.Units %}
            <li>{{u.Type.Name}} (id {{u.Id}})</li>
            {% endfor %}
        </ul>

        <h2>Buildings</h2>
        <ul>
            {% for b in Land.Buildings %}
            <li>{{b.Type.Name}} (id {{b.Id}})</li>
            {% endfor %}
        </ul>
    </div>

    <div>
        <h2>Stocks</h2>
        <table>
            <tr>
                <td>Base Capacity</td>
                {% for r in Land.Stock.Base %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>

            <th>
            <td>Multipliers</td>
            </th>
            <tr>
                <td>Buildings</td>
                {% for r in Land.Stock.Buildings.Mult %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Knowledge</td>
                {% for r in Land.Stock.Knowledge.Mult %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Troops</td>
                {% for r in Land.Stock.Troops.Mult %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>

            <th>
            <td>Bonus</td>
            </th>
            <tr>
                <td>Buildings</td>
                {% for r in Land.Stock.Buildings.Plus %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Knowledge</td>
                {% for r in Land.Stock.Knowledge.Plus %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Troops</td>
                {% for r in Land.Stock.Troops.Plus %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>

            <th>
            <td>Total</td>
            </th>
            <tr>
                <td>Actual Capacity</td>
                {% for r in Land.Stock.Actual %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Usage</td>
                {% for r in Land.Stock.Usage %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
        </table>
    </div>

    <div>
        <h2>Resources</h2>
        <table>
            <tr>
                <td>Base Production</td>
                {% for r in Land.Production.Base %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>

            <th>
            <td>Multipliers</td>
            </th>
            <tr>
                <td>Buildings</td>
                {% for r in Land.Production.Buildings.Mult %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Knowledge</td>
                {% for r in Land.Production.Knowledge.Mult %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Troops</td>
                {% for r in Land.Production.Troops.Mult %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>

            <th>
            <td>Bonus</td>
            </th>
            <tr>
                <td>Buildings</td>
                {% for r in Land.Production.Buildings.Plus %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Knowledge</td>
                {% for r in Land.Production.Knowledge.Plus %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
            <tr>
                <td>Troops</td>
                {% for r in Land.Production.Troops.Plus %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>

            <th>
            <td>Total</td>
            </th>
            <tr>
                <td>Actual Production</td>
                {% for r in Land.Production.Actual %}
                <td>{{r}}</td>
                {% endfor %}
            </tr>
        </table>
    </div>
</section>
{% include "footer.tpl" %}