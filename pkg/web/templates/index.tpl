{% include "header.tpl" %}

<div>
    <p>Online management RPG game</p>
    <form action="/action/login" method="post">
        <input type="text" name="email" value=""/>
        <input type="password" name="password" value=""/>
        <input type="submit" value="Enter"/>
    </form>
</div>

{% include "footer.tpl" %}