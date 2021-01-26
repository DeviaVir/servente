var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href').split("/")[1] == window.location.pathname.split("/")[1]) {
		link.classList.add("live");
		break;
	}
}

function addField(el){
	const [settings, attributes] = discoverTypes()

	console.log(settings)
	console.log(attributes)

	var parent = el.parentElement
	for (var i = 0; i < parent.childNodes.length; i++) {
		if (parent.childNodes[i].className !== undefined && parent.childNodes[i].className.indexOf("input-container") > -1) {
			var container = parent.childNodes[i]
			var type = "default"
			if (container.attributes.length > 0 && container.attributes["attr-name"] !== undefined) {
				type = container.attributes["attr-name"].nodeValue
			}

			var inputIdentifier = document.createElement("input");
			inputIdentifier.type = "text";
			inputIdentifier.name = type + "[identifier][]"

			container.appendChild(inputIdentifier);
		}
	}
}

function discoverTypes() {
	// settings
	var settings
	var sTholder = document.getElementById("settingsTypes")
	if (sTholder !== undefined && sTholder.attributes.length > 0 && sTholder.attributes["value"] !== undefined) {
		settings = sTholder.attributes["value"].nodeValue.split(",")
	}

	// attributes
	var attributes
	var aTholder = document.getElementById("attributesTypes")
	if (aTholder !== undefined && aTholder.attributes.length > 0 && aTholder.attributes["value"] !== undefined) {
		attributes = aTholder.attributes["value"].nodeValue.split(",")
	}

	return [settings, attributes]
}
