// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

const installAction = function (obj, cls, action) {
    let elements = obj.getElementsByClassName('clickable')
    for (i = 0; i < elements.length; i++) {
        let element = elements[i];
        element.onclick = action
    }
};

const drawCircle = function (svg1, cell) {
    let c = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    if (cell.city != null && cell.city != 0) {
        c.setAttribute("class", "city clickable");
    } else {
        c.setAttribute("class", "cell clickable");
    }
    c.setAttribute("id", cell.id);
    c.setAttribute("cx", cell.x);
    c.setAttribute("cy", cell.y);
    c.setAttribute("r", 5);
    svg1.appendChild(c);
    return c;
};

const drawLine = function (svg1, src, dst) {
    let l = document.createElementNS("http://www.w3.org/2000/svg", "line");
    l.setAttribute("class", "road");
    l.setAttribute("x2", src.x);
    l.setAttribute("y2", src.y);
    l.setAttribute("x1", dst.x);
    l.setAttribute("y1", dst.y);
    svg1.appendChild(l);
    return l;
};

const drawMap = function (svg1, map, onClickPosition) {
    map.roads.forEach(function (road, idx, tab) {
        let src = map.cells[road.src]
        let dst = map.cells[road.dst]
        let r = drawLine(svg1, src, dst);
    });
    Object.keys(map.cells).forEach(function (key, idx, tab) {
        let cell = map.cells[key]
        let c = drawCircle(svg1, cell);
        if (onClickPosition != null) {
            c.onclick = function (e) {
                console.log("onClickPosition")
                onClickPosition(key);
            };
        }
    });
    return map;
};

const patchWithArmies = function (svg1, map, armies, onClickArmy) {
    const pad = 5;
    armies.forEach(function (army, idx, tab) {
        let c = map.cells[army.cell];
        let cDom = document.getElementById(army.cell);
        let a = document.createElementNS("http://www.w3.org/2000/svg", "use");
        a.setAttribute("fill", "black");
        a.setAttribute("stroke", "black");
        a.setAttribute("stroke-width", 2);
        let r = cDom.getAttribute("r");
        a.setAttribute("x", c.x - r + pad);
        a.setAttribute("y", c.y - r + pad);
        a.setAttributeNS("http://www.w3.org/1999/xlink", "href", "#shield");
        svg1.appendChild(a);
        if (onClickArmy != null) {
            c.onclick = function (e) {
                onClickArmy(army.cell, army.id);
            };
        }
    });
    return map;
};

const patchWithCities = function (svg1, name, map, onClickCity, onClickArmy) {
    return fetch('/map/cities?id=' + name)
        .then(response => {
            return response.json();
        })
        .then(data => {
            Object.keys(data).forEach(function (key, idx, tab) {
                let city = data[key];
                let c = svg1.getElementById(city.cell);
                c.setAttribute("r", 23);
                if (onClickCity != null) {
                    c.onclick = function (e) {
                        console.log("onClickCity")
                        onClickCity(city.cell, city.id, city.name);
                    };
                }
            })
            return map
        })
};

const drawMapWithCities = function (svg1, name, onClickPosition, onClickCity) {
    // Step 1: draw the map itself, define the graph with nodes (a.k.a cells) and
    // vertices (a.k.a. roads).
    return fetch('/map/region?id=' + name)
        .then(response => {
            return response.json();
        })
        .then(map => {
            drawMap(svg1, map, onClickPosition);
            // Step 2: alter the representation of the cells with a City
            return patchWithCities(svg1, name, map, onClickCity);
        });
};

const drawMapWithArmies = function (svg1, name, armies, onClickArmy) {
    return drawMapWithCities(svg1, name)
        .then(map => {
            return patchWithArmies(svg1, map, armies, onClickArmy);
        });
}

const hightlightCell = function (svg1, id) {
    let c = svg1.getElementById(id);
    console.log(c);
    if (c != null) {
        c.setAttribute("class", c.getAttribute("class") + " here");
    }
}
