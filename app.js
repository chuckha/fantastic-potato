var map = new maplibregl.Map({
    container: 'map',
    center: [137, 38],
    zoom: 5
});

map.addSource('japan', {
    type: 'geojson',
    data: prefectures,
    generateId: true
});

map.addLayer({
    'id': 'prefectures',
    'type': 'fill',
    'source': 'japan',
    'layout': {},
    'paint': {
        'fill-color': [
            'case',
            ['boolean', ['feature-state', 'correct'], false],
            'rgba(70, 192, 138, 0.19)',
            ['boolean', ['feature-state', 'clicked'], false],
            'rgba(182, 92, 138, 0.19)',
            'rgba(222, 222, 222, 0.19)'
        ],
        'fill-outline-color': '#111',

    }
});

var source = map.getSource('japan');
console.log(source)


// add hiragana
const data = source._data.features.map(feat => ({ja: feat.properties.name_ja, fu: feat.properties.hiragana}))
const quiz = shuffle(data)
console.log(quiz)

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

class Game {
    constructor(data, label) {
        this.data = data;
        this.label = label;
        this.updateLabel();
    }

    updateLabel() {
        this.label.innerHTML = '<ruby>'+data[0].ja+'<rt>'+data[0].fu+'</rt></ruby>';
    }

    correct(input) {
        if (input === this.data[0].ja) {
            this.data.shift();
            this.updateLabel();
            return true;
        }
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
            source: 'japan',
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
            source: 'japan',
            id:id
        }, {
            clicked: false
        })
    }
}

const game = new Game(quiz, document.getElementById('current'));
const filler = new Filler();
const correct = [];

map.on('click', 'prefectures', function (e) {
    // ignore clicks on already correctly guessed prefectures
    if (correct.includes(e.features[0].id)) {
        return;
    }
    const clickedon = e.features[0].properties.name_ja;
    if (game.correct(clickedon)) {
        correct.push(e.features[0].id)
        map.setFeatureState({
            source: 'japan',
            id: e.features[0].id
        }, {
            correct: true
        })
        filler.reset();
    } else {
        filler.fill(e.features[0].id)
    }
});
// game?
// get an array of all prefectures
// shuffle array

// for each item in array
// update the HTML
// - user clicks on wrong prefecture
// -- prefecture turns red
// - user clicks on correct prefecture
// -- prefecture turns green
// -- all red prefectures return to original color
// -- next item