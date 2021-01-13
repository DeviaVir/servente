var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href').split("/")[1] == window.location.pathname.split("/")[1]) {
		link.classList.add("live");
		break;
	}
}
