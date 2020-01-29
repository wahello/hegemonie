{% include "header.tpl" %}

    {% for c in User.Characters %}
    <div>
        <h2>{{c.Name}}</h2>
        <p><a href="/game/character?cid={{c.Id}}">Play!</a></p>
    </div>
    {% endfor %}

{% include "footer.tpl" %}