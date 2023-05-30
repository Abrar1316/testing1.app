function updatePinProject(userId, pName, value) {
    console.log(userId)
    fetch(`/${userId}/pinproject/${value}/${pName}`, {
      method: 'GET',
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Failed to update pinning project');
        }
        window.location.href = `/${userId}/days/30`;
      })
      .catch(error => {
        console.error(error);
      });
  }
  