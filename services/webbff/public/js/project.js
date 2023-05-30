if (typeof buttons === 'undefined') {
  const buttons = document.querySelectorAll('[id^="pin-button-"]');
  buttons.forEach(function (button) {
    const projectId = button.id.split('-').pop();
    const pinValue = document.getElementById(`pin-value-${projectId}`);

    button.addEventListener('mouseenter', function () {
      pinValue.style.opacity = 1;
    });

    button.addEventListener('mouseleave', function () {
      pinValue.style.opacity = 0;
    });
  });
}


function openUpdateModal(id, userId) {
  var modal = document.getElementById("updateModal"+id);
  var form = document.getElementById("updateprojectsform"+id);
  form.setAttribute("method", "POST");
  form.action = "/" + userId + "/projects/" + id + "/settings/";
  modal.style.display = "block";
}

function closeUpdateModal(id) {
  document.getElementById("updateModal"+id).style.display = "none";
}

// window.onclick = function (event) {
//   if (event.target == updateModal) {
//     closeUpdateModal();
//   }
// }


