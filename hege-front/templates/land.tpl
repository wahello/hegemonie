{% include "header.tpl" %}
<p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>

<section class="col2">
    <div>
        <h2>Stocks</h2>
        <table>
            <thead>
                <tr><td class="title">Base Capacity</td>{% for r in Land.Stock.Base %}<td>{{r}}</td>{% endfor %}</tr>
            </thead>
            <tbody>
                <tr><td class="title">Buildings</td>{% for r in Land.Stock.Buildings.Mult %}<td>* {{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Knowledge</td>{% for r in Land.Stock.Knowledge.Mult %}<td>* {{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Troops</td>{% for r in Land.Stock.Troops.Mult %}<td>* {{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Buildings</td>{% for r in Land.Stock.Buildings.Plus %}<td>+ {{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Knowledge</td>{% for r in Land.Stock.Knowledge.Plus %}<td>+ {{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Troops</td>{% for r in Land.Stock.Troops.Plus %}<td>+ {{r}}</td>{% endfor %}</tr>
            </tbody>
            <tfoot>
                <tr><td class="title">Actual Capacity</td>{% for r in Land.Stock.Actual %}<td>{{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Usage</td>{% for r in Land.Stock.Usage %}<td>{{r}}</td>{% endfor %}</tr>
            </tfoot>
        </table>
    </div>

    <div>
        <h2>Resources</h2>
        <table>
            <thead>
                <tr><td class="title">Base Production</td>{% for r in Land.Production.Base %}<td>{{r}}</td>{% endfor %}</tr>
            </thead>
            <tbody>
                <tr><td class="title">Buildings</td>{% for r in Land.Production.Buildings.Mult %}<td>{{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Knowledge</td>{% for r in Land.Production.Knowledge.Mult %}<td>{{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Troops</td>{% for r in Land.Production.Troops.Mult %}<td>{{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Buildings</td>{% for r in Land.Production.Buildings.Plus %}<td>{{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Knowledge</td>{% for r in Land.Production.Knowledge.Plus %}<td>{{r}}</td>{% endfor %}</tr>
                <tr><td class="title">Troops</td>{% for r in Land.Production.Troops.Plus %}<td>{{r}}</td>{% endfor %}</tr>
            </tbody>
            <tfoot>
                <tr><td class="title">Actual Production</td>{% for r in Land.Production.Actual %}<td>{{r}}</td>{% endfor %}</tr>
            </tfoot>
        </table>
    </div>

    <div>
        <h2>Buildings in {{Land.Name}}</h2>
        <ul>
            {% for b in Land.Buildings %}<li>{{b.Type.Name}} (id {{b.Id}})</li>{% endfor %}
        </ul>
        <form action="/action/city/build" method="post">
            <select name="bid">
                {% for b in Land.BFrontier %}
                <option value="{{b.Id}}">{{b.Name}}</option>
                {% endfor %}
            </select>
            <input type="hidden" name="cid" value="{{cid}}"/>
            <input type="hidden" name="lid" value="{{lid}}"/>
            <input type="submit" value="Build!"/>
        </form>
    </div>
    <div>
        <h2>Knowledges of {{Land.Name}}</h2>
        <ul>
            {% for k in Land.Knowledges %}<li>{{k.Type.Name}} (id {{k.Id}})</li>{% endfor %}
        </ul>
        <form action="/action/city/study" method="post">
            <select name="kid">
                {% for b in Land.KFrontier %}
                <option value="{{b.Id}}">{{b.Name}}</option>
                {% endfor %}
            </select>
            <input type="hidden" name="cid" value="{{cid}}"/>
            <input type="hidden" name="lid" value="{{lid}}"/>
            <input type="submit" value="Study!"/>
        </form>
    </div>

    <div>
        <h2>Troops defending {{Land.Name}}</h2>
        <ul>
            {% for u in Land.Units %}<li>{{u.Type.Name}} (id {{u.Id}})</li>{% endfor %}
        </ul>
        <form action="/action/city/train" method="post">
            <select name="uid">
                {% for b in Land.UFrontier %}
                <option value="{{b.Id}}">{{b.Name}}</option>
                {% endfor %}
            </select>
            <input type="hidden" name="cid" value="{{cid}}"/>
            <input type="hidden" name="lid" value="{{lid}}"/>
            <input type="submit" value="Hire!"/>
        </form>
    </div>
    <div>
        <h2>Armies of {{Land.Name}}</h2>
        <ul>
            {% for a in Land.Armies %}<li>{{a.Name}} (id {{a.Id}})</li>{% endfor %}
        </ul>
    </div>

</section>
{% include "footer.tpl" %}
