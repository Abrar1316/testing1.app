function psOpenModel() {
  var modal = document.getElementById("psmyModel");
  modal.style.display = "block";
}

function psCloseModel() {
  var modal = document.getElementById("psmyModel");
  modal.style.display = "none";
}
window.onclick = function(event) {
  if (event.target == modal) {
    psCloseModel();
  }
}