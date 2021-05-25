var editElement = document.getElementById('edit');
editElement.contentEditable = true;
editElement.focus();


var originalContent = editElement.innerHTML;

function onChange() {
  if (editElement.innerHTML !== originalContent ) {
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

/*
$('edit').on("change keyup paste click", function(){
  console.log("Wassup");
})
*/