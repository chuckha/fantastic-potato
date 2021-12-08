const country = "japan";
const tl = "japanese";

const mapId = 'map';
const sourceId = 'country';
const divisionId = 'division';
var map = new maplibregl.Map({
    container: mapId,
});

fetch(`/api/get-country-data?country=${country}&target_language=${tl}`)
.then(resp => resp.json())
.then(callback)

const el = document.getElementById("korea")
el.addEventListener("click", () => {
    clear();
    init();
    fetch(`/api/get-country-data?country=korea&target_language=korean`)
    .then(resp => resp.json())
    .then(callback)
});

function init() {
    map = new maplibregl.Map({
        container: mapId,
    });
}

function clear() {
    document.getElementById('error-clicks').innerText = 0;
    map.removeLayer(divisionId);
    map.removeSource(sourceId);
    map.remove();
}

function callback(obj) {
    const geojson = JSON.parse(obj.geojson);
    map.setZoom(obj.default_zoom);
    map.setCenter([obj.center_lon, obj.center_lat]);

    map.addSource(sourceId, {
        type: 'geojson',
        data: geojson,
        generateId: true
    });

    const correctColor = 'rgba(70, 192, 138, 0.19)';
    const incorrectColor = 'rgba(182, 92, 138, 0.19)';
    const defaultColor = 'rgba(222, 222, 222, 0.19)';
    const outlineColor = '#111';

    map.addLayer({
        'id': divisionId,
        'type': 'fill',
        'source': sourceId,
        'layout': {},
        'paint': {
            'fill-color': [
                'case',
                ['boolean', ['feature-state', 'correct'], false],
                correctColor,
                ['boolean', ['feature-state', 'clicked'], false],
                incorrectColor,
                defaultColor
            ],
            'fill-outline-color': outlineColor,
        }
    });

    var source = map.getSource(sourceId);

    // fetch the geojson
    // do the game stuff client side
    // report stats back to the backend

    // What do i need to render a map of Japan? The GeoJSON.
    // What does the frontend do with it? all this stuff.
    // Does the frontend need to do anything else?

    // Gets a list of feature names and furigana
    // shuffles the list
    // Game starts

    // add hiragana
//    console.log(source._data.features)
    const data = source._data.features.map(extractName)
    const quiz = shuffle(data)

    function shuffle(data) {
        var tmp;
        for (let i = data.length-1; i >= 0; i--) {
            var j = getRandomInt(i)
            tmp = data[j];
            data[j] = data[i];
            data[i] = tmp;
        }
        return data
    }

    function getRandomInt(max) {
        return Math.floor(Math.random() * max);
    }

    const game = new Game(country, quiz, document.getElementById('current'), document.getElementById('error-clicks'));
    const filler = new Filler();
    const correct = [];

    map.on('click', divisionId, function (e) {
        // ignore clicks on already correctly guessed prefectures
        if (correct.includes(e.features[0].id)) {
            return;
        }
        // ignore clicks on already clicked prefectures
        if (e.features[0].state.clicked === true) {
            return;
        }
        const clickedon = e.features[0].properties.name_tl;
        if (game.correct(clickedon)) {
            correct.push(e.features[0].id)
            map.setFeatureState({
                source: sourceId,
                id: e.features[0].id
            }, {
                correct: true
            })
            filler.reset();
        } else {
            filler.fill(e.features[0].id)
        }
    });
}

const extractName = (feat) => {
    let name = `<ruby>${feat.properties.name_tl}`
    if (feat.properties.ruby_text) {
        name = name + `<rt>${feat.properties.ruby_text}</rt>`
    }
    name = name + '</ruby>'
    return {
        name: feat.properties.name_tl,
        nameHTML: name
    }
}

class Filler {
    constructor() {
        this.ids = [];
    }
    fill(id) {
        this.ids.push(id)
        map.setFeatureState({
            source: sourceId,
            id: id
        }, {
            clicked: true
        })
    }
    reset() {
        this.ids.forEach(id => this.unfill(id));
        this.ids.length = 0;
    }
    unfill(id) {
        map.setFeatureState({
            source: sourceId,
            id:id
        }, {
            clicked: false
        })
    }
}

class Game {
    constructor(country, data, labelEl, clicksEl) {
        this.country = country
        this.data = data;
        this.label = labelEl;
        this.updateLabel();
        this.erroneousClicks = 0;
        this.clicksEl = clicksEl;
    }

    updateLabel() {
        this.label.innerHTML = this.data[0].nameHTML;
    }

    victory() {
        localStorage.setItem(Date.now(), {
            "erroneous clicks": this.erroneousClicks,
            "country": this.country
        });
        this.label.innerHTML = 'Good job! Refresh to play again!';
    }

    updateErroneousClicks() {
        this.erroneousClicks++;
        this.clicksEl.innerText = this.erroneousClicks;
    }

    correct(input) {
        if (input === this.data[0].name) {
            this.data.shift();
            if (this.data.length == 0) {
                this.victory();
                return true;
            }
            this.updateLabel();
            return true;
        }
        this.updateErroneousClicks();
        return false;
    }
}



// domain 1: application supports email auth
// domain 2: statistics
//           GameStarted<>
//           WrongGuessEvent<GameID: id, Finding: xyz, Guessed: abc>
//           IdentifiedEvent<GameID: id, Identified: abc>
//           GameFinished<>
//           With these events we can ask questions like:
//           - average misses per game
//           - which states are confused for other states
//           - games played, etc
// domain 3: game data
//           GetData<Country: xyz>
//           core here is an app that gets geojson data from <somewhere>
//           it turns that geojson data into quiz-friendly data and returns it
//           what is quiz friendly data?

