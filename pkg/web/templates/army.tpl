{% include "header.tpl" %}
<script src="/static/hege-map.js"></script>
<script>
// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

window.addEventListener("load", function() {
  let svg1 = document.getElementById('interactive-map');
  let armies = [{"id":{{aid}}, "cell":{{Army.Location}}}];
  drawMapWithCities(svg1, "calaquyr",
    function (position) {
      document.getElementById('Location').value = position;
      //document.getElementById('CityId').value = null;
      document.getElementById('CityName').value = null;
    },
    function (position, cityId, cityName) {
      document.getElementById('Location').value = position;
      //document.getElementById('CityId').value = cityId;
      document.getElementById('CityName').value = cityName;
    })
    .then(map => {
      hightlightCell(svg1, {{Land.Location}});
      return map;
    })
    .then(map => {
      return patchWithArmies(svg1, map, armies);
    })
    .catch(err => { console.log(err); });
});
</script>
{% include "map.tpl" %}

<div><h2>Actions</h2>
<form action="/action/army/command" method="post">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="location" id="Location" value=""/>
    Target: <input type="text" id="CityName" value=""/><br/>
    <fieldset>
        <legend>Action</legend>
        <input type="radio" name="action" value="move" checked> Move<br/>
        <input type="radio" name="action" value="wait"/> Wait<br/>
        <input type="radio" name="action" value="attack"> Attack<br/>
        <input type="radio" name="action" value="defend"> Defend<br/>
        <input type="radio" name="action" value="massacre"/> Massacre<br/>
        <input type="radio" name="action" value="break"/> Break<br/>
        <input type="radio" name="action" value="overlord"/> Overlord<br/>
        <input type="radio" name="action" value="deposit"/> Deposit<br/>
        <input type="radio" name="action" value="disband"/> Disband<br/>
    </fieldset>
    <input type="submit" value="Go"/>
</form>
</div>

<div><h2>Commands</h2>
{% for cmd in Commands %}
    <p>{{cmd.Order}}:
    {% if cmd.CommandId == 0 %}
    Take selfies on
    {% elif cmd.CommandId == 1 %}
    Disband on
    {% elif cmd.CommandId == 2 %}
    Wait on
    {% elif cmd.CommandId == 3 %}
    Move to
    {% elif cmd.CommandId == 4 %}
    Attack and Overlord
    {% elif cmd.CommandId == 5 %}
    Attacke and Break a building on
    {% elif cmd.CommandId == 6 %}
    Attack and Massacre on
    {% elif cmd.CommandId == 7 %}
    Drop resources at
    {% elif cmd.CommandId == 8 %}
    Disband in
    {% else %}
    ?
    {% endif %}
    {{cmd.CityName}} (id {{cmd.CityId}}, located at {{cmd.Location}})</p>
{% endfor %}
</div>

<div><h2>Payload</h2>
    <p>
    Gold ({{Army.Stock.R1}}),
    Cereals ({{Army.Stock.R2}}),
    Livestock ({{Army.Stock.R3}}),
    Wood ({{Army.Stock.R4}}),
    Stone ({{Army.Stock.R5}})
    </p>
</div>

<div><h2>Enroll</h2>
    {% for u in Land.Assets.Units %}
    {% if u.Ticks == 0 %}
    <p>{{u.Type.Name}} (id {{u.Id}}) Health({{u.Health}}/{{u.Type.Health}})</p>
    {% endif %}
    {% endfor %}
</div>

<div><h2>Units</h2>
    {% for u in Army.Units %}
    <p>{{u.Type.Name}} (id {{u.Id}}) Health({{u.Health}}/{{u.Type.Health}})</p>
    {% endfor %}
</div>

{% include "footer.tpl" %}
