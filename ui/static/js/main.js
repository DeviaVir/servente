var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href').split("/")[1] == window.location.pathname.split("/")[1]) {
		link.classList.add("live");
		break;
	}
	// services <=> service (hackish)
	if (link.getAttribute('href').split("/")[1] == "services" && window.location.pathname.split("/")[1] == "service") {
		link.classList.add("live");
		break;
	}
}

function addField(el){
	const [settings, attributes] = discoverTypes();

	var parent = el.parentElement;
	for (var i = 0; i < parent.childNodes.length; i++) {
		if (parent.childNodes[i].className !== undefined && parent.childNodes[i].className.indexOf("input-container") > -1) {
			var container = parent.childNodes[i];
			var type = "default";
			if (container.attributes.length > 0 && container.attributes["attr-name"] !== undefined) {
				type = container.attributes["attr-name"].nodeValue;
			}

			var itemContainer = document.createElement("div")
			itemContainer.className = "clear item"

			var inputIdentifier = document.createElement("input");
			inputIdentifier.type = "text";
			inputIdentifier.placeholder = "Identifier";
			inputIdentifier.name = `${type}[identifier][]`;

			var selectType = document.createElement("select");
			selectType.name = `${type}[type][]`;

			var selectedTypes = (type == "settings" ? settings : attributes)
			for (ii = 0; ii < selectedTypes.length; ii++) {
				var option = document.createElement("option");
				option.value = selectedTypes[ii];
				option.innerHTML = selectedTypes[ii];
				selectType.appendChild(option);
			}

			itemContainer.appendChild(inputIdentifier);
			if (type == "settings") {
				var inputValue = document.createElement("input");
				inputValue.type = "text";
				inputValue.placeholder = "Value";
				inputValue.name = `${type}[value][]`;
				itemContainer.appendChild(inputValue);
			}
			itemContainer.appendChild(selectType);

			container.appendChild(itemContainer)
		}
	}
}

function discoverTypes() {
	// settings
	var settings;
	var sTholder = document.getElementById("settingsTypes");
	if (sTholder !== undefined && sTholder.attributes.length > 0 && sTholder.attributes["value"] !== undefined) {
		settings = sTholder.attributes["value"].nodeValue.split(",");
	}

	// attributes
	var attributes;
	var aTholder = document.getElementById("attributesTypes");
	if (aTholder !== undefined && aTholder.attributes.length > 0 && aTholder.attributes["value"] !== undefined) {
		attributes = aTholder.attributes["value"].nodeValue.split(",");
	}

	return [settings, attributes];
}
