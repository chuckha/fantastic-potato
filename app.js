var map = new maplibregl.Map({
    container: 'map',
    center: [137, 38],
    zoom: 3
});

    map.addSource('japan', {
        type: 'geojson',
        data: prefectures
    });
    map.addLayer({
        'id': 'prefectures',
        'type': 'fill',
        'source': 'japan',
        'layout': {},
        'paint': {
            'fill-color': '#088',
            'fill-opacity': 0.8
        }
    });
