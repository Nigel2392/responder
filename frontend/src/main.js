import * as _ from "./wasm/wasm_exec";

var go = new Go();

// Fetch and instantiate the WebAssembly file
function FetchAndInstantiate(url, importObject) {
    return fetch(url).then(response => response.arrayBuffer()
    ).then(bytes => WebAssembly.instantiate(bytes, importObject)
    ).then(results => results.instance );
};

let wasm_link_elem = document.getElementById("wasm_main");

let mod = FetchAndInstantiate(wasm_link_elem.href, go.importObject);

document.addEventListener("DOMContentLoaded", function(event) {
    mod.then(function(instance) { go.run(instance); })
});

export var app = document.getElementById("app");
