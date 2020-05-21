{% include "header.tpl" %}

<div><h2>Production</h2>
    <form action="/action/produce" method="post"><input type="submit" value="Produce"/></form>
</div>

<div><h2>Movement</h2>
    <form action="/action/move" method="post"><input type="submit" value="Movement"/></form>
</div>

<div><h2>Scoreboard</h2>
    <table>
        <thead>
        <tr>
            <td>Score</td>
            <td>Name</td>
            <td>Cult</td>
            <td>Chaos</td>
            <td>Alignment</td>
            <td>Ethny</td>
            <td>Politics</td>
        </tr>
        </thead>
        <tbody>
        {% for s in Scores %}
        <tr>
            <td>{{s.Score}}</td>
            <td>{{s.Name}}</td>
            <td>{{s.Cult}}</td>
            <td>{{s.Chaos}}</td>
            <td>{{s.Alignment}}</td>
            <td>{{s.Ethny}}</td>
            <td>{{s.Politics}}</td>
        </tr>
        {% endfor %}
        </tbody>
    </table>
</div>

<div>
    {% include "map.tpl" %}
</div>

{% include "footer.tpl" %}