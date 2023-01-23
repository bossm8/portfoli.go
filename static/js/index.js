(function () {
    let nameField = document.getElementById("my-name");
    doTheHarlemShake(nameField);
})();

function shake(element) {

    var magnitude = 20;
    var counter = 1;
    var shakes = 20;
    var decrease = magnitude/shakes;
  
    var randInt = (min, max) => {
      return Math.floor(Math.random() * (max - min + 1)) + min;
    };

    var _shake = function () {
        if (counter < shakes) {
            magnitude -= decrease;

            let x = randInt(-magnitude, magnitude);
            let y = randInt(-magnitude, magnitude);

            element.style.transform = 'translate(' + x + 'px, ' + y + 'px)';

            counter++;

            requestAnimationFrame(_shake);
        } else {
            element.style.transform = 'translate(0,0)';
        }
    }
    _shake();
}

function doTheHarlemShake(element) {
    let strings = [
        "Hi, my name is, what?",
        "My name is, who?",
        "My name is ...",
        "Chka-Chka ...",
        element.textContent
    ];
    
    var animate = function (idx) {
        if (strings.length === idx) {
            return;
        }
        shake(element);
        element.textContent = strings[idx];
        setTimeout(animate, 1000, idx+1);
    }
    animate(0);
}