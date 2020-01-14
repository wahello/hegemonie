{% include "header.tpl" %}

    <div>
        <h2>Characters</h2>
        <ul>
            {% for c in User.Characters %}
            <li><a href="/game/character?cid={{c.Id}}">{{c.Name}}</a></li>
            {% endfor %}
        </ul>
    </div>
    <div>
        <h2>Admin</h2>
        <form action="/action/produce" method="post"><input type="submit" value="Produce"/></form>
        <form action="/action/move" method="post"><input type="submit" value="Movement"/></form>
    </div>

    <div class="large">
        <h2>Scoreboard</h2>
        <table>
            <thead>
            <tr>
                <td>Id</td>
                <td>Name</td>
                <td>Score</td>
            </tr>
            </thead>
            <tbody>
            {% for s in Score %}
            <tr>
                <td>{{s.Id}}</td>
                <td>{{s.Name}}</td>
                <td>{{s.Score}}</td>
            </tr>
            {% endfor %}
            </tbody>
        </table>
    </div>

{% include "footer.tpl" %}