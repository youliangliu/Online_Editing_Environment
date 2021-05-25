(function ( doc ) {

var editElement = document.getElementById('edit');
editElement.contentEditable = true;
editElement.focus();


var originalContent = editElement.innerHTML;
console.log(originalContent)

function onChange() {
  if (editElement.innerHTML !== originalContent ) {
      originalContent = editElement.innerHTML;
    console.log("something changed");
  }
}

editElement.onkeyup = function () {
  onChange();
};
editElement.onchange = function () {
  onChange();
};
editElement.onpaste = function () {
  onChange();
};
editElement.onclick = function () {
  onChange();
};


})( document );