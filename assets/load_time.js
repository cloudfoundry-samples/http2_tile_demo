document.getElementById("http").innerHTML = performance.getEntries()[0]?.nextHopProtocol;
window.addEventListener("load", function() {
    window.setTimeout(function() {
        document.getElementById("time").innerHTML = (window.performance.timing.loadEventEnd - window.performance.timing.navigationStart);
    }, 0);
});
