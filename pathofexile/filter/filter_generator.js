// Generates item filter code for Path of Exile game.

// Open https://poe.ninja/challenge/base-types. Paste the code below in console (F12) and it will generate filter rules. Change these two values below to fit your needs:
var minCountConfidence = 10; // poe.ninja considers < 5 as low confidence. < 10 as medium. >= 10 normal
var cheapBasesMinPrice = 6;
var expensiveBasesMinPrice = 10;

// I'm sorry you have to read this messy code, I was extremely sleep deprived and depressed when I wrote it. But hey! it works!
fetch(`https://poe.ninja/api/data/itemoverview?league=Delve&type=BaseType`).then(function(response){return response.json();}).then(function(data){
var items = data.lines;
console.log(
	generateRules(items, cheapBasesMinPrice, expensiveBasesMinPrice, "0 255 0", "PlayAlertSound 16 300") +
	"\n\n" + 
	generateRules(items, expensiveBasesMinPrice, -1, "255 0 0", `CustomAlertSound "LakadMatataaag_NormalinNormalin.mp3"`));
}).catch(function(err){console.log(err);});

// minChaosValue is inclusive. We use >= when comparing it
// maxChaosValueExclusive is compared using < only. Pass -1 for no limit
function generateRules(items, minChaosValue, maxChaosValueExclusive, borderColor, soundLine) {
	// If the base is worth picking up even without elder/shaper, we dont need to repeat it on the filter with shaper/elder
	// If the base is worth picking up at iLvl 82, we don't need to repeat it at 83+ on filter. Put it 82+
	var filter = {}; // "Lion Pelt:":89, "Lion Pelt:Shaper":87,
	for (var item of items.values()) {
		if (item.chaosValue < minChaosValue || item.count < minCountConfidence || 
			(maxChaosValueExclusive > -1 && item.chaosValue >= maxChaosValueExclusive)) continue;
		item.variant = item.variant || "";
		var filterKey = item.name + ":";
		var filterKeyWithVariant = filterKey + item.variant;
		var ilvlOnFilter = parseInt(filter[filterKey] || 999);
		var ilvlOnFilterWithVariant = parseInt(filter[filterKeyWithVariant] || 999);
		// Base item, not Elder or Shaper
		if (item.variant == "" && item.levelRequired < ilvlOnFilter) {
			filter[filterKey] = item.levelRequired; // Add base to filter
			// If current Shaper or Elder filter is higher or equal than this base, delete them since we can just filter by all
			if ((filter[filterKey+"Elder"] || 999) <= item.levelRequired) delete filter[filterKey+"Elder"];
			if ((filter[filterKey+"Shaper"] || 999) <= item.levelRequired) delete filter[filterKey+"Shaper"];
		}
		// Shaper or Elder bases
		if (item.variant != "" && item.levelRequired < ilvlOnFilterWithVariant) {
			filter[filterKeyWithVariant] = item.levelRequired; // Add variant base to filter
		}
	}
	// Now filter looks like: {"Steel Ring:Elder:": 83, "Steel Ring:Shaper:": 82, "Thicket Bow:Elder:": 82} and we'll reorganize it 
	// to facilitate making filter rules: {"82":{"Elder": ["Thicket Bow"], "Shaper": ["Steel Ring"]},"83":"Elder":["Steel Ring"]}}
	var sorted = {};
	for (var base in filter) {
		var ilvl = filter[base];
		var [name, variant] = base.split(':');
		variant = variant || "Normal";
		if (!sorted.hasOwnProperty(ilvl)) sorted[ilvl] = {};
		if (!sorted[ilvl].hasOwnProperty(variant)) sorted[ilvl][variant] = [];
		sorted[ilvl][variant].push('"' + name + '"');
	}
	var rules = `
# Highlighting bases that are worth at least ${minChaosValue}c`;
	for (var ilvl in sorted) {
		for (var variant in sorted[ilvl]) {
			var variantLine = (variant == "Normal" ? "" : (variant == "Elder" ? "\nElderItem True" : "\nShaperItem True"));
			rules = rules + `
Show${variantLine}
ItemLevel >= ${ilvl}
Rarity <= Rare
BaseType ${sorted[ilvl][variant].join(" ")}
SetTextColor 255 255 255
SetBorderColor ${borderColor}
${soundLine}
SetBackgroundColor 0 128 0 96
SetFontSize 50
DisableDropSound
MinimapIcon 0 Green Triangle`
		}
	}
	// Now that we have clean structured data, we can generate the filter rules
	return rules;
}