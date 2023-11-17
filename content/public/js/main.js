function headerNavListeners() {
    const navToggle = document.querySelector('header>nav>.nav-toggle');
    const nav = document.querySelector('header>nav');
    navToggle.addEventListener('click', function (e) {
        console.debug('click', e);
        nav.classList.toggle('open');
    });
    const links = document.querySelectorAll('header>nav>a');
    console.debug('links', links);
    links.forEach(function (link) {
        console.debug('link', link);
        link.addEventListener('click', function (e) {
            console.debug('click', e);
            e.preventDefault();
            links.forEach(function (link) {
                setNotActive(link);
            })
            setActive(link);
        });
    });
}


function setActive(link) {
    link.classList.add('active');
    link.tabIndex = -1;
}


function setNotActive(link) {
    link.classList.remove('active');
    link.tabIndex = 0;
}

function init() {
    console.debug('init');
    headerNavListeners();
}



init();