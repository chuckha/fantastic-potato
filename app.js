const canvas = document.getElementById('map');
const width = canvas.clientWidth;
const height = canvas.clientHeight;
/** @type {CanvasRenderingContext2D} */
const ctx = canvas.getContext('2d');
console.log(canvas.clientWidth, canvas.clientHeight);

// this is a list of features (islands and shit)
for (let i=0; i < data.features[0].geometry.coordinates.length; i++) {
    const feature = data.features[0].geometry.coordinates[i]

    ctx.beginPath();
    const start2 = latLonToOffsets(feature[0][0][1], feature[0][0][0], width, height);
    ctx.moveTo(start2.x, start2.y);
    for (let j=1; j < feature[0].length; j++) {
        const point2 = latLonToOffsets(feature[0][j][1], feature[0][j][0], width, height);
        console.log(point2.x)
        ctx.lineTo(point2.x, point2.y);
    }
    ctx.stroke();
}
console.log("stroked")

function getCursorPosition(canvas, event) {
    const rect = canvas.getBoundingClientRect()
    const x = event.clientX - rect.left
    const y = event.clientY - rect.top
    console.log("x: " + x + " y: " + y)
}

canvas.addEventListener('mousedown', function(e) {
    getCursorPosition(canvas, e)
})

function latLonToOffsets(latitude, longitude, mapWidth, mapHeight) {
    const FE = -110; // false easting
    const FN = -30;
    const radius = mapWidth / (Math.PI/4 );

    const latRad = degreesToRadians(latitude + FN);
    const lonRad = degreesToRadians(longitude + FE);

    const x = lonRad * radius;

    const yFromEquator = radius * Math.log(Math.tan(Math.PI / 4 + latRad / 2));
    const y = mapHeight / 2 - yFromEquator;
    console.log("mapheight", mapHeight, "yFromEquator", yFromEquator )
    return { x, y };
  }

  function degreesToRadians(deg) {
      return deg * (Math.PI / 180);
  }