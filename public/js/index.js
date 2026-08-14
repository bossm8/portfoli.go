// Copyright (c) 2023, Boss Marco <bossm8+portfoligo@hotmail.com>
// All rights reserved.

(function () {
    let nameField = document.getElementById("my-name");
    var running = false;
    var hasPlayed = false;
    nameField.addEventListener("mouseover", () => {
        if (hasPlayed)
            return;
        hasPlayed = true;
        doTheHarlemShake(nameField);
    });

    let originalHTML = nameField.innerHTML;
    let strings = [
        "Hi, my name is, what?",
        "My name is, who?",
        "My name is ...",
        "Chka-Chka ...",
    ];

    var doTheHarlemShake = function(element) {
        // This method takes an already running flag, which determines if the animation
        // is already running (whaat...?) The first time it get's called this will be false
        // because it is set to true only after anmiate was called the first time with this value.
        // Once animate finishes all animations of the current run (handled all texts) it will set
        // the flag to false again so another mouse enter triggers the animation again.
        let animate = function (idx, alreadyRunning) {
            if (alreadyRunning)
                return;
            shake(element);
            if (strings.length === idx) {
                // restore via innerHTML so markup (e.g. <strong>) inside the
                // original name isn't lost - textContent would strip it
                element.innerHTML = originalHTML;
                running = false;
                return;
            }
            element.textContent = strings[idx];
            setTimeout(animate, 1000, idx+1, alreadyRunning);
        }
        animate(0, running);
        running = true;
    };
})();

function shake(element) {

    let magnitude = 20;
    let counter = 1;
    let shakes = 20;
    let decrease = magnitude/shakes;

    const randInt = (min, max) => {
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
