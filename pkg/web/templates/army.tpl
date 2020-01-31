{% include "header.tpl" %}
<div><h2>Armies</h2>
    <ul>{% for a in Land.Assets.Armies %}
        <li>{{a.Name}} (id {{a.Id}})</li>{% endfor %}
    </ul>
</div>

<div><h2>Disband</h2>
    <p>Cancel the Army and give both its freight and its troops to the local City.
    The action only works if there is a City on the local position of the Army.</p>
    <form action="/action/army/disband" method="post">
        <input type="hidden" name="aid" value="{{a.Id}}"/>
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <input type="submit" value="Disband!"/>
    </form>
</div>

<div><h2>Cancellation</h2>
    <p>Cancel the Army and return both its freight and its troops back to the owner City.
    The action only works if the army is on the position of its owner City.</p>
    <form action="/action/army/cancel" method="post">
        <input type="hidden" name="aid" value="{{a.Id}}"/>
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <input type="submit" value="Disband!"/>
    </form>
</div>

<div><h2>Attack</h2>
    <p>Move to the given City and attack it. will you dare?</p>
    <form action="/action/army/command" method="post">
        <input type="hidden" name="action" value="attack"/>
        <input type="hidden" name="aid" value="{{a.Id}}"/>
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <input type="submit" value="Attack!"/>
    </form>
</div>

<div><h2>Defend</h2>
    <p>Move to the given City and join its defence, or wait there for an assault to start against it.</p>
    <form action="/action/army/command" method="post">
        <input type="hidden" name="action" value="defend"/>
        <input type="hidden" name="aid" value="{{a.Id}}"/>
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <input type="submit" value="Defend!"/>
    </form>
</div>

<div><h2>Move</h2>
    <p>Just move by thegiven City.</p>
    <form action="/action/army/move" method="post">
        <input type="hidden" name="aid" value="{{a.Id}}"/>
        <input type="hidden" name="cid" value="{{Character.Id}}"/>
        <input type="hidden" name="lid" value="{{Land.Id}}"/>
        <input type="submit" value="Move!"/>
    </form>
</div>

{% include "footer.tpl" %}
