<!--
Copyright (C) 2018-2019 Hegemonie's AUTHORS
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
-->
<!DOCTYPE html>
<html lang="fr">
<head><meta charset="UTF-8"><title>Hegemonie</title>
<meta name="author" content="Jean-Francois Smigielski"/>
<meta name="theme-color" content="#FFF"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<meta name="description" content="${description}"/>
<link rel="stylesheet" href="/static/style.css"/>
</head>
<body>
    <header>
        <h1>{{Title}}</h1>
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

        {% if cid %}<a href="/game/character?cid={{ cid }}">Character</a>{% endif %}
        {% if userid %}<a href="/game/user">User</a>{% endif %}
        {% if User.Admin %}<a href="/game/admin">Admin</a>{% endif %}
        {% if userid %}<a href="/action/logout">Log-Out</a>{% endif %}
    </nav>
    <aside>
        <p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>
    </aside>
    <main>