{% include "header.tpl" %}
<script src="/static/hege-map.js"></script>
<script>
// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
window.addEventListener("load", function() {
  let svg1 = document.getElementById('interactive-map');
  let armies = [
  {% for a in Land.Assets.Armies %}
      {"id":{{a.Id}}, "cell":{{a.Location}}},
  {% endfor %}
  ];
  drawMapWithArmies(svg1, "calaquyr", armies)
    .then(map => {
        hightlightCell(svg1, {{Land.Location}});
        return map;
    })
    .catch(err => { console.log(err); });
});
</script>
{% include "map.tpl" %}
<div><h2>Armies</h2>
    <ul>{% for a in Land.Assets.Armies %}
        <li><a href="/game/army?cid={{cid}}&lid={{lid}}&aid={{a.Id}}">{{a.Name}}</a> (id {{a.Id}})</li>{% endfor %}
    </ul>
</div>
<div><h2>Local Units</h2>
    {% for u in Land.Assets.Units %}
    <p>
    {{u.Type.Name}} (id {{u.Id}})
    Health({{u.Health}} / {{u.Type.Health}})
    ETA({{u.Ticks}} ticks)
    </p>
    {% endfor %}
</div>
<div><h2>Buildings</h2>
    {% for b in Land.Assets.Buildings %}
    <p>
    {{b.Type.Name}} (id {{b.Id}})
    ETA({{b.Ticks}} ticks)
    </p>
    {% endfor %}
</div>
<div><h2>Knowledge</h2>
    {% for k in Land.Assets.Knowledges %}
    <p>
    {{k.Type.Name}} (id {{k.Id}})
    ETA({{k.Ticks}} ticks)
    </p>
    {% endfor %}
</div>
{% include "footer.tpl" %}
