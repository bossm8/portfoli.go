// Copyright (c) 2023, Boss Marco <bossm8+portfoligo@hotmail.com>
// All rights reserved.

(function () {
    const form = document.getElementById('contact-form');
    const button = document.getElementById("runaway");

    form.addEventListener('submit', validate);

    function validate(event) {
        if (!form.checkValidity()) {
            event.preventDefault()
            event.stopPropagation()
            button.addEventListener('mouseover', run);
            button.addEventListener('click', run);
        }
        this.classList.add('was-validated')
    }

    const randInt = (max) => {
      return Math.floor(Math.random() * (max + 1));
    };

    function run() {
        if (form.checkValidity()) {
            button.removeEventListener('click', run);
            button.removeEventListener('mouseover', run);
            if (button.style.position === 'fixed') {
                // stop any in-flight dodge, otherwise its still-running
                // animation frames keep overwriting left/top right after
                // we reset them below
                anime.remove(button);
                button.style.position = '';
                button.style.left = '';
                button.style.top = '';
                button.style.zIndex = '';
            }
            return
        }

        if (button.style.position !== 'fixed') {
            // anchor it exactly where it already visually is before undocking
            // it from the document flow, otherwise switching straight to
            // fixed makes it jump to its raw in-flow position, which can
            // land outside the viewport entirely
            const rect = button.getBoundingClientRect();
            button.style.position = 'fixed';
            button.style.zIndex = '60';
            button.style.left = `${rect.left}px`;
            button.style.top = `${rect.top}px`;
        }

        // roam the whole visible viewport, not just its original form column
        const maxLeft = Math.max(0, window.innerWidth - button.offsetWidth);
        const maxTop = Math.max(0, window.innerHeight - button.offsetHeight);
        const left = randInt(maxLeft);
        const top = randInt(maxTop);
        anime({
            targets: button,
            left: `${left}px`,
            top: `${top}px`,
            easing: "easeOutCirc"
        }).play();
    }
})();
