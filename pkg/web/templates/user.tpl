{% include "header.tpl" %}

<div>
    <h2>Cities</h2>
    <ul>{% for char in User.Characters %}{% for city in char.Cities %}
        <li><a href="/game/land/overview?cid={{char.Id}}&lid={{city.Id}}">{{city.Name}}</a>
            as <em>{{char.Name}}</em></li>{% endfor %}{% endfor %}
    </ul>
</div>

<div>
    <h2>Profile</h2>
    <ul>
        <li>Name: {{User.Name}}</li>
        <li>Id: {{User.Id}}</li>
        <li>Admin: {{User.Admin}}</li>
        <li>E-Mail: {{User.Mail}}</li>
    </ul>
</div>

<div>
    <h2>Characters</h2>
    <ul>{% for char in User.Characters %}
        <li><a href="/game/character?cid={{char.Id}}">{{char.Name}}</a></li>{% endfor %}
    </ul>
</div>

{% include "footer.tpl" %}