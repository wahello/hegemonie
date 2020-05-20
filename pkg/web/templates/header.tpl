<!DOCTYPE html>
<html lang="fr">
<head><meta charset="UTF-8"><title>{{Title}}</title>
<meta name="author" content="Jean-Francois Smigielski"/>
<meta name="theme-color" content="#FFF"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<meta name="description" content="Hegemone {{Land.Name}}"/>
<link rel="icon" type="image/png" href="https://www.hegemonie.be/static/favicon.png"/>
<link rel="shortcut icon" type="image/png" href="https://www.hegemonie.be/static/favicon.png"/>
<link rel="apple-touch-icon" type="image/png" href="https://www.hegemonie.be/static/favicon-apple.png"/>
<link rel="stylesheet" type="text/css" href="/static/style.css"/>
</head>
<body>
    <header>
        {% if lid %}
        <h1>{{Land.Name}}</h1>
        <h2>{{Character.Name}}</h2>
        {% elif cid %}
        <h1>{{Character.Name}}</h1>
        <h2>{{User.Name}}</h2>
        {% elif userid %}
        <h1>{{User.Name}}</h1>
        {% else %}
        <h1>{{Title}}</h1>
        {% endif %}
    </header>
    <nav>
        {% if lid %}
        <a href="/game/land/overview?cid={{ cid }}&lid={{ lid }}">Overview</a>
        <a href="/game/land/budget?cid={{ cid }}&lid={{ lid }}">Budget</a>
        <a href="/game/land/units?cid={{ cid }}&lid={{ lid }}">Troops</a>
        <a href="/game/land/armies?cid={{ cid }}&lid={{ lid }}">Armies</a>
        <a href="/game/land/buildings?cid={{ cid }}&lid={{ lid }}">Building</a>
        <a href="/game/land/knowledges?cid={{ cid }}&lid={{ lid }}">Science</a>
        <br/>
        {% endif %}

        {% if userid %}
        <a href="/game/user">User</a>
            {% if User.Admin %}
            <a href="/game/admin">Admin</a>
            {% endif %}
        <a href="/action/logout">Log-Out</a>
        {% endif %}

        {% if cid %}
        <a href="/game/character?cid={{ cid }}">Character</a>
        {% endif %}

    </nav>
    <aside>
        <p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>
    </aside>
    <main>