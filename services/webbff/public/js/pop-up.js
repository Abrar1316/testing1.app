var modal = document.getElementById("myModal");

function openModal() {
  document.getElementById("myModal").style.display = "block";
}
function closeModal() {
  document.getElementById("myModal").style.display = "none";
}

window.onclick = function(event) {
  if (event.target == modal) {
    closeModal();
  }
}

async function deleteProject(userId, projectId) {
  var confirmBox = document.createElement("div");
  confirmBox.classList.add("confirm-box");
  confirmBox.innerHTML = `<p><b>Are you sure you want to delete this project?</b></p>
                          <button class="confirm-btn">Confirm</button>
                          <button class="cancel-btn">Cancel</button>`;
  document.body.appendChild(confirmBox);

  return new Promise(resolve => {
    document.querySelector(".confirm-btn").addEventListener("click", () => {
      document.body.removeChild(confirmBox);
      resolve(true);
      fetch(`/${userId}/projects/${projectId}/delete/`, {
        method: "GET",
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to delete project");
          }
          // redirect to the projects page for the given user ID
          window.location.href = `/${userId}/projects`;
        })
        .catch((error) => {
          console.error(error);
          // handle the error and display a message to the user
        });
    });

    document.querySelector(".cancel-btn").addEventListener("click", () => {
      document.body.removeChild(confirmBox);
      resolve(false);
    });
  });
}

function updatePinProject(userId, projectId, value) {
  console.log(userId)
  fetch(`/${userId}/projects/${projectId}/pinproject/${value}/`, {
    method: 'GET',
  })
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to update pinning project');
      }
      window.location.href = `/${userId}/projects`;
    })
    .catch(error => {
      console.error(error);
    });
}

async function updateActiveProject(userId, projectId, active){
var confirmBox = document.createElement("div");
  confirmBox.classList.add("confirm-box");
  confirmBox.innerHTML = `<p><b>Are you sure?</b></p>
                          <button class="confirm-btn">Confirm</button>
                          <button class="cancel-btn">Cancel</button>`;
  document.body.appendChild(confirmBox);

  return new Promise(resolve => {
    document.querySelector(".confirm-btn").addEventListener("click", () => {
      document.body.removeChild(confirmBox);
      resolve(true);
      fetch(`/${userId}/projects/${projectId}/activeproject/${active}/`, {
        method: "GET",
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to deactivate project");
          }
          // redirect to the projects page for the given user ID
          window.location.href = `/${userId}/projects`;
        })
        .catch((error) => {
          console.error(error);
          // handle the error and display a message to the user
        });
    });

    document.querySelector(".cancel-btn").addEventListener("click", () => {
      document.body.removeChild(confirmBox);
      resolve(false);
    });
  });

}

function deleteAwsCredentials(userId, projectId) {
  const deleteBtn = document.getElementById("disconnect-btn");
  deleteBtn.disabled = true; // disable the button

  // create confirmation box
  const confirmBox = document.createElement("div");
  confirmBox.classList.add("confirm-box");
  confirmBox.innerHTML = `<p><b>Are you sure you want to disconnect?</b></p>
                            <button class="confirm-btn">Confirm</button>
                            <button class="cancel-btn">Cancel</button>`;
  document.body.appendChild(confirmBox);

  return new Promise(resolve => {
    document.querySelector(".confirm-btn").addEventListener("click", () => {
      document.body.removeChild(confirmBox);
      // send delete request
      fetch(`/${userId}/projects/${projectId}/deleteAwsCreds/`, {
        method: "DELETE"
      })
      .then(response => {
        if (!response.ok) {
          throw new Error("Failed to delete credentials");
        }
        // redirect to the home page
        window.location.href = window.location.href;
      })
      .catch(error => {
        console.error(error);
        // handle the error and display a message to the user
      });
      resolve(true);
    });

    document.querySelector(".cancel-btn").addEventListener("click", () => {
      document.body.removeChild(confirmBox);
      deleteBtn.disabled = false; // enable the button
      resolve(false);
    });
  });
}