// input is a GeoJSON
// need to know the center (potentially)

const mapId = 'map';
const centerLon = 137;
const centerLat = 38;
const defaultZoom = 5;

var map = new maplibregl.Map({
    container: mapId,
    center: [centerLon, centerLat],
    zoom: defaultZoom
});

const sourceId = 'country'

map.addSource(sourceId, {
    type: 'geojson',
    data: geojson,
    generateId: true
});

const divisionId = 'division';
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

const extractName = (feat) => ({
    name: feat.properties.name_ja,
    nameHTML: `<ruby>feat.properties.name_ja<rt>feat.properties.hiragana</rt></ruby>`
})

// add hiragana
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

// game data is a list of objects with a .nameHTML property

class Game {
    constructor(data, labelEl, clicksEl) {
        this.data = data;
        this.label = labelEl;
        this.updateLabel();
        this.erroneousClicks = 0;
        this.clicksEl = clicksEl;
    }

    updateLabel() {
        this.label.innerHTML = data[0].nameHTML;
    }

    victory() {
        localStorage.setItem(Date.now(), this.erroneousClicks);
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

const game = new Game(quiz, document.getElementById('current'), document.getElementById('error-clicks'));
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
    const clickedon = e.features[0].properties.name_ja;
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