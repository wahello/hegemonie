{% include "header.tpl" %}
<div><h2>Knowledge</h2>
    <ul>{% for k in Land.Assets.Knowledges %}
        <li>{{k.Type.Name}} (id {{k.Id}})</li>{% endfor %}
    </ul>
</div>
<div><h2>Learn</h2>
    <form action="/action/city/study" method="post">
        <select name="kid">{% for b in Land.Evol.KFrontier %}
            <option value="{{b.Id}}">{{b.Name}}</option>{% endfor %}
        </select>
        <input type="hidden" name="cid" value="{{cid}}"/>
        <input type="hidden" name="lid" value="{{lid}}"/>
        <input type="submit" value="Start!"/>
    </form>
</div>
{% include "footer.tpl" %}
