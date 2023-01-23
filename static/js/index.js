(function () {
    let nameField = document.getElementById("my-name");
    doTheHarlemShake(nameField);
})();

function shake(element) {

    let magnitude = 20;
    let counter = 1;
    let shakes = 20;
    let decrease = magnitude/shakes;
  
    let randInt = (min, max) => {
      return Math.floor(Math.random() * (max - min + 1)) + min;
    };

    let _shake = function () {
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
    
    let animate = function (idx) {
        if (strings.length === idx) {
            return;
        }
        shake(element);
        element.textContent = strings[idx];
        setTimeout(animate, 1000, idx+1);
    }
    setTimeout(animate, 1000, 0);
}