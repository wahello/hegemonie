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
      document.getElementById('Position').value = position;
      document.getElementById('CityId').value = null;
      document.getElementById('CityName').value = null;
    },
    function (position, cityId, cityName) {
      document.getElementById('Position').value = position;
      document.getElementById('CityId').value = cityId;
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
<form>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="position" id="Position" value=""/>
    <input type="hidden" name="cityId" id="CityId" value=""/>
    <input type="text" id="CityName" value=""/>
    <table class="action-set">
        <tbody>
            <tr>
                <td><input type="submit" value="Move"/></td>
                <td><p>Just move there. Wait'n see.</p></td>
            </tr>
            <tr>
                <td><input type="submit" value="Attack"/></td>
                <td><p>Move to the given City and attack it. Will you dare?</p></td>
            </tr>
            <tr>
                <td><input type="submit" value="Defend"/></td>
                <td><p>Move to the given City and join its defence, or wait there for an assault to start against it.</p></td>
            </tr>
        </tbody>
        <tfoot>
            <tr>
                <td><input type="submit" value="Disband"/></td>
                <td><p>Cancel the Army and give both its freight and its troops
                 to the local City.</p></td>
            </tr>
            <tr>
                <td><input type="submit" value="Cancel"/></td>
                <td><p>Cancel the Army and return both its freight and its troops
                 back to {{Land.Name}}. The action only works if the army is at home.</p></td>
            </tr>
        </tfoot>
    </table>
</form>
</div>

<div><h2>Commands</h2>
{% for cmd in Commands %}
<p>{{cmd.Target}}</p>
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
